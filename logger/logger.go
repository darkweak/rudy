package logger

import "go.uber.org/zap"

var Logger, _ = zap.NewDevelopment()
