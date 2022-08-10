package gotest

import (
	"fmt"
	"testing"
)

func TestBotPost(t *testing.T) {

}

func run() (code int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic")
			code = 2
			return
		}
		if code != 0 {
			fmt.Println("err")
			return
		}
		fmt.Println("ok")
		code = 0
	}()
	fmt.Println("run")
	return
}
