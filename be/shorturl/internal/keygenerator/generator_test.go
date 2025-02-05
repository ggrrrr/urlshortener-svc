package keygenerator

import (
	"fmt"
	"testing"
)

func TestGenerate(t *testing.T) {
	testG := &Generator{}

	for i := 0; i < 100; i++ {
		key := testG.Generate()
		fmt.Printf("%s \n", key)
	}
}
