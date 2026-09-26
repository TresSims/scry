package facts

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMemFactCollection(t *testing.T) {
	actual := &MemInfo{}
	err := ReadProcFileIntoStruct("../../testdata/facter/facts/proc_meminfo", actual)
	if err != nil {
		t.Errorf("Mem read failed: %s", err)

		return
	}

	expected := &MemInfo{
		Total: "31988108 kB",
		Free:  "389032 kB",
	}

	if !cmp.Equal(actual, expected) {
		t.Errorf("actual values do not match expected.\nactual: %+v\nexpected: %+v", actual, expected)
	}
}
