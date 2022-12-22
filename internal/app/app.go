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
