// Package tribonacci contains the core logic for calculating tribonacci sequence numbers
package tribonacci

import (
	"math"
)

// Recurrent calculation of sequence members is CPU and/or memory heavy.
// So we can use function like Binet's Fibonacci number formula instead:
// https://ru.wikipedia.org/wiki/%D0%A7%D0%B8%D1%81%D0%BB%D0%B0_%D1%82%D1%80%D0%B8%D0%B1%D0%BE%D0%BD%D0%B0%D1%87%D1%87%D0%B8
// We can precalculate some constant arguments for increase performance:
var aplus = math.Pow((19 + 3*math.Sqrt(33)), 1.0/3)
var aminus = math.Pow((19 - 3*math.Sqrt(33)), 1.0/3)
var b = math.Pow((586 + 102*math.Sqrt(33)), 1.0/3)
var c = (aplus + aminus + 1) / 3.0
var d = 3 * b / (math.Pow(b, 2) - 2*b + 4)

// Number function calculates N'th tribonacci number
func Number(n uint64) float64 {
	// we use 1-indexed values
	if n < 3 {
		return 0
	}
	return math.Round(math.Pow(c, float64(n-2)) * d)
}
