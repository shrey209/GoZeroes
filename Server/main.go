package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"golang.org/x/sys/unix"
)

const maxEvents = 1024

func main() {
	port := ":8080"

	listener, err := net.Listen("tcp", port)
	if err != nil {
		slog.Error("failed to initialize the listener")
		os.Exit(1)
	}
	defer listener.Close()

	slog.Info("server started at" + port)

	tcpListener := listener.(*net.TCPListener)
	listenerFd, err := tcpListener.File()
	if err != nil {
		slog.Error("failed to create the file descriptor")
		os.Exit(1)
	}
	defer listenerFd.Close()

	epfd, err := unix.EpollCreate1(0)
	if err != nil {
		slog.Error("failed to create epoll instance")
		os.Exit(1)
	}
	defer unix.Close(epfd)

	event := &unix.EpollEvent{Events: unix.EPOLLIN, Fd: int32(listenerFd.Fd())}
	if err := unix.EpollCtl(epfd, unix.EPOLL_CTL_ADD, int(listenerFd.Fd()), event); err != nil {
		slog.Error("failed to add the file descriptor")
		os.Exit(1)
	}

	events := make([]unix.EpollEvent, maxEvents)

	for {
		n, err := unix.EpollWait(epfd, events, -1)
		if err != nil {
			slog.Error("failed to wait for epoll events", err)
			continue
		}

		for i := 0; i < n; i++ {
			fd := int(events[i].Fd)
			if fd == int(listenerFd.Fd()) {
				conn, err := listener.Accept()
				if err != nil {
					slog.Error("failed to accept connection", err)
					continue
				}

				tcpConn := conn.(*net.TCPConn)
				connFd, err := tcpConn.File()
				if err != nil {
					slog.Error("failed to get connection file descriptor", err)
					conn.Close()
					continue
				}
				defer connFd.Close()

				connEvent := &unix.EpollEvent{Events: unix.EPOLLIN, Fd: int32(connFd.Fd())}
				if err := unix.EpollCtl(epfd, unix.EPOLL_CTL_ADD, int(connFd.Fd()), connEvent); err != nil {
					slog.Error("failed to add connection to epoll", err)
					conn.Close()
				}
				slog.Info("new connection accepted", "address", conn.RemoteAddr().String())
			} else {
				//read
				buffer := make([]byte, 1024)
				n, err := unix.Read(fd, buffer)
				if err != nil || n == 0 {
					slog.Info("connection closed or read error", "fd", fd)
					unix.Close(fd)
					continue
				}
				//tokens
				str := string(buffer[:n])
				token, err := ParseRESP(str)
				if err != nil {
					slog.Error("error parsing")
					data := []byte("ERR" + err.Error())
					writeToConnection(data, fd)
					unix.Close(fd)
				}
				//execution
				fmt.Println(token)
				response, err := handleCommand(token)
				if err != nil {
					slog.Error("error in executing")
					data := []byte("ERR" + err.Error())
					writeToConnection(data, fd)
					unix.Close(fd)
				}

				data := []byte(response)
				writeToConnection(data, fd)
			}
		}
	}
}

func writeToConnection(data []byte, fd int) error {
	totalWritten := 0

	for totalWritten < len(data) {
		n, err := unix.Write(fd, data[totalWritten:])
		if err != nil {
			slog.Error("Failed to write data to connection", "file_descriptor", fd, "error", err)
			unix.Close(fd)
			return err
		}
		totalWritten += n
	}

	return nil
}
