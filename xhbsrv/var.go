package xhbsrv

import (
	"errors"
)

const (
	defaultPackSizePos      = 0
	defaultPackMaxLen       = 100 * 1024
	defaultTimeOutSecond    = 2 * 60
	defaultMaxConn          = 100 * 1000
	defaultMaxSendQueueSize = 100 // 单个发送的消息队列长度
	defaultReadMiliSecond   = 100 // 每次读操作超时ms
)

var (
	ErrClosedActiveByPeer = errors.New("对方主动关闭")
	ErrClosedTimeOut      = errors.New("超时")     // 超时
	ErrClosedSelf         = errors.New("自己手动关闭") // 自己手动关闭
	ErrClosedSrvFull      = errors.New("服务满载")   // 服务满载
	ErrCloseSrvDown       = errors.New("服务停止")

	ErrSizeInvalid = errors.New("头部大小错误") // 头部大小错误

	ErrConnDoNotExist = errors.New("未找到该连接")
	ErrConnQueueFull  = errors.New("发送队列已满")
	ErrConnQueueClose = errors.New("发送队列已经销毁")
)
