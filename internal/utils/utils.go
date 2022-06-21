package utils

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
)

func CalcFileSHA256(fname string) (res string, err error) {
	f, err := os.Open(fname)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Println(err)
		return "", err
	}

	str := fmt.Sprintf("%x", h.Sum(nil))
	return str, nil
}
