package memory

import (
	"log"
	"math/rand"
	"time"
)

type RAM interface {
	FreeMemory() int
}

type Samsung struct {
}

func (s *Samsung) FreeMemory() int {
	rand.Seed(time.Now().UnixNano())
	max := 1024 * 8
	min := 1024 * 3
	usage := rand.Intn(max-min) + min
	log.Printf("Rand usage of RAM: %v", usage)
	return max - usage
}
