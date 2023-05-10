package container

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/roadrunner-server/endure/v2"
	"github.com/roadrunner-server/errors"
	"github.com/rumorsflow/rumors/v2/internal/common"
	"github.com/rumorsflow/rumors/v2/internal/config"
	"github.com/rumorsflow/rumors/v2/internal/logger"
	"github.com/rumorsflow/rumors/v2/pkg/errs"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	errors.Separator = " |> "
}

type Container struct {
	EnvPrefix string
	CfgFile   string
	Version   string
}

func (c *Container) Run(plugins ...any) error {
	const op = errors.Op("container_run")

	containerCfg, err := NewConfig(c.CfgFile, c.EnvPrefix)
	if err != nil {
		return errors.E(op, err)
	}

	cfg := &config.Plugin{
		Path:    c.CfgFile,
		Prefix:  c.EnvPrefix,
		Timeout: containerCfg.GracePeriod,
		Version: c.Version,
	}

	endureOptions := []endure.Options{
		endure.GracefulShutdownTimeout(containerCfg.GracePeriod),
	}

	if containerCfg.PrintGraph {
		endureOptions = append(endureOptions, endure.Visualize())
	}

	cont := endure.New(containerCfg.LogLevel, endureOptions...)

	if err = cont.RegisterAll(append(plugins, cfg, &logger.Plugin{})...); err != nil {
		return errors.E(op, err)
	}

	if err = cont.Init(); err != nil {
		return errors.E(op, err)
	}

	errCh, err := cont.Serve()
	if err != nil {
		return errors.E(op, err)
	}

	oss, stop := make(chan os.Signal, 5), make(chan struct{}, 1) //nolint:gomnd
	signal.Notify(oss, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-oss

		stop <- struct{}{}

		<-oss
		fmt.Println("exit forced")
		os.Exit(1)
	}()

	for {
		select {
		case e := <-errCh:
			if e.Error == common.Success {
				return cont.Stop()
			}

			err1 := fmt.Errorf("error: %w\nplugin: %s", e.Error, e.VertexID)
			err2 := cont.Stop()

			return errs.Append(err1, err2)
		case <-stop:
			_, _ = color.New(color.FgHiBlue, color.Bold).Fprintf(os.Stderr, "stop signal received, grace timeout is: %0.f seconds\n", containerCfg.GracePeriod.Seconds())

			return cont.Stop()
		}
	}
}
