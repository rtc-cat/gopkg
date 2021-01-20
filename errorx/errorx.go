package errorx

import "log"

const (
	ErrCodeNotSet ErrCode = iota
	ErrCodeUnexpected
)

const (
	ErrLogLevelInfo  ErrLogLevel = "info"
	ErrLogLevelDebug ErrLogLevel = "debug"
	ErrLogLevelWarn  ErrLogLevel = "warn"
	ErrLogLevelError ErrLogLevel = "error"
)

type (
	ErrCode      int    // 自定义错误码
	ErrOperation string // 自定义操作
	ErrLogLevel  string // 自定义日志级别
)

// Errorx 带有错误码和错误操作的扩展类型
type Errorx struct {
	ErrCode      ErrCode
	ErrOperation ErrOperation
	ErrLogLevel  ErrLogLevel // 自定义日志级别
	Err          error       // 被封装的错误
}

// Error 实现error接口
func (err *Errorx) Error() string {
	if err.Err == nil {
		return ""
	}
	return err.Err.Error()
}

// E 创建Errorx,参数可以为ErrCode,ErrMessage,ErrOperation
func E(args ...interface{}) error {
	e := &Errorx{}
	for _, argi := range args {
		switch arg := argi.(type) {
		case ErrCode:
			e.ErrCode = arg
		case int:
			e.ErrCode = ErrCode(arg)
		case string:
			e.ErrOperation = ErrOperation(arg)
		case ErrOperation:
			e.ErrOperation = arg
		case ErrLogLevel:
			e.ErrLogLevel = arg
		case error:
			e.Err = arg
		default:
			panic("unknown errorx field")
		}
	}
	return e
}

// Code 递归寻找第一个非0的错误码
func Code(err error) ErrCode {
	e, ok := err.(*Errorx)
	if !ok {
		return ErrCodeUnexpected
	}
	if e.ErrCode != ErrCodeNotSet {
		return e.ErrCode
	}
	return Code(e.Err)
}

// Operations 返回操作的stack,递归合并操作信息
func Operations(err error) []ErrOperation {
	e, ok := err.(*Errorx)
	if !ok {
		return nil
	}
	result := []ErrOperation{e.ErrOperation}
	subErr, ok := e.Err.(*Errorx)
	if !ok {
		return result
	}
	result = append(result, Operations(subErr)...)
	return result
}

// LogLevel 获取日志级别
// Deprecated 应该自己实现
func LogLevel(err error) ErrLogLevel {
	errx, ok := err.(*Errorx)
	if !ok {
		return ErrLogLevelError
	}
	switch errx.ErrLogLevel {
	case ErrLogLevelInfo:
		return ErrLogLevelInfo
	case ErrLogLevelDebug:
		return ErrLogLevelDebug
	case ErrLogLevelWarn:
		return ErrLogLevelWarn
	case "":
		return LogLevel(errx.Err)
	default:
		return ErrLogLevelError
	}
}

// Log 简单实现如何分级打印错误信息
// Deprecated 应该自己实现
func Log(err error) {
	errx, ok := err.(*Errorx)
	if !ok {
		log.Println(err)
		return
	}
	switch errx.ErrLogLevel {
	case "info":
		log.Printf("[info ] %v: %v", Operations(errx), errx)
	case "debug":
		log.Printf("[debug] %v: %v", Operations(errx), errx)
	case "warn", "warning":
		log.Printf("[warn ] %v: %v", Operations(errx), errx)
	default:
		log.Printf("[error] %v: %v", Operations(errx), errx)
	}
}
