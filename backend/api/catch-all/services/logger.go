package services

import (
	"fmt"
	"log"
)

func ErrorLog(msg string) {
	log.Println(fmt.Sprintf("ERROR: %s", msg))
}

func InfoLog(msg string) {
	log.Println(fmt.Sprintf("INFO: %s", msg))
}
