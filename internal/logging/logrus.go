package logging

import (
	"github.com/rs/zerolog/log"
	"github.com/sirupsen/logrus"
)

// LogrusFormatter is a logrus formatter
type LogrusFormatter struct{}

// Format renders a single log entry from logrus entry to zerolog
func (f *LogrusFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	fields := map[string]interface{}(entry.Data)
	switch entry.Level {
	case logrus.ErrorLevel:
		log.Error().Fields(fields).Msg(entry.Message)
	case logrus.WarnLevel:
		log.Warn().Fields(fields).Msg(entry.Message)
	case logrus.DebugLevel:
		log.Debug().Fields(fields).Msg(entry.Message)
	case logrus.TraceLevel:
		log.Trace().Fields(fields).Msg(entry.Message)
	default:
		log.Info().Fields(fields).Msg(entry.Message)
	}
	return nil, nil
}
