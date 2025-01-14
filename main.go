package main

import (
	"bufio"
	"fmt"
	"github.com/Omotolani98/tcp-server-chat/pkg"
	"net"
	"strings"
)

var clientManager = pkg.NewClientManager()

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
	_, _ = conn.Write([]byte("Welcome to the chat server!\nEnter your name: "))
	reader := bufio.NewReader(conn)
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	clientManager.AddClient(conn, name)

	pkg.BroadcastSystemNotification(fmt.Sprintf("%s has joined the chat!", name), conn, "you have joined the chat :)", clientManager)
	fmt.Printf("Client %s connected\n", name)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Client %s disconnected\n", name)
			clientManager.RemoveClient(conn)
			pkg.BroadcastSystemNotification(fmt.Sprintf("%s has joined the chat!", name), conn, "you have joined the chat :)", clientManager)
			return
		}

		//broadcast the message
		message = strings.TrimSpace(message)
		if message == "" {
			continue
		}

		if strings.HasPrefix(message, "/create") {
			pkg.CreateGroup(message, name, conn)
		} else if strings.HasPrefix(message, "/join") {
			pkg.JoinGroup(message, name, conn)
		} else if strings.HasPrefix(message, "/groups") {
			pkg.ListGroups(conn)
		} else if strings.HasPrefix(message, "/members") {
			pkg.ListGroupMembers(message, conn)
		} else if strings.HasPrefix(message, "/gmsg") {
			pkg.GroupMessage(message, name, conn)
		} else if strings.HasPrefix(message, "/msg") {
			pkg.HandlePrivateMessage(message, name, conn)
		} else if strings.HasPrefix(message, "/users") {
			pkg.ListActiveUsers(conn)
		} else {
			pkg.Broadcast(fmt.Sprintf("%s", message), conn, clientManager)
		}
	}
}
