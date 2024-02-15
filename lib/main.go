package util

import (
	"log"
)

func Check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
