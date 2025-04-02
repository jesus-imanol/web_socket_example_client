package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/twinj/uuid"
	"github.com/gorilla/websocket"
)
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // Allow all origins * 
    },
}
type Notification struct {
    Title string `json:"title"`
    Body  string `json:"body"`
    Token string `json:"token"`
}
type Message struct {
    Event string `json:"event"`
    Data  interface {}
}
var connections = make(map[string]*websocket.Conn)
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Print("Error upgrading to WebSocket:", err)
        return
    }
    defer conn.Close()
    fmt.Println("Cliente conectado")

    uuid := uuid.NewV4().String()
    err  = conn.WriteJSON(Message{
        Event: "token",
        Data: uuid,
    })
    connections[uuid] = conn

    for {
        var msg Message
        err := conn.ReadJSON(&msg)
        if err != nil {
            log.Println("Error reading JSON:", err)
            break
        }

        switch msg.Event {
        case "notification":
            jsonData, err := json.Marshal(msg.Data)
            if err != nil {
                log.Println("Error marshaling JSON:", err)
                break
            }
            var notification Notification
            err = json.Unmarshal(jsonData, &notification)
            if err != nil {
                log.Println("Error unmarshaling JSON:", err)
                break
            }
            conn := connections[notification.Token]
            err = conn.WriteJSON(Message{
                Event: "notification",
                Data: notification,
            })
            if err!= nil {
                log.Println("Error writing JSON:", err)
                break
            }
            log.Println("Notificacion enviada")


        }


    }
}
func main() {
    http.HandleFunc("/ws", handleWebSocket)
    log.Println("Server running")

    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Println("Error iniciando server:", err)
    }
}