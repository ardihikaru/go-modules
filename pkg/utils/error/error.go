package error

import "go.uber.org/zap"

// FatalOnError stops the application if any fatal error occurs
func FatalOnError(err error, msg string) {
	if err != nil {
		// prepares a basic logger
		log, _ := zap.NewProduction(zap.WithCaller(false))
		log.Error("fatal error.", zap.Error(err))
		zap.S().Fatalf("%s:%s", msg, err)
	}
}
