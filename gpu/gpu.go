package gpu

import (
	"log"
	"math/rand"
	"time"
)

type GPU interface {
	Usage() int // percentage of the GPU
}

type Nvidia struct {
}

func (g *Nvidia) Usage() int {
	rand.Seed(time.Now().UnixNano())
	usage := rand.Intn(100)
	log.Printf("GPU Rand usage: %v", usage)
	return usage
}
