package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/coder/websocket"
)

func main() {
	//mux := http.NewServeMux()

	//mux.HandleFunc("/hello", hello)

	l, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		panic(err)
	}
	s := http.Server{
		Handler: wsServer{},
	}
	s.Serve(l)

	fmt.Println("serving on http://localhost:8000")
	//http.ListenAndServe("localhost:8000", mux)
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world!"))
}

type wsServer struct{}

func (s wsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, _ := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols:   []string{"echo"},
		OriginPatterns: []string{"*"},
	})

	for {
		// Doing the logic in the separate function, then invoking that function
		// in an infinite for{} seems to allow it to run not just once.
		// Probably ServeHTTP just gets called one time, upon setting up the
		// websocket.
		err := echo(c)
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
			websocket.CloseStatus(err) == websocket.StatusGoingAway {
			return
		} else if err != nil {
			return
		}
	}
}

func echo(c *websocket.Conn) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	typ, r, err := c.Reader(ctx)
	if err != nil {
		return err
	}

	w, err := c.Writer(ctx, typ)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, r)
	if err != nil {
		return fmt.Errorf("failed to io.Copy: %w", err)
	}

	err = w.Close()
	return err
}
