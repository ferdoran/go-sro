package utils

import "math"

func BaseMinAttack(strOrInt int) int {
	return int(0.45 * float32(strOrInt))
}

func BaseMaxAttack(strOrInt int) int {
	return int(0.65 * float32(strOrInt))
}

func BaseDef(strOrInt int) int {
	return int(0.4 * float32(strOrInt))
}

func PhyBalance(level, strOrInt int) int {
	return int(100 - (100 * 2 / 3 * ((28 + float32(level)*4) - float32(strOrInt)) / (28 + float32(level)*4)))
}

func MagBalance(level, strOrInt int) int {
	return int(100*float32(strOrInt)/28 + float32(level)*4)
}

func BaseHPOrMP(level, strOrInt int) int {
	return int(math.Pow(1.02, float64(level-1)) * float64(strOrInt) * 10)
}
