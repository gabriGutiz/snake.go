package utils

import (
	"log"
)

func Assert(cond bool, msg string) {
	if !cond {
		log.Fatalf("Runtime assertion error: %s", msg)
	}
}
