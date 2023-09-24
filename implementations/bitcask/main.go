package main

import (
  "fmt"
  "go-bitcask/datastore"
)

func main() {
  fmt.Println("Hello")
  datastore.Append("name", "pranay")
}
