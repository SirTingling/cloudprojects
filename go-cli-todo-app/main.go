package main

import (
// "fmt"
)

func main() {
	todos := Todos{}
	todos.add("Buy milk")
	todos.add("Buy bread")
	todos.toggle(0)
	todos.print()
}
