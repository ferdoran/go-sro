package opcode

const (
	AgentLoginTokenRequest          uint16 = 0xE100
	AgentGameserverForPlayerRequest uint16 = 0xE101
	AgentPlayerJoinRequest          uint16 = 0xE102

	AgentGameserverForPlayerResponse uint16 = 0xF101
	AgentPlayerJoinResponse          uint16 = 0xF102

	AgentLogoutRequest        uint16 = 0x7005
	AgentLogoutResponse       uint16 = 0xB005
	AgentLogoutSuccess        uint16 = 0x300A
	AgentLogoutCancelRequest  uint16 = 0x7006
	AgentLogoutCancelResponse uint16 = 0xB006

	CharacterDataBegin uint16 = 0x34A5
	CharacterDataBody  uint16 = 0x3013
	CharacterDataEnd   uint16 = 0x34A6

	ItemOperationRequest  uint16 = 0x7034
	ItemOperationResponse uint16 = 0xB034
)
