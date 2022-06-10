package opcode

const (
	EntityGroupSpawnBegin uint16 = 0x3017
	EntityGroupSpawnData  uint16 = 0x3019
	EntityGroupSpawnEnd   uint16 = 0x3018
	EntitySingleSpawn     uint16 = 0x3015

	EntityUpdatePosition      uint16 = 0x3028
	EntityAnimationPickup     uint16 = 0x3036
	EntityEquipItem           uint16 = 0x3038
	EntityUnequipItem         uint16 = 0x3039
	EntityUpdateStats         uint16 = 0x303D
	EntityRemoveOwnership     uint16 = 0x304D
	EntityUpdatePoints        uint16 = 0x304E
	EntityAnimationPromote    uint16 = 0x3054
	EntityUpdateExperience    uint16 = 0x3056
	EntityUpdateStatus        uint16 = 0x3057
	EntityDamageEffect        uint16 = 0x3058
	EntityItemUsed            uint16 = 0x305C
	EntityEmotion             uint16 = 0x3091
	EntityUpdateMovementState uint16 = 0x30BF
	EntityUpdateMovementSpeed uint16 = 0x30D0
	EntityUpdateHwanLevel     uint16 = 0x30DF
	EntityUpdateAttackSpeed   uint16 = 0x3200
	EntityMask                uint16 = 0x3207
	EntityMovementRequest     uint16 = 0x7021
	EntityMovementResponse    uint16 = 0xB021
	MovementPositionUpdate    uint16 = 0xB023

	EntitySelectRequest  uint16 = 0x7045
	EntitySelectResponse uint16 = 0xB045
)
