// Package utils provides utilirize for market assignments.
package utils

import (
	"fmt"
	"os"
	"time"
)

//KeyVal struct for storing key and value of map
type KeyVal struct {
	key string
	val int
}

//Data proves struct for the format
type Data struct {
	ID          string
	User        string
	RequestedAt time.Time
	StartedAt   time.Time
	FinishedAt  time.Time
	Deleted     bool
	ExitCode    int
	Size        int64
}

//Openfile from the provided file name as input.
func Openfile(name string) *os.File {
	file, err := os.Open(name)
	if err != nil {
		fmt.Println("Unable to open file", name)
	}
	return file

}
