package exceptions

import "testing"

func TestDeliberateFailure(t *testing.T) {
	t.Fatal("deliberate CI breakage")
}
