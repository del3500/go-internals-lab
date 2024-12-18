package ch03

import (
	"io"
	"net"
	"testing"
)

// A listener in port 12.0.0.1 using a random port

func TestListener(t *testing.T) {
	// net.Listen function accept a network type, and an IP
	// address and port separated by colon.
	listener, err := net.Listen("tcp", ":")
	if err != nil {
		t.Fatal(err)
	}

	defer func() { _ = listener.Close }()

	t.Logf("bound to %q", listener.Addr())
}

func TestDial(t *testing.T) {
	// Create a listener on a random port
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	// A channel to allow communication between the main test function
	// and the goroutines
	done := make(chan struct{})
	go func() {
		// Before a goroutine exits, this defer function will
		// send a signal into the done channel indicating
		// that the goroutine completed its task.
		defer func() { done <- struct{}{} }()

		for {
			// Listener.Accept() call waits for incoming connection
			conn, err := listener.Accept()
			if err != nil {
				t.Log(err)
				return
			}

			go handleConnection(conn, done, t)
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	conn.Close()
	<-done
	listener.Close()
	<-done
}

func handleConnection(c net.Conn, done chan struct{}, t *testing.T) {
	defer func() {
		c.Close()
		done <- struct{}{}
	}()

	buf := make([]byte, 1024)
	for {
		n, err := c.Read(buf)
		if err != nil {
			if err != io.EOF {
				t.Error(err)
			}
			return
		}
		t.Logf("received: %q", buf[:n])
	}
}
