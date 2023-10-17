package util

import "fmt"

func SafeGo(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
				return
			}
		}()
		f()
	}()
}
