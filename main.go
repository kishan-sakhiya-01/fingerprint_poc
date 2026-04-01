package main

import (
	"fmt"

	"github.com/kishan-sakhiya-01/fingerprint_poc/deviceid"
)

func main() {
	fp, hash, err := deviceid.Compute()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(hash)
	fmt.Println(fp)
}
