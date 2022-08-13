package webmon

import "net"

type TCPProbe struct {
	addr string
}

func NewTCPProbe(addr string) (probe *TCPProbe, err error) {
	return &TCPProbe{addr: addr}, nil
}

func (p *TCPProbe) Ping() error {
	conn, err := net.Dial("tcp", p.addr)
	if err != nil {
		return err
	}

	defer conn.Close()

	return nil
}
