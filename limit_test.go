package oval

import (
	"fmt"
	"testing"
	"time"
)

func Test(t *testing.T) {
	for i := 0; i < 5; i++ {
		got := Limited("name", 2, 3)
		time.Sleep(1 * time.Second)
		fmt.Println(i, got)
	}
	UnLimited("name")
	for i := 0; i < 5; i++ {
		got := Limited("name", 10, 3)
		time.Sleep(1 * time.Second)
		fmt.Println(i, got)
	}
}
