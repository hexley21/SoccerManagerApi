package zap_logger

import (
	"io"
	"os"

	"github.com/hexley21/soccer-manager/pkg/config"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var loggerLevelMap = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
}

func getLoggerLevel(lvl string) zapcore.Level {
	level, exist := loggerLevelMap[lvl]
	if !exist {
		return zapcore.DebugLevel
	}

	return level
}

var echoLevelMap = map[zapcore.Level]log.Lvl{
	zapcore.DebugLevel: log.DEBUG,
	zapcore.InfoLevel:  log.INFO,
	zapcore.WarnLevel:  log.WARN,
	zapcore.ErrorLevel: log.ERROR,
}

func getEchoLoggerLevel(lvl zapcore.Level) log.Lvl {
	level, exist := echoLevelMap[lvl]
	if !exist {
		return log.DEBUG
	}

	return level
}

type ZapLogger struct {
	sugar  *zap.SugaredLogger
	writer io.Writer
}

func New(logPath string, cfg config.Logging, isProd bool) *ZapLogger {
	logWriter := zapcore.AddSync(os.Stdout)

	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		log.Fatal(err)
	}

	fileWriter := zapcore.AddSync(logFile)

	var encoderCfg zapcore.EncoderConfig
	var encoder zapcore.Encoder

	if isProd {
		encoderCfg = zap.NewProductionEncoderConfig()
		encoderCfg.NameKey = "[SERVICE]"
		encoderCfg.TimeKey = "[TIME]"
		encoderCfg.LevelKey = "[LEVEL]"
		encoderCfg.FunctionKey = "[CALLER]"
		encoderCfg.CallerKey = "[LINE]"
		encoderCfg.MessageKey = "[MESSAGE]"
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
		encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder
		encoderCfg.EncodeName = zapcore.FullNameEncoder
		encoderCfg.EncodeDuration = zapcore.StringDurationEncoder
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	} else {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
		encoderCfg.NameKey = "[SERVICE]"
		encoderCfg.TimeKey = "[TIME]"
		encoderCfg.LevelKey = "[LEVEL]"
		encoderCfg.FunctionKey = "[CALLER]"
		encoderCfg.CallerKey = "[LINE]"
		encoderCfg.MessageKey = "[MESSAGE]"
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderCfg.EncodeName = zapcore.FullNameEncoder
		encoderCfg.EncodeDuration = zapcore.StringDurationEncoder
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderCfg.EncodeCaller = zapcore.FullCallerEncoder
		encoderCfg.ConsoleSeparator = " | "
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	}

	var cfgLevel zapcore.Level
	if isProd {
		cfgLevel = getLoggerLevel(cfg.ProdLogLevel)
	} else {
		cfgLevel = getLoggerLevel(cfg.DevLogLevel)
	}

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, logWriter, zap.NewAtomicLevelAt(cfgLevel)),
		zapcore.NewCore(encoder, fileWriter, zap.NewAtomicLevelAt(cfgLevel)),
	)

	var options []zap.Option

	if cfg.CallerEnabled {
		options = append(options, zap.AddCaller())
		options = append(options, zap.AddCallerSkip(1))
	}

	return &ZapLogger{
		sugar:  zap.New(core, options...).Sugar(),
		writer: io.MultiWriter(os.Stdout, logFile),
	}
}

func jsonToFields(j log.JSON) []interface{} {
	fields := make([]interface{}, 0, len(j)*2)
	for k, v := range j {
		fields = append(fields, k, v)
	}
	return fields
}

func (z *ZapLogger) Output() io.Writer {
	return z.writer
}

func (z *ZapLogger) SetOutput(w io.Writer) {
	z.writer = w
}

func (z *ZapLogger) Prefix() string {
	return ""
}

func (z *ZapLogger) SetPrefix(p string) {}

func (z *ZapLogger) Level() log.Lvl {
	return getEchoLoggerLevel(z.sugar.Level())
}

func (z *ZapLogger) SetLevel(v log.Lvl) {}

func (z *ZapLogger) SetHeader(h string) {}

func (z *ZapLogger) Print(i ...interface{}) {
	z.sugar.Info(i...)
}

func (z *ZapLogger) Printf(format string, args ...interface{}) {
	z.sugar.Infof(format, args...)
}

func (z *ZapLogger) Printj(j log.JSON) {
	z.sugar.Infow("", jsonToFields(j)...)
}

func (z *ZapLogger) Println(i ...interface{}) {
	z.sugar.Info(i...)
}

func (z *ZapLogger) Debug(i ...interface{}) {
	z.sugar.Debug(i...)
}

func (z *ZapLogger) Debugf(format string, args ...interface{}) {
	z.sugar.Debugf(format, args...)
}

func (z *ZapLogger) Debugj(j log.JSON) {
	z.sugar.Debugw("", jsonToFields(j)...)
}

func (z *ZapLogger) Info(i ...interface{}) {
	z.sugar.Info(i...)
}

func (z *ZapLogger) Infof(format string, args ...interface{}) {
	z.sugar.Infof(format, args...)
}

func (z *ZapLogger) Infoj(j log.JSON) {
	z.sugar.Infow("", jsonToFields(j)...)
}

func (z *ZapLogger) Warn(i ...interface{}) {
	z.sugar.Warn(i...)
}

func (z *ZapLogger) Warnf(format string, args ...interface{}) {
	z.sugar.Warnf(format, args...)
}

func (z *ZapLogger) Warnj(j log.JSON) {
	z.sugar.Warnw("", jsonToFields(j)...)
}

func (z *ZapLogger) Error(i ...interface{}) {
	z.sugar.Error(i...)
}

func (z *ZapLogger) Errorf(format string, args ...interface{}) {
	z.sugar.Errorf(format, args...)
}

func (z *ZapLogger) Errorj(j log.JSON) {
	z.sugar.Errorw("", jsonToFields(j)...)
}

func (z *ZapLogger) Fatal(i ...interface{}) {
	z.sugar.Fatal(i...)
}

func (z *ZapLogger) Fatalf(format string, args ...interface{}) {
	z.sugar.Fatalf(format, args...)
}

func (z *ZapLogger) Fatalj(j log.JSON) {
	z.sugar.Fatalw("", jsonToFields(j)...)
}

func (z *ZapLogger) Panic(i ...interface{}) {
	z.sugar.Panic(i...)
}

func (z *ZapLogger) Panicf(format string, args ...interface{}) {
	z.sugar.Panicf(format, args...)
}

func (z *ZapLogger) Panicj(j log.JSON) {
	z.sugar.Panicw("", jsonToFields(j)...)
}
