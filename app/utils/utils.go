package utils

import (
	"errors"

	"github.com/matoous/go-nanoid/v2"
)

var runes = []rune("0123456789abcdefghijklmnopABCDERG")

func NanoId(size int) (string,error) {
	id, err := gonanoid.Generate(string(runes),size)

	if err != nil {
		return "", errors.New("size must be positive integer")
	}

	return string(id), nil
}
