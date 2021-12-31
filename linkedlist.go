package orderbookgo

import (
	"sync"
	"log"
)

//double linked list implementation

func NewOrderbookList() *OrderbookList{
	return &OrderbookList{


		buysidemutex: &sync.RWMutex{},
		sellsidemutex: &sync.RWMutex{},
	}

}


type OrderbookList struct {
	buySideHeadNode  *OrderbookLEntry
	buysidedepth int	
	
	sellSideHeadNode *OrderbookLEntry
	sellsidedepth int	
	
	buysidemutex  *sync.RWMutex
	sellsidemutex *sync.RWMutex

	lastSeq interface{}
}

func (orderbook *OrderbookList) InsertBuySide(price float64, quantity float64, seq interface{}) {

	orderbook.buysidemutex.Lock()
	defer orderbook.buysidemutex.Unlock()
	
	newEntry := newEntry(price, quantity, seq)

	orderbook.lastSeq = seq

	//base case
	if orderbook.buySideHeadNode == nil && quantity > 0{
		orderbook.buySideHeadNode = newEntry
		orderbook.buysidedepth++
		return
	}

	//second base case. Entry is the new best bid to buy
	if price > orderbook.buySideHeadNode.Price && quantity > 0 {
		newEntry.childNode = orderbook.buySideHeadNode
		orderbook.buySideHeadNode = newEntry
		orderbook.buysidedepth++	
		return
	}

	//in this case the order needs to be removed
	if quantity == 0 {
		//find the order
		var tmpnodeParent *OrderbookLEntry 	
		tmpnodeParent = nil	
		tmpnode := orderbook.buySideHeadNode
		for ;tmpnode != nil ; tmpnodeParent, tmpnode = tmpnode, tmpnode.childNode {
			if tmpnode.Price == price {
				if tmpnodeParent != nil {
					tmpnodeParent.childNode = tmpnode.childNode
				} else {
					orderbook.buySideHeadNode = tmpnode.childNode //this is a base case where the headnode is being removed
				}
				orderbook.buysidedepth--
				return
			}
		} 
		return
	}

	var tmpnodeParent *OrderbookLEntry
	tmpnode := orderbook.buySideHeadNode
	for ; tmpnode != nil && tmpnode.Price >= price; tmpnodeParent, tmpnode = tmpnode, tmpnode.childNode  {
		
		if tmpnode.Price == price {
			log.Print("price is the same")	
			tmpnode.Quantity = quantity
			tmpnode.Seq = seq
			return
		}
	}

	if tmpnode == nil {
		tmpnodeParent.childNode = newEntry
		orderbook.buysidedepth++
		return
	}

	tmpnodeParent.childNode = newEntry
	newEntry.childNode = tmpnode	
	orderbook.buysidedepth++	
	return


}

//returns a
func (ob *OrderbookList) GetBuyside() []Entry {
	
	if ob.buySideHeadNode == nil {	
		return nil
	}

	entries := make([]Entry, ob.buysidedepth)

	for i, node := 0, ob.buySideHeadNode; i < ob.buysidedepth; i++ {
		entries[i] = node.Entry	
		node = node.childNode
	}

	return entries

}


func (orderbook *OrderbookList) InsertAskSide(price float64, quantity float64, seq interface{})  {

	orderbook.sellsidemutex.Lock()
	defer orderbook.sellsidemutex.Unlock()

	orderbook.lastSeq = seq	

	newEntry := newEntry(price, quantity, seq)



	//base case
	if orderbook.sellSideHeadNode == nil {
		orderbook.sellSideHeadNode = newEntry
		orderbook.sellsidedepth++
		return 
	}

	//second base case. Entry is the new best bid to buy
	if price < orderbook.sellSideHeadNode.Price && quantity > 0{
		newEntry.childNode = orderbook.sellSideHeadNode
		orderbook.sellSideHeadNode = newEntry
		orderbook.sellsidedepth++	
		return 
	}

	//in this case the order needs to be removed
	if quantity == 0 {
		//find the order
		var tmpnodeParent *OrderbookLEntry 	
		tmpnodeParent = nil	
		tmpnode := orderbook.sellSideHeadNode
		for ;tmpnode != nil ; tmpnodeParent, tmpnode = tmpnode, tmpnode.childNode {
			if tmpnode.Price == price {
				if tmpnodeParent != nil {
					tmpnodeParent.childNode = tmpnode.childNode
				} else {
					orderbook.sellSideHeadNode = tmpnode.childNode //this is a base case where the headnode is being removed
				}
				orderbook.sellsidedepth--
				return
			}
		} 
		return
	}

	var tmpnodeParent *OrderbookLEntry
	tmpnode := orderbook.sellSideHeadNode
	for ; tmpnode != nil && tmpnode.Price <= price; tmpnodeParent, tmpnode = tmpnode, tmpnode.childNode  {
		
		if tmpnode.Price == price {
			log.Print("price is the same")	
			tmpnode.Quantity = quantity
			tmpnode.Seq = seq
			return
		}
	}

	if tmpnode == nil {
		tmpnodeParent.childNode = newEntry
		orderbook.sellsidedepth++
		return
	}

	tmpnodeParent.childNode = newEntry
	newEntry.childNode = tmpnode	
	orderbook.sellsidedepth++	
	return


}

func (ob *OrderbookList) GetSellside() []Entry {
	
	if ob.sellSideHeadNode == nil {	
		return nil
	}

	entries := make([]Entry, ob.sellsidedepth)

	for i, node := 0, ob.sellSideHeadNode; i < ob.sellsidedepth; i++ {
		entries[i] = node.Entry	
		node = node.childNode
	}

	return entries

}



func newEntry(price float64, quantity float64, seq interface{}) *OrderbookLEntry {

	return &OrderbookLEntry{
		Entry{price, quantity, seq},
		nil,
	}

}

type OrderbookLEntry struct {
	Entry
	childNode  *OrderbookLEntry
}


func (orderbook *OrderbookList) GetSpread() float64 {
	if orderbook.sellSideHeadNode == nil || orderbook.buySideHeadNode == nil {
		return 0
	}	
	
	return orderbook.sellSideHeadNode.Price - orderbook.buySideHeadNode.Price

}

func (orderbook *OrderbookList) GetMidPrice() float64 {
	if orderbook.sellSideHeadNode == nil || orderbook.buySideHeadNode == nil {
		return 0
	}
		return (orderbook.buySideHeadNode.Price + orderbook.sellSideHeadNode.Price) / 2.0

}



//returns the total dollar amount of buy orders to a given price but not including the price.
func (orderbook *OrderbookList) GetSellLiquidity(price float64) float64 {
	total := 0.0

	for node := orderbook.buySideHeadNode; node != nil && node.Price > price; node = node.childNode {
		total += node.Quantity * node.Price
	}

	return total

}

func (orderbook *OrderbookList) GetBuyLiquidity(price float64) float64 {
	total := 0.0

	for node := orderbook.sellSideHeadNode; node != nil && node.Price < price; node = node.childNode {
		total += node.Quantity * node.Price
	}

	return total

}


func (ob *OrderbookList) GetLastSequence() interface{} {


	return ob.lastSeq
}