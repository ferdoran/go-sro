package utils

func ItemClassToDegree(itemClass int) int {
	return ((itemClass - 1) / 3) + 1
}

func ItemClassToDegreeTier(itemClass int) int {
	return itemClass - (3 * ((itemClass - 1) / 3)) - 1
}
