package patterns

/*
snippet demonstrating worker pool pattern.

The worker pool has a simple mechanism, there is a queue and there are a fixed
number of workers. As work comes to queue, its size is increased. As work gets
send to workers, the size of the queue decreases. When queue reaches its
maximum size, we stop accepting work.

This is a very common pattern in real-life scenarios. You have a limited
number of resources for a task, and you want to balance out the incoming work
to existing resources.
*/

import (
	"fmt"
	"net"
	"os"
	"regexp"
)

var r, _ = regexp.Compile("GET (.+) HTTP/1.1\r\n")

func handleHttpRequest(conn net.Conn) {
	buff := make([]byte, 1024)

	size, _ := conn.Read(buff)
	if r.Match(buff[:size]) {
		// if request is valid proceed with work
		file, err := os.ReadFile(fmt.Sprintf("../asset/%s",
			r.FindSubmatch(buff[:size])[1]))
		if err == nil {
			conn.Write([]byte(fmt.Sprintf("HTTP/1.1/200 OK\r\nContent-Length: %d\r\n\r\n", len(file))))
			conn.Write(file)
		} else {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n<html>Not Found</html>"))
		}
	} else {
		conn.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n\r\n"))
	}
	conn.Close()
}

func StartHttpWorkers(n int, incomingConnections <-chan net.Conn) {
	for i := 0; i < n; i++ {
		go func() {
			for request := range incomingConnections {
				handleHttpRequest(request)
			}
		}()
	}
}

func WorkerPoolMain() {
	incomingConn := make(chan net.Conn)
	nbWorkers := 3
	StartHttpWorkers(nbWorkers, incomingConn)
	server, _ := net.Listen("tcp", "localhost:8888")
	defer server.Close()
	for {
		conn, _ := server.Accept()
		select {
		case incomingConn <- conn: // send work to queue
		default:
			fmt.Println("Server is busy")
			conn.Write([]byte("HTTP/1.1 429 Too Many Requests\r\n\r\n" +
				"<html>Busy</html>\n"))
			conn.Close()
		}
	}
}
