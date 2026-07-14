package pkg

import (
	"sync"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	// 只需要它的 init() 先把 config/*.yaml（含 loggers.yaml）灌進全域 viper，
	// 這裡不使用 bootstrap.CONFIG，分類設定改用 viper 動態讀取。
	_ "example/internal/bootstrap"
)

// LogConfig 對應 config/loggers.yaml 底下每一個業務分類的設定。
type LogConfig struct {
	Directory  string `mapstructure:"directory"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

const sLogDefaultCategory = "default"

var oLoggerCache sync.Map // sCategory(string) -> *zap.Logger

// Logger 依業務分類（middleware / controller / cron / sdk / service ... 對應
// config/loggers.yaml 的 key）取得專屬的 *zap.Logger，各分類獨立目錄、
// 同一分類再依 level（debug/info/warn/error）拆成不同檔案。
// 分類在設定檔裡不存在時，退回 default 分類的設定。
// 同一分類重複呼叫會共用同一個 *zap.Logger，不會重複開檔案 handle。
func Logger(sCategory string) *zap.Logger {
	if oExisting, bOk := oLoggerCache.Load(sCategory); bOk {
		return oExisting.(*zap.Logger)
	}

	oLogger := newCategoryLogger(sCategory)
	oActual, _ := oLoggerCache.LoadOrStore(sCategory, oLogger)
	return oActual.(*zap.Logger)
}

func newCategoryLogger(sCategory string) *zap.Logger {
	oConfig := loadLogConfig(sCategory)

	oEncoderConfig := zap.NewProductionEncoderConfig()
	oEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	oEncoder := zapcore.NewJSONEncoder(oEncoderConfig)

	aCores := []zapcore.Core{
		zapcore.NewCore(oEncoder, logWriter(oConfig, sCategory, "debug"),
			zap.LevelEnablerFunc(func(iLevel zapcore.Level) bool { return iLevel == zap.DebugLevel })),
		zapcore.NewCore(oEncoder, logWriter(oConfig, sCategory, "info"),
			zap.LevelEnablerFunc(func(iLevel zapcore.Level) bool { return iLevel == zap.InfoLevel })),
		zapcore.NewCore(oEncoder, logWriter(oConfig, sCategory, "warn"),
			zap.LevelEnablerFunc(func(iLevel zapcore.Level) bool { return iLevel == zap.WarnLevel })),
		zapcore.NewCore(oEncoder, logWriter(oConfig, sCategory, "error"),
			zap.LevelEnablerFunc(func(iLevel zapcore.Level) bool { return iLevel >= zap.ErrorLevel })), // Error, DPanic, Panic, Fatal 全收
	}

	return zap.New(zapcore.NewTee(aCores...), zap.AddCaller())
}

func loadLogConfig(sCategory string) LogConfig {
	var oConfig LogConfig

	if err := viper.UnmarshalKey("loggers."+sCategory, &oConfig); err == nil && oConfig.Directory != "" {
		return oConfig
	}

	// 分類不存在或設定不完整時退回 default
	_ = viper.UnmarshalKey("loggers."+sLogDefaultCategory, &oConfig)
	return oConfig
}

func logWriter(oConfig LogConfig, sCategory string, sLevel string) zapcore.WriteSyncer {
	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   oConfig.Directory + "/" + sCategory + "-" + sLevel + ".log",
		MaxSize:    oConfig.MaxSize,
		MaxAge:     oConfig.MaxAge,
		MaxBackups: oConfig.MaxBackups,
		Compress:   oConfig.Compress,
	})
}
