package main

import (
  "time"
  "sort"
  "fmt"
)

type TimeWindowedList struct {
  Windows       map[TimeHash]TimeWindow
  DurationType  time.Duration
  MaxDurations  int
}

type TimeWindow struct {
  Entries       []TimeWindowEntry
}

type TimeWindowEntry interface {
  TStamp()      time.Time
}

type TimeHash int64
type ByTimeHash []TimeHash
func (a ByTimeHash) Len() int           { return len(a) }
func (a ByTimeHash) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTimeHash) Less(i, j int) bool { return a[i] < a[j] }

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
  var all []TimeWindowEntry

  for windowTime := range wl.Windows {
    window := wl.Windows[windowTime]

    for i := range window.Entries {
      all = append(all, window.Entries[i])
    }
  }

  return all
}

func (wl *TimeWindowedList) DisplayContents() {
  wl.ExpireOldEntries()

  fmt.Println("----")
  for windowTime := range wl.Windows {
    window := wl.Windows[windowTime]

    fmt.Printf("%d - %d\n", windowTime, len(window.Entries))
  }

  fmt.Println("----\n")
}
