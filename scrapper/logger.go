package scrapper

import (
	"log"
	"os"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "scrapper: ", log.Ldate|log.Ltime|log.Lshortfile)
}
