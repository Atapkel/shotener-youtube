package random

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRandomString(t *testing.T) {
	testCases := []struct {
		name string
		size int
	}{
		{name: "size = 1", size: 1},
		{name: "size = 2", size: 3},
		{name: "size = 3", size: 3},
		{name: "size = 4", size: 4},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			str1 := NewRandomString(tc.size)
			str2 := NewRandomString(tc.size)
			fmt.Println(str2 + " " + str1)
			assert.Len(t, str1, tc.size)
			assert.Len(t, str2, tc.size)

			assert.NotEqual(t, str1, str2)
		})
	}
}
