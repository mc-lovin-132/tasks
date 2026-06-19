package app

import (
	"context"
	"net"
	"tasks/config"
	"tasks/internal/infrastructure/delivery/handlers"
	"tasks/internal/infrastructure/delivery/interceptors"
	"tasks/internal/infrastructure/repository"
	"tasks/internal/service"
	"tasks/pb"

	// "github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type App struct {
	cfg    *config.Config
	logger *zap.Logger
}

func New(cfg *config.Config, logger *zap.Logger) *App {
	return &App{
		cfg:    cfg,
		logger: logger,
	}
}

func (a *App) Start(ctx context.Context) error {
	db, err := sqlx.Connect("postgres", a.cfg.DSN())
	if err != nil {
		return err
	}
	defer func() {
		err := db.Close()
		if err != nil {
			a.logger.Error("err while closing db connection", zap.Error(err))
			return
		}
		a.logger.Info("db connection successfuly closed")
	}()
	err = repository.RunMigrations(db.DB)
	if err != nil {
		return err
	}
	a.logger.Info("successfuly init db connection")
	repo := repository.New(db)
	srvc := service.New(repo, &userServiceMock{}, &eventSendenerMock{})
	cron := service.NewDeadlineCrone(&eventSendenerMock{}, repo)
	handler := handlers.New(srvc)

	lis, err := net.Listen("tcp", a.cfg.Addr())
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.NewLoggingInterceptor(a.logger)),
	)
	pb.RegisterTaskServiceServer(grpcServer, handler)
	a.logger.Info("gRPC server listening", zap.String("port", a.cfg.Port))

	errG, gCtx := errgroup.WithContext(ctx)

	errG.Go(func() error {
		<-gCtx.Done()
		grpcServer.GracefulStop()
		return gCtx.Err()
	})

	errG.Go(func() error {
		err := cron.Start(gCtx)
		if err != nil {
			a.logger.Error("cron run failed", zap.Error(err))
			return err
		}
		a.logger.Info("cron started successfuly")
		return nil
	})

	errG.Go(func() error {
		if err := grpcServer.Serve(lis); err != nil {
			a.logger.Error("grpc server run failed", zap.Error(err))
			return err
		}
		return nil
	})

	return errG.Wait()
}
