package main

import (
	"fmt"
	"log"
)

const initErr = "initialization error in step %s: %w"

func handleErr(step string, err error) {
	if err != nil {
		log.Fatal(fmt.Errorf(initErr, step, err))
	}
}

func main() {

}
