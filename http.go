package main

import (
    "log"
    "strconv"
    "net/http"
    "golang.org/x/net/websocket"
    "io"
    "fmt"
)

const channelBufSize = 100

type WebSocketCommand struct {
    Command string `json:"command"`
}

type WebSocketMessage struct {
    Type    string `json:"type"`
    Payload interface{} `json:"payload"`
}

type WebSocketServer struct {
    api          *Api
    clients      map[int]*WebSocketClient
    addClient    chan *WebSocketClient
    removeClient chan *WebSocketClient
    done         chan bool
    error        chan error
}

type WebSocketClient struct {
    id         int
    connection *websocket.Conn
    server     *WebSocketServer
    messages   chan *WebSocketMessage
    done       chan bool
}

func serve(api *Api, port int, path string) {
    go websocketServe(api)
    httpServe(port, path)
}

func httpServe(port int, path string) {
    log.Println("Serving static content    :", path)
    http.Handle("/", http.FileServer(http.Dir(path)))
    log.Fatal(http.ListenAndServe(":" + strconv.Itoa(port), nil))
}

func websocketServe(api *Api) {
    var clientId int = 0

    server := &WebSocketServer{
        api,
        make(map[int]*WebSocketClient),
        make(chan *WebSocketClient),
        make(chan *WebSocketClient),
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

        clientId++
        client := &WebSocketClient{
            clientId,
            ws,
            server,
            make(chan *WebSocketMessage, channelBufSize),
            make(chan bool),
        }
        server.addClient <- client
        client.Listen(api)
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

        case message := <-server.api.stream:
            server.broadcast(toWebSocketMessage(message))

        case err := <-server.error:
            log.Println("Error:", err.Error())

        case <-server.done:
            return
        }
    }
}

func (server *WebSocketServer) broadcast(message *WebSocketMessage) {
    if message != nil {
        for _, client := range server.clients {
            client.Write(message)
        }
    }
}

func (client *WebSocketClient) Listen(api *Api) {
    go client.listenWrite()
    client.listenRead(api)
}

func (client *WebSocketClient) listenWrite() {
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

func (client *WebSocketClient) listenRead(api *Api) {
    for {
        select {

        case <-client.done:
            client.server.removeClient <- client
            client.done <- true
            return

        default:
            var message WebSocketCommand
            err := websocket.JSON.Receive(client.connection, &message)
            if err == io.EOF {
                client.done <- true
            } else if err != nil {
                client.server.error <- err
            } else {
                log.Println("Command [", client.id, "]", message.Command)
                api.Command(message.Command, client)
            }
        }
    }
}

func (client *WebSocketClient) Write(message *WebSocketMessage) {
    if message != nil {
        select {
        case client.messages <- message:
        default:
            client.server.removeClient <- client
            err := fmt.Errorf("Client %d is already disconnected.", client.id)
            client.server.error <- err
        }
    }
}

func (client *WebSocketClient) Reply(message interface{}) {
    client.Write(toWebSocketMessage(message))
}

func toWebSocketMessage(message interface{}) *WebSocketMessage {
    switch t := message.(type) {
    default:
        log.Printf("Unexpected message type %T", t)
        return nil
    case Execution:
        return &WebSocketMessage{"execution", message }
    case ExecutionStart:
        return &WebSocketMessage{"execution-start", message }
    case ExecutionFinish:
        return &WebSocketMessage{"execution-finish", message }
    case ExecutionLog:
        return &WebSocketMessage{"execution-log", message }
    }
}
