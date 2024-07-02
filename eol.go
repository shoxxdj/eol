package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"encoding/json"
	"os"
	"github.com/olekukonko/tablewriter"
	"strings"
)

type ProductDetail struct {
	Cycle  string  `json:"cycle"`
	Lts bool `json:"lts"` 
	ReleaseDate string `json:"releaseDate"` 
	Support string `json:"support"` 
	Eol interface{} `json:"eol"` 
	Latest string `json:"latest"` 
	Link string `json:"link"` 
}

func main() {
	productFlag := flag.String("p", "", "Product")
	cycleFlag := flag.String("c", "", "Cycle")
	formatFlag := flag.String("f","inline","Ouptut Format : table,inline(default)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "EOL.Date : a binary to fetch the API of endoflife.date (v1.0.1)\n")

		flag.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(os.Stderr, "\t-%v: %v\n", f.Name, f.Usage) // f.Name, f.Value
		})
	}

	// Parse the command-line flags.
	flag.Parse()

	// Read products from stdin and add them to the slice.
	var input string 
	var product string
	var cycle string

	cycle = *cycleFlag

	if len(*productFlag) == 0 {
		for {
			_, err := fmt.Scanln(&input)
			if err != nil {
				break
			}
			if input != "" {
				if(strings.Contains(input,":")){
					subStrings := strings.Split(input,":")
					product = subStrings[0]
					cycle=subStrings[1]
				}else{
					product = input
				}
			}
		}
	}else{
		product = *productFlag
	}

	// Should handle the stdin values "product:cycle" 

	
	var resp *http.Response
	var err error
	var items []ProductDetail

	if(cycle==""){
		resp, err = http.Get("https://endoflife.date/api/"+product+".json")
	}else{
		resp, err = http.Get("https://endoflife.date/api/"+product+"/"+cycle+".json")
	}

	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	
	if(cycle==""){
		if err := decoder.Decode(&items); err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}
	}else{
		var item ProductDetail
		if err := decoder.Decode(&item); err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}
		items = []ProductDetail{item}
	}

	//fmt.Println(*items)

	if(*formatFlag=="table"){
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Cycle", "LTS", "Release Date","Support","EOL","Latest","Link"})
		for _, item := range items {
			var eol string
			if str, ok := item.Eol.(string); ok {
				// 'Eol' is a string, so we can use it directly.
				eol = str
			} else if b, ok := item.Eol.(bool); ok {
				// 'Eol' is a boolean, so we need to convert it to a string.
				eol = fmt.Sprintf("%t", b)
			}
			if(cycle!=""){
				item.Cycle=cycle; //display bug
			}
			table.Append([]string{item.Cycle, strconv.FormatBool(item.Lts), item.ReleaseDate,item.Support,eol,item.Latest,item.Link})
		}
		table.Render()
	} else if *formatFlag=="inline"{
		for _, item:= range items{
			var eol string
			if str, ok := item.Eol.(string); ok {
				// 'Eol' is a string, so we can use it directly.
				eol = str
			} else if b, ok := item.Eol.(bool); ok {
				// 'Eol' is a boolean, so we need to convert it to a string.
				eol = fmt.Sprintf("%t", b)
			}
			if(cycle!=""){
				item.Cycle=cycle; //display bug
			}
			fmt.Printf("Cycle : %s, Lts : %s, ReleaseDate : %s, Support : %s, Eol : %s, Latest : %s, Link: %s\n",item.Cycle, strconv.FormatBool(item.Lts), item.ReleaseDate,item.Support,eol,item.Latest,item.Link)
		}
	} else if *formatFlag=="json"{

	}
}

func generateRandomValue() string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomValue := make([]byte, 10)
	for i := range randomValue {
		randomValue[i] = charset[rand.Intn(len(charset))]
	}
	return string(randomValue)
}
