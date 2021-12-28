package orderbookgo

import "sync"



type Entry struct {
	Price float64
	Quantity float64
	Seq interface{}
}

type EntryStats struct {

	Distance float64 //represents a percentage decimal of the distance from the head node. 
	OrdersBetween int //the actual orders between the new order and the headnode

}

var entryPool = sync.Pool{}

func GetEntry() *Entry {

	entry := entryPool.Get()
	if entry != nil {
		return entry.(*Entry)
	}
	return &Entry{}
}

func ReleaseEntry(e *Entry) {

	entryPool.Put(e)

}