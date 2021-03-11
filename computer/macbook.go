package computer

import (
	"errors"

	"github.com/jurabek/table-mock-sample/cpu"
	"github.com/jurabek/table-mock-sample/gpu"
	"github.com/jurabek/table-mock-sample/memory"
)

var CpuUtilizationError = errors.New("cpu utilized more than expected")
var MemoryUsageError = errors.New("not enough memory")
var GpuUsageError = errors.New("gpu usage error")

type Diagnoseable interface {
	Diagnose(cpuThreshold, gpuThreshold, memoryThreshold int) error
}

type MacBook struct {
	cpu    cpu.CPU
	gpu    gpu.GPU
	memory memory.RAM
}

func NewMacBook(cpu cpu.CPU, gpu gpu.GPU, memory memory.RAM) *MacBook {
	return &MacBook{
		cpu:    cpu,
		gpu:    gpu,
		memory: memory,
	}
}

func (m *MacBook) Diagnose(cpuThreshold, gpuThreshold, memoryThreshold int) error {
	if (m.cpu.Usage()) > cpuThreshold {
		return CpuUtilizationError
	}

	if m.gpu.Usage() > gpuThreshold {
		return GpuUsageError
	}

	if m.memory.FreeMemory() <= memoryThreshold {
		return MemoryUsageError
	}

	return nil
}
