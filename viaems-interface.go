package main

import (
	"fmt"
)

func main() {
	target, err := OpenTCPInterface("localhost:1234")
	if err != nil {
		fmt.Println(err)
		return
	}

	updates := target.GetStatusUpdates()
	go func() {
		for {
			u := <-updates
			fmt.Println(u)
		}
	}()

	x, err := target.ListTables()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(x)
	}

	t, err := target.GetTable("ve")
	fmt.Println(t)
	fmt.Println(err)

	select {}

}
