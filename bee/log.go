package bee

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"time"
)

var Logger *zap.SugaredLogger

// LogCfg 日志配置
type logConfig struct {
	FilePath      string `mapstructure:"filePath"`      // 日志文件路径
	FileName      string `mapstructure:"fileName"`      // 日志文件名
	MaxSize       int    `mapstructure:"maxSize"`       // 单文件最大体积(M),超过后自动切分
	MaxBackups    int    `mapstructure:"maxBackups"`    // 保留历史文件的最大个数
	MaxAge        int    `mapstructure:"maxAge"`        // 保留天数
	IsJsonEncoder bool   `mapstructure:"isJsonEncoder"` // 是否使用JSON编码器
}

// setDefault 设置默认值
func (l *logConfig) setDefault() {
	if l.FilePath == "" {
		l.FilePath = "log"
	}
	if l.FileName == "" {
		l.FileName = "app"
	}
	if l.MaxSize == 0 {
		l.MaxSize = 2
	}
	if l.MaxBackups == 0 {
		l.MaxBackups = 10
	}
	if l.MaxAge == 0 {
		l.MaxAge = 7
	}
}

// MagicLog 日志
type MagicLog struct {
	cfg    *logConfig
	LogKey string // 日志配置前缀key
}

func (m *MagicLog) initLogger() {
	// 配置
	var config logConfig
	if err := mapstructure.Decode(viper.GetStringMap(m.LogKey), &config); err != nil {
		panic(errors.New("日志配置转换失败。"))
	}
	config.setDefault()
	m.cfg = &config

	// 创建核心日志组件
	core := zapcore.NewCore(m.getEncoder(), zapcore.NewMultiWriteSyncer(m.getWriteSyncer(), zapcore.AddSync(os.Stdout)), zapcore.InfoLevel)
	Logger = zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel)).Sugar()
}

// getEncoder 获取编码器
func (m *MagicLog) getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.LevelKey = "L"
	encoderConfig.TimeKey = "T"
	encoderConfig.CallerKey = "C"
	encoderConfig.MessageKey = "M"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(fmt.Sprintf("[%s]", time.Local().Format("2006/01/02 15:04:05.00000")))
	}
	if m.cfg.IsJsonEncoder {
		return zapcore.NewJSONEncoder(encoderConfig)
	}

	return zapcore.NewConsoleEncoder(encoderConfig)
}

// getWriteSyncer 写入日志配置
func (m *MagicLog) getWriteSyncer() zapcore.WriteSyncer {
	// 设置日志存放目录及名称
	stSeparator := string(filepath.Separator)
	logFile := fmt.Sprintf("%s%s%s_%s.log", m.cfg.FilePath, stSeparator, m.cfg.FileName, time.Now().Format(time.DateOnly))
	// 设置文件切割,需用到第三方包(lumberjack),zap本身不支持
	lumberjackSyncer := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    m.cfg.MaxSize,
		MaxBackups: m.cfg.MaxBackups,
		MaxAge:     m.cfg.MaxAge,
		Compress:   false,
	}

	return zapcore.AddSync(lumberjackSyncer)
}
