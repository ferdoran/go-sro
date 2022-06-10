package utils

func Int16ToXAndZ(region int16) (x, z int) {
	x = int(region & 0xFF)
	z = int(region >> 8 & 0xFF)
	return
}

func XAndZToInt16(x, z byte) (regionId int16) {
	regionId += int16(x)
	regionId += int16(z) << 8
	return
}
