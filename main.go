package main

import (
  "time"
  "fmt"
)

type MyEntry struct {
  HappenedAt    time.Time
}

func(e MyEntry) TStamp() time.Time {
  return e.HappenedAt
}

func main() {
  windowedList := NewTimeWindowedList(time.Second, 5)

  for {
    newEntry := MyEntry{HappenedAt: time.Now()}
    windowedList.Add(newEntry)

    fmt.Printf("Current Time: %d\n", time.Now().Unix())
    windowedList.DisplayContents()
  }
}
