package App

import "fmt"

type Cleaner func()
type CleanessStack []func()

var fullStackCleaning CleanessStack

func RegisterCleaner(cleaner Cleaner) {
	if fullStackCleaning == nil {
		fullStackCleaning = CleanessStack{}
	}
	fullStackCleaning = append([]func(){cleaner}, fullStackCleaning...)
}

func Cleanup()  {
	Log("Cleaning up")
	for _, function := range fullStackCleaning {
		function()
	}
	fmt.Println("-- All cleaners done -- ")
}