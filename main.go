package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/gorilla/mux"
)

var scriptPath string
var statusMap map[int64]*Status = make(map[int64]*Status)

type Status struct {
	// started int64
	// ended int64
	isRunning bool
	err       error
}

func main() {
	scriptPath = ""

	setupServer(":8080")
}

func setupServer(port string) {
	r := mux.NewRouter()

	r.HandleFunc("/start", startScript)
	r.HandleFunc("/startAndWait", startScriptAndWait)
	r.HandleFunc("/status/{id}", runStatus)

	http.ListenAndServe(port, r)
}

func startScript(w http.ResponseWriter, r *http.Request) {
	var id int64 = 1

	go runScript(id)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, id)
}

func startScriptAndWait(w http.ResponseWriter, r *http.Request) {
	var id int64 = 99

	stat, err := runScript(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	data, _ := json.Marshal(stat)
	w.Write(data)
}

func runStatus(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.ParseInt(idStr, 10, 64)
	stat := statusMap[id]
	data, _ := json.Marshal(stat)

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func runScript(id int64) (*Status, error) {
	stat := &Status{isRunning: true,
		err: nil,
	}
	statusMap[id] = stat

	log.Printf("START [%d] %s\n", id, "sleep 5")
	cmd := exec.Command("sleep", "5")
	err := cmd.Run()
	if err != nil {
		log.Printf("END [%d] %s %s\n", id, "sleep 5", err.Error())
	} else {
		log.Printf("END [%d] %s\n", id, "sleep 5")
	}

	stat.err = err
	stat.isRunning = false

	return stat, err
}
