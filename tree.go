package orderbookgo

import "sync"


//currently this function assumes sorted data in ascending order
func BuildTree(entries []Entry) *OrderbookTree {
	if entries == nil {
		return &OrderbookTree{
			mutex: &sync.RWMutex{},
		}
	}

	node := buildTrees(entries)

	return &OrderbookTree{
		mutex: &sync.RWMutex{},
		headNode: node,
	}

}

func buildTrees(entries []Entry) *orderbookTentry {

	mid := len(entries) / 2 //gets the middle index of the slice

	//base cases
	if len(entries) < 3 {
		
		
		switch len(entries) {

		case 0: 
			//idk how this could happen
		case 1:
			return &orderbookTentry{
				Entry: entries[0],
			}
		case 2:
			child := &orderbookTentry{
				Entry: entries[0],

			}
			return &orderbookTentry{
				Entry: entries[1],
				leftChild: child,
			}
		}	
	}
	

	return &orderbookTentry {
		Entry: entries[mid], 
		leftChild: buildTrees(entries[:mid]),
		rightChild: buildTrees(entries[mid + 1:]),
	}


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
	Entry	
	leftChild *orderbookTentry //less than
	rightChild *orderbookTentry //greater than 
	 
}

func ( entry *orderbookTentry) insert(newEntry *orderbookTentry)  {
	if entry.Price == newEntry.Price {
		entry.Quantity = newEntry.Quantity
		entry.Seq = newEntry.Seq
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

