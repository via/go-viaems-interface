package main

import (
	"fmt"
  "github.com/via/go-viaems-interface/pkg/viaems"
  "encoding/json"
)

func displayThings(tg viaems.StatusTarget) {
  updates := tg.GetStatusUpdates()
  for {
    select {
    case status := <-updates:
      js, _ := json.Marshal(status)
      fmt.Println(string(js))
    }
  }
}


func main() {
	target, err := viaems.OpenTCPInterface("localhost:1234")
	if err != nil {
		fmt.Println(err)
		return
	}

	x, err := target.ListTables()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(x)
	}

	t, err := target.GetTable("ve")
	fmt.Println(t)
	fmt.Println(err)

  go displayThings(target)

	select {}

}
