package application

import (
	"context"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"io"
	"net/http"
	"os"
	"os/signal"
	"quic_upload/api/internal/config"
	"quic_upload/api/internal/config/certificate"
	"quic_upload/api/internal/controller"
	"quic_upload/api/internal/service"
	"quic_upload/api/internal/storage"
	"quic_upload/api/pkg/logging"
	"quic_upload/api/pkg/s3Storage"
	"syscall"
	"time"
)

type Application struct {
	cfg         *config.Config
	router      http.Handler
	http3Server *http3.Server
	s3          S3Storage
	logger      *logging.Logger
}

type S3Storage struct {
	s3Storage *s3.S3
	bucket    string
}

func NewApplication(_ context.Context, configFile string, logger *logging.Logger) (app *Application, err error) {
	logger.Infof("init config...")
	cfg, err := config.NewConfig(configFile)
	if err != nil {
		return nil, errors.Wrap(err, "fail to initialize config")
	}

	logger.Infof("init s3 storage...")
	sessionS3, err := s3Storage.NewStorageS3(cfg.S3.Host, cfg.S3.Region, cfg.S3.AccessKey, cfg.S3.SecretKey)
	if err != nil {
		return nil, errors.Wrap(err, "fail to initialize s3-storage")
	}
	s3Bucket := cfg.S3.Bucket

	uploadStorage := storage.NewStorage(sessionS3, s3Bucket)
	uploadService := service.NewService(uploadStorage)
	uploadController := controller.NewController(uploadService)

	routes := uploadController.InitRoutes()

	return &Application{
		cfg: cfg,
		s3: S3Storage{
			s3Storage: sessionS3,
			bucket:    s3Bucket,
		},
		router: routes,
		logger: logger,
	}, nil

}

func (a *Application) StartHttp3Server() error {
	a.http3Server = &http3.Server{
		Handler:    a.router,
		Addr:       a.cfg.Server.Addr,
		QuicConfig: &quic.Config{},
	}
	a.logger.Infof("http server started on :%s", a.cfg.Server.Addr)
	go a.GracefulShutdown([]os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM}, a.http3Server)

	return a.http3Server.ListenAndServeTLS(certificate.GetCertificatePaths())
}

func (a *Application) Shutdown(_ context.Context) error {
	if err := a.http3Server.CloseGracefully(time.Duration(a.cfg.Server.ShutdownTimeout) * time.Second); err != nil {
		return errors.Wrap(err, "fail to gracefully db shutdown")
	}
	return nil
}

func (a *Application) GracefulShutdown(signals []os.Signal, closeItems ...io.Closer) {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, signals...)

	sig := <-sigc
	a.logger.Infoln("--- shutdown application ---")
	a.logger.Infof("Caught signal %s. Shutting down...", sig)

	for _, closer := range closeItems {
		time.Sleep(time.Duration(a.cfg.Server.ShutdownTimeout) * time.Second)
		if err := closer.Close(); err != nil {
			a.logger.Errorf("failed to close %v: %v", closer, err)
		}

		if err := a.Shutdown(context.Background()); err != nil {
			a.logger.Errorf("failed to shutdown: %v", err)
		}
	}
}
