package server

type UserContext struct {
	UserID       uint32
	GameServerID byte
	ShardID      uint16
	CharName     string
	UniqueID     uint32
	RefObjId     uint32
	Username     string
}
