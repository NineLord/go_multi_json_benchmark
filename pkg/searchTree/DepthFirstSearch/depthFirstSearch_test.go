package DepthFirstSearch

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary
var messyJson map[string]interface{}

func setup() error {
	return json.Unmarshal([]byte(`{
            "a": {
                "b": [
                    0,
                    0.5,
                    "shimi"
                ],
                "c": [
                    null
                ]
            },
            "d": [
                [
                    1,
                    "hey"
                ],
                [
                    "lol",
                    "lol"
                ]
            ],
            "e": {
                "f": {
                    "g": 2
                },
                "h": [
                    3,
                    true
                ]
            }
        }`), &messyJson)
}

func TestMain(m *testing.M) {
	if err := setup(); err != nil {
		panic(err)
	}
	code := m.Run()
	// shutdown()
	os.Exit(code)
}

func expectToFind(assertions *assert.Assertions, valueToFind interface{}) {
	assertions.True(Run(messyJson, valueToFind), "Expected to find: %+v", valueToFind)
}

func expectNotToFind(assertions *assert.Assertions, valueToNotFind interface{}) {
	assertions.False(Run(messyJson, valueToNotFind), "Expected to not find %+v", valueToNotFind)
}

func TestRunShouldFind(t *testing.T) {
	assertions := assert.New(t)

	for letter := 'a'; letter <= 'h'; letter++ {
		expectToFind(assertions, string(letter))
	}

	for number := 0; number <= 3; number++ {
		expectToFind(assertions, float64(number))
	}

	expectToFind(assertions, true)
	expectToFind(assertions, 0.5)
	expectToFind(assertions, "shimi")
	expectToFind(assertions, "hey")
	expectToFind(assertions, "lol")
}

func TestRunShouldNotFind(t *testing.T) {
	assertions := assert.New(t)

	expectNotToFind(assertions, false)
	expectNotToFind(assertions, float64(4))
	expectNotToFind(assertions, 0)
	expectNotToFind(assertions, 1.5)
	expectNotToFind(assertions, string('i'))
	expectNotToFind(assertions, "Hello")
}
