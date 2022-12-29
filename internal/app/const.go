/*
 * Copyright © 2021 - 2022 vity <vityme@icloud.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 */

package app

/*
 * 全局定义
 */

// 审核状态
const (
	AuditNone   int8 = 0 // 未提交审核
	AuditWait   int8 = 1 // 已提交审核
	AuditPass   int8 = 2 // 已通过审核
	AuditReject int8 = 3 // 已驳回审核
	AuditCancel int8 = 9 // 已撤销审核
)

// 内置用户ID
const (
	PlatformID uint64 = 1 // 平台ID
	GuestID    uint64 = 0 // 访客ID
)

// 用户类型
const (
	UserTypeEmployee int8 = 1 // 平台用户，工作人员用户
	UserTypeMember   int8 = 2 // 注册用户，外部注册用户
)
