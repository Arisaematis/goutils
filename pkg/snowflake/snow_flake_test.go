package snowflake

import (
	"fmt"
	"testing"
)

func TestGetSnowFlake(t *testing.T) {
	id := GetWorkId()
	fmt.Println(id)
}
