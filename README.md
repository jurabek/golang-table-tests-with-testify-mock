# golang-table-tests-with-testify-mock
This repo shows how to test and mock multiple dependencies using table-driven tests with testify-mock in Golang


## Problem
Sometimes you need to mock multiple dependencies for specific test cases, imagine we have a type `MacBook` and it depends on `CPU`, `GPU` and `RAM` and we have a `Diagnose` method that returns an error if some of the dependencies hit the threshold.

```
type MacBook struct {
	cpu    cpu.CPU
	gpu    gpu.GPU
	memory memory.RAM
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
```
as you can see we have 3 if condtions, to be able to test happy case that returns no error, we need to write at least 4 unit tests and mock all dependencies.

## Solution


## Links to read
[Table-Driven tests in Go](https://github.com/golang/go/wiki/TableDrivenTests)

[Testify Mocking](https://github.com/stretchr/testify#mock-package)
