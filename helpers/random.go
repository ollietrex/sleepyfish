package helpers

import (
	"crypto/rand"
	"fmt"
)

//Generate a random string of chars to be used as test data
func RandomString(strSize int) string {
	dictionary := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	var bytes = make([]byte, strSize)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}

	return string(bytes)
}

func RandomEmail() string {
	return fmt.Sprintf("%s.%s@%s.com", RandomString(5), RandomString(6), RandomString(10))
}

func RandomName() string {
	return fmt.Sprintf("%s %s", RandomString(5), RandomString(6))
}
