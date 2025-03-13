package errcode

import "net/http"

type errCode int64

type Status struct {
	Code errCode `json:"code"`
	Msg  string  `json:"msg"`
}

const (
	AllOK errCode = 10000 // 成功

	RedisRunTimeError errCode = 20000 // Redis运行时错误
	RedisNilError     errCode = 20004 // Redis返回空值

	SQLRunTimeError errCode = 30000 // SQL运行时错误
	SQLNilError     errCode = 30004 // SQL返回空值

	PermissionDenied errCode = 40000 // 当前用户权限不足

	ProgramError errCode = 50000 // 程序运行出错

)

var errorMsg = map[errCode]string{
	AllOK:             "成功",
	RedisRunTimeError: "Redis运行时错误",
	RedisNilError:     "Redis返回空值",
	SQLRunTimeError:   "SQL运行时错误",
	SQLNilError:       "SQL返回空值",

	PermissionDenied: "当前用户权限不足",
	ProgramError:     "程序运行出错",
}

var httpCode = map[errCode]int{
	AllOK:             http.StatusOK,
	RedisRunTimeError: http.StatusInternalServerError,
	RedisNilError:     http.StatusNotFound,
	SQLRunTimeError:   http.StatusInternalServerError,
	SQLNilError:       http.StatusNotFound,
	PermissionDenied:  http.StatusUnauthorized,
	ProgramError:      http.StatusInternalServerError,
}

func GetMsg(code errCode) string {
	return errorMsg[code]
}

func GetHttpCode(code errCode) int {
	return httpCode[code]
}

func GetStatus(code errCode) Status {
	return Status{Code: code, Msg: errorMsg[code]}
}
