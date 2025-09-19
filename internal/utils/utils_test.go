package utils

import (
	"github.com/brianvoe/gofakeit/v7"
	"testing"
)

func TestFormatFloat2Float(t *testing.T) {
	//t.Log(FormatFloat2Float(1.234567, 1))

	for i := 0; i < 10; i++ {

		num := gofakeit.IntRange(30, 99)
		f := FormatFloat2Float(float64(num), 1)
		t.Logf("num: %d, f: %.1f", num, f)

	}

}
