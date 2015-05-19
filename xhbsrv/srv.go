package xhbsrv

import (
	. "fmt"
	"net"
	"sync"
	"time"
)

type Srv struct {
	PackSizePos      int
	PackMaxLen       int
	TimeOutSecond    int
	MaxConn          int // 最多支持链接
	MaxSendQueueSize int // 单个发送的消息队列长度
	stop             bool
	stopchan         chan bool // 服务停止
	lis              net.Listener

	numlock    *sync.Mutex //   connmap connincrid
	connmap    map[int]*xhbconn
	connincrid int // 自增id

	cb CallBack
}

func New() *Srv {
	return &Srv{
		defaultPackSizePos,
		defaultPackMaxLen,
		defaultTimeOutSecond,
		defaultMaxConn,
		defaultMaxSendQueueSize,

		false,
		make(chan bool, 1),
		nil,

		&sync.Mutex{},
		make(map[int]*xhbconn),
		0,

		&defaultcb{},
	}
}

func (s *Srv) Set(sizepos, maxlen, timeout, maxconn int) {
	s.PackSizePos = sizepos
	s.PackMaxLen = maxlen
	s.TimeOutSecond = timeout
	s.MaxConn = maxconn
}

func (s *Srv) run() {
	defer s.lis.Close()
	for !s.stop {
		conn, err := s.lis.Accept()
		if err != nil {
			return
		}
		connid := s.incrconnid()
		connaddr := conn.RemoteAddr()

		if s.ConnNum() > s.MaxConn {
			conn.Close()
			s.cb.OnConnClosed(connaddr, connid, ErrClosedSrvFull)
		} else {
			newconn := &xhbconn{
				conn,
				connid,
				s,
				make(chan []byte, s.MaxSendQueueSize),
				make(chan bool, 1),
				&sync.Once{},
			}
			s.addmap(connid, newconn)

			if s.cb.OnNewConn(connaddr, connid) == false {
				conn.Close()
				s.delmap(connid)
				continue
			}

			go handleconnrecv(newconn)
			go handleconnsend(newconn)
		}

	}
}

// 重新启动服务 应该重新的New,然后RunOn
func (s *Srv) RunOn(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.lis = l
	s.stop = false

	go s.run()

	return nil
}

func (s *Srv) Stop() {
	if s.stop == true {
		return
	}
	s.stop = true
	close(s.stopchan)
	c := time.Tick(1 * time.Second)
	for w := range c {
		_ = w
		n := s.ConnNum()
		Println("current connection number:", n)
		if n <= 0 {
			break
		}
	}
}
func (s *Srv) incrconnid() int {
	s.numlock.Lock()
	defer s.numlock.Unlock()
	s.connincrid++
	return s.connincrid
}
func (s *Srv) addmap(connid int, c *xhbconn) {
	s.numlock.Lock()
	defer s.numlock.Unlock()
	s.connmap[connid] = c
}
func (s *Srv) delmap(connid int) {
	s.numlock.Lock()
	defer s.numlock.Unlock()
	if _, exist := s.connmap[connid]; exist {
		delete(s.connmap, connid)
	}
}
func (s *Srv) ConnNum() int {
	s.numlock.Lock()
	defer s.numlock.Unlock()
	return len(s.connmap)
}
func (s *Srv) CloseConn(connid int) {
	s.numlock.Lock()
	defer s.numlock.Unlock()

	if conn, exist := s.connmap[connid]; exist {
		delete(s.connmap, connid)
		conn.close(ErrClosedSelf)
	}
}
func (s *Srv) Addr(connid int) net.Addr {
	s.numlock.Lock()
	defer s.numlock.Unlock()

	if conn, exist := s.connmap[connid]; exist {
		return conn.conn.RemoteAddr()
	} else {
		return nil
	}
}

func (s *Srv) SendCmd(connid int, cmd int, body []byte) error {
	msg, err := s.netpack(cmd, body)
	if err != nil {
		return err
	}

	return s.Send(connid, msg)
}
func (s *Srv) Send(connid int, data []byte) error {
	s.numlock.Lock()
	conn, exist := s.connmap[connid]
	s.numlock.Unlock()
	if exist {
		return conn.send(data)
	}

	return ErrConnDoNotExist
}

func (s *Srv) Attach(c net.Conn) int {
	connid := s.incrconnid()
	s.attach(c, connid)
	return connid
}
func (s *Srv) attach(c net.Conn, connid int) {
	s.numlock.Lock()
	defer s.numlock.Unlock()

	newconn := &xhbconn{
		c,
		connid,
		s,
		make(chan []byte, s.MaxSendQueueSize),
		make(chan bool, 1),
		&sync.Once{},
	}
	s.connmap[connid] = newconn

	go handleconnrecv(newconn)
	go handleconnsend(newconn)
}

// 立即返回connid, 连接完成后会回调.
func (s *Srv) Connect(addr string, timeout time.Duration) (int, error) {
	tcpaddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return 0, err
	}
	connid := s.incrconnid()

	go func(addr net.Addr, timeout time.Duration, connid int) {
		conn, err := net.DialTimeout("tcp", addr.String(), timeout)
		if err == nil {
			s.attach(conn, connid)
		}
		s.cb.OnConnectRemote(addr, connid, err)
	}(tcpaddr, timeout, connid)
	return connid, nil
}

// 阻塞,不会得到回调
func (s *Srv) ConnectBlock(addr string, timeout time.Duration) (int, error) {
	_, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return 0, err
	}
	connid := s.incrconnid()
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err == nil {
		s.attach(conn, connid)
		return connid, nil
	}

	return 0, err
}
