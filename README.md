# golang-table-tests-with-testify-mock
This repo shows how to test and mock multiple dependencies using table-driven tests with testify-mock in Golang


## Problem
In unit testing and mocking world we need to mock multiple dependencies for specific test cases, (e.g) imagine we have a  `struct MacBook` and it depends on `CPU`, `GPU` and `RAM` and we have a `Diagnose` method that returns an error if some of the dependencies `used resources` hits the threshold.
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
example above we have 3 if conditions, in order to test and get happy result that returns no error, we should have to write at least 4 unit tests or more and mock all dependencies.

## Solution
```
func TestMacbookDiagnose(t *testing.T) {
	type depFields struct {
		cpu    *mockedCPU
		gpu    *mockedGPU
		memory *mockedRAM
	}

	type args struct {
		cpuThreshold, gpuThreshold, memoryThreshold int
	}

	tests := []struct {
		name string
		in   *args
		out  error

		on     func(*depFields)
		assert func(*depFields)
	}{
	
	 	// Test cases...
		
	 	{
			name: "when all thresholds not hit return nil",
			in:   &args{50, 90, 1000},
			out:  nil,
			on: func(df *depFields) {
				df.cpu.On("Usage").Return(40)           // 40% CPU usage less than cpuThreshold
				df.gpu.On("Usage").Return(50)           // 50% gpu usage less than gpuThreshold
				df.memory.On("FreeMemory").Return(2000) // 2000 MB free memory left so it is larger than 1000 mb threshold
			},
			assert: func(t *testing.T, df *depFields) {
				df.cpu.AssertNumberOfCalls(t, "Usage", 1)
				df.gpu.AssertNumberOfCalls(t, "Usage", 1)
				df.memory.AssertNumberOfCalls(t, "FreeMemory", 1)
			},
		},
	}
```
basically in the table-driven testing we will define slice of structs that represents test cases, in the example above we have struct that contains callback functions `on: func(*depFields)` and `assert: func(*testing.T, *depFields)` when we call those methods `type depFields struct{}` will be passed into it in order to `mock` or `assert`
```
for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			f := &depFields{
				&mockedCPU{},
				&mockedGPU{},
				&mockedRAM{},
			}
			mb := NewMacBook(f.cpu, f.gpu, f.memory)
			if tt.on != nil {
				tt.on(f)
			}
			// act
			err := mb.Diagnose(tt.in.cpuThreshold, tt.in.gpuThreshold, tt.in.memoryThreshold)

			// assert
			if err != tt.out {
				t.Errorf("got %v, want %v", err, tt.out)
			}
			if tt.assert != nil {
				tt.assert(t, f)
			}
		})
	}
```

example above we are iterating all test cases and running sub tests, on each case we creating `f := &depFields{}` of mocks and passing it onto `on` function to prepare mocks, at the end of the test we also checking `assertion` of mock calls, this helps us to assert number of calls that mocked methods are called inside the `Diagnose` method. 

## Links to read
[Table-Driven tests in Go](https://github.com/golang/go/wiki/TableDrivenTests)

[Testify Mocking](https://github.com/stretchr/testify#mock-package)
