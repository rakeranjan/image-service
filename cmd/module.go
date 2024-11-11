package main

import (
	"context"

	"github.com/rakeranjan/image-service/internal/configuration"
	"github.com/rakeranjan/image-service/module"
	"github.com/rakeranjan/image-service/utils"
)

type ServiceModule interface {
	StartProcess() error
}

func NewModule(ctx context.Context, conf *configuration.Config) ServiceModule {
	switch conf.Module {
	case utils.USER_SERVICE:
		{
			return module.NewUserService(ctx, conf)
		}
	case utils.IMAGE_SERVICE:
		{
			return module.NewImageService(ctx, conf)
		}
	}
	return nil
}
