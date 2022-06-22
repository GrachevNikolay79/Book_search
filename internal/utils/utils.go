package utils

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"time"
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
