package main
/*
import (
	"fmt"
	"io/ioutil"
	"flag"
)

func main(){
	fptr := flag.String("fpath", "data", "file path to csv file")
	flag.Parse()
	//data, err := ioutil.ReadFile("data.csv")
	data, err := ioutil.ReadFile(*fptr)
	if err !=nil{
		fmt.Println("Error in reading file", err)
		return
	}
	fmt.Println("Contents of the file", string(data))


} */
import (
	"fmt"
	"flag"
	"bufio"
	"os"
	"time"
	"encoding/csv"
	"io"
	"sort"
)


type keyVal struct{
	key string
	val int

} 
type Data struct {
    ID   string
	User string
	RequestedAt string
	StartedAt string
	FinishedAt string
	Deleted string
	ExitCode string
	Size string

}

func count(d [] Data)map[string]int{
	UserCount := make(map[string]int)
	for _, n := range d{
		user := n.User
		_, present := UserCount[user]
		if present{
			UserCount[user] +=1
		}else{
			UserCount[user] =1
		}

	}
	return UserCount

}



// func total_builds(start, end, duration tim.Time){

// 	start := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
// 	end := time.Date(2000, 1, 1, 0, 16, 0, 0, time.UTC)

// 	difference := end.Sub(start)
// 	min, _ := time.ParseDuration("15m")
// 	fmt.Printf("difference = %v\n", difference)
// 	if min >= difference{
// 		fmt.Printf("diff nad min are same = %v\n", difference)
// 	}



// }

//func total_builds(start, end time.Time)
 
// type Data struct{
// 	ID 
// }

func main(){
	fptr := flag.String("fpath", "data.csv", "path to csv file")
	cptr := flag.Int("top", 1, "path to csv file")
	flag.Parse()
	

	f, err := os.Open(*fptr)
	if err !=nil{
		fmt.Println("Unable to read the file", err)
		return
	}

	// defer func(){
	// 	if err = f.Close();err !=nil{
	// 		fmt.Println("Unable to close the file")
	// 	}

	// }()

	// s := bufio.NewScanner(f)
	// //min, _ := time.ParseDuration("15m")
	// s1:="2018-11-29T08:30:20-05:00"
	// t, _ := time.Parse(time.RFC3339, s1)
	
	// fmt.Println("duration", t)
	// bb()
	// for s.Scan(){
	// 	//fmt.Println(s.Text())

	// }
	reader := csv.NewReader(bufio.NewReader(f))
	var data []Data 
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		}

		data = append(data, Data{
			ID : line[0],
			User :line[1],
			RequestedAt: line[2],
			StartedAt  : line[3],
			FinishedAt : line[4],
			Deleted  :line[5],
			ExitCode :line[6],
			Size     :line[7],
		})
		
	}
	
	fmt.Println(time.Parse(time.RFC3339, data[0].StartedAt))

	user_count:= count(data)
	fmt.Println(user_count)
	users := make([]string, 0, len(user_count))
	func(){
		
		for k:= range user_count{
			users = append(users, k)
		} 
	}()

	var kv []keyVal

	for k, v := range user_count{
		kv = append(kv, keyVal{k, v})

	}

	sort.Slice(kv, func(i, j int) bool{
		return kv[i].val > kv[j].val
	})
	for _, v :=range kv{
		fmt.Printf("%s, %d\n", v.key, v.val)
		*cptr --
		if *cptr == 0{break}

	}



}