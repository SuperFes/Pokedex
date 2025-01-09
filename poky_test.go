package main

import (
	"testing"
	"time"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  Hello  World  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "CARRY bridge lEft piNg",
			expected: []string{"carry", "bridge", "left", "ping"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("Expected %v, got %v", c.expected, actual)

			t.Failed()
		} else {
			t.Logf("Expected %v, got %v", len(c.expected), len(actual))
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf("Expected %v, got %v", expectedWord, word)

				t.Failed()
			} else {
				t.Logf("Expected %v, got %v", expectedWord, word)
			}
		}
	}
}

func TestCache(t *testing.T) {
	c := make(Caches)

	c.Set("test", []byte("test"), 1*time.Second)

	if len(c) != 1 {
		t.Errorf("Expected 1, got %v", len(c))

		t.Failed()
	} else {
		t.Logf("Expected 1, got %v", len(c))
	}

	time.Sleep(2 * time.Second)

	c.Set("test2", []byte("test2"), 1*time.Second)

	if len(c) != 1 {
		t.Errorf("Expected 1, got %v", len(c))

		t.Failed()
	} else {
		t.Logf("Expected 1, got %v", len(c))
	}

	if _, ok := c.Get("test"); ok {
		t.Errorf("Expected false, got true")

		t.Failed()
	} else {
		t.Logf("Expected false, got true")
	}

	if _, ok := c.Get("test2"); !ok {
		t.Errorf("Expected true, got false")

		t.Failed()
	} else {
		t.Logf("Expected true, got false")
	}
}
