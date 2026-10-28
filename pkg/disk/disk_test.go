package disk

import (
	"testing"
)

func TestListDisks(t *testing.T) {
	disks, err := ListDisks()
	if err != nil {
		t.Fatalf("ListDisks failed: %v", err)
	}

	if len(disks) == 0 {
		t.Errorf("expected at least one disk enumerated, got 0")
	}

	jsonStr, err := FormatDisksJSON()
	if err != nil {
		t.Fatalf("FormatDisksJSON failed: %v", err)
	}

	if len(jsonStr) == 0 {
		t.Errorf("expected non-empty JSON disk output")
	}
}
