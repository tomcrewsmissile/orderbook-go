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
	if orderbook.buySideHeadNode == nil {
		orderbook.buySideHeadNode = newEntry
		orderbook.buysidedepth++
		return
	}

	//second base case. Entry is the new best bid to buy
	if price > orderbook.buySideHeadNode.Price {
		newEntry.childNode = orderbook.buySideHeadNode
		orderbook.buySideHeadNode.parentNode = newEntry	
		orderbook.buySideHeadNode = newEntry
		orderbook.buysidedepth++	
		return
	}

	tmpnode := orderbook.buySideHeadNode	
	for ; tmpnode.childNode != nil && tmpnode.Price > price; tmpnode = tmpnode.childNode {}

		//check to see if price update
		if tmpnode.Price == price {
			
			//in this case the order needs to be removed	
			if quantity == 0 {
				log.Print(tmpnode.parentNode)
				tmpnode.parentNode.childNode = tmpnode.childNode
				tmpnode.childNode.parentNode = tmpnode.parentNode
				return
			}	
			
			
			//update the corresponding node	
			tmpnode.Quantity = quantity
			tmpnode.Seq = seq
			return
		}


	//check to see if the node needs to be inserted at the end
	if tmpnode.childNode == nil && tmpnode.Price > price {
		tmpnode.childNode = newEntry
		newEntry.parentNode = tmpnode
		orderbook.buysidedepth++	
		return
	}

	//insertion into the list
	newEntry.parentNode = tmpnode.parentNode	
	newEntry.childNode = tmpnode	

	//reassign pointers
	newEntry.childNode.parentNode = newEntry
	newEntry.parentNode.childNode = newEntry

	orderbook.buysidedepth++	

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
	if price < orderbook.sellSideHeadNode.Price {
		newEntry.childNode = orderbook.sellSideHeadNode
		orderbook.sellSideHeadNode.parentNode = newEntry	
		orderbook.sellSideHeadNode = newEntry
		orderbook.sellsidedepth++	
		return 
	}

	tmpnode := orderbook.sellSideHeadNode	
	for ; tmpnode.childNode != nil && tmpnode.Price < price; tmpnode = tmpnode.childNode {}

		//check to see if price update
		if tmpnode.Price == price {
			
			if quantity == 0 { //in this case, the order needs to be removed
				log.Print(tmpnode.parentNode)	
				
				tmpnode.parentNode.childNode = tmpnode.childNode
				tmpnode.childNode.parentNode = tmpnode.parentNode
				return
			}  	
			
			//update the corresponding node	
			tmpnode.Quantity = quantity
			tmpnode.Seq = seq
			return
		}


	//check to see if the node needs to be inserted at the end
	if tmpnode.childNode == nil && tmpnode.Price < price {
		tmpnode.childNode = newEntry
		newEntry.parentNode = tmpnode
		orderbook.sellsidedepth++	
		return
	}

	//insertion into the list
	newEntry.parentNode = tmpnode.parentNode	
	newEntry.childNode = tmpnode	

	//reassign pointers
	newEntry.childNode.parentNode = newEntry
	newEntry.parentNode.childNode = newEntry

	orderbook.sellsidedepth++	

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
		nil,
	}

}

type OrderbookLEntry struct {
	Entry
	parentNode *OrderbookLEntry
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