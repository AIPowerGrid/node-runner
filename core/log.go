package core

import (
	"os"

	zap "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var sugar *zap.SugaredLogger = nil

func GetEnv() string {
	env := os.Getenv("GO_ENV")
	if env == "" {
		return "development"
	} else {
		return env
	}

}

// Get get log
func GetLogger() *zap.SugaredLogger {
	if sugar != nil {
		return sugar
	}
	var logger *zap.Logger
	GOENV := GetEnv()
	if GOENV == "development" {
		// fmt.Println("goenv is dev", appSettings.GOENV)
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.DisableStacktrace = true
		// config.DisableCaller = true
		config.Level.SetLevel(zap.DebugLevel)
		logger, _ = config.Build()
		// logger, _ = zap.NewDevelopment(config)
	} else {
		// fmt.Println("goenv is", appSettings.GOENV)
		config := zap.NewProductionConfig()
		config.DisableStacktrace = true
		config.DisableCaller = true
		config.Level.SetLevel(zap.DebugLevel)
		logger, _ = config.Build()
	}
	// defer logger.Sync() // flushes buffer, if any
	sugar = logger.Sugar()
	return sugar
}
