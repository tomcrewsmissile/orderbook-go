package orderbookgo



type Entry struct {
	Price float64
	Quantity float64
	seq interface{}
}

type EntryStats struct {

	Distance float64 //represents a percentage decimal of the distance from the head node. 
	OrdersBetween int //the actual orders between the new order and the headnode

}