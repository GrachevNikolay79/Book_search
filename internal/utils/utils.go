package utils

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func CalcFileSHA256(fname string) (res string, size int64, err error) {
	f, err := os.Open(fname)
	if err != nil {
		log.Println(err)
		return "", 0, err
	}
	defer f.Close()

	h := sha256.New()
	size, err = io.Copy(h, f)
	if err != nil {
		log.Println(err)
		return "", size, err
	}

	str := fmt.Sprintf("%x", h.Sum(nil))
	return str, size, nil
}

func DoWithTries(fn func() error, attemtps int, delay time.Duration) (err error) {
	for attemtps > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attemtps--

			continue
		}

		return nil
	}

	return
}

func MaxI(a, b int) int {
	if a > b {
		return a
	}
	return b
}
