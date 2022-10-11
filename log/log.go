package log

import (
	"os"
	"path/filepath"
	"strings"

	"ds-agent-watcher/utils"

	"github.com/ian-kent/go-log/appenders"
	"github.com/ian-kent/go-log/layout"
	"github.com/ian-kent/go-log/levels"
	"github.com/ian-kent/go-log/log"
	"github.com/ian-kent/go-log/logger"
)

//GetLogger - get logger
func GetLogger(level string, agentType string, agentID string) logger.Logger {
	logger := log.Logger("logger")
	pattern := "[%d]" + " [%p]" + " [" + agentType + "]" + " [" + agentID + "]" + " %m"
	mode := strings.ToUpper(level)
	loggingLevel := levels.INFO
	switch mode {
	case "INFO":
		loggingLevel = levels.INFO
		break
	case "DEBUG":
		loggingLevel = levels.DEBUG
		break
	case "PROD":
		loggingLevel = levels.INFO
		break
	default:
		loggingLevel = levels.INFO
		break
	}
	logger.SetLevel(loggingLevel)
	layoutPattern := layout.Pattern(pattern)
	consoleAppender := appenders.Console()
	consoleAppender.SetLayout(layoutPattern)
	if !(IsKubernetesEnv()) {
		logFileSize := 1024 * 1024 * 5
		d, _, _ := utils.GetExecutablePathAndName()
		rollingFileAppender := appenders.RollingFile(filepath.Join(d, "..", "log", "sentinel.log"), true)
		rollingFileAppender.MaxFileSize = int64(logFileSize)
		rollingFileAppender.MaxBackupIndex = 4
		rollingFileAppender.SetLayout(layoutPattern)
		logger.SetAppender(appenders.Multiple(layoutPattern, consoleAppender, rollingFileAppender))
		return logger
	}
	logger.SetAppender(consoleAppender)
	return logger
}

//IsKubernetesEnv - is kubernetes environment
func IsKubernetesEnv() bool {
	if os.Getenv("KUBERNETES_SERVICE_PORT") != "" && os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		return true
	}
	return false
}
