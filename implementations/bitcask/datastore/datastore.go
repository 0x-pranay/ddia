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
  // if err := os.WriteFile("file.txt", []byte)
  err := os.Mkdir(DATA_DIR, os.ModePerm)
  if err != nil && !os.IsExist(err) {
    panic(err)
  }

  fo, err := os.OpenFile("db_data/log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
  if err != nil {
    panic(err)
  }

  defer func() {
    if err := fo.Close(); err != nil {
      panic(err)
    }
  }()
  
  // buf := make([]byte, 1024)
  // for {
    // write that chunk
    // if _, err := fo.Write(buf[:n])

  // }
  // _, err := fmt.Fprintln(fo, key, "," value)
  // if err != nil {
  //   panic(err)
  // }
  data := key + "," + value + "\n"
  fmt.Println(data)
  // if _, err := fo.WriteString(data); err != nil {
  //   panic(err)
  // }

  if _, err := fo.Write([]byte(data)); err != nil {
    fo.Close()
    panic(err)
  }
}

//appen('name', 'pranay')
