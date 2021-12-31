package orderbookgo

import (
	"testing"
)

func getTestLLOrderbook() *OrderbookList {


	list := NewOrderbookList()
	
	list.InsertAskSide(10, 2, nil)
	list.InsertAskSide(13, 2, nil)
	list.InsertAskSide(17, 2, nil)
	list.InsertAskSide(12, 2, nil)
	list.InsertAskSide(11, 2, nil)
	list.InsertAskSide(9, 3, nil)
	list.InsertAskSide(14, 4, nil)
	list.InsertAskSide(11, 3, nil)	
	//list.InsertBuySide(11, 0, nil)	
	list.InsertAskSide(12, 0 , nil)	
	list.InsertAskSide(9, 0, nil)	
	list.InsertAskSide(9, 1, nil)	
	list.InsertAskSide(10, 0, nil)	
	list.InsertAskSide(17, 0, nil)	
	list.InsertAskSide(20, 0, nil)	
	list.InsertAskSide(1, 0, nil)	
	return list

} 


func TestOrderbookLinkedListEntry(t *testing.T) {

	list := getTestLLOrderbook()

	for node := list.sellSideHeadNode; node != nil ; node = node.childNode {
		t.Log(node)

	}

	t.Fail()

}