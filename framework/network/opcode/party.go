package opcode

const (
	PartyCreateRequest               uint16 = 0x7060
	PartyKickRequest                 uint16 = 0x7063
	PartyUpdateResponse              uint16 = 0x3864
	PartyCreateResponse              uint16 = 0xB060
	PartyMatchingFormRequest         uint16 = 0x7069
	PartyMatchingFormResponse        uint16 = 0xB069
	PartyMatchingUpdateRequest       uint16 = 0x706A
	PartyMatchingUpdateResponse      uint16 = 0xB06A
	PartyMatchingDeleteRequest       uint16 = 0x706B
	PartyMatchingDeleteResponse      uint16 = 0xB06B
	PartyMatchingListRequest         uint16 = 0x706C
	PartyMatchingListResponse        uint16 = 0xB06C
	PartyMatchingJoinRequest         uint16 = 0x706D
	PartyMatchingJoinResponse        uint16 = 0xB06D
	PartyMatchingPlayerJoinRequest   uint16 = 0x306E
	PartyMatchingPlayerJoinResponse  uint16 = 0x706D
	PartyMemberCountResponse         uint16 = 0xB067
	PartyCreatedFromMatchingResponse uint16 = 0x3065
)
