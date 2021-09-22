package main

import (
	"context"
	"time"

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
	ticker := time.NewTicker(time.Duration(1) * time.Second)

	for _ = range ticker.C {
		if err := FileStorage.MakeBucket(ctx, "qwerty"); err != nil {
			logrus.Fatalf("error make bucket: %s", err.Error())
		}
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
