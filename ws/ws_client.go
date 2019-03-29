package ws

import (
	"github.com/gorilla/websocket"
	"fmt"

	"github.com/spf13/viper"
	"net/http"
	"os"
)

func StartListeningMessagesAsClient() {

	dialer := websocket.Dialer{
		Subprotocols:    []string{},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	url := fmt.Sprintf("ws://%s:%d/ws",viper.GetString("app.host"),viper.GetInt("app.port"))
	header := http.Header{"Accept-Encoding": []string{"gzip"}}

	conn,_, err := dialer.Dial(url, header)
	if err != nil {
		fmt.Println("failed to send: ", err.Error())
		os.Exit(1)
	}

	for {
		_,message,err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("Error::: %s\n", err.Error())
			break
		}

		fmt.Println("message : ",string(message))
	}

	defer conn.Close()
}

func StartSendMessagesAsClient(res http.ResponseWriter,req *http.Request) {

	message := req.FormValue("message")
	if message == "" {
		message = "ping"
	}

	dialer := websocket.Dialer{
		Subprotocols:    []string{},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	url := fmt.Sprintf("ws://%s:%d/ws",viper.GetString("app.host"),viper.GetInt("app.port"))
	header := http.Header{"Accept-Encoding": []string{"gzip"}}

	conn,_, err := dialer.Dial(url, header)
	if err != nil {
		fmt.Println("failed to connect : ", err.Error())
		os.Exit(1)
	}

	err = conn.WriteMessage(websocket.TextMessage,[]byte (message))
	if err != nil {
		fmt.Println("failed to send: ", err.Error())
	}

	defer conn.Close()
}
