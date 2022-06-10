package opcode

const (
	StallCreateRequest         uint16 = 0x70B1
	StallCreateResponse        uint16 = 0xB0B1
	StallDestroyRequest        uint16 = 0x70B2
	StallDestroyResponse       uint16 = 0xB0B2
	StallTalkRequest           uint16 = 0x70B3
	StallTalkResponse          uint16 = 0xB0B3
	StallBuyRequest            uint16 = 0x70B4
	StallBuyResponse           uint16 = 0xB0B4
	StallLeaveRequest          uint16 = 0x70B5
	StallLeaveResponse         uint16 = 0xB0B5
	StallUpdateRequest         uint16 = 0x70BA
	StallUpdateResponse        uint16 = 0xB0BA
	StallEntityActionResponse  uint16 = 0x30B7
	StallEntityCreateResponse  uint16 = 0x30B8
	StallEntityDestroyResponse uint16 = 0x30B9
	StallEntityNameResponse    uint16 = 0x30BB
)
