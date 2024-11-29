package util

import (
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

// Take the list of blocking users and split them into chunks
func SliceChunker(input []string, chunkSize int) (output [][]string) {
	for i := 0; i < len(input); i += chunkSize {
		end := i + chunkSize
		if end > len(input) {
			end = len(input)
		}
		output = append(output, input[i:end])
	}
	return output
}

// ReturnLogLevels : Return a comma-delimited list of log levels
func ReturnLogLevels(levelMap map[string]logrus.Level) string {
	logLevels := make([]string, 0, len(levelMap))
	for k := range levelMap {
		logLevels = append(logLevels, k)
	}
	sort.Strings(logLevels)

	return strings.Join(logLevels, ", ")
}

// ConfigureLogger : Configure the logger
func ConfigureLogger(logLevel logrus.Level, nocolorFlag bool) (logger *logrus.Logger) {
	disableColors := false
	if nocolorFlag {
		disableColors = true
	}
	logger = &logrus.Logger{
		Out:   os.Stderr,
		Level: logLevel,
		Formatter: &prefixed.TextFormatter{
			DisableColors:    disableColors,
			DisableTimestamp: true,
			TimestampFormat:  "2006-01-02 15:04:05",
			FullTimestamp:    true,
			ForceFormatting:  false,
		},
	}
	logger.SetLevel(logLevel)

	return logger
}

func StripNonPrintable(s string) string {
	re, _ := regexp.Compile(`[^\x00-\x7F]+`)
	return re.ReplaceAllString(s, "")
}
