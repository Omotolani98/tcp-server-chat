package pkg

import (
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

var clientManager = NewClientManager()

func Broadcast(message string, sender net.Conn, clientManager *ClientManager) {
	timestamp := time.Now().Format("15:04:05")
	allClients := clientManager.GetAllClients()
	senderName, _ := clientManager.GetClientName(sender)

	for client := range allClients {
		if client == sender {
			_, _ = client.Write([]byte(fmt.Sprintf("[%s] me: %s\n", timestamp, message)))
		} else {
			_, _ = client.Write([]byte(fmt.Sprintf("[%s] %s: %s\n", timestamp, senderName, message)))
		}
	}
}

func BroadcastSystemNotification(message string, sender net.Conn, senderMessage string, clientManager *ClientManager) {
	timestamp := time.Now().Format("15:04:05")
	allClients := clientManager.GetAllClients()

	for client := range allClients {
		if client == sender {
			_, _ = client.Write([]byte(fmt.Sprintf("[%s] %s\n", timestamp, senderMessage)))
		} else {
			_, _ = client.Write([]byte(fmt.Sprintf("[%s] %s\n", timestamp, message)))
		}
	}
}

func HandlePrivateMessage(message, senderName string, senderConn net.Conn) {
	parts := strings.SplitN(message, " ", 3)
	if len(parts) < 3 {
		_, _ = senderConn.Write([]byte("Invalid format. Use /msg <username> <message>\n"))
		return
	}

	username := parts[1]
	privateMessage := parts[2]

	allClients := clientManager.GetAllClients()

	for client, name := range allClients {
		if name == username {
			timestamp := time.Now().Format("15:04:05")
			_, _ = client.Write([]byte(fmt.Sprintf("[%s] [Private] %s: %s \r\n", timestamp, senderName, privateMessage)))
			_, _ = senderConn.Write([]byte(fmt.Sprintf("[%s] [To] %s: %s \r\n", timestamp, username, privateMessage)))
			return
		}
	}

	_, _ = senderConn.Write([]byte("User not found\n"))
}

func ListActiveUsers(conn net.Conn) {
	allClients := clientManager.GetAllClients()

	_, _ = conn.Write([]byte("Active users:\n"))
	for _, name := range allClients {
		_, _ = conn.Write([]byte("- " + name + "\n"))
	}
}
