package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
)

var scriptPath string
var statusMap map[int]*Status = make(map[int]*Status)

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
	http.HandleFunc("/start", startScript)
	http.HandleFunc("/startAndWait", startScriptAndWait)
	http.HandleFunc("/status", runStatus)
	http.ListenAndServe(port, nil)
}

func startScript(w http.ResponseWriter, r *http.Request) {
	id := 1

	go runScript(id)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "")
}

func startScriptAndWait(w http.ResponseWriter, r *http.Request) {
	id := 99

	stat, err := runScript(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
	} else {
		data, _ := json.Marshal(stat)
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func runStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "")
}

func runScript(id int) (*Status, error) {
	stat := &Status{isRunning: true,
		err: nil,
	}
	statusMap[id] = stat

	cmd := exec.Command("sleep", "5")
	err := cmd.Run()

	stat.err = err
	stat.isRunning = false

	return stat, err
}
