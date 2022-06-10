package network

import "sync"

const DefaultSeed uint32 = 0x9ABFB3B6

func generateValue(value uint32) (uint32, uint32) {
	newVal := value
	complement1 := ^uint32(1)
	for i := 0; i < 32; i++ {
		v := newVal
		v = (v >> 2) ^ newVal
		v = (v >> 2) ^ newVal
		v = (v >> 1) ^ newVal
		v = (v >> 1) ^ newVal
		v = (v >> 1) ^ newVal
		newVal = ((newVal>>1)|(newVal<<31))&complement1 | (v & 1)
	}
	return newVal, newVal
}

type MessageSequence struct {
	byte0 byte
	byte1 byte
	byte2 byte
	mutex *sync.Mutex
}

func (ms *MessageSequence) Next() byte {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	value := ms.byte2 * (^ms.byte0 + ms.byte1)
	ms.byte0 = value ^ (value >> 4)
	return ms.byte0
}

func NewMessageSequence(seed uint32) MessageSequence {
	var mut0 uint32
	if seed != 0 {
		mut0 = seed
	} else {
		mut0 = DefaultSeed
	}
	mut0, mut1 := generateValue(mut0)
	mut0, mut2 := generateValue(mut0)
	mut0, mut3 := generateValue(mut0)
	mut0, _ = generateValue(mut0)

	byte1 := byte((mut1 & 0xFF) ^ (mut2 & 0xFF))
	if byte1 == 0 {
		byte1 = 1
	}

	byte2 := byte((mut0 & 0xFF) ^ (mut3 & 0xFF))
	if byte2 == 0 {
		byte2 = 1
	}

	byte0 := byte2 ^ byte1

	return MessageSequence{byte0, byte1, byte2, &sync.Mutex{}}
}
