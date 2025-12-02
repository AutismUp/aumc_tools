package mcprops

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return path
}

func TestLoadProperties(t *testing.T) {
	dir := t.TempDir()
	content := strings.Join([]string{
		"# comment",
		"key1=value1",
		"key2 = value2",
		"key3=value=with=equals",
		"",
		"# another comment",
		"key4=value4",
	}, "\n")

	path := writeTempFile(t, dir, "server.properties", content)

	config, err := LoadProperties(path)
	if err != nil {
		t.Fatalf("LoadProperties returned error: %v", err)
	}

	tests := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value=with=equals",
		"key4": "value4",
	}

	for key, expected := range tests {
		if got, ok := config.Properties[key]; !ok || got != expected {
			t.Fatalf("expected %s=%s, got %s", key, expected, got)
		}
	}
}

func TestUpdateProperty(t *testing.T) {
	cfg := &MCConfig{}
	cfg.UpdateProperty("key1", "value1")
	cfg.UpdateProperty("key2", "value2")
	cfg.UpdateProperty("key1", "newvalue")

	if cfg.Properties["key1"] != "newvalue" {
		t.Fatalf("expected key1 to be updated")
	}

	if len(cfg.order) != 2 {
		t.Fatalf("expected order length 2, got %d", len(cfg.order))
	}
}

func TestWriteProperties(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "server.properties")

	cfg := &MCConfig{
		Properties: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
		order: []string{"key1", "key2"},
	}

	cfg.UpdateProperty("key3", "value3")

	if err := cfg.WriteProperties(path); err != nil {
		t.Fatalf("WriteProperties returned error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "#Minecraft server properties") {
		t.Fatalf("missing header")
	}

	expectedLines := []string{"key1=value1", "key2=value2", "key3=value3"}
	for _, line := range expectedLines {
		if !strings.Contains(content, line) {
			t.Fatalf("missing line %s", line)
		}
	}
}
