package bot

import "testing"

func TestCommandMatchUsesExactCommandToken(t *testing.T) {
	if !StartCommand.Match("/start now") {
		t.Fatal("expected /start with args to match start command")
	}

	if StartCommand.Match("/start123") {
		t.Fatal("expected /start123 not to match start command")
	}
}

func TestParseCommandNormalizesInput(t *testing.T) {
	got := ParseCommand("LIST_PRODUCTS?category=abc page")
	if got != ListProductsCommand {
		t.Fatalf("expected %q, got %q", ListProductsCommand, got)
	}
}
