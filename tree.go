package orderbookgo

import "sync"



func BuildTree([]orderbookTentry) {

}


type OrderbookTree struct {
	headNode *orderbookTentry
	mutex *sync.RWMutex
}





func (ob *OrderbookTree) Insert(entry orderbookTentry) {
	ob.mutex.Lock()
	ob.headNode.insert(&entry)


	ob.mutex.Unlock()
}




type orderbookTentry struct {
	price float64
	quantity float64	
	seq interface{}	
	leftChild *orderbookTentry //less than
	rightChild *orderbookTentry //greater than 
	 
}

func ( entry *orderbookTentry) insert(newEntry *orderbookTentry)  {
	if entry.price == newEntry.price {
		entry.quantity = newEntry.quantity
		entry.seq = newEntry.seq
	} 


}


func (entry *orderbookTentry) getDepth() int {

	leftDepth := 0
	rightDepth := 0
	
	if entry.leftChild != nil {
		leftDepth = entry.leftChild.getDepth()
	}	
	if entry.rightChild != nil {
		rightDepth = entry.rightChild.getDepth() 
	}

	if leftDepth > rightDepth {
		return 1 + leftDepth
	}

	return 1 + rightDepth 
}

