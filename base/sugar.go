package base

import (
	"fmt"
	"sync"

	"github.com/gobuffalo/envy"
	"go.uber.org/zap"
)

// 全局变量，是一个 pointer
var zapLogger *zap.Logger
var sugar *zap.SugaredLogger
var loggerOnce sync.Once

// Sugar 返回全局指针
func Sugar() *zap.SugaredLogger {
	loggerOnce.Do(func() {
		var err error

		// 为全局变量赋值
		zapLogger, err = newZapLogger()
		if err != nil {
			fmt.Printf("%+v\n", err)
			panic("zapLogger build failed.")
		}

		// 为全局变量赋值
		sugar = zapLogger.Sugar()
	})

	return sugar
}

// ZapSync 善后工作
// from doc: Sync calls the underlying Core's Sync method, flushing any buffered log entries.
// Applications should take care to call Sync before exiting.
func ZapSync() {
	defer zapLogger.Sync()
}

// newZapLogger 设置 logger 参数
func newZapLogger() (*zap.Logger, error) {
	var cfg zap.Config
	env := envy.Get("GO_ENV", "development")

	if env == "development" {
		// 默认是开发环境
		cfg = zap.NewDevelopmentConfig()
	} else {
		// 如果是 release
		cfg = zap.NewProductionConfig()
	}

	// 在 web 应用下，outputs 的配置是相同的
	cfg.OutputPaths = []string{
		"stdout",
	}

	return cfg.Build()
}
