package main

import (
	"fmt"
	"proj3/concurrent"
)

func main() {
	dequeu := concurrent.NewUnBoundedDEQueue()
	// fmt.Println("Dequeu is empty: ", dequeu.IsEmpty())
	// fmt.Println("Dequeu contains ", dequeu.Size(), " elements")

	dequeu.PushBottom(1)
	dequeu.PushBottom(2)
	dequeu.PushBottom(3)

	// fmt.Println("Dequeu is empty: ", dequeu.IsEmpty())
	// fmt.Println("Dequeu contains ", dequeu.Size(), " elements")

	dequeu.Show()

	fmt.Println("------------------")

	fmt.Println("poptop : ", dequeu.PopTop())
	fmt.Println("pop bottom : ", dequeu.PopBottom())

	// fmt.Println("Dequeu is empty: ", dequeu.IsEmpty())
	// fmt.Println("Dequeu contains ", dequeu.Size(), " elements")

	dequeu.Show()

}
