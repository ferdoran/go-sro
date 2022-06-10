package security

import (
	utils2 "github.com/ferdoran/go-sro/framework/utils"
)

func G_pow_X_mod_P(generator, private, prime uint32) uint32 {
	var result int64 = 1
	var mult = int64(generator)

	if private == 0 {
		return 1
	}

	for private != 0 {
		if (private & 1) > 0 {
			result = (mult * result) % int64(prime)
		}

		private >>= 1
		mult = (mult * mult) % int64(prime)
	}
	return uint32(result)
}

func KeyTransformValue(value []byte, key uint32, keyByte byte) []byte {
	value[0] ^= byte(uint32(value[0]) + ((key) >> uint32(0) & uint32(0xFF)) + uint32(keyByte))
	value[1] ^= byte(uint32(value[1]) + ((key) >> uint32(8) & uint32(0xFF)) + uint32(keyByte))
	value[2] ^= byte(uint32(value[2]) + ((key) >> uint32(16) & uint32(0xFF)) + uint32(keyByte))
	value[3] ^= byte(uint32(value[3]) + ((key) >> uint32(24) & uint32(0xFF)) + uint32(keyByte))

	value[4] ^= byte(uint32(value[4]) + ((key) >> uint32(0) & uint32(0xFF)) + uint32(keyByte))
	value[5] ^= byte(uint32(value[5]) + ((key) >> uint32(8) & uint32(0xFF)) + uint32(keyByte))
	value[6] ^= byte(uint32(value[6]) + ((key) >> uint32(16) & uint32(0xFF)) + uint32(keyByte))
	value[7] ^= byte(uint32(value[7]) + ((key) >> uint32(24) & uint32(0xFF)) + uint32(keyByte))
	return value
}

func CalculateKey(commonSecret, secret1, secret2 uint32) []byte {
	newKey := utils2.ConcatTwoUint32ToByteArray(secret1, secret2)
	return KeyTransformValue(newKey, commonSecret, byte(commonSecret&3))
}

func CalculateChallenge(commonSecret, secret1, secret2 uint32) []byte {
	newKey := utils2.ConcatTwoUint32ToByteArray(secret1, secret2)
	return KeyTransformValue(newKey, commonSecret, byte(secret1&7))
}
