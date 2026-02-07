//nolint:testpackage // Tests validate internal config form helpers directly.
package tui

import (
	"testing"

	"github.com/mholtzscher/ugh/internal/config"
)

func TestConfigFormApplyToConfigParsesBoolAndTheme(t *testing.T) {
	base := config.Config{Version: config.DefaultVersion}
	form := configFormState{
		values: configFormValues{
			syncURL:     "libsql://example",
			authToken:   "token",
			syncOnWrite: "true",
			theme:       "ansi-light",
		},
	}

	updated, err := form.applyToConfig(base)
	if err != nil {
		t.Fatalf("applyToConfig returned error: %v", err)
	}

	if !updated.DB.SyncOnWrite {
		t.Fatal("expected sync_on_write to be true")
	}
	if updated.UI.Theme != "ansi-light" {
		t.Fatalf("expected theme ansi-light, got %q", updated.UI.Theme)
	}
}

func TestConfigFormApplyToConfigRejectsInvalidBool(t *testing.T) {
	form := configFormState{values: configFormValues{syncOnWrite: "nope"}}
	_, err := form.applyToConfig(config.Config{Version: config.DefaultVersion})
	if err == nil {
		t.Fatal("expected error for invalid sync_on_write")
	}
}

func TestMaskToken(t *testing.T) {
	if got := maskToken(""); got != "-" {
		t.Fatalf("expected '-' for empty token, got %q", got)
	}
	if got := maskToken("abcd"); got != "********" {
		t.Fatalf("expected fallback mask for short token, got %q", got)
	}
	if got := maskToken("abcdefgh1234"); got != "abcd...1234" {
		t.Fatalf("unexpected masked token: %q", got)
	}
}
