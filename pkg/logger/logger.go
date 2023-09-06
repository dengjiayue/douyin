package logger

import (
	"douyin/internal/gateway/config"
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

var (
	defaultLogger Logger
	zapLogger     *zap.Logger
)

type Log struct {
	Env        string `yaml:"env"`
	Path       string `yaml:"path"`
	Encoding   string `yaml:"encoding"`
	MaxSize    int    `yaml:"max_size"`
	MaxAge     int    `yaml:"max_age"`
	MaxBackups int    `yaml:"max_backups"`
}

func Init(cfg interface{}) {
	var logCfg config.Log

	// 默认值
	if cfg == nil {
		logCfg = config.Log{
			Env:        "dev",
			Path:       "",
			Encoding:   "console",
			MaxSize:    100,
			MaxAge:     30,
			MaxBackups: 30,
		}
	} else {
		// 强行转换
		buf, err := yaml.Marshal(cfg)
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(buf, &logCfg)
		if err != nil {
			panic(err)
		}

	}

	// 日志输出文件还是控制台
	var writer zapcore.WriteSyncer
	if logCfg.Path != "" {
		writer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), getLogFileWriter(logCfg))
	} else {
		writer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout))
	}
	core := zapcore.NewCore(
		buildEncoder(logCfg),
		writer,
		enableLevel(logCfg),
	)

	// 增加堆栈打印 （caller跳过日志库工具类的代码行）
	zapLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	defaultLogger = zapLogger.Sugar()
}

func Sync() {
	_ = zapLogger.Sync()
}

func GetZapLogger() *zap.Logger {
	return zapLogger
}

func buildEncoder(logCfg config.Log) zapcore.Encoder {
	// 区分环境配置
	var encoderConfig zapcore.EncoderConfig
	if logCfg.Env == "prd" {
		encoderConfig = zap.NewProductionEncoderConfig()
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	}

	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("[2006-01-02 15:04:05.000]") // zapcore.ISO8601TimeEncoder

	// 输出格式json还是console
	if logCfg.Encoding == "json" {
		encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// 输出日志等级
func enableLevel(logCfg config.Log) zapcore.Level {
	if logCfg.Env == "prd" {
		return zapcore.InfoLevel
	}
	return zapcore.DebugLevel
}

// 轮转存文件
func getLogFileWriter(logCfg config.Log) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   logCfg.Path,
		MaxSize:    logCfg.MaxSize,    // 单个文件最大尺寸，单位M （默认值100M）
		MaxBackups: logCfg.MaxBackups, // 最多保留备份个数 (跟MaxAge都不配，则保留全部)
		MaxAge:     logCfg.MaxAge,     // 最大时间，默认单位 day
		LocalTime:  true,              // 使用本地时间
		Compress:   false,
	}
	_ = lumberJackLogger.Rotate()
	return zapcore.AddSync(lumberJackLogger)
}
