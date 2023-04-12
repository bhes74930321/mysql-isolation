package main

import (

	"fmt"

	"example.com/data"

)

func main() {
	data.Init()
	
	fmt.Println("Example 1: Read Uncommitted")
	data.ReadUncommitted()

	fmt.Println("Example 2: Read Committed")
	data.ReadCommitted()

	fmt.Println("Example 3: Repeatable Read")
	data.RepeatableRead()

	fmt.Println("Example 4: Serializable")
	data.Serializable()
}
