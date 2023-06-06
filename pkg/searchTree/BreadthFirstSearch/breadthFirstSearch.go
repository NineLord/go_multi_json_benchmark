package BreadthFirstSearch

import (
	"github.com/NineLord/go_multi_json_benchmark/pkg/utils/Vector"
	"reflect"
)

func Run(tree interface{}, value interface{}) bool {
	currentNodes := Vector.NewVector[interface{}](0)
	currentNodes.Push(tree)
	nextLevelNodes := Vector.NewVector[interface{}](0)

	for currentNodes.Len() != 0 {
		currentNode := currentNodes.Pop()

		if currentNode == nil {
			if currentNode == value {
				return true
			}
		} else {
			switch reflect.TypeOf(currentNode).Kind() {
			case reflect.Map:
				currentNode := currentNode.(map[string]interface{})
				for nodeKey, nodeValue := range currentNode {
					if nodeKey == value {
						return true
					}
					nextLevelNodes.Push(nodeValue)
				}
			case reflect.Slice:
				currentNode := currentNode.([]interface{})
				for index := 0; index < len(currentNode); index++ {
					nodeValue := currentNode[index]
					nextLevelNodes.Push(nodeValue)
				}
			default:
				if currentNode == value {
					return true
				}
			}
		}

		if currentNodes.Len() == 0 {
			currentNodes = nextLevelNodes
			nextLevelNodes = Vector.NewVector[interface{}](0)
		}
	}

	return false
}
