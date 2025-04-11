package errcode

// 统一错误码定义
// 错误码规则：
// 1. 错误码为5位数字
// 2. 第1位表示错误级别：1为系统级错误，2为业务级错误
// 3. 第2-3位表示服务模块：10为通用，11为用户服务，12为短链接服务
// 4. 第4-5位表示具体错误码

// 系统级错误码 (10000-19999)
const (
	// 通用系统错误 (10000-10999)
	ServerError        = 10000 // 服务器内部错误
	InvalidParams      = 10001 // 参数错误
	Unauthorized       = 10002 // 未授权
	NotFound           = 10003 // 资源不存在
	TooManyRequests    = 10004 // 请求过多
	Timeout            = 10005 // 请求超时
	ServiceUnavailable = 10006 // 服务不可用

	// 数据库错误 (11000-11999)
	DBError           = 11000 // 数据库错误
	DBConnectionError = 11001 // 数据库连接错误
	DBQueryError      = 11002 // 数据库查询错误
	DBInsertError     = 11003 // 数据库插入错误
	DBUpdateError     = 11004 // 数据库更新错误
	DBDeleteError     = 11005 // 数据库删除错误

	// 缓存错误 (12000-12999)
	CacheError           = 12000 // 缓存错误
	CacheConnectionError = 12001 // 缓存连接错误
	CacheSetError        = 12002 // 缓存设置错误
	CacheGetError        = 12003 // 缓存获取错误
	CacheDeleteError     = 12004 // 缓存删除错误
)

// 业务级错误码 (20000-29999)
const (
	// 用户服务错误码 (21000-21999)
	UserNotFound      = 21000 // 用户不存在
	UserAlreadyExists = 21001 // 用户已存在
	PasswordIncorrect = 21002 // 密码错误
	TokenInvalid      = 21003 // Token无效
	TokenExpired      = 21004 // Token过期
	LoginFailed       = 21005 // 登录失败
	RegisterFailed    = 21006 // 注册失败
	LogoutFailed      = 21007 // 登出失败

	// 短链接服务错误码 (22000-22999)
	ShortlinkCreateFailed = 22000 // 创建短链接失败
	ShortlinkNotFound     = 22001 // 短链接不存在
	ShortlinkExpired      = 22002 // 短链接已过期
	ShortlinkInvalid      = 22003 // 短链接无效
	BatchCreateFailed     = 22004 // 批量创建短链接失败
	EmptyURLList          = 22005 // URL列表为空
	TopLinksQueryFailed   = 22006 // 获取热门链接失败
)

// 错误码与HTTP状态码的映射
var ErrCodeToHTTPStatus = map[int]int{
	// 系统级错误
	ServerError:        500,
	InvalidParams:      400,
	Unauthorized:       401,
	NotFound:           404,
	TooManyRequests:    429,
	Timeout:            504,
	ServiceUnavailable: 503,

	// 数据库错误
	DBError:           500,
	DBConnectionError: 500,
	DBQueryError:      500,
	DBInsertError:     500,
	DBUpdateError:     500,
	DBDeleteError:     500,

	// 缓存错误
	CacheError:           500,
	CacheConnectionError: 500,
	CacheSetError:        500,
	CacheGetError:        500,
	CacheDeleteError:     500,

	// 用户服务错误
	UserNotFound:      404,
	UserAlreadyExists: 409,
	PasswordIncorrect: 401,
	TokenInvalid:      401,
	TokenExpired:      401,
	LoginFailed:       401,
	RegisterFailed:    500,
	LogoutFailed:      401,

	// 短链接服务错误
	ShortlinkCreateFailed: 500,
	ShortlinkNotFound:     404,
	ShortlinkExpired:      410,
	ShortlinkInvalid:      400,
	BatchCreateFailed:     500,
	EmptyURLList:          400,
	TopLinksQueryFailed:   500,
}

// 错误码对应的错误信息
var ErrCodeMessages = map[int]string{
	// 系统级错误
	ServerError:        "服务器内部错误",
	InvalidParams:      "参数错误",
	Unauthorized:       "未授权",
	NotFound:           "资源不存在",
	TooManyRequests:    "请求过多",
	Timeout:            "请求超时",
	ServiceUnavailable: "服务不可用",

	// 数据库错误
	DBError:           "数据库错误",
	DBConnectionError: "数据库连接错误",
	DBQueryError:      "数据库查询错误",
	DBInsertError:     "数据库插入错误",
	DBUpdateError:     "数据库更新错误",
	DBDeleteError:     "数据库删除错误",

	// 缓存错误
	CacheError:           "缓存错误",
	CacheConnectionError: "缓存连接错误",
	CacheSetError:        "缓存设置错误",
	CacheGetError:        "缓存获取错误",
	CacheDeleteError:     "缓存删除错误",

	// 用户服务错误
	UserNotFound:      "用户不存在",
	UserAlreadyExists: "用户已存在",
	PasswordIncorrect: "密码错误",
	TokenInvalid:      "Token无效",
	TokenExpired:      "Token过期",
	LoginFailed:       "登录失败",
	RegisterFailed:    "注册失败",
	LogoutFailed:      "登出失败",

	// 短链接服务错误
	ShortlinkCreateFailed: "创建短链接失败",
	ShortlinkNotFound:     "短链接不存在",
	ShortlinkExpired:      "短链接已过期",
	ShortlinkInvalid:      "短链接无效",
	BatchCreateFailed:     "批量创建短链接失败",
	EmptyURLList:          "URL列表为空",
	TopLinksQueryFailed:   "获取热门链接失败",
}

// Error 定义错误结构体
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewError 创建一个新的错误
func NewError(code int) *Error {
	message, ok := ErrCodeMessages[code]
	if !ok {
		message = "未知错误"
	}

	return &Error{
		Code:    code,
		Message: message,
	}
}

// WithMessage 设置自定义错误信息
func (e *Error) WithMessage(message string) *Error {
	e.Message = message
	return e
}

// HTTPStatusCode 获取对应的HTTP状态码
func (e *Error) HTTPStatusCode() int {
	code, ok := ErrCodeToHTTPStatus[e.Code]
	if !ok {
		return 500 // 默认返回500
	}
	return code
}
