package core

import (
	"fmt"

	"github.com/sunmi-OS/gocore/v2/utils"

	"github.com/sunmi-OS/gocore/v2/api/ecode"
)

func WrapError(err error, format string, a ...any) error {
	if err == nil {
		return ecode.NewV2(ecode.SystemErrorCode, fmt.Sprintf(format, a...))
	}

	e2 := ecode.FromError(err)
	newMsg := fmt.Sprintf(format, a...) + ", " + e2.Message()
	return ecode.NewV2(e2.Code(), newMsg)
}

// WrapDebugError 仅在非线上环境才添加该错误信息
func WrapDebugError(err error, format string, a ...any) error {
	if utils.IsRelease() {
		return err
	}
	return WrapError(err, "[DebugError]"+format, a...)
}
