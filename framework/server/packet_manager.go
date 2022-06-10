package server

import (
	"sync"
)

type PacketManager struct {
	queues map[uint16]chan PacketChannelData
	mutex  sync.Mutex
}

func (pm *PacketManager) GetQueue(opcode uint16) chan PacketChannelData {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	if queue, exists := pm.queues[opcode]; exists {
		return queue
	} else {
		queue := make(chan PacketChannelData, 128)
		pm.queues[opcode] = queue
		return queue
	}
}

var PacketManagerInstance = &PacketManager{queues: make(map[uint16]chan PacketChannelData)}
