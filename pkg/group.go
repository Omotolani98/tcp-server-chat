package pkg

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

var (
	groups    = make(map[string]map[net.Conn]string)
	groupsMux sync.Mutex
)

func CreateGroup(message, username string, conn net.Conn) {
	parts := strings.Split(message, " ")
	if len(parts) != 2 {
		_, _ = conn.Write([]byte("Invalid format. Use /create <group_name> (e.g. new_grp)\n"))
		return
	}

	groupName := "grp/" + parts[1]
	groupsMux.Lock()
	defer groupsMux.Unlock()

	if _, exists := groups[groupName]; exists {
		_, _ = conn.Write([]byte("Group already exists.\n"))
		return
	}

	groups[groupName] = make(map[net.Conn]string)
	_, _ = conn.Write([]byte("Group created successfully.\n"))
}

func JoinGroup(message, username string, conn net.Conn) {
	parts := strings.Split(message, " ")
	if len(parts) != 2 {
		_, _ = conn.Write([]byte("Invalid format. Use /join <group_name>\n"))
		return
	}

	groupName := "grp/" + parts[1]
	groupsMux.Lock()
	defer groupsMux.Unlock()

	group, exists := groups[groupName]
	if !exists {
		_, _ = conn.Write([]byte("Group does not exist.\n"))
		return
	}

	group[conn] = username
	_, _ = conn.Write([]byte(fmt.Sprintf("Joined group %s.\n", groupName)))
}

func ListGroups(conn net.Conn) {
	groupsMux.Lock()
	defer groupsMux.Unlock()

	if len(groups) == 0 {
		_, _ = conn.Write([]byte("No groups available.\n"))
		return
	}

	_, _ = conn.Write([]byte("Available groups:\n"))
	for groupName := range groups {
		_, _ = conn.Write([]byte("- " + groupName + "\n"))
	}
}

func ListGroupMembers(message string, conn net.Conn) {
	parts := strings.Split(message, " ")
	if len(parts) != 2 {
		_, _ = conn.Write([]byte("Invalid format. Use /members <group_name>\n"))
		return
	}

	groupName := "grp/" + parts[1]
	groupsMux.Lock()
	defer groupsMux.Unlock()

	group, exists := groups[groupName]
	if !exists {
		_, _ = conn.Write([]byte("Group does not exist.\n"))
		return
	}

	_, _ = conn.Write([]byte(fmt.Sprintf("Members of %s:\n", groupName)))
	for _, username := range group {
		_, _ = conn.Write([]byte("- " + username + "\n"))
	}
}

func GroupMessage(message, senderName string, senderConn net.Conn) {
	parts := strings.SplitN(message, " ", 3)
	if len(parts) != 3 {
		_, _ = senderConn.Write([]byte("Invalid format. Use /gmsg <group_name> <message>\n"))
		return
	}

	groupName := "grp/" + parts[1]
	groupMessage := parts[2]

	groupsMux.Lock()
	defer groupsMux.Unlock()

	group, exists := groups[groupName]
	if !exists {
		_, _ = senderConn.Write([]byte("Group does not exist.\n"))
		return
	}

	timestamp := time.Now().Format("15:04:05")
	for conn, _ := range group {
		if conn == senderConn {
			_, _ = conn.Write([]byte(fmt.Sprintf("[%s] me (to %s): %s\n", timestamp, groupName, groupMessage)))
		} else {
			_, _ = conn.Write([]byte(fmt.Sprintf("[%s] %s (in %s): %s\n", timestamp, senderName, groupName, groupMessage)))
		}
	}
}
