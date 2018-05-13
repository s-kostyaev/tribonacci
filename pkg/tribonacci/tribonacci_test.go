package tribonacci

import (
	"fmt"
	"testing"
)

// https://ru.wikipedia.org/wiki/%D0%A7%D0%B8%D1%81%D0%BB%D0%B0_%D1%82%D1%80%D0%B8%D0%B1%D0%BE%D0%BD%D0%B0%D1%87%D1%87%D0%B8
// leading zero was added for simplify indexing inside test
var results = []float64{0, 0, 0, 1, 1, 2, 4, 7, 13, 24, 44, 81, 149, 274, 504, 927, 1705, 3136, 5768, 10609, 19513, 35890, 66012, 121415, 223317, 410744, 755476, 1389537, 2555757, 4700770, 8646064, 15902591, 29249425, 53798080, 98950096, 181997601, 334745777}

type testData struct {
	name string
	n    uint64
	want float64
}

func TestNumber(t *testing.T) {
	tests := make([]testData, 0, len(results)-2)
	for i := 1; i < len(results)-1; i++ {
		tests = append(tests, testData{
			name: fmt.Sprintf("%v number", i),
			n:    uint64(i),
			want: results[i],
		})
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Number(tt.n)
			if got != tt.want {
				t.Errorf("Number() = %v, want %v", got, tt.want)
			}
		})
	}
}
