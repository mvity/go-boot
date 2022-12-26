/*
 * Copyright © 2021 - 2022 vity <vityme@icloud.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 */

package app

// Context 标识定义
const (
	GormContext string = "Ctx_Gorm" // Gorm Context 标识
	GinContext  string = "Ctx_Gin"  // Gin Context 标识
)

// Gin 标识定义
const (
	GinTime    = "Gin_Time" // Gin请求时间，服务端接收到请求时间
	GinBody    = "Gin_Body" // Gin请求内容，请求体解密后内容
	GinData    = "Gin_Data" // Gin请求内容，请求体原始内容，未解密
	GinLogger  = "Gin_Log"  // 是否记录请求日志
	GinEncrypt = "Gin_Aes"  // 是否进行AES加密
	GinUserID  = "Gin_Uid"  // 当前请求用户ID
)

// Gin 错误定义
const (
	GinSignError        int8 = 10 // 签名校验错误
	GinTimeError        int8 = 11 // 客户端时间错误
	GinVersionError     int8 = 12 // 客户端版本过低
	GinServerClosed     int8 = 13 // 服务端暂停服务
	GinParamError       int8 = 14 // 参数校验失败
	GinAuthError        int8 = 15 // 身份鉴权失败
	GinActionError      int8 = 16 // 无接口操作权限
	GinDataError        int8 = 17 // 无数据操作权限
	GinNotFound         int8 = 18 // 未找到匹配接口
	GinMustParam        int8 = 19 // 缺少必要参数
	GinApiStop          int8 = 20 // 接口已停用
	GinApiLock          int8 = 21 // 请求业务处理中
	GinSysError         int8 = 22 // 系统处理异常
	GinTransactionError int8 = 99 // 数据处理异常
)
