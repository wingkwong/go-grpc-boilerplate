package logger

import (
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Log              *zap.Logger
	customTimeFormat string
	onceInit         sync.Once
)

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(customTimeFormat))
}

// Overview of Level:
// -1 : DebugLevel logs are typically voluminous, and are usually disabled in production.
// 0: InfoLevel is the default logging priority.
// 1: WarnLevel logs are more important than Info, but don't need individual human review.
// 2: ErrorLevel logs are high-priority. If an application is running smoothly, it shouldn't generate any error-level logs.
// 3: DPanicLevel logs are particularly important errors. In development the logger panics after writing the message.
// 4: PanicLevel logs a message, then panics.
// 5: FatalLevel logs a message, then calls os.Exit(1).

func Init(level int, timeFormat string) error {
	var err error

	onceInit.Do(func() {
		globalLevel := zapcore.Level(level)
		highPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
			return level >= zapcore.ErrorLevel
		})
		lowPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
			return level >= globalLevel && level < zapcore.ErrorLevel
		})
		consoleInfos := zapcore.Lock(os.Stdout)
		consoleErrors := zapcore.Lock(os.Stderr)

		var useCustomTimeFormat bool
		ecfg := zap.NewProductionEncoderConfig()
		if len(timeFormat) > 0 {
			customTimeFormat = timeFormat
			ecfg.EncodeTime = customTimeEncoder
			useCustomTimeFormat = true
		}
		consoleEncoder := zapcore.NewJSONEncoder(ecfg)

		core := zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
			zapcore.NewCore(consoleEncoder, consoleInfos, lowPriority),
		)

		Log = zap.New(core)
		zap.RedirectStdLog(Log)

		if !useCustomTimeFormat {
			Log.Warn("[WARN] Time format for logger is not provided - use zap default")
		}
	})

	return err
}
