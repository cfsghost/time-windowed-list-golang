package main

import (
  "time"
  "sort"
  "fmt"
)

type TimeWindowEntry interface {
  TStamp()      time.Time
}

type TimeWindow struct {
  Entries       []TimeWindowEntry
}

type TimeHash int64
type ByTimeHash []TimeHash
func (a ByTimeHash) Len() int           { return len(a) }
func (a ByTimeHash) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTimeHash) Less(i, j int) bool { return a[i] < a[j] }

type TimeWindowedList struct {
  Windows       map[TimeHash]TimeWindow
  DurationType  time.Duration
  MaxDurations  int
}

func NewTimeWindowedList(durType time.Duration, maxDur int) TimeWindowedList {
  l := TimeWindowedList{}
  l.Windows = make(map[TimeHash]TimeWindow)
  l.DurationType = durType
  l.MaxDurations = maxDur

  return l
}

func (wl *TimeWindowedList) ExpireOldEntries() {
  windowCount := len(wl.Windows)
  windowOverflow := windowCount - wl.MaxDurations

  if windowOverflow > 0 {
      windowTimes := make([]TimeHash, 0, len(wl.Windows))
      for windowTime := range wl.Windows {
        windowTimes = append(windowTimes, windowTime)
      }

      sort.Sort(ByTimeHash(windowTimes))
      for i := 0; i < windowOverflow; i++ {
        delete(wl.Windows, windowTimes[i])
      }
  }
}

func (wl *TimeWindowedList) Add(entry TimeWindowEntry) {
  entryTime := TimeHash(entry.TStamp().Truncate(wl.DurationType).Unix())

  if window, ok := wl.Windows[entryTime]; ok {
    wl.Windows[entryTime] = TimeWindow{Entries: append(window.Entries, entry)}
  } else {
    wl.Windows[entryTime] = TimeWindow{Entries: []TimeWindowEntry{entry}}
  }

  wl.ExpireOldEntries()
}

func (wl *TimeWindowedList) All() []TimeWindowEntry {
  wl.ExpireOldEntries()
  startTime := time.Now().Add(-time.Duration(wl.MaxDurations) * wl.DurationType).Unix()

  var all []TimeWindowEntry
  for windowTime := range wl.Windows {
    window := wl.Windows[windowTime]

    for i := range window.Entries {
      entry := window.Entries[i]

      if(entry.TStamp().Unix() > startTime) {
        all = append(all, entry)
      }
    }
  }

  return all
}

func (wl *TimeWindowedList) DisplayContents() {
  wl.ExpireOldEntries()
  startTime := time.Now().Add(-time.Duration(wl.MaxDurations) * wl.DurationType).Unix()

  fmt.Println("----")
  for windowTime := range wl.Windows {
    window := wl.Windows[windowTime]

    var allInBucket []TimeWindowEntry
    for i := range window.Entries {
      entry := window.Entries[i]

      if(entry.TStamp().Unix() > startTime) {
        allInBucket = append(allInBucket, entry)
      }
    }

    fmt.Printf("%d - %d\n", windowTime, len(allInBucket))
  }

  fmt.Printf("\nTotal in list: %d\n", len(wl.All()))
  fmt.Println("----\n")
}
