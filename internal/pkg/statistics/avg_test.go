package statistics

import (
	"fmt"
	"testing"
)

func TestAvgWeighted(t *testing.T) {
	a1 := []int{2, 2, 4, 3, 6, 7}
	a2 := []int{4, 4, 4, 4, 2, 1, 9, 5}

	fmt.Println(AvgWeighted(a1))
	fmt.Println(AvgWeighted(a2))
	t.Fail()
}
