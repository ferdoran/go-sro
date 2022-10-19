package navmeshv2

type RtNavmeshEdgeLink struct {
	LinkedObjID     int16
	LinkedObjEdgeID int16
	EdgeID          int16
}

type RtNavmeshInstObj struct {
	RtNavmeshInstBase
	Region    Region
	WorldID   int
	EdgeLinks []RtNavmeshEdgeLink
}
