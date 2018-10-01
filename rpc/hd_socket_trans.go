/**********************************************************\
|                                                          |
|                          hprose                          |
|                                                          |
| Official WebSite: http://www.hprose.com/                 |
|                   http://www.hprose.org/                 |
|                                                          |
\**********************************************************/
/**********************************************************\
 *                                                        *
 * rpc/hd_socket_trans.go                                 *
 *                                                        *
 * hprose half duplex socket transport for Go.            *
 *                                                        *
 * LastModified: Jan 7, 2017                              *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import (
	"net"
	"time"
)

type halfDuplexSocketTransport struct {
	idleTimeout time.Duration
	poolSize    int
	createConn  func() (net.Conn, error)
}

func newHalfDuplexSocketTransport() (hd *halfDuplexSocketTransport) {
	hd = new(halfDuplexSocketTransport)
	hd.idleTimeout = 0
	hd.poolSize = 0
	return
}

func (hd *halfDuplexSocketTransport) setCreateConn(createConn func() (net.Conn, error)) {
	hd.createConn = createConn
}

// IdleTimeout returns the conn pool idle timeout of hprose socket client
func (hd *halfDuplexSocketTransport) IdleTimeout() time.Duration {
	return hd.idleTimeout
}

// SetIdleTimeout sets the conn pool idle timeout of hprose socket client
func (hd *halfDuplexSocketTransport) SetIdleTimeout(timeout time.Duration) {
	hd.idleTimeout = timeout
}

// MaxPoolSize returns the max conn pool size of hprose socket client
func (hd *halfDuplexSocketTransport) MaxPoolSize() int {
	return hd.poolSize
}

// SetMaxPoolSize sets the max conn pool size of hprose socket client
func (hd *halfDuplexSocketTransport) SetMaxPoolSize(size int) {
	hd.poolSize = size
}

func (hd *halfDuplexSocketTransport) close() {
}

func (hd *halfDuplexSocketTransport) closeConn(conn net.Conn) {
	conn.Close()
}

func (hd *halfDuplexSocketTransport) sendAndReceive(
	data []byte, context *ClientContext) ([]byte, error) {
	conn, err := hd.createConn()

	if err == nil {
		err = hdSendData(conn, data)
	}
	if err == nil {
		data, err = hdRecvData(conn, data)
	}
	if err == nil {
		err = conn.SetDeadline(time.Time{})
	}
	hd.closeConn(conn)
	if err != nil {
		return nil, err
	}
	return data, nil
}
