package computer

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

type mockedCPU struct {
	mock.Mock
}

func (m *mockedCPU) Usage() int {
	args := m.Called()
	return args.Get(0).(int)
}

type mockedGPU struct {
	mock.Mock
}

func (m *mockedGPU) Usage() int {
	args := m.Called()
	return args.Get(0).(int)
}

type mockedRAM struct {
	mock.Mock
}

func (m *mockedRAM) FreeMemory() int {
	args := m.Called()
	return args.Get(0).(int)
}

func TestMacbook(t *testing.T) {
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
		{
			name: "when CPU usage larger than CPU threshold diagnose return CpuUtilizationError",
			in:   &args{50, 60, 1000},
			out:  CpuUtilizationError,
			on: func(df *depFields) {
				df.cpu.On("Usage").Return(60) // 60% CPU usage
			},
			assert: func(df *depFields) {
				df.cpu.AssertNumberOfCalls(t, "Usage", 1)
			},
		},
		{
			name: "when GPU usage larger than GPU threshold diagnose return GpuUsageError",
			in:   &args{50, 90, 1000},
			out:  GpuUsageError,
			on: func(df *depFields) {
				df.cpu.On("Usage").Return(40) // 40% CPU usage less than cpuThreshold
				df.gpu.On("Usage").Return(95) // 95% gpu usage larger than gpuThreshold
			},
			assert: func(df *depFields) {
				df.cpu.AssertNumberOfCalls(t, "Usage", 1)
				df.gpu.AssertNumberOfCalls(t, "Usage", 1)
			},
		},

		{
			name: "when Free memory less than memory threshold diagnose return MemoryUsageError",
			in:   &args{50, 90, 1000},
			out:  MemoryUsageError,
			on: func(df *depFields) {
				df.cpu.On("Usage").Return(40)          // 40% CPU usage less than cpuThreshold
				df.gpu.On("Usage").Return(50)          // 50% gpu usage less than gpuThreshold
				df.memory.On("FreeMemory").Return(900) // 900 MB free memory left so it is less than 1000 mb threshold
			},
			assert: func(df *depFields) {
				df.cpu.AssertNumberOfCalls(t, "Usage", 1)
				df.gpu.AssertNumberOfCalls(t, "Usage", 1)
				df.memory.AssertNumberOfCalls(t, "FreeMemory", 1)
			},
		},

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
}
