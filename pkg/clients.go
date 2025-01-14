package pkg

import (
	"net"
	"sync"
)

type ClientManager struct {
	clients    map[net.Conn]string
	clientsMux sync.Mutex
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		clients: make(map[net.Conn]string),
	}
}

func (cm *ClientManager) AddClient(conn net.Conn, name string) {
	cm.clientsMux.Lock()
	defer cm.clientsMux.Unlock()
	cm.clients[conn] = name
}

func (cm *ClientManager) RemoveClient(conn net.Conn) {
	cm.clientsMux.Lock()
	defer cm.clientsMux.Unlock()
	delete(cm.clients, conn)
}

func (cm *ClientManager) GetClientName(conn net.Conn) (string, bool) {
	cm.clientsMux.Lock()
	defer cm.clientsMux.Unlock()
	name, exists := cm.clients[conn]
	return name, exists
}

func (cm *ClientManager) GetAllClients() map[net.Conn]string {
	cm.clientsMux.Lock()
	defer cm.clientsMux.Unlock()

	clientsCopy := make(map[net.Conn]string)
	for conn, name := range cm.clients {
		clientsCopy[conn] = name
	}
	return clientsCopy
}
