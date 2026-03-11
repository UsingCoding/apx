package landlock

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	ll "github.com/landlock-lsm/go-landlock/landlock"

	"github.com/UsingCoding/apx/internal/sandbox"
)

func TestFsRules_SystemPathsPresent(t *testing.T) {
	rules, err := fsRules(context.Background(), sandbox.Policy{}, discardLogger())
	assert.NoError(t, err)

	strs := rulesStr(rules)

	// Verify count covers at least system paths that exist on this machine.
	assert.Greater(t, len(rules), 0)

	// Check a representative subset of system RO paths.
	for _, p := range []string{"/usr", "/etc"} {
		if _, statErr := os.Stat(p); statErr == nil {
			assert.True(t, containsPath(strs, p), "expected rule for system RO path %s", p)
		}
	}

	// Check a representative subset of system RW paths.
	for _, p := range []string{"/tmp", "/dev"} {
		if _, statErr := os.Stat(p); statErr == nil {
			assert.True(t, containsPath(strs, p), "expected rule for system RW path %s", p)
		}
	}
}

func TestFsRules_WorkdirPresent(t *testing.T) {
	wd, err := os.Getwd()
	assert.NoError(t, err)
	if err != nil {
		return
	}

	rules, err := fsRules(context.Background(), sandbox.Policy{}, discardLogger())
	assert.NoError(t, err)
	if err != nil {
		return
	}

	strs := rulesStr(rules)
	assert.True(t, containsPath(strs, wd), "expected rule for working directory %s", wd)
}

func TestFsRules_PolicyROPathsIncluded(t *testing.T) {
	// Use /tmp which is guaranteed to exist on any Linux/macOS machine.
	rules, err := fsRules(context.Background(), sandbox.Policy{
		Filesystem: sandbox.Filesystem{
			ROPaths: []string{"/tmp"},
		},
	}, discardLogger())
	assert.NoError(t, err)
	if err != nil {
		return
	}

	strs := rulesStr(rules)
	assert.True(t, containsPath(strs, "/tmp"), "expected rule for policy RO path /tmp")
}

func TestFsRules_PolicyRWPathsIncluded(t *testing.T) {
	rules, err := fsRules(context.Background(), sandbox.Policy{
		Filesystem: sandbox.Filesystem{
			RWPaths: []string{"/tmp"},
		},
	}, discardLogger())
	assert.NoError(t, err)
	if err != nil {
		return
	}

	strs := rulesStr(rules)
	assert.True(t, containsPath(strs, "/tmp"), "expected rule for policy RW path /tmp")
}

func TestFsRules_FullDiskReadAccess(t *testing.T) {
	rules, err := fsRules(context.Background(), sandbox.Policy{
		Filesystem: sandbox.Filesystem{
			FullDiskReadAccess: true,
		},
	}, discardLogger())
	assert.NoError(t, err)
	if err != nil {
		return
	}

	strs := rulesStr(rules)
	assert.True(t, containsPath(strs, "/"), "expected RODirs('/') rule when FullDiskReadAccess is set")
}

func TestFsRules_NonExistentPathSkipped(t *testing.T) {
	const ghost = "/nonexistent/apx/test/path"

	rules, err := fsRules(context.Background(), sandbox.Policy{
		Filesystem: sandbox.Filesystem{
			ROPaths: []string{ghost},
		},
	}, discardLogger())
	assert.NoError(t, err)
	if err != nil {
		return
	}

	strs := rulesStr(rules)
	assert.False(t, containsPath(strs, ghost), "non-existent path should be skipped")
}

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// rulesStr returns the fmt.Sprintf representation of each rule.
func rulesStr(rules []ll.FSRule) []string {
	out := make([]string, len(rules))
	for i, r := range rules {
		out[i] = r.String()
	}
	return out
}

// containsPath returns true when at least one rule string contains path.
func containsPath(strs []string, path string) bool {
	for _, s := range strs {
		if strings.Contains(s, path) {
			return true
		}
	}
	return false
}
