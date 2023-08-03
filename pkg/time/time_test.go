package time

import (
	"fmt"
	"testing"
)

func TestTimeStringToGoTime(t *testing.T) {
	time := StringToGoTime("2023-07-19")
	fmt.Println(time)
}
