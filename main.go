package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/sha3"
	//"github.com/ethereum/go-ethereum/crypto/sha3"
)

func main() {
	var h string
	var err error
	switch len(os.Args) {
	case 1:
		h, err = parseReader(os.Stdin)
	case 2:
		h, err = parseString(os.Args[1])
	default:
		err = fmt.Errorf("Must provide either a single argument, or piped data")
	}

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(h)
	}
}

func parseString(s string) (string, error) {
	hash := sha3.NewLegacyKeccak256()
	hash.Write([]byte(s))
	buf := hash.Sum(nil)
	return hex.EncodeToString(buf), nil
}

func parseReader(r io.Reader) (string, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}

	if info.Mode()&os.ModeCharDevice != 0 || info.Size() <= 0 {
		return "", fmt.Errorf(`Invalid pipe data.\nUsage: echo "Error(string)" | keccak256`)
	}

	hash := sha3.NewLegacyKeccak256()
	reader := bufio.NewReader(r)
	reader.WriteTo(os.Stdout)
	reader.WriteTo(hash)
	buf := hash.Sum(nil)
	return hex.EncodeToString(buf), nil
}
