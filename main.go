package main

import (
	"github.com/jurabek/table-mock-sample/computer"
	"github.com/jurabek/table-mock-sample/cpu"
	"github.com/jurabek/table-mock-sample/gpu"
	"github.com/jurabek/table-mock-sample/memory"
)

func main() {
	intel := &cpu.Intel{}
	nvidia := &gpu.Nvidia{}
	samsung := &memory.Samsung{}
	_ = computer.NewMacBook(intel, nvidia, samsung)
}
