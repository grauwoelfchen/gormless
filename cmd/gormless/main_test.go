package main

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestPathToID(t *testing.T) {
	tests := map[string]struct {
		str      string
		expected string
	}{
		"empty": {
			"",
			"",
		},
		"unknown": {
			"unknown",
			"unknown",
		},
		"directory": {
			filepath.Join("/", "path", "to", "dir"),
			"dir",
		},
		"invalid file": {
			filepath.Join("/", "path", "to", "file.go"),
			"file.go",
		},
		"plugin": {
			filepath.Join("/", "path", "to", "plugin.so"),
			"plugin",
		},
		"replace so": {
			filepath.Join("/", "path", "to", "awesome.so"),
			"awesome",
		},
		"replace .so once": {
			filepath.Join("/", "path", "to", "awesome.so.so"),
			"awesome.so",
		},
		"expected migration plugin file": {
			filepath.Join("/", "path", "to", "20190504_create_users.so"),
			"20190504_create_users",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := pathToID(tc.str)
			if tc.expected != result {
				t.Errorf("pathToID(%v) -> %v, want: %v", tc.str, result, tc.expected)
			}
		})
	}
}

func TestValidateAction(t *testing.T) {
	tests := map[string]struct {
		str      string
		expected bool
	}{
		"empty": {
			"",
			false,
		},
		"unknown": {
			"unknown",
			false,
		},
		"valid: commit": {
			"commit",
			true,
		},
		"valid: migrate": {
			"migrate",
			true,
		},
		"valid: revert": {
			"revert",
			true,
		},
		"valid: rollback": {
			"rollback",
			true,
		},
		"valid: run": {
			"run",
			true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := validateAction(tc.str)
			if tc.expected != result {
				t.Errorf(
					"validateAction(%v) -> %v, want: %v",
					tc.str,
					result,
					tc.expected,
				)
			}
		})
	}
}

func TestValidateMigrationDirectory(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("invalid caller")
	}
	rootDir := filepath.Join(filepath.Dir(filename), "..", "..")

	tests := map[string]struct {
		str      string
		expected bool
	}{
		"no such directory": {
			"",
			false,
		},
		"not a directory (file)": {
			filepath.Join(rootDir, "README.md"),
			false,
		},
		"system directory": {
			"/dev",
			true,
		},
		"a directory": {
			filepath.Join(rootDir, "example"),
			true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := validateMigrationDirectory(tc.str)
			if tc.expected != result {
				t.Errorf(
					"validateMigrationDirectory(%v) -> %v, want: %v",
					tc.str,
					result,
					tc.expected,
				)
			}
		})
	}
}
