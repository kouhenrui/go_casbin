package test

import (
	"errors"
	"go_casbin/internal/logger"
	"testing"
)

func TestLogger(t *testing.T) {
	logger.Init(nil)
	var tt=map[string]interface{}{
		"tt":"test",
	}
	logger.ErrorWithObject("test",tt,errors.New("test"))
	// logger.Info("test",logger.String("test", "test"))
	// logger.InfoMap("map",tt)
	// logger.InfoObject("object",tt)
}