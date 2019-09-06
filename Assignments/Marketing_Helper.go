package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"time"
)

type keyVal struct {
	key string
	val int
}
type Data struct {
	ID          string
	User        string
	RequestedAt string
	StartedAt   string
	FinishedAt  string
	Deleted     string
	ExitCode    string
	Size        string
}

func count(d []Data) map[string]int {
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

func totalBuilds(duration string, start string, end string, data []Data) {
	fmt.Println("In side totalBuilds", duration)
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

				startedat, _ := time.Parse(time.RFC3339, ele.StartedAt)
				endat, _ := time.Parse(time.RFC3339, ele.FinishedAt)
				if inTimeSpan(startTime, endTime, time.Time(startedat)) && inTimeSpan(startTime, endTime, time.Time(endat)) {
					fmt.Println("ID: ", ele.ID)
					count++
				} else {
					fmt.Println("NOT")

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

				startedat, _ := time.Parse(time.RFC3339, ele.StartedAt)
				endat, _ := time.Parse(time.RFC3339, ele.FinishedAt)
				if inTimeSpan(startTime, endTime, time.Time(startedat)) && inTimeSpan(startTime, endTime, time.Time(endat)) {
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

	f, err := os.Open(*fptr)
	if err != nil {
		fmt.Println("Unable to read the file", err)
		return
	}

	defer func() {
		if err = f.Close(); err != nil {
			fmt.Println("Unable to close the file")
		}

	}()

	reader := csv.NewReader(bufio.NewReader(f))
	var data []Data
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		}

		data = append(data, Data{
			ID:          line[0],
			User:        line[1],
			RequestedAt: line[2],
			StartedAt:   line[3],
			FinishedAt:  line[4],
			Deleted:     line[5],
			ExitCode:    line[6],
			Size:        line[7],
		})

	}

	fmt.Println(time.Parse(time.RFC3339, data[0].StartedAt))

	userCount := count(data)
	fmt.Println(userCount)
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
