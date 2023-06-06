package DepthFirstSearch

import (
	"reflect"
)

func Run(node interface{}, value interface{}) bool {
	if node == nil { // Has to do this check since `reflect.TypeOf(nil) == nil`
		return node == value
	}
	switch reflect.TypeOf(node).Kind() {
	case reflect.Map:
		node := node.(map[string]interface{})
		for nodeKey, nodeValue := range node {
			if nodeKey == value || Run(nodeValue, value) {
				return true
			}
		}
	case reflect.Slice:
		/*node := reflect.ValueOf(node)
		for index := 0; index < node.Len(); index++ {
			nodeValue := node.Index(index)
			if Run(nodeValue, value) {
				return true
			}
		}*/
		/*node := node.([]interface{})
		for nodeValue := range node {
			if Run(nodeValue, value) {
				return true
			}
		}*/
		node := node.([]interface{})
		for index := 0; index < len(node); index++ {
			nodeValue := node[index]
			if Run(nodeValue, value) {
				return true
			}
		}
	default:
		if node == value {
			return true
		}
	}
	return false
}
