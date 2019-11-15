package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type record struct {
	City         string
	State        string
	Timezone     string
	StateCapital bool `json:"state_capital"`
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	f, err := ioutil.ReadFile("./us.json")
	check(err)
	records := make([]record, 0)
	err = json.Unmarshal(f, &records)
	check(err)

	type request struct {
		City  string
		State string
	}

	http.HandleFunc("/timezone", func(w http.ResponseWriter, r *http.Request) {
		req := request{}
		err := json.NewDecoder(r.Body).Decode(&req)
		check(err)
		for _, r := range records {
			if r.State == req.State && r.City == req.City {
				_, err = w.Write([]byte(fmt.Sprintf(`{"timezone":"%s"}`, r.Timezone)))
				check(err)
			}
		}

	})

	log.Fatal(http.ListenAndServe(":8888", nil))
}
