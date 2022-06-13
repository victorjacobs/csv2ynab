package csv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColumnIndices(t *testing.T) {
	tests := map[string]struct {
		header                []string
		columnNamesList       [][]string
		expectedColumnIndices []int
	}{
		"a": {
			header: []string{"a", "b", "c"},
			columnNamesList: [][]string{
				{"x", "y"},
				{"z", "b"},
			},
			expectedColumnIndices: []int{-1, 1},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			columnIndices := columnIndices(tt.header, tt.columnNamesList...)

			assert.Equal(t, tt.expectedColumnIndices, columnIndices)
		})
	}
}
