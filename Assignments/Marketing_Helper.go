package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"sort"
	"strconv"
	"time"

	"github.com/amitkhomane/GoLearning/Assignments/utils"
)

type keyVal struct {
	key string
	val int
}

func count(d []utils.Data) map[string]int {
	UserCount := make(map[string]int)
	for _, n := range d {
		user := n.User
		_, present := UserCount[user]
		if present {
			UserCount[user]++
		} else {
			UserCount[user] = 1
		}

	}
	return UserCount

}

func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

func totalBuilds(duration string, start string, end string, data []utils.Data) {
	done := make(chan bool, 1)
	if duration != "" {
		d, _ := time.ParseDuration(duration)
		startTime := time.Now().Add(-d)
		fmt.Println("startTime", startTime)
		endTime := time.Now()
		fmt.Println("endTime", endTime)

		go func() {
			count := 0
			for _, ele := range data {

				if inTimeSpan(startTime, endTime, time.Time(ele.StartedAt)) && inTimeSpan(startTime, endTime, time.Time(ele.FinishedAt)) {
					fmt.Println("ID: ", ele.ID)
					count++
				}
			}
			fmt.Printf("Total %d builds executed in given time ", count)
			done <- true

		}()
		<-done

	} else {
		startTime, _ := time.Parse("2006-01-02 15:04:05", start)
		endTime, _ := time.Parse("2006-01-02 15:04:05", end)
		fmt.Println("Start:", startTime)
		fmt.Println("end:", endTime)
		go func() {
			count := 0
			for _, ele := range data {

				if inTimeSpan(startTime, endTime, time.Time(ele.StartedAt)) && inTimeSpan(startTime, endTime, time.Time(ele.FinishedAt)) {
					fmt.Println("ID: ", ele.ID)
					count++
				}
			}
			fmt.Printf("Total %d builds executed in given time ", count)
			done <- true

		}()
		<-done

	}

}

func main() {
	fptr := flag.String("fpath", "data.csv", "path to csv file")
	cptr := flag.Int("top", 0, "list top users")
	dptr := flag.String("d", "", "time duration window")
	fromptr := flag.String("from", "", "time from")
	toptr := flag.String("to", "", "time till")

	flag.Parse()

	f := utils.Openfile(*fptr)

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println("Unable to close file")
		}
	}()

	readData := csv.NewReader(bufio.NewReader(f))

	var data []utils.Data
	var record utils.Data
	for {
		line, error := readData.Read()
		if error == io.EOF {
			break
		}
		record.ID = line[0]
		record.User = line[1]
		record.RequestedAt, _ = time.Parse(time.RFC3339, line[2])
		record.StartedAt, _ = time.Parse(time.RFC3339, line[3])
		record.FinishedAt, _ = time.Parse(time.RFC3339, line[4])
		record.Deleted, _ = strconv.ParseBool(line[5])
		record.ExitCode, _ = strconv.Atoi(line[6])
		record.Size, _ = strconv.ParseInt(line[7], 10, 32)

		data = append(data, record)

	}

	userCount := count(data)
	users := make([]string, 0, len(userCount))
	func() {

		for k := range userCount {
			users = append(users, k)
		}
	}()

	var kv []keyVal

	for k, v := range userCount {
		kv = append(kv, keyVal{k, v})

	}

	sort.Slice(kv, func(i, j int) bool {
		return kv[i].val > kv[j].val
	})
	for _, v := range kv {
		fmt.Printf("%s, %d\n", v.key, v.val)
		*cptr--
		if *cptr == 0 {
			break
		}

	}

	// Duration calculation
	if *dptr != "" || *fromptr != "" || *toptr != "" {
		totalBuilds(*dptr, *fromptr, *toptr, data)
	}
}
