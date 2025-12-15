package logging

import (
	"github.com/sirupsen/logrus"
	// prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func Init() {
	SetFormatter(logrus.StandardLogger())
	logrus.SetLevel(logrus.DebugLevel)
}

// SetFormatter sets the logrus logger formatter.
func SetFormatter(logger *logrus.Logger) {
	// Set JSON formatter with custom field names
	logger.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyTime:  "time",
			logrus.FieldKeyMsg:   "message",
		},
	})
	// If not in local environment, use prefixed text formatter for better readability
	// if isLocal, _ := strconv.ParseBool(viper.GetString("local-env")); !isLocal {
	// 	logger.SetFormatter(&prefixed.TextFormatter{
	// 		ForceFormatting: true,
	// 	})
	// }
}
