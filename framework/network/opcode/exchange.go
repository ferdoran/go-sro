package opcode

const (
	ExchangeStartRequest  uint16 = 0x7081
	ExchangeStartResponse uint16 = 0xB081

	ExchangeConfirmRequest  uint16 = 0x7082
	ExchangeConfirmResponse uint16 = 0xB082

	ExchangeApproveRequest  uint16 = 0x7083
	ExchangeApproveResponse uint16 = 0xB083

	ExchangeCancelRequest  uint16 = 0x7084
	ExchangeCancelResponse uint16 = 0xB084

	ExchangeStartedResponse     uint16 = 0x3085
	ExchangeConfirmedResponse   uint16 = 0x3086
	ExchangeApprovedResponse    uint16 = 0x3087
	ExchangeCancelledResponse   uint16 = 0x3088
	ExchangeUpdateResponse      uint16 = 0x3089
	ExchangeUpdateItemsResponse uint16 = 0x308C
)
