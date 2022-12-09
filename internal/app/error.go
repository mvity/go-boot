package app

import "fmt"

// ApiError 接口错误
type ApiError struct {
	ErrCode int8
	Message string
	Origin  error
}

func (e *ApiError) Error() string {
	if e.Origin != nil {
		return fmt.Sprintf("ApiError : [ %v ] %v , Origin : %v", e.ErrCode, e.Message, e.Origin)
	} else {
		return fmt.Sprintf("ApiError : [ %v ] %v", e.ErrCode, e.Message)
	}
}

// MySQLError 数据库错误
type MySQLError struct {
	Message string
	SQL     string
	Origin  error
}

func (e *MySQLError) Error() string {
	if e.Origin != nil {
		return fmt.Sprintf("MySQLError : %v , SQL : [ %v ] , Origin : %v", e.Message, e.SQL, e.Origin)
	} else {
		return fmt.Sprintf("MySQLError : %v , SQL : [ %v ]", e.Message, e.SQL)
	}
}

// RedisError Redis 错误
type RedisError struct {
	Message string
	Command string
	Origin  error
}

func (e *RedisError) Error() string {
	if e.Origin != nil {
		return fmt.Sprintf("RedisError : %v , Command : [ %v ] , Origin : %v", e.Message, e.Command, e.Origin)
	} else {
		return fmt.Sprintf("RedisError : %v , Command : [ %v ]", e.Message, e.Command)
	}
}
