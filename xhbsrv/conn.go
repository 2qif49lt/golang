package xhbsrv

import (
	//	"errors"
	"io"
	"net"
	"sync"
	"time"
)

type xhbconn struct {
	conn   net.Conn
	connid int
	srv    *Srv

	squeue  chan []byte
	stop    chan bool
	clsonce *sync.Once
}

func (c *xhbconn) close(err error) {
	c.clsonce.Do(func() {
		addr := c.conn.RemoteAddr()
		c.srv.delconn(c.connid)
		close(c.stop)
		c.conn.Close()
		c.srv.cb.OnConnClosed(addr, c.connid, err)
	})
}
func (c *xhbconn) send(buff []byte) error {
	select {
	case c.squeue <- buff:
	case <-c.stop:
		return ErrConnQueueClose
	default:
		return ErrConnQueueFull
	}
	return nil
}
func handleconnrecv(c *xhbconn) {
	addr := c.conn.RemoteAddr()
	var errcb error = nil
	defer c.close(errcb)
	buff := make([]byte, c.srv.PackMaxLen)
	datalen := 0

	//	timeoutcount := 0
	//	timeoutmax := c.srv.TimeOutSecond * 1000 / defaultReadMiliSecond

	for {
		select {
		case <-c.srv.stopchan:
			errcb = ErrCloseSrvDown
		case <-c.stop:
			errcb = ErrClosedSelf
		default:
			c.conn.SetReadDeadline(time.Now().Add(time.Duration(c.srv.TimeOutSecond) * time.Second))
			n, e := c.conn.Read(buff[datalen:])

			if e != nil {
				if e == io.EOF {
					errcb = ErrClosedActiveByPeer
				} else if neterr, ok := e.(net.Error); ok && neterr.Timeout() {
					//		timeoutcount++
					//				if timeoutcount > timeoutmax {
					errcb = ErrClosedTimeOut
					//			}
				} else {
					errcb = e
				}
			}

			if errcb != nil {
				return
			}
			//		timeoutcount = 0

			datalen += n
			databeg := 0
			dataend := datalen

			for {
				pack, e := c.srv.netunpack(buff[databeg:dataend])
				if e != nil {
					errcb = ErrClosedTimeOut
					return
				}
				if pack == nil {
					break
				} else {
					datalen -= len(pack)
					databeg += len(pack)

					if bhdle := c.srv.cb.OnRecvMsg(addr, c.connid, pack); bhdle == false {
						errcb = ErrClosedSelf
						return
					}
				}
			}
			if databeg != 0 {
				copy(buff, buff[databeg:dataend])
			}
		}
	}

}

func handleconnsend(c *xhbconn) {
	var errcb error = nil
	defer c.close(errcb)

	for {
		select {
		case <-c.srv.stopchan:
			errcb = ErrCloseSrvDown
		case <-c.stop:
			errcb = ErrClosedSelf
		case data, ok := <-c.squeue:
			{
				if !ok || data == nil {
					return
				}
				nlen := len(data)
				nwrt := 0
				for nwrt < nlen {
					per, err := c.conn.Write(data[nwrt:])
					if err != nil {
						errcb = err
						break
					}
					nwrt += per
				}
			}
		}

		if errcb != nil {
			break
		}
	}
}
