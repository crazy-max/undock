package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	_ "time/tzdata"

	"github.com/alecthomas/kong"
	"github.com/crazy-max/undock/internal/app"
	"github.com/crazy-max/undock/internal/config"
	"github.com/crazy-max/undock/internal/logging"
	"github.com/rs/zerolog/log"
)

var (
	undock  *app.Undock
	cli     config.Cli
	version = "dev"
	meta    = config.Meta{
		ID:     "undock",
		Name:   "Undock",
		Desc:   "Extract contents of a container image in a local folder",
		URL:    "https://github.com/crazy-max/undock",
		Author: "CrazyMax",
	}
)

func main() {
	var err error
	runtime.GOMAXPROCS(runtime.NumCPU())

	meta.Version = version
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

	// Logging
	logging.Configure(cli)

	// Handle os signals
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt, SIGTERM)
	go func() {
		sig := <-channel
		undock.Close()
		log.Warn().Msgf("caught signal %v", sig)
		os.Exit(0)
	}()

	// Init
	if undock, err = app.New(meta, cli); err != nil {
		log.Fatal().Err(err).Msg("cannot initialize undock")
	}

	// Start
	if err = undock.Start(); err != nil {
		log.Fatal().Stack().Err(err).Send()
	}
}
