package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/howeyc/gopass"
)

func readSecretInput(msg string, masked bool) (string, error) {
	var (
		input []byte
		err   error
	)

	fmt.Print(strings.TrimSpace(msg) + " ")

	if !masked {
		input, err = gopass.GetPasswdMasked()
	} else {
		input, err = gopass.GetPasswd()
	}
	return string(input[:]), err
}

func readInput(msg string, r ...io.ReadCloser) (string, error) {
	var (
		rd io.ReadCloser
	)

	if len(r) == 0 {
		rd = os.Stdin
	} else if len(r) >= 1 {
		rd = r[0]
	}

	fmt.Print(strings.TrimSpace(msg) + " ")

	reader := bufio.NewReader(rd)
	str, err := reader.ReadString('\n')
	return strings.TrimSpace(str), err
}

func strSliceContains(slice []string, key string) bool {
	for _, s := range slice {
		if s == key {
			return true
		}
	}
	return false
}
