package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	_ "time/tzdata"

	"github.com/alecthomas/kong"
	"github.com/crazy-max/undock/internal/app"
	"github.com/crazy-max/undock/internal/config"
	"github.com/crazy-max/undock/internal/logging"
	"github.com/rs/zerolog/log"
)

var version = "dev"

func main() {
	if err := run(); err != nil {
		log.Fatal().Stack().Err(err).Send()
	}
}

func run() error {
	cli := config.Cli{}
	meta := config.Meta{
		ID:      "undock",
		Name:    "Undock",
		Desc:    "Extract contents of a container image in a local folder",
		URL:     "https://github.com/crazy-max/undock",
		Author:  "CrazyMax",
		Version: version,
	}

	meta.UserAgent = fmt.Sprintf("%s/%s go/%s %s", meta.ID, meta.Version, runtime.Version()[2:], strings.Title(runtime.GOOS)) //nolint:staticcheck // ignoring "SA1019: strings.Title is deprecated", as for our use we don't need full unicode support

	_ = kong.Parse(&cli,
		kong.Name(meta.ID),
		kong.Description(fmt.Sprintf("%s. More info: %s", meta.Desc, meta.URL)),
		kong.UsageOnError(),
		kong.Vars{
			"version": version,
		},
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}))

	logging.Configure(cli)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	undock, err := app.New(meta, cli)
	if err != nil {
		return err
	}
	return undock.Start(ctx)
}
