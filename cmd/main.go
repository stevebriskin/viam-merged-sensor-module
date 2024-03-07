package main

import (
	"context"

	mymodule "github.com/stevebriskin/viam-merged-sensor-module"
	"go.viam.com/utils"

	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/module"
)

func main() {
	utils.ContextualMain(mainWithArgs, module.NewLoggerFromArgs("sds011"))
}

func mainWithArgs(ctx context.Context, args []string, logger logging.Logger) error {
	mod, err := module.NewModuleFromArgs(ctx, logger)
	if err != nil {
		return err
	}
	if err := mod.AddModelFromRegistry(ctx, sensor.API, mymodule.Model); err != nil {
		return err
	}

	if err := mod.Start(ctx); err != nil {
		return err
	}
	defer mod.Close(ctx)
	<-ctx.Done()
	return nil
}
