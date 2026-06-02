package main

import (
	"fmt"

	"github.com/HarryCoburn/blog-aggregator/internal/config"
)

type state struct {
	state *config.Config
}

func main() {
	file, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	file.SetUser("harry")
	file, err = config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", file)

}
