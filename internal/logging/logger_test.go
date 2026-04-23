package logging

import (
	"testing"

	"github.com/crazy-max/undock/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestConfigureSetsGlobalLevels(t *testing.T) {
	oldLogger := log.Logger
	oldGlobalLevel := zerolog.GlobalLevel()
	oldLogrusLevel := logrus.GetLevel()
	oldFormatter := logrus.StandardLogger().Formatter
	defer func() {
		log.Logger = oldLogger
		zerolog.SetGlobalLevel(oldGlobalLevel)
		logrus.SetLevel(oldLogrusLevel)
		logrus.SetFormatter(oldFormatter)
	}()

	Configure(config.Cli{
		LogLevel: "debug",
		LogJSON:  true,
	})

	assert.Equal(t, zerolog.DebugLevel, zerolog.GlobalLevel())
	assert.Equal(t, logrus.DebugLevel, logrus.GetLevel())
	_, ok := logrus.StandardLogger().Formatter.(*LogrusFormatter)
	assert.True(t, ok)
}
