package orderbookgo

import (
	"fmt"
	//"log"
	"testing"
	"strconv"
)

func getTestEntries() []Entry {



	entries := make([]Entry, 100, 200)

	for i, _ := range entries {
		entries[i] = Entry{Price: float64(i + 10), Quantity: 1}
	} 

	return entries

}


func TestBuildTrees(t *testing.T) {

	entries := getTestEntries()

	treeTop := buildTrees(entries)

	DisplayTree(treeTop)



}

func DisplayTree(top *orderbookTentry) {

	node := top
	
	strconv.FormatFloat(node.Price, 'f', -1, 64)
	printChildren(node)

}

func printChildren(node *orderbookTentry) {

	children := make([]*orderbookTentry, 0, 10)
	children = append(children, node.leftChild)
	children = append(children, node.rightChild)	
	for ;len(children) > 0 ; {
		//prints children	
		for i := 0; i < len(children); i++ {
			
			if children[i] == nil {
				fmt.Print("  ")
				continue
			}	
			
			fmt.Print(strconv.FormatFloat(children[i].Price, 'f', -1, 64) + "  ")
		} 
		fmt.Print("\n")

		//adds children of children to new slice

		newChildren := make([]*orderbookTentry, 0, len(children) * 2)
			
		notnil := false
		for i := 0; i < len(children); i++ {
		
			
			if children[i] == nil {
				newChildren = append(newChildren, nil, nil)
			}

			if children[i].leftChild != nil {
				notnil = true
				newChildren = append(newChildren, children[i].leftChild)
			}

			if children[i].rightChild != nil {
				notnil = true
				newChildren = append(newChildren, children[i].rightChild)
			}

		}

		if !notnil {
			break
		}
	
		children = newChildren
		} 


}