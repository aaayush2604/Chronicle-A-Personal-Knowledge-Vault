package config

import "testing"

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.PageSize != 10 {
		t.Fatalf("expected PageSize 10, got %d", cfg.PageSize)
	}

	if !cfg.ShowBanner {
		t.Fatalf("expected ShowBanner=true")
	}

	if len(cfg.Paths) != 0 {
		t.Fatalf("expected empty Paths")
	}
}

func TestNormalizeLogPathFile(t *testing.T) {
	path := normalizeLogPath("/tmp/chronicle.log")

	if path != "/tmp/chronicle.log" {
		t.Fatalf("unexpected path: %s", path)
	}
}

func TestDeduplicate(t *testing.T) {
	input := []string{
		"a",
		"b",
		"a",
		"c",
		"b",
	}

	out := deduplicate(input)

	if len(out) != 3 {
		t.Fatalf("expected 3 unique paths, got %d", len(out))
	}
}

func TestVersionExists(t *testing.T) {
	if Version == "" {
		t.Fatalf("version should not be empty")
	}
}
