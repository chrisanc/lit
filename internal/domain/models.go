package domain

import (
	"os"
)

/*
 * language.go: Definition of all the types used in the parsing
 */

// Point is used to save the position of a certain line of code (feedback usages)
type Point struct {
	Row, Column uint
}

// FunctionData - > struct made to register all the data returned by the parser and save it
type FunctionData struct {
	Name                                  string
	TotalParams, Complexity, InvalidNames uint
	StartPosition                         Point
	StartByte, EndByte, Size              uint
	Feedback                              string
}

// Directory saves up the dir name with its content
type Directory struct {
	DirName string
	Content []os.DirEntry
}

// AddInitialData Method to append the initial data into a FunctionData "object"
func (f *FunctionData) AddInitialData(name string, totalParams, startByte, endByte, size uint, startPos Point) {
	f.Name = name
	f.TotalParams = totalParams
	f.StartByte = startByte
	f.EndByte = endByte
	f.Size = size
	f.StartPosition = startPos
}

func (f *FunctionData) UpdateInvalidNames() {
	f.InvalidNames++
}

// IsTargetInRange validates the range given by another function to validate it's written on the same byte range
func (f *FunctionData) IsTargetInRange(startByte, endByte uint) bool {
	return f.StartByte <= startByte && f.EndByte >= endByte
}

// SetFunctionFeedback loops through the feedback map and sets up the right feedback
func (f *FunctionData) SetFunctionFeedback(messages map[string][]Message) {
	for key, value := range messages {
		var msg string
		for _, item := range value {
			if f.getValue(key) >= item.MinValue {
				msg = item.Message
			}
		}
		f.Feedback += msg
	}
}

// getValue is a helper function used to get the determined integer based on a key
func (f *FunctionData) getValue(key string) uint {
	dict := map[string]uint{
		"parameters": f.TotalParams,
		"complexity": f.Complexity,
		"size":       f.Size,
	}

	return dict[key]
}
