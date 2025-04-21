package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps the zap logger with convenient methods
type Logger struct {
	*zap.Logger
}

var (
	once   sync.Once
	logger *Logger
)

// LogLevel represents the application log level
type LogLevel string

const (
	// Debug level logging
	Debug LogLevel = "debug"
	// Info level logging
	Info LogLevel = "info"
	// Warn level logging
	Warn LogLevel = "warn"
	// Error level logging
	Error LogLevel = "error"
)

// getLogLevel converts a string log level to zapcore.Level
func getLogLevel(level LogLevel) zapcore.Level {
	switch level {
	case Debug:
		return zapcore.DebugLevel
	case Info:
		return zapcore.InfoLevel
	case Warn:
		return zapcore.WarnLevel
	case Error:
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// InitLogger initializes the logger with the specified log level
func InitLogger(logLevel LogLevel, serviceName string) *Logger {
	once.Do(func() {
		// Create encoder config
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.TimeKey = "timestamp"
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

		// Create encoder and level
		level := getLogLevel(logLevel)
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			level,
		)

		// Create logger with specified options
		zapLogger := zap.New(
			core,
			zap.AddCaller(),
			zap.AddCallerSkip(1),
			zap.Fields(
				zap.String("service", serviceName),
			),
		)

		logger = &Logger{zapLogger}
	})

	return logger
}

// Get returns the singleton logger instance
func Get() *Logger {
	if logger == nil {
		// Initialize with default values if not already initialized
		return InitLogger(Info, "grpc-multi-tenant")
	}
	return logger
}

// With adds fields to the logger
func (l *Logger) With(fields map[string]interface{}) *Logger {
	if len(fields) == 0 {
		return l
	}

	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}

	return &Logger{l.Logger.With(zapFields...)}
}

// WithField adds a single field to the logger
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{l.Logger.With(zap.Any(key, value))}
}

// Debug logs a message at the debug level
func (l *Logger) Debug(msg string, fields map[string]interface{}) {
	if fields == nil {
		l.Logger.Debug(msg)
		return
	}
	l.With(fields).Logger.Debug(msg)
}

// Info logs a message at the info level
func (l *Logger) Info(msg string, fields map[string]interface{}) {
	if fields == nil {
		l.Logger.Info(msg)
		return
	}
	l.With(fields).Logger.Info(msg)
}

// Warn logs a message at the warn level
func (l *Logger) Warn(msg string, fields map[string]interface{}) {
	if fields == nil {
		l.Logger.Warn(msg)
		return
	}
	l.With(fields).Logger.Warn(msg)
}

// Error logs a message at the error level
func (l *Logger) Error(msg string, fields map[string]interface{}) {
	if fields == nil {
		l.Logger.Error(msg)
		return
	}
	l.With(fields).Logger.Error(msg)
}

// Fatal logs a message at the fatal level and then calls os.Exit(1)
func (l *Logger) Fatal(msg string, fields map[string]interface{}) {
	if fields == nil {
		l.Logger.Fatal(msg)
		return
	}
	l.With(fields).Logger.Fatal(msg)
}
