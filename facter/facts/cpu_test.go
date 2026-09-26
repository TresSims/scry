package facts

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// TestCPUFactCollection tests reading a sample proc file
func TestCPUFactCollection(t *testing.T) {
	actual := &CPUInfo{}
	err := ReadProcFileIntoStruct("../../testdata/facter/facts/proc_cpuinfo", actual)
	if err != nil {
		t.Errorf("CPU read failed: %s", err)

		return
	}

	expected := &CPUInfo{
		Name:  "AMD Ryzen 7 9800X3D 8-Core Processor",
		Cores: 8,
		Clock: 5224.175,
	}

	if !cmp.Equal(actual, expected) {
		t.Errorf("actual values do not match expected.\nactual: %+v\nexpected: %+v", actual, expected)
	}
}
