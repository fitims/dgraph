package main

import (
	"fmt"
	"log"
	"sort"
	"strconv"
)

type uidMap struct {
	lastUID uint64
	uids    map[string]uint64
}

func newUIDMap() *uidMap {
	return &uidMap{
		lastUID: 1,
		uids:    map[string]uint64{},
	}
}

func (m *uidMap) uid(str string) uint64 {

	hint, err := strconv.ParseUint(str, 10, 64)
	if err == nil {
		uid, ok := m.uids[str]
		if ok {
			if uid == hint {
				return uid
			} else {
				log.Fatalf("bad node hint: %v", str)
			}
		} else {
			m.uids[str] = hint
			return hint
		}
	}

	uid, ok := m.uids[str]
	if ok {
		return uid
	}
	m.lastUID++
	m.uids[str] = m.lastUID
	return m.lastUID
}

func (m *uidMap) lease() uint64 {
	// lastUID => lease
	//    9999 => 10001
	//   10000 => 10001
	//   10001 => 10001
	//   10002 => 20001
	//   10003 => 20001
	if m.lastUID <= 2 {
		return 10001
	} else {
		return (m.lastUID-2)/10000*10000 + 10001
	}
}

func (m *uidMap) logState() {
	type entry struct {
		uid uint64
		str string
	}
	entries := make([]entry, 0, len(m.uids))
	for str, uid := range m.uids {
		entries = append(entries, entry{uid, str})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].uid < entries[j].uid })
	fmt.Printf("UID Map: (%d entries)\n", len(entries))
	for _, entry := range entries {
		fmt.Printf("%d:%q\n", entry.uid, entry.str)
	}
}
