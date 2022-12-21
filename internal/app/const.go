package app

/*
 * 全局定义
 */

// 审核状态
const (
	AuditNone   int8 = 0 // 未提交审核
	AuditWait        = 1 // 已提交审核
	AuditPass        = 2 // 已通过审核
	AuditReject      = 3 // 已驳回审核
	AuditCancel      = 9 // 已撤销审核
)

// 内置用户ID
const (
	PlatformID uint64 = 1 // 平台用户ID
	GuestID           = 0 // 访客用户ID
)

// 用户类型
const (
	UserPlatform int8 = 1 // 平台用户，系统内建用户
	UserMember        = 2 // 注册用户，外部注册用户
)
