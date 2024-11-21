package docker

import "testing"

// test for fixed and not fixed parent

func TestParent(t *testing.T) {
	parent()
}

func TestParentFixed(t *testing.T) {
	parentFixed()
}

func FuzzEvent(f *testing.F) {
	testinputs := []int{5, 0, 50}

	for _, tc := range testinputs {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}
	f.Fuzz(func(t *testing.T, in int) {
		// we don't use the input in any way -> hinfällig
		parent()

	})
}
