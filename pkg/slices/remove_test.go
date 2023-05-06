package slices

import (
	"reflect"
	"testing"
)

func TestRemove(t *testing.T) {
	testCases := []struct {
		name     string
		s        []string
		x        string
		expected []string
	}{
		{
			name:     "remove one item from slice with multiple items",
			s:        []string{"apple", "banana", "cherry"},
			x:        "banana",
			expected: []string{"apple", "cherry"},
		},
		{
			name:     "remove only item from slice with one item",
			s:        []string{"apple"},
			x:        "apple",
			expected: []string{},
		},
		{
			name:     "remove item not in slice",
			s:        []string{"apple", "banana", "cherry"},
			x:        "durian",
			expected: []string{"apple", "banana", "cherry"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := Remove(tc.s, tc.x)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("expected %v but got %v", tc.expected, actual)
			}
		})
	}
}
