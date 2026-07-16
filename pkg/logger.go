package pkg

import (
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Module string

const (
	Default         Module = "default"
	Cron            Module = "cron"
	Consumer        Module = "consumer"
	Websocket       Module = "websocket"
	HttpAdmin       Module = "http-admin"
	HttpApp         Module = "http-app"
	HttpThird       Module = "http-third"
	MiddlewareAdmin Module = "middleware-admin"
	MiddlewareApp   Module = "middleware-app"
	MiddlewareThird Module = "middleware-third"
	Client          Module = "client"
	FacadeGame      Module = "facade-game"
	FacadeTable     Module = "facade-table"
	FacadeRegister  Module = "facade-register"
	ResourceLogic   Module = "resource-logic"
	ResourceModel   Module = "resource-model"
	Repository      Module = "repository"
	Sdk             Module = "sdk"
	Publisher       Module = "publisher"
)

var (
	cache = make(map[Module]*zap.Logger)
	mu    sync.Mutex
)

// Get 取得指定模組的 Logger
func Logger(module Module) *zap.Logger {
	mu.Lock()
	defer mu.Unlock()

	if logger, ok := cache[module]; ok {
		return logger
	}

	dir := filepath.Join("runtime", "log", string(module))
	_ = os.MkdirAll(dir, 0755)

	newWriter := func(filename string) zapcore.WriteSyncer {
		return zapcore.AddSync(&lumberjack.Logger{
			Filename:   filepath.Join(dir, filename),
			MaxSize:    100, // MB
			MaxBackups: 10,
			MaxAge:     30, // Days
			Compress:   true,
		})
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// app.log (Info 以上全部)
	appCore := zapcore.NewCore(
		jsonEncoder,
		newWriter("index.log"),
		zap.LevelEnablerFunc(func(level zapcore.Level) bool {
			return level >= zap.InfoLevel
		}),
	)

	// warn.log (Warn 以上)
	warnCore := zapcore.NewCore(
		jsonEncoder,
		newWriter("warn.log"),
		zap.LevelEnablerFunc(func(level zapcore.Level) bool {
			return level >= zap.WarnLevel
		}),
	)

	// error.log (Error 以上)
	errorCore := zapcore.NewCore(
		jsonEncoder,
		newWriter("error.log"),
		zap.LevelEnablerFunc(func(level zapcore.Level) bool {
			return level >= zap.ErrorLevel
		}),
	)

	// Console（帶顏色，只影響終端輸出，不影響寫檔的 json log）
	consoleEncoderConfig := zap.NewDevelopmentEncoderConfig()
	consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(consoleEncoderConfig),
		zapcore.AddSync(os.Stdout),
		zap.DebugLevel,
	)

	logger := zap.New(
		zapcore.NewTee(
			appCore,
			warnCore,
			errorCore,
			consoleCore,
		),
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	).With(
		zap.String("module", string(module)),
	)

	cache[module] = logger

	return logger
}
