package math

import "testing"

func TestAdd(t *testing.T) {

	got := Add(2, 2)
	want := 4

	if got != want {
		t.Errorf("got %d, wanted %d", got, want)
	}
}
