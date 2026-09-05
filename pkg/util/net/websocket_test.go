package net

import (
	"errors"
	"net"
	"testing"
	"time"

	"golang.org/x/net/websocket"
)

func TestWebsocketListenerCloseUnblocksAccept(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	wl := NewWebsocketListener(ln)
	defer wl.Close()
	done := make(chan error, 1)
	go func() { _, err := wl.Accept(); done <- err }()
	if err := wl.Close(); err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-done:
		if !errors.Is(err, ErrWebsocketListenerClosed) {
			t.Fatalf("unexpected accept error: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Accept remained blocked after listener Close")
	}
}

func TestWebsocketListenerCloseRejectsUnacceptedUpgrade(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	wl := NewWebsocketListener(ln)
	defer wl.Close()
	c, err := websocket.Dial("ws://"+ln.Addr().String()+FrpWebsocketPath, "", "http://localhost/")
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	if err := wl.Close(); err != nil {
		t.Fatal(err)
	}
	_ = c.SetReadDeadline(time.Now().Add(time.Second))
	var buf [1]byte
	_, err = c.Read(buf[:])
	if err == nil {
		t.Fatal("unaccepted websocket remained open")
	}
	var timeout net.Error
	if errors.As(err, &timeout) && timeout.Timeout() {
		t.Fatalf("unaccepted upgraded connection leaked: %v", err)
	}
}

func TestWebsocketListenerClosePreservesAcceptedConnection(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	wl := NewWebsocketListener(ln)
	defer wl.Close()
	c, err := websocket.Dial("ws://"+ln.Addr().String()+FrpWebsocketPath, "", "http://localhost/")
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	accepted, err := wl.Accept()
	if err != nil {
		t.Fatal(err)
	}
	defer accepted.Close()
	if err := wl.Close(); err != nil {
		t.Fatal(err)
	}
	_ = c.SetDeadline(time.Now().Add(time.Second))
	_ = accepted.SetDeadline(time.Now().Add(time.Second))
	written := make(chan error, 1)
	go func() { _, err := accepted.Write([]byte("x")); written <- err }()
	var buf [1]byte
	if _, err := c.Read(buf[:]); err != nil {
		t.Fatal(err)
	}
	if buf[0] != 'x' {
		t.Fatalf("unexpected data: %q", buf)
	}
	if err := <-written; err != nil {
		t.Fatal(err)
	}
}
