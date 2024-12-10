package supervisor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"golang.org/x/time/rate"
)

func Start() {
	l, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		panic(err)
	}
	s := http.Server{
		Handler: wsServer{},
	}
	s.Serve(l)

	fmt.Println("serving on http://localhost:8000")
}

func hello(w http.ResponseWriter) {
	w.Write([]byte("hello world!"))
}

type wsServer struct{}

func (s wsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, _ := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols:   []string{"echo"},
		OriginPatterns: []string{"*"},
	})

	msgType, err := greetClient(c)
	if err != nil {
		fmt.Println("GREETING ERROR:", err)
		return
	}

	rateLimiter := rate.NewLimiter(rate.Every(time.Second*5), 10)

	for {
		// Doing the logic in the separate function, then invoking that function
		// in an infinite for{} seems to allow it to run not just once.
		// Probably ServeHTTP just gets called one time, upon setting up the
		// websocket.
		err := rateLimiter.Wait(context.Background())
		if err != nil {
			fmt.Println("rate limiter error:", err)
			return
		}

		err = writeNodeData(c, msgType)
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
			websocket.CloseStatus(err) == websocket.StatusGoingAway {
			fmt.Println("Closed:", websocket.CloseStatus(err))
			return
		} else if err != nil {
			fmt.Println("write node data error:", err)
			return
		}
	}
}

type clientMsg struct {
	nodes []int
}

func greetClient(c *websocket.Conn) (websocket.MessageType, error) {
	msgType, r, err := c.Reader(context.Background())
	if err != nil {
		fmt.Println(err)
		return msgType, err
	}

	clientBuffer := new(bytes.Buffer)
	clientBuffer.ReadFrom(r)
	clientData := clientMsg{}
	json.Unmarshal(clientBuffer.Bytes(), &clientData)
	fmt.Println("client data:", clientData)

	nodes, err := writeNodes(c, msgType)
	if err != nil {
		fmt.Println("NODES ERROR:", err)
		return msgType, err
	}
	fmt.Println(nodes)

	return msgType, nil
}

func writeNodes(c *websocket.Conn, typ websocket.MessageType) ([]int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	c.Write(ctx, typ, []byte("{ \"nodes\": [0, 1, 2] }"))
	return []int{0, 1, 2}, nil
}

func writeNodeData(c *websocket.Conn, typ websocket.MessageType) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	data, err := getDataFromNode()
	if err != nil {
		fmt.Println(err)
		return err
	}

	c.Write(ctx, typ, data)
	return nil
}

func getDataFromNode() ([]byte, error) {
	url := "http://localhost:8001/"
	client := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, getErr := client.Do(req)
	if getErr != nil {
		return nil, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return nil, err
	}

	return body, err
}
