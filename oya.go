package main

import (
	"fmt"
	"os"

	kasia "github.com/ziutek/kasia.go"
)

type Ctx struct {
	H, W string
}

func main() {
	ctx := map[string]string{"H": "Hello", "W": "world"}

	tpl, err := kasia.Parse("$H $W $$O!\n")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = tpl.Run(os.Stdout, ctx)
	if err != nil {
		fmt.Println(err)
	}
}
