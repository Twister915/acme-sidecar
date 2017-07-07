package main

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"os"
)

func init() {
	log.SetHandler(text.New(os.Stdout))
}
