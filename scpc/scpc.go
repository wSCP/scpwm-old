package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

var (
	socketEnv       string = "SCPWM_SOCKET"
	socketPathTpl   string = "/tmp/scpwm%s_%d_%d-socket"
	socketPath      string
	timeoutDuration int = 3

	logger = log.New(os.Stderr, "[SCPC] ", log.Ldate|log.Lmicroseconds)
)

func defaulSocketPath() string {
	var socketPath = os.Getenv(socketEnv)
	if socketPath == "" {
		socketPath = fmt.Sprintf(socketPathTpl, "", 0, 0)
	}
	return socketPath
}

func removable(index int, remove []int) bool {
	for _, r := range remove {
		if index == r {
			return true
		}
	}
	return false
}

func send(args []string) []byte {
	var r []int
	for i, v := range args {
		switch v {
		case "-p", "-e":
			r = append(r, i, i+1)
		}
	}

	var send []string

	for i, v := range args {
		if !removable(i, r) {
			send = append(send, v)
		}
	}

	idx := len(send) - 1

	b := new(bytes.Buffer)
	for i, v := range send {
		b.WriteString(v)
		if i != idx {
			b.WriteRune(' ')
		}
	}
	return b.Bytes()
}

func outcome(r Response, sent []byte) {
	logger.Println(fmt.Sprintf("euclid returned: `%s` to sent message: `%s`", r, sent))
}

type Response int

const (
	unknown Response = iota
	failure
	syntax
	length
	timeout
	success
)

func (r Response) String() string {
	switch r {
	case unknown:
		return "unknown"
	case failure:
		return "failure"
	case syntax:
		return "syntax"
	case length:
		return "length"
	case timeout:
		return "timeout"
	case success:
		return "success"
	}
	return ""
}

var toConstResponse = map[string]Response{
	"unknown": unknown,
	"failure": failure,
	"syntax":  syntax,
	"length":  length,
	"timeout": timeout,
	"success": success,
}

func response(c io.Reader, sent []byte) Response {
	buf := make([]byte, 1024)
	t := time.After(3 * time.Second)
	for {
		select {
		case <-t:
			logger.Println("timeout, no response from euclid after %d seconds", timeoutDuration)
			return timeout
		default:
			n, err := c.Read(buf[:])
			if err != nil {
				logger.Println(err)
				return failure
			}
			if resp, ok := toConstResponse[string(buf[0:n])]; ok {
				switch resp {
				case failure, syntax, length:
					outcome(resp, sent)
					return resp
				case success:
					return resp
				default:
					outcome(unknown, sent)
					return resp
				}
			}
		}
	}
	return timeout
}

func exit(c *net.UnixConn, code int) {
	c.Close()
	os.Remove("/tmp/scpc")
	os.Exit(code)
}

func main() {
	t := "unix"
	laddr := net.UnixAddr{"/tmp/scpc", t}
	conn, err := net.DialUnix(t, &laddr, &net.UnixAddr{socketPath, t})
	if err != nil {
		panic(err)
	}

	sent := send(os.Args[1:])
	_, err = conn.Write(sent)
	if err != nil {
		panic(err)
	}

	rsp := response(conn, sent)
	switch rsp {
	case success:
		exit(conn, 0)
	case failure, length, timeout:
		exit(conn, 1)
	case syntax:
		exit(conn, 2)
	default:
		exit(conn, 3)
	}
}

func init() {
	flag.StringVar(&socketEnv, "e", socketEnv, "Read the socket from the given env variable")
	flag.StringVar(&socketPath, "p", defaulSocketPath(), "Read the socket from the given path.")
	flag.IntVar(&timeoutDuration, "t", timeoutDuration, "Wait only specified seconds euclid's response, default is 3")
	flag.Parse()
}
