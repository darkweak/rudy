// Package logger wrap the logger
package logger

import "go.uber.org/zap"

// Logger is the global logger.
var Logger, _ = zap.NewDevelopment()
