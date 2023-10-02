package main

import (
	"encoding/json"
	"fmt"
	// "go-bitcask/datastore"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type Bitcask struct {
  mu    sync.RWMutex

  // An append only file in which new data is inserted
  activeFile  *os.File

  // map to store all the previous log segments
  segments  map[int]*os.File

  // last segment index
  segmentIndex int

  // Hashmap which stores key and value pairs
  hintMap map[string]DataLocation

  // hintMap written to the disk for next process start
  hintFile *os.File

  // maxFileSize of the active log file in KB
  maxDataFileSize int64

  // maxFileSize of the segment files in Kilobytes
  maxSegmentSize int64 

  // Database direcoty in which data is stored
  dataDir string
}

type DataLocation struct {
  SegmentIndex    int
  Offset          int64
  ValueSize       int
}

func (b *Bitcask) Initialize(dataDir string) (*Bitcask, error) {
  // b.mu.Lock()
  // defer b.mu.Unlock()

  b.maxDataFileSize = 1 // in Kilobytes
  b.maxSegmentSize = 2 // in Kilobytes
  b.dataDir = dataDir

  if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
    return nil, err
  }

  // Initialize the active log file

  activeFile, err := b.createFile("active.log")
  if err != nil {
    return nil, err
  }
  b.activeFile = activeFile

  // Initialize segment map
  b.segments = make(map[int]*os.File)

  // Initialize segment segmentIndex
  b.segmentIndex = 0

  // load hint map from hint file 

  hintFile, err := b.createFile("hintfile")
  if err != nil {
    return nil, err
  }

  b.hintFile = hintFile
  if err := b.loadHintMap(); err != nil {
    return nil, err
  }

  return b, nil
}

func (b *Bitcask) createFile (filename string) (*os.File, error) {
  filePath := filepath.Join(b.dataDir, filename)
  return os.OpenFile(filePath, os.O_CREATE | os.O_RDWR, 0644)
}

func (b *Bitcask) loadHintMap() error  {
  // b.mu.Lock()
  // defer b.mu.Unlock()

  data, err := io.ReadAll(b.hintFile)
  if err != nil {
    return err
  }

  fmt.Println("Size of the hintfile", len(data))
  if len(data) == 0 {
    b.hintMap = make(map[string]DataLocation)
    return nil
  }
  if err := json.Unmarshal(data, &b.hintMap); err != nil {
    return err
  }

  return nil
}

func (b *Bitcask) writeToDataFile(key, value string) (int64, error)  {
  b.mu.Lock()
  defer b.mu.Unlock()

  // return 1123141, nil

  // Check if active file size exceeded

  if fileInfo, err := b.activeFile.Stat(); err == nil && fileInfo.Size() > b.maxDataFileSize {
    if err := b.rotateActiveFile(); err != nil {
      return 0, err
    }
  }

  offset, err := b.activeFile.Seek(0, io.SeekEnd)
  if err != nil {
    return 0, err
  }

  data := fmt.Sprintf("%s:%s\n", key, value)
  _, err = b.activeFile.Write([]byte(data))
  if err != nil {
    return 0, err
  }

  return offset, nil
  
}

func (b *Bitcask) rotateActiveFile() error {
  b.mu.Lock()
  defer b.mu.Unlock()
  
  // close the current active file
  if err := b.activeFile.Close(); err != nil {
    return err
  }

  // Rename the active file to a segment file
  segmentFilename := fmt.Sprintf("%d.log", b.segmentIndex)
  segmentPath := filepath.Join(b.dataDir, segmentFilename)

  err := os.Rename(filepath.Join(b.dataDir, "active.log"), segmentPath) 
  if err != nil {
    return err
  }

  // open a new active file
  activeFile, err := b.createFile("active.log")
  if err != nil {
    return err
  }
  b.activeFile = activeFile
  
  b.segments[b.segmentIndex] = activeFile

  //Increment the segment Index
  b.segmentIndex++

  b.hintMap = make(map[string]DataLocation)

  // persist hint map to hint file
  if err := b.persistHintMap(); err != nil {
    return err
  }

  return nil
}

func (b *Bitcask) persistHintMap() error {
  b.mu.Lock()
  defer b.mu.Unlock()

  // Marshall the hint map to json

  data, err := json.Marshal(b.hintMap)
  if err != nil {
    return err
  }

  // Write json data to hint file
  if err := os.WriteFile(filepath.Join(b.dataDir, "hintfile"), data, 0644); err != nil {
    return err
  }

  return nil
}

func (b *Bitcask) Get(key string) (string, error) {
  b.mu.RLock()
  defer b.mu.RUnlock()

  location, ok := b.hintMap[key]
  if !ok {
    return "", fmt.Errorf("key not found")
  }

  value := make([]byte, location.ValueSize)
  _, err := b.activeFile.ReadAt(value, location.Offset)
  if err != nil {
    return "", err
  }

  return string(value), nil
}

func (b *Bitcask) Put(key string, value string) error {
  b.mu.Lock()
  defer b.mu.Unlock()

  offset, err := b.writeToDataFile(key, value)
  if err != nil {
    return err
  }

  // update hint map
  b.hintMap[key] = DataLocation{
    SegmentIndex: 0,
    Offset: offset,
    ValueSize: len(value),
  }

  if err := b.persistHintMap(); err != nil {
    return err
  }

  return nil
}

func (b *Bitcask) Close() {
  b.mu.Lock()
  defer b.mu.Unlock()
  
  b.activeFile.Close()
  b.hintFile.Close()
}

func main() {
  // fmt.Println("Hello")
  // datastore.Append("name", "pranay")
  db, err := (&Bitcask{}).Initialize("db_dir")
  if err != nil {
    fmt.Println("Error Initialize database:", err)
    return
  }

  defer db.Close()

  err = db.Put("key1", "value1")
  if err != nil {
    fmt.Println("Error putting data:", err)
    return
  }

  // Get
  value, err := db.Get("key1")
  if err != nil {
    fmt.Println("Error getting data: ", err)
    return
  }

  fmt.Println("Value:", value)
}
