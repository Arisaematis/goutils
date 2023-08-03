package slice

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	s := []string{"hello", "world", "hello", "golang", "hello", "ruby", "php", "java"}
	fmt.Printf("RemoveRepByMap(s): %v\n", RemoveRepByMap(s))
	fmt.Printf("RemoveRepByLoop(s): %v\n", RemoveRepByLoop(s))
}
