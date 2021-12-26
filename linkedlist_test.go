package orderbookgo

import (
	"testing"
)

func getTestLLOrderbook() *OrderbookList {


	list := NewOrderbookList()
	
	list.InsertBuySide(10, 2, nil)
	list.InsertBuySide(13, 2, nil)
	list.InsertBuySide(17, 2, nil)
	list.InsertBuySide(12, 2, nil)
	list.InsertBuySide(11, 2, nil)
	list.InsertBuySide(9, 3, nil)
	list.InsertBuySide(14, 4, nil)
	list.InsertBuySide(11, 3, nil)	
	return list

} 


func TestOrderbookLinkedListEntry(t *testing.T) {

	list := getTestLLOrderbook()

	for node := list.buySideHeadNode; node != nil ; node = node.childNode {
		t.Log(node)

	}

	t.Fail()

}