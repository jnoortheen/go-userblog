package main

import (
	"io/ioutil"
	"log"
	"github.com/gobuffalo/plush"
	"fmt"
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
