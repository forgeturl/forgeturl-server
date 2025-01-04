package maths

import (
	"math/rand"
)

const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// Base58Encode encodes an int64 to a Base58 string
func Base58Encode(num int64) string {
	if num == 0 {
		return string(alphabet[0])
	}

	// Convert the int64 to bytes
	var result []byte
	for num > 0 {
		remainder := num % 58
		num /= 58
		result = append(result, alphabet[remainder]) // Append the corresponding character
	}

	// Reverse the result to get the correct order
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}

// GenPageID 如果是int32最长是7个字符，如果是int64最长是12个字符
func GenPageID(prefix string) string {
	// ri := rand.Int63()
	ri := rand.Int31()
	return prefix + Base58Encode(int64(ri))
}
