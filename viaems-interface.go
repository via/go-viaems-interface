package main

import (
  "fmt"
)

func main() {
	fmt.Println("Hello, World\n")

	target, err := OpenTCPInterface("localhost:1234")
	if err != nil {
		fmt.Println(err)
		return
	}

  go func() {
    for {
      _ = <-target.updates
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
