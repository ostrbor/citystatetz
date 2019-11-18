package main

import (
	"encoding/json"
	"fmt"
	"github.com/gobuffalo/packr"
	"log"
	"net/http"
	"os"
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
	box := packr.NewBox("./static")
	f, err := box.Find("./static/us.json")
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
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		for _, r := range records {
			if r.State == req.State && r.City == req.City {
				_, err = w.Write([]byte(fmt.Sprintf(`{"timezone":"%s"}`, r.Timezone)))
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(err.Error()))
				}
			}
		}
		w.WriteHeader(http.StatusNotFound)
	})

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8888"
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
