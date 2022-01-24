package logging

import (
	"github.com/rs/zerolog/log"
	"github.com/sirupsen/logrus"
)

// LogrusFormatter is a logrus formatter
type LogrusFormatter struct{}

// Format renders a single log entry from logrus entry to zerolog
func (f *LogrusFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	switch entry.Level {
	case logrus.ErrorLevel:
		log.Error().Fields(entry.Data).Msg(entry.Message)
	case logrus.WarnLevel:
		log.Warn().Fields(entry.Data).Msg(entry.Message)
	case logrus.DebugLevel:
		log.Debug().Fields(entry.Data).Msg(entry.Message)
	case logrus.TraceLevel:
		log.Trace().Fields(entry.Data).Msg(entry.Message)
	default:
		log.Info().Fields(entry.Data).Msg(entry.Message)
	}
	return nil, nil
}
