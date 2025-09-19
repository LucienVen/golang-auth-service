package data_gen

import (
	"testing"
)

func TestGenerateBatchChineseNames(t *testing.T) {
	res := GenerateBatchChineseNames(10)
	t.Log(res)
}

func TestGenTimestamp(t *testing.T) {
	input := "202306"
	res, err := GenTimestamp(input)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(res)
	}
}
