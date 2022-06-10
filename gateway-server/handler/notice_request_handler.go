package handler

import (
	"github.com/ferdoran/go-sro/framework/network"
	"github.com/ferdoran/go-sro/framework/network/opcode"
	"github.com/ferdoran/go-sro/framework/server"
	"time"

	"github.com/ferdoran/go-sro/gateway-server/db"
)

type NoticeRequestHandler struct {
	channel chan server.PacketChannelData
}

func InitNoticeRequestHandler() {
	handler := NoticeRequestHandler{channel: server.PacketManagerInstance.GetQueue(opcode.NoticeRequest)}
	go handler.Handle()
}

func (h *NoticeRequestHandler) Handle() {
	for {
		packet := <-h.channel
		p := network.EmptyPacket()
		p.MessageID = 0x600D
		p.WriteByte(1)
		p.WriteUInt16(1)
		p.WriteUInt16(opcode.NoticeResponse)
		packet.Session.Conn.Write(p.ToBytes())

		notices := db.GetNotices()
		p = network.EmptyPacket()
		p.MessageID = 0x600D
		p.WriteByte(0)
		p.WriteByte(byte(len(notices)))
		for _, notice := range notices {
			p.WriteString(notice.Subject)
			p.WriteString(notice.Article)
			p.WriteUInt16(uint16(notice.Ctime.Year()))
			p.WriteUInt16(uint16(notice.Ctime.Month()))
			p.WriteUInt16(uint16(notice.Ctime.Day()))
			p.WriteUInt16(uint16(notice.Ctime.Hour()))
			p.WriteUInt16(uint16(notice.Ctime.Minute()))
			p.WriteUInt16(uint16(notice.Ctime.Second()))
			p.WriteUInt32(uint32((notice.Ctime.Nanosecond()) / int(time.Millisecond)))
		}

		packet.Session.Conn.Write(p.ToBytes())
	}
}
