package main

import (
	"fmt"
	"os/user"
)

func main() {
	user, _ := user.Current()

	fmt.Printf("%v", user)
}
