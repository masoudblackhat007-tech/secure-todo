package logging

import (
	"fmt"

	"go.uber.org/zap"
)

// Security Logging Policy:
// - هیچ secret، password، token، API key، DSN کامل یا هر دادهٔ حساس دیگری نباید در لاگ چاپ شود.
// - فقط متای غیرحساس مجاز است (مثلاً طول رشته، نوع، یا یک شناسهٔ هش‌شده).
// - مثال‌های صریحِ ممنوع: JWT_SECRET، DB_DSN کامل، access/refresh tokens، session IDs، کدهای یک‌بار مصرف
var logger *zap.Logger

// Init initializes the global logger based on environment.
// env == "production"  -> zap.NewProduction()
// env != "production" -> zap.NewDevelopment()
func Init(env string) (*zap.Logger, error) {
	var (
		l   *zap.Logger
		err error
	)

	if env == "production" {
		l, err = zap.NewProduction(zap.AddCaller())
	} else {
		l, err = zap.NewDevelopment(zap.AddCaller())
	}

	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	logger = l
	return logger, nil
}

// L returns the initialized global logger.
// It panics if Init has not been called.
func L() *zap.Logger {
	if logger == nil {
		panic("logger not initialized: call logging.Init(env) before using logging.L()")
	}
	return logger
}
