package datastore

import (
  // "log"
  "os"
  "fmt"
  // "strings"
)

// type Datastore interface {
//   append()
// }
//
// func (ds *Datastore) append() int {
//   return 1
// }
//
// func (ds *Datastore) read(key string) string {
//  
//   return 'some value'
// }



const DATA_DIR = "db_data"

func Append(key string, value string) {

  err := os.Mkdir(DATA_DIR, os.ModePerm)
  if err != nil && !os.IsExist(err) {
    // todo: create a directory in default location
    panic(err)
  }

  // todo: make sure this file is written by only one process at a time.
  fo, err := os.OpenFile("db_data/active.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
  if err != nil {
    panic(err)
  }

  defer func() {
    if err := fo.Close(); err != nil {
      panic(err)
    }
  }()
  
  // inserting a single entry

  data := key + "," + value + "\n"
  fmt.Println(data)
  if _, err := fo.Write([]byte(data)); err != nil {
    fo.Close()
    panic(err)
  }
}

//appen('name', 'pranay')
