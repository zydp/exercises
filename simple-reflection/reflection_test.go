package main

import (
	"fmt"
	"testing"
	"time"
)

func Test_PrintMethods(t *testing.T) {
	PrintMethods(time.Now())
	t.Error("abcdefg")
	fmt.Println(t.Name())
	t.SkipNow()
}
