package main

import (
	"fmt"
	"sort"
	"testing"
)

func Test_randInRangeMaxValue(t *testing.T) {
	type args struct {
		min  int
		max  int
		size int
	}
	tests := []struct {
		args
	}{
		{args{min: 5, max: 5, size: 1}},
		{args{min: 1, max: 10, size: 1000}},
		{args{min: 1, max: 2, size: 1_000_000}},
		{args{min: 0, max: 42, size: 100_000_000}},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("min=%d, max=%d, size=%d", tt.min, tt.max, tt.size)
		t.Run(name, func(t *testing.T) {
			randInts := make([]int, tt.size)
			for i := 0; i < tt.size; i++ {
				randInts[i] = randInRange(tt.min, tt.max)
			}
			sort.Ints(randInts)
			maxInArr := randInts[len(randInts)-1]
			if maxInArr > tt.max {
				t.Errorf("want max: %d, got max: %d", tt.max, maxInArr)
			}
		})
	}
}
