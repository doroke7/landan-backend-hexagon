package helper

import (
	"example/internal/bootstrap"
	"fmt"
	"reflect"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

/*

由於 loggerhlepr 是全container 同一個，
即便 Logger *zap.Logger 沒有注入，也會維持一份

*/

type LoggerHelper struct {
	Logger *zap.Logger
	*AbstractHelper
}

func NewLoggerHelper(oAbstractHelper *AbstractHelper) *LoggerHelper {

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 時間格式：2026-03-17T...
	oEncoder := zapcore.NewJSONEncoder(encoderConfig)

	oDebugEnabler := zap.LevelEnablerFunc(func(iLevel zapcore.Level) bool { return iLevel == zap.DebugLevel })
	oInfoEnabler := zap.LevelEnablerFunc(func(iLevel zapcore.Level) bool { return iLevel == zap.InfoLevel })
	oWarnEnabler := zap.LevelEnablerFunc(func(iLevel zapcore.Level) bool { return iLevel == zap.WarnLevel })
	oErrorEnabler := zap.LevelEnablerFunc(func(iLevel zapcore.Level) bool { return iLevel >= zap.ErrorLevel }) // Error, DPanic, Panic, Fatal 全收

	// oValueLogger := mapping(bootstrap.CONFIG.LOGGERS, sDirectory)

	// if LOGGER, bOk := oValueLogger.(struct {
	// 	DIRECTORY   string `mapstructure:"directory"`
	// 	MAX_SIZE    int    `mapstructure:"max_size"`
	// 	MAX_BACKUPS int    `mapstructure:"max_backups"`
	// 	MAX_AGE     int    `mapstructure:"max_age"`
	// 	COMPRESS    bool   `mapstructure:"compress"`
	// }); bOk {
	// 	aCores := []zapcore.Core{
	// 		zapcore.NewCore(oEncoder, writer(LOGGER.DIRECTORY+"/"+sFilename+"-debug.log", LOGGER.MAX_SIZE, LOGGER.MAX_BACKUPS, LOGGER.COMPRESS), oDebugEnabler),
	// 		zapcore.NewCore(oEncoder, writer(LOGGER.DIRECTORY+"/"+sFilename+"-info.log", LOGGER.MAX_SIZE, LOGGER.MAX_BACKUPS, LOGGER.COMPRESS), oInfoEnabler),
	// 		zapcore.NewCore(oEncoder, writer(LOGGER.DIRECTORY+"/"+sFilename+"-warn.log", LOGGER.MAX_SIZE, LOGGER.MAX_BACKUPS, LOGGER.COMPRESS), oWarnEnabler),
	// 		zapcore.NewCore(oEncoder, writer(LOGGER.DIRECTORY+"/"+sFilename+"-error.log", LOGGER.MAX_SIZE, LOGGER.MAX_BACKUPS, LOGGER.COMPRESS), oErrorEnabler),
	// 	}
	// 	return &LoggerHelper{
	// 		Logger: zap.New(zapcore.NewTee(aCores...), zap.AddCaller()),
	// 	}
	// }

	LOGGER := bootstrap.CONFIG.LOGGERS.DEFAULT

	aCores := []zapcore.Core{
		zapcore.NewCore(oEncoder, writer(LOGGER.DIRECTORY+"/"+"index-debug.log", LOGGER.MAX_SIZE, LOGGER.MAX_BACKUPS, LOGGER.COMPRESS), oDebugEnabler),
		zapcore.NewCore(oEncoder, writer(LOGGER.DIRECTORY+"/"+"index-info.log", LOGGER.MAX_SIZE, LOGGER.MAX_BACKUPS, LOGGER.COMPRESS), oInfoEnabler),
		zapcore.NewCore(oEncoder, writer(LOGGER.DIRECTORY+"/"+"index-warn.log", LOGGER.MAX_SIZE, LOGGER.MAX_BACKUPS, LOGGER.COMPRESS), oWarnEnabler),
		zapcore.NewCore(oEncoder, writer(LOGGER.DIRECTORY+"/"+"index-error.log", LOGGER.MAX_SIZE, LOGGER.MAX_BACKUPS, LOGGER.COMPRESS), oErrorEnabler),
	}
	return &LoggerHelper{
		Logger:         zap.New(zapcore.NewTee(aCores...), zap.AddCaller()),
		AbstractHelper: oAbstractHelper,
	}

}

func mapping(oConfig interface{}, sSubDirectory string) interface{} {

	v := reflect.ValueOf(oConfig)

	// 2. 為了防呆且匹配你的結構體命名，將輸入轉為大寫
	fieldName := strings.ToUpper(sSubDirectory)

	// 3. 根據字串名稱尋找欄位
	field := v.FieldByName(fieldName)

	var oInterface interface{}

	// 4. 檢查欄位是否存在
	if !field.IsValid() {
		fmt.Printf("警告: 找不到層級 [%s], 將返回預設配置\n", fieldName)
		// 如果找不到，通常建議返回 DEFAULT 欄位作為 Fallback
		oInterface = v.FieldByName("DEFAULT").Interface()
	}

	if field.IsValid() {
		oInterface = field.Interface()

	}

	return oInterface
}

func writer(sFilename string, iMaxSize int, iMaxBackups int, bCompress bool) zapcore.WriteSyncer {

	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   sFilename,
		MaxSize:    iMaxSize,    // 50-100MB 視你的硬碟配置
		MaxBackups: iMaxBackups, // 保留幾份舊檔案
		Compress:   bCompress,   // 百萬級日誌建議開啟壓縮
	})
}
