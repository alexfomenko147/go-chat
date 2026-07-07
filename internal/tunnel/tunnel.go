package tunnel

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
)

func RunServer(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	defer ln.Close()

	nextID := 1
	var mu sync.Mutex

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		go func(c net.Conn) {
			defer c.Close()

			var localPort uint32
			if err := binary.Read(c, binary.BigEndian, &localPort); err != nil {
				return
			}

			mu.Lock()
			id := nextID
			nextID++
			publicPort := 20000 + id
			mu.Unlock()

			if err := binary.Write(c, binary.BigEndian, uint32(publicPort)); err != nil {
				return
			}

			fmt.Printf("tunnel: %s -> port %d\n", c.RemoteAddr(), publicPort)

			clientAddr := c.RemoteAddr().(*net.TCPAddr)
			for {
				_, err := c.Write([]byte{1})
				if err != nil {
					return
				}

				dataConn, err := net.Dial("tcp", clientAddr.String())
				if err != nil {
					return
				}

				go func() {
					io.Copy(dataConn, c)
					dataConn.Close()
				}()
				go func() {
					io.Copy(c, dataConn)
					dataConn.Close()
				}()
			}
		}(conn)
	}
}

func RunClient(serverAddr string, localPort int) (int, error) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return 0, fmt.Errorf("connect to tunnel server: %w", err)
	}

	if err := binary.Write(conn, binary.BigEndian, uint32(localPort)); err != nil {
		conn.Close()
		return 0, fmt.Errorf("send local port: %w", err)
	}

	var publicPort uint32
	if err := binary.Read(conn, binary.BigEndian, &publicPort); err != nil {
		conn.Close()
		return 0, fmt.Errorf("read public port: %w", err)
	}

	go func() {
		defer conn.Close()
		buf := make([]byte, 1)
		for {
			if _, err := io.ReadFull(conn, buf); err != nil {
				return
			}

			go func() {
				dataConn, err := net.Dial("tcp", serverAddr)
				if err != nil {
					return
				}
				defer dataConn.Close()

				localConn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", localPort))
				if err != nil {
					return
				}
				defer localConn.Close()

				var wg sync.WaitGroup
				wg.Add(2)
				go func() {
					io.Copy(dataConn, localConn)
					wg.Done()
				}()
				go func() {
					io.Copy(localConn, dataConn)
					wg.Done()
				}()
				wg.Wait()
			}()
		}
	}()

	return int(publicPort), nil
}
