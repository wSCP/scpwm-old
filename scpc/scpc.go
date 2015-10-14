package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

var (
	socketEnv     string = "SCPWM_SOCKET"
	socketPathTpl string = "/tmp/scpwm%s_%d_%d-socket"
	socketPath    string

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

func outcome(response string, sent []byte) {
	logger.Println(fmt.Sprintf("euclid returned: `%s` to sent message: `%s`", response, sent))
}

func response(c io.Reader, sent []byte) {
	buf := make([]byte, 1024)
	for {
		n, err := c.Read(buf[:])
		if err != nil {
			logger.Println(err)
			return
		}
		resp := string(buf[0:n])
		if resp != "success" {
			switch resp {
			case "failure", "syntax", "unknown", "length":
				outcome(resp, sent)
				return
			default:
				outcome("UNCATEGORIZABLE", sent)
				return
			}
		}
	}
}

func main() {
	t := "unix"
	laddr := net.UnixAddr{"/tmp/scpc", t}
	conn, err := net.DialUnix(t, &laddr, &net.UnixAddr{socketPath, t})
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	defer os.Remove("/tmp/scpc")

	sent := send(os.Args[1:])
	_, err = conn.Write(sent)
	if err != nil {
		panic(err)
	}

	response(conn, sent)
}

func init() {
	flag.StringVar(&socketEnv, "e", socketEnv, "Read the socket from the given env variable")
	flag.StringVar(&socketPath, "p", defaulSocketPath(), "Read the socket from the given path.")
	flag.Parse()
}
