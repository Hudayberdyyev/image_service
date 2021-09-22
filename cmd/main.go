package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Hudayberdyyev/image_service/pkg/handler"
	"github.com/Hudayberdyyev/image_service/server"
	"github.com/Hudayberdyyev/image_service/storage"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	ctx := context.Background()
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}
	FileStorage, err := storage.NewStorage(storage.Config{
		Endpoint:       viper.GetString("storage.host") + ":" + viper.GetString("storage.port"),
		AccessKeyId:    viper.GetString("storage.username"),
		SecretAccesKey: viper.GetString("storage.password"),
		UseSSL:         viper.GetBool("storage.sslmode"),
	})

	if err != nil {
		logrus.Fatalf("error initializing minio storage: %s", err.Error())
	}

	// if err := FileStorage.MakeBucket(ctx, storage.LogoBucket); err != nil {
	// 	logrus.Printf("error make bucket(%s): %s", storage.LogoBucket, err.Error())
	// }

	// if err := FileStorage.UploadAuthorsLogo(ctx, storage.LogoBucket); err != nil {
	// 	logrus.Printf("error upload authors logo: %s", err.Error())
	// }
	handlers := handler.NewHandler(FileStorage)
	srv := new(server.Server)
	go func() {
		if err := srv.Run(viper.GetString("server.ip"), viper.GetString("server.port"), viper.GetString("server.protocol"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logrus.Print("Server started ...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logrus.Print("Server Shutting Down")

	if err := srv.Shutdown(ctx); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
