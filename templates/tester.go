package main

import (
	"fmt"
	"github.com/gobuffalo/plush"
	"io/ioutil"
	"log"
)

func main() {
	c := plush.NewContext()
	h, err := plush.Render(html(), c)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(h)
}

func html() string {
	st, err := ioutil.ReadFile("/home/noor/src/go/src/muserblog/templates/application.html")
	if err != nil {
		log.Fatal(err)
	}
	return string(st)
}
