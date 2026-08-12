package driver

import (
	"testing"

	"github.com/container-storage-interface/spec/lib/go/csi"
)

// A volume ID becomes a filesystem path, so anything that escapes the mount
// root has to be refused rather than cleaned up.
func TestVolumeDirRejectsEscapes(t *testing.T) {
	for _, id := range []string{
		"", "..", ".", "../etc", "a/b", "/abs", "pvc-1/../..", ".hidden", "a\x00b",
	} {
		if _, err := volumeDir("/mnt/etcfs", id); err == nil {
			t.Errorf("volumeDir accepted %q, want rejection", id)
		}
	}

	got, err := volumeDir("/mnt/etcfs", "pvc-3f0c9a12-4b1e")
	if err != nil {
		t.Fatalf("volumeDir rejected a valid ID: %v", err)
	}
	if want := "/mnt/etcfs/pvc-3f0c9a12-4b1e"; got != want {
		t.Errorf("volumeDir = %q, want %q", got, want)
	}
}

func TestValidateCapabilityRejectsBlock(t *testing.T) {
	block := &csi.VolumeCapability{
		AccessType: &csi.VolumeCapability_Block{Block: &csi.VolumeCapability_BlockVolume{}},
	}
	if err := validateCapability(block); err == nil {
		t.Error("validateCapability accepted a block volume")
	}

	mount := &csi.VolumeCapability{
		AccessType: &csi.VolumeCapability_Mount{Mount: &csi.VolumeCapability_MountVolume{}},
		AccessMode: &csi.VolumeCapability_AccessMode{
			Mode: csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER,
		},
	}
	if err := validateCapability(mount); err != nil {
		t.Errorf("validateCapability rejected a multi-writer mount volume: %v", err)
	}
	if err := validateCapabilities(nil); err == nil {
		t.Error("validateCapabilities accepted an empty capability list")
	}
}
