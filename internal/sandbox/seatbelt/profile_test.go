package seatbelt

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/UsingCoding/apx/internal/sandbox"
)

func TestSnapshot_Default(t *testing.T) {
	prof, _, err := makeProfile(sandbox.Policy{Network: sandbox.Network{Deny: true}})
	assert.NoError(t, err)
	writeOrCompareSnapshot(t, "default", prof)
}

func TestSnapshot_FullDiskReadAccess(t *testing.T) {
	prof, _, err := makeProfile(sandbox.Policy{Filesystem: sandbox.Filesystem{FullDiskReadAccess: true}, Network: sandbox.Network{Deny: true}})
	assert.NoError(t, err)
	writeOrCompareSnapshot(t, "fs_full_disk_read", prof)
}

func TestSnapshot_NoCache(t *testing.T) {
	// HOME is used to compute cache dir; the profile itself references a param, so content is stable
	t.Setenv("HOME", filepath.Join(t.TempDir(), "home"))
	prof, _, err := makeProfile(sandbox.Policy{Filesystem: sandbox.Filesystem{NoCache: true}, Network: sandbox.Network{Deny: true}})
	assert.NoError(t, err)
	writeOrCompareSnapshot(t, "fs_no_cache", prof)
}

func TestSnapshot_ROPaths(t *testing.T) {
	prof, _, err := makeProfile(sandbox.Policy{Filesystem: sandbox.Filesystem{ROPaths: []string{"/any/a", "/any/b"}}, Network: sandbox.Network{Deny: true}})
	assert.NoError(t, err)
	writeOrCompareSnapshot(t, "fs_ro_paths", prof)
}

func TestSnapshot_RWPaths(t *testing.T) {
	prof, _, err := makeProfile(sandbox.Policy{Filesystem: sandbox.Filesystem{RWPaths: []string{"/tmp/a", "/tmp/b"}}, Network: sandbox.Network{Deny: true}})
	assert.NoError(t, err)
	writeOrCompareSnapshot(t, "fs_rw_paths", prof)
}

func TestSnapshot_NetworkAllowed(t *testing.T) {
	prof, _, err := makeProfile(sandbox.Policy{Network: sandbox.Network{Deny: false}})
	assert.NoError(t, err)
	writeOrCompareSnapshot(t, "network_allowed", prof)
}

func snapshotPath(name string) string {
	return filepath.Join("testdata", name+".golden")
}

func writeOrCompareSnapshot(t *testing.T, name, got string) {
	t.Helper()
	p := snapshotPath(name)
	if os.Getenv("UPDATE_SNAPSHOTS") == "1" {
		err := os.MkdirAll(filepath.Dir(p), 0o755)
		assert.NoErrorf(t, err, "mkdir testdata")
		if err == nil {
			err = os.WriteFile(p, []byte(got), 0o600)
			assert.NoErrorf(t, err, "write snapshot")
		}
		return
	}
	want, err := os.ReadFile(p)
	assert.NoErrorf(t, err, "read snapshot %s", p)
	if err != nil {
		return
	}
	assert.Equalf(t, string(want), got, "snapshot mismatch for %s", name)
}
