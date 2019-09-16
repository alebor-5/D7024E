package main

import (
	"fmt"
)

type RPC struct {
	
}

func SendPONG(target Contact){
	fmt.Println(target.Address)
}
