package errcode

import (
	"encoding/json"
	"fmt"
	"path"
	"runtime"
)

type AppError struct {
	code     int    `json:"code"`
	msg      string `json:"msg"`
	cause    error  `json:"cause"`
	occurred string `json:"occurred"` // 保存由底层错误导致AppErr发生时的位置
}

// 实现error接口  Error方法变成error 类型
func (e *AppError) Error() string {
	if e == nil {
		return ""
	}

	formattedErr := struct {
		Code     int    `json:"code"`
		Msg      string `json:"msg"`
		Cause    string `json:"cause"`
		Occurred string `json:"occurred"`
	}{
		Code:     e.Code(),
		Msg:      e.Msg(),
		Occurred: e.occurred,
	}

	if e.cause != nil {
		formattedErr.Cause = e.cause.Error()
	}
	errByte, _ := json.Marshal(formattedErr)
	return string(errByte)
}

func (e *AppError) String() string {
	return e.Error()
}

func (e *AppError) Code() int {
	return e.code
}

func (e *AppError) Msg() string {
	return e.msg
}

// WithCause 在逻辑执行中出现错误, 比如dao层返回的数据库查询错误
// 可以在领域层返回预定义的错误前附加上导致错误的基础错误。
// 如果业务模块预定义的错误码比较详细, 可以使用这个方法, 反之错误码定义的比较笼统建议使用Wrap方法包装底层错误生成项目自定义Error
// 并将其记录到日志后再使用预定义错误码返回接口响应
func (e *AppError) WithCause(err error) *AppError {
	e.cause = err
	e.occurred = getAppErrOccurredInfo()
	return e
}

// newError 创建新的应用错误实例
// 参数说明：
//   - code: 错误码，必须大于等于0
//   - msg: 错误信息
//
// 返回值：
//   - *AppError: 返回新创建的错误实例
//
// 说明：
//  1. 错误码必须大于等于0
//  2. 错误码不能重复，重复时会触发panic
//  3. 新创建的错误码会被记录到全局错误码映射中
func newError(code int, msg string) *AppError {
	// 检查错误码是否有效
	if code < 0 {
		panic("错误码必须大于等于0")
	}
	// 检查错误码是否重复
	if _, exists := codes[code]; exists {
		panic(fmt.Sprintf("错误码 %d 已存在，请使用其他错误码", code))
	}

	// 记录新的错误码
	codes[code] = struct{}{}

	// 创建并返回错误实例
	return &AppError{
		code: code,
		msg:  msg,
	}
}

//func newError(code int, msg string) *AppError {
//	if code > -1 {
//		if _, duplicated := codes[code]; duplicated {
//			panic(fmt.Sprintf("预定义错误码 %d 不能重复, 请检查后更换", code))
//		}
//		codes[code] = struct{}{}
//	}
//	return &AppError{code: code, msg: msg}
//}

// Wrap 用于逻辑中包装底层函数返回的error 和 WithCause 一样都是为了记录错误链条
// 该方法生成的error 用于日志记录, 返回响应请使用预定义好的error
//
//	func Wrap(msg string, err error) *AppError {
//		if err == nil {
//			return nil
//		}
//		appErr := &AppError{code: -1, msg: msg, cause: err}
//		appErr.occurred = getAppErrOccurredInfo()
//		return appErr
//	}
//
// ，当你拿到一个error不确定它该是什么错误，你就用这个Wrap方法包装
func Wrap(msg string, err error) *AppError {
	if err == nil {
		return nil
	}
	appErr := &AppError{code: -1, msg: msg, cause: err}
	appErr.occurred = getAppErrOccurredInfo()
	return appErr
}

// getAppErrOccurredInfo 获取项目中调用Wrap或者WithCause方法时的程序位置, 方便排查问题
func getAppErrOccurredInfo() string {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return ""
	}
	file = path.Base(file)
	funcName := runtime.FuncForPC(pc).Name()
	triggerInfo := fmt.Sprintf("func: %s, file: %s, line: %d", funcName, file, line)
	return triggerInfo
}
