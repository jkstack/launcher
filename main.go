package main

import (
	"fmt"
	"net/http"
)

const port = 19999

func main() {
	fmt.Println("Hello World")

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
