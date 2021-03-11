package cpu

import (
	"log"
	"math/rand"
	"time"
)

type CPU interface {
	Usage() int // percentage of CPU usage
}

type Intel struct {
}

func (i *Intel) Usage() int {
	rand.Seed(time.Now().UnixNano())
	usage := rand.Intn(100)
	log.Printf("CPU Rand usage: %v", usage)
	return usage
}
