package main

import (
	"os/user"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	spew.Dump(user.Current())
}
