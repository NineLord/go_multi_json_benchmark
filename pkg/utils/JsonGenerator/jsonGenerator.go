package JsonGenerator

import (
	"fmt"
	"github.com/NineLord/go_multi_json_benchmark/pkg/utils"
	"github.com/NineLord/go_multi_json_benchmark/pkg/utils/JsonType"
	"github.com/NineLord/go_multi_json_benchmark/pkg/utils/Randomizer"
	"reflect"
	"strings"
)

type jsonGenerator struct {
	characterPoll    []rune
	numberOfLetters  uint
	depth            uint
	numberOfChildren uint
}

func newJsonGenerator(characterPoll string, numberOfLetters uint, depth uint, numberOfChildren uint) *jsonGenerator {
	result := makeJsonGenerator(characterPoll, numberOfLetters, depth, numberOfChildren)
	return &result
}

func makeJsonGenerator(characterPoll string, numberOfLetters uint, depth uint, numberOfChildren uint) jsonGenerator {
	return jsonGenerator{
		characterPoll:    []rune(characterPoll),
		numberOfLetters:  numberOfLetters,
		depth:            depth,
		numberOfChildren: numberOfChildren,
	}
}

func GenerateJson(characterPoll string, numberOfLetters uint, depth uint, numberOfChildren uint) (map[string]interface{}, error) {
	return newJsonGenerator(characterPoll, numberOfLetters, depth, numberOfChildren).generateFullTree()
}

func (jsonGenerator *jsonGenerator) generateFullTree() (map[string]interface{}, error) {
	result := make(map[string]interface{})

	currentNodes := utils.NewVector2[interface{}](1)
	currentNodes.Push(result)
	nextLevelNodes := utils.NewVector2[interface{}](0)
	lastLevel := jsonGenerator.depth - 1

	for level := uint(0); level < jsonGenerator.depth; level++ {

		for currentNodes.Len() != 0 {
			currentNode := currentNodes.Pop()

			switch currentNodeType := reflect.TypeOf(currentNode); {
			case currentNodeType.Kind() == reflect.Map:
				currentNode := currentNode.(map[string]interface{})
				for _nodeCount := uint(0); _nodeCount < jsonGenerator.numberOfChildren; _nodeCount++ {
					var childNodeValue interface{}
					if level == lastLevel {
						childNodeValue = JsonType.GetRandomLeafJson()
					} else {
						childNodeValue = JsonType.GetRandomNoneLeafJson(jsonGenerator.numberOfChildren)
					}
					currentNode[jsonGenerator.getRandomNodeName()] = childNodeValue
					nextLevelNodes.Push(childNodeValue)
				}
			case strings.HasPrefix(currentNodeType.String(), "*utils.Vector"):
				currentNode, _ := currentNode.(*utils.Vector[interface{}])
				for _nodeCount := uint(0); _nodeCount < jsonGenerator.numberOfChildren; _nodeCount++ {
					var childNodeValue interface{}
					if level == lastLevel {
						childNodeValue = JsonType.GetRandomLeafJson()
					} else {
						childNodeValue = JsonType.GetRandomNoneLeafJson(jsonGenerator.numberOfChildren)
					}
					currentNode.Push(childNodeValue)
					nextLevelNodes.Push(childNodeValue)
				}
			default:
				panic(fmt.Sprintf("generateFullTree unknown JSON type: %s", currentNodeType.String()))
			}
		}

		currentNodes = nextLevelNodes
		nextLevelNodes = utils.NewVector2[interface{}](0)
	}

	return result, nil
}

// #region Helper methods

func (jsonGenerator *jsonGenerator) getRandomNodeCharacter() *rune {
	index := Randomizer.GetRandomIndexFromArray(jsonGenerator.characterPoll)
	char := jsonGenerator.characterPoll[index]
	return &char
}

func (jsonGenerator *jsonGenerator) getRandomNodeName() string {
	var stringBuilder strings.Builder
	for count := uint(0); count < jsonGenerator.numberOfLetters; count++ {
		stringBuilder.WriteRune(*jsonGenerator.getRandomNodeCharacter())
	}
	return stringBuilder.String()
}

// #endregion
