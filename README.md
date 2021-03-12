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
			assert: func(df *depFields) {
				df.cpu.AssertNumberOfCalls(t, "Usage", 1)
				df.gpu.AssertNumberOfCalls(t, "Usage", 1)
				df.memory.AssertNumberOfCalls(t, "FreeMemory", 1)
			},
		},
	}
```
in the table-driven testing we gonna define slice of struct that represents test cases, in the example above we have struct that contains callback functions `on: func(*depFields)` and `assert: func(*depFields)` this functions will be called when call test cases and dependencies will be passed that we wanted to `mock` and `assertion calls`

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
				tt.assert(f)
			}
		})
	}
```

we iterate all our test cases and run child test functions, on each case we creating dependency of mocks and passing it onto `on` function that prepares us mocks, at the end we also checking test cases `assertion` this helps us to assert number calls that mocked methods are called and etc. 

## Links to read
[Table-Driven tests in Go](https://github.com/golang/go/wiki/TableDrivenTests)

[Testify Mocking](https://github.com/stretchr/testify#mock-package)
