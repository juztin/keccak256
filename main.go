package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/sha3"
)

func main() {
	isHex := flag.Bool("x", false, "input is hexadecimal")
	flag.Parse()

	var b []byte
	var err error
	switch len(os.Args) {
	case 1:
		b, err = parseReader(os.Stdin, *isHex)
	default:
		b, err = parseString(os.Args[len(os.Args)-1], *isHex)
	}

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(hex.EncodeToString(b))
	}
}

func parseString(s string, isHex bool) ([]byte, error) {
	var b []byte
	var err error
	if isHex {
		b, err = hex.DecodeString(s)
	} else {
		b = []byte(s)
	}

	hash := sha3.NewLegacyKeccak256()
	if err == nil {
		_, err = hash.Write(b)
	}
	return hash.Sum(nil), nil
}

func parseReader(r io.Reader, isHex bool) ([]byte, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}

	if info.Mode()&os.ModeCharDevice != 0 || info.Size() <= 0 {
		return nil, fmt.Errorf("Invalid pipe data.\n\nUsage: echo \"Error(string)\" | keccak256")
	}

	if isHex {
		r = hex.NewDecoder(r)
	}

	reader := bufio.NewReader(r)
	hash := sha3.NewLegacyKeccak256()
	for {
		l, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		hash.Write(l)
	}
	return hash.Sum(nil), err
}
