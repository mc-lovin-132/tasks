package service

import (
	"context"
	"fmt"
	"time"
)

type deadlineCron struct {
	eventSendener eventSendener
	repo          repository
	interval      time.Duration
}

func NewDeadlineCrone(eventSendener eventSendener, repo repository) *deadlineCron {
	return &deadlineCron{
		eventSendener: eventSendener,
		repo:          repo,
	}
}

func (c *deadlineCron) isDeadlineSoon(deadline time.Time) bool {
	remaining := time.Until(deadline)
	if remaining <= c.interval {
		return true
	}
	return false
}

func (c *deadlineCron) Start(ctx context.Context) error {
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			timeoutCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
			defer cancel()
			tasks, err := c.repo.List(timeoutCtx, nil, nil, nil)
			if err != nil {
				fmt.Println(err.Error())
				return err
			}
			for _, task := range tasks {
				if c.isDeadlineSoon(task.Deadline) {
					err := c.eventSendener.DeadlineSoon(task, c.interval)
					if err != nil {
						return err
					}
				}
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
