package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/net/websocket"
)

const channelBufSize = 100

type wsCommand struct {
	Command string `json:"command"`
}

type wsMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type wsServer struct {
	service      *executionService
	clients      map[int]*wsClient
	addClient    chan *wsClient
	removeClient chan *wsClient
	done         chan bool
	error        chan error
}

type wsClient struct {
	id         int
	connection *websocket.Conn
	server     *wsServer
	messages   chan *wsMessage
	done       chan bool
}

func serve(service *executionService, port int, path string) {
	go websocketServe(service)
	httpServe(port, path)
}

func httpServe(port int, path string) {
	http.Handle("/", http.FileServer(http.Dir(path)))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

func websocketServe(service *executionService) {
	var clientID int

	server := &wsServer{
		service,
		make(map[int]*wsClient),
		make(chan *wsClient),
		make(chan *wsClient),
		make(chan bool),
		make(chan error),
	}

	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				server.error <- err
			}
		}()

		clientID++
		client := &wsClient{
			clientID,
			ws,
			server,
			make(chan *wsMessage, channelBufSize),
			make(chan bool),
		}
		server.addClient <- client
		client.listen(service)
	}

	pattern := "/websocket"
	log.Println("Creating websocket handler:", pattern)
	http.Handle(pattern, websocket.Handler(onConnected))

	for {
		select {

		case client := <-server.addClient:
			log.Println("Add new websocket client:", client.id)
			server.clients[client.id] = client
			log.Println("Number of connected clients:", len(server.clients))

		case client := <-server.removeClient:
			log.Println("Remove websocket client:", client.id)
			delete(server.clients, client.id)

		case message := <-server.service.stream:
			server.broadcast(toWebSocketMessage(message))

		case err := <-server.error:
			log.Println("Error:", err.Error())

		case <-server.done:
			return
		}
	}
}

func (server *wsServer) broadcast(message *wsMessage) {
	if message != nil {
		for _, client := range server.clients {
			client.write(message)
		}
	}
}

func (client *wsClient) listen(service *executionService) {
	go client.listenWrite()
	client.listenRead(service)
}

func (client *wsClient) listenWrite() {
	for {
		select {
		case message := <-client.messages:
			websocket.JSON.Send(client.connection, message)

		case <-client.done:
			client.server.removeClient <- client
			client.done <- true
			return
		}
	}
}

func (client *wsClient) listenRead(service *executionService) {
	for {
		select {

		case <-client.done:
			client.server.removeClient <- client
			client.done <- true
			return

		default:
			var message wsCommand
			err := websocket.JSON.Receive(client.connection, &message)
			if err == io.EOF {
				client.done <- true
			} else if err != nil {
				client.server.error <- err
			} else {
				log.Println("Command [", client.id, "]", message.Command)
				service.command(message.Command, client)
			}
		}
	}
}

func (client *wsClient) write(message *wsMessage) {
	if message != nil {
		select {
		case client.messages <- message:
		default:
			client.server.removeClient <- client
			err := fmt.Errorf("client %d is already disconnected: ", client.id)
			client.server.error <- err
		}
	}
}

func (client *wsClient) Reply(message interface{}) {
	client.write(toWebSocketMessage(message))
}

func toWebSocketMessage(message interface{}) *wsMessage {
	switch t := message.(type) {
	default:
		log.Printf("Unexpected message type %T", t)
		return nil
	case execution:
		return &wsMessage{"execution", message}
	case executionStart:
		return &wsMessage{"execution-start", message}
	case executionFinish:
		return &wsMessage{"execution-finish", message}
	case executionLog:
		return &wsMessage{"execution-log", message}
	}
}
