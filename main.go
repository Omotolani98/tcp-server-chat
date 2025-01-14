package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

var (
	clients    = make(map[net.Conn]string)
	clientsMux sync.Mutex
)

func main() {
	server, err := net.Listen("tcp", ":5080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer server.Close()

	fmt.Println("Listening on :5080")

	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	//get client name
	conn.Write([]byte("Welcome to the chat server!\nEnter your name: "))
	reader := bufio.NewReader(conn)
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	clientsMux.Lock()
	clients[conn] = name
	clientsMux.Unlock()

	broadcastSystemNotification(fmt.Sprintf("%s has joined the chat!", name), conn, "you have joined the chat :)")
	fmt.Printf("Client %s connected\n", name)

	for {
		//_, _ = conn.Write([]byte("me: "))
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Client %s disconnected\n", name)
			clientsMux.Lock()
			delete(clients, conn)
			clientsMux.Unlock()
			broadcastSystemNotification(fmt.Sprintf("%s has left the chat!", name), conn, "you have left the chat :(")
			return
		}

		//broadcast the message
		message = strings.TrimSpace(message)
		if message == "" {
			continue
		}

		if strings.HasPrefix(message, "/msg") {
			handlePrivateMessage(message, name, conn)
		} else if strings.HasPrefix(message, "/users") {
			listActiveUsers(conn)
		} else {
			//timestamp := time.Now().Format("15:04:05")
			broadcast(fmt.Sprintf("%s", message), conn)
		}
	}
}

func broadcastSystemNotification(message string, sender net.Conn, senderMessage string) {
	clientsMux.Lock()
	defer clientsMux.Unlock()

	timestamp := time.Now().Format("15:04:05")
	for client := range clients {
		if client == sender {
			_, _ = client.Write([]byte(fmt.Sprintf("[%s] %s\n", timestamp, senderMessage)))
		} else {
			_, _ = client.Write([]byte(fmt.Sprintf("[%s] %s\n", timestamp, message)))
		}
	}
}

func broadcast(message string, sender net.Conn) {
	clientsMux.Lock()
	defer clientsMux.Unlock()

	timestamp := time.Now().Format("15:04:05")
	senderName := clients[sender]
	for client := range clients {
		if client == sender {
			_, _ = client.Write([]byte(fmt.Sprintf("[%s] me: %s\n", timestamp, message)))
		} else {
			_, _ = client.Write([]byte(fmt.Sprintf("[%s] %s: %s\n", timestamp, senderName, message)))
		}
	}
}

func handlePrivateMessage(message, senderName string, senderConn net.Conn) {
	parts := strings.SplitN(message, " ", 3)
	if len(parts) < 3 {
		_, _ = senderConn.Write([]byte("Invalid format. Use /msg <username> <message>\n"))
		return
	}

	username := parts[1]
	privateMessage := parts[2]

	clientsMux.Lock()
	defer clientsMux.Unlock()

	for client, name := range clients {
		if name == username {
			timestamp := time.Now().Format("15:04:05")
			_, _ = client.Write([]byte(fmt.Sprintf("[%s] [Private] %s: %s \r\n", timestamp, senderName, privateMessage)))
			_, _ = senderConn.Write([]byte(fmt.Sprintf("[%s] [To] %s: %s \r\n", timestamp, username, privateMessage)))
			return
		}
	}

	_, _ = senderConn.Write([]byte("User not found\n"))
}

func listActiveUsers(conn net.Conn) {
	clientsMux.Lock()
	defer clientsMux.Unlock()

	_, _ = conn.Write([]byte("Active users:\n"))
	for _, name := range clients {
		_, _ = conn.Write([]byte("- " + name + "\n"))
	}
}
