package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
)

var scriptPath string
var statusMap map[int64]*Status = make(map[int64]*Status)
var currId int64

type Status struct {
	Started     time.Time
	Ended       time.Time
	IsRunning   bool
	ErrorString string
	Err         error
}

func main() {
	port := flag.String("port", "8080", "http port")
	readScriptPath()

	setupServer(":" + *port)
}

func readScriptPath() {
	flag.StringVar(&scriptPath, "scriptPath", "UNDEFINED", "path to the script")
	flag.Parse()
	if scriptPath == "UNDEFINED" {
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func setupServer(port string) {
	r := mux.NewRouter()

	r.HandleFunc("/start", startScript)
	r.HandleFunc("/startAndWait", startScriptAndWait)
	r.HandleFunc("/status/{id}", runStatus)

	http.ListenAndServe(port, r)
}

func startScript(w http.ResponseWriter, r *http.Request) {
	var id = atomic.AddInt64(&currId, 1)

	go runScript(id)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, id)
}

func startScriptAndWait(w http.ResponseWriter, r *http.Request) {
	var id = atomic.AddInt64(&currId, 1)

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
	if stat != nil {
		data, _ := json.Marshal(stat)

		w.WriteHeader(http.StatusOK)
		w.Write(data)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func runScript(id int64) (*Status, error) {
	stat := &Status{
		Started:   time.Now(),
		IsRunning: true,
		Err:       nil,
	}
	statusMap[id] = stat

	log.Printf("START [%d] %s\n", id, scriptPath)
	err := runFile(id)
	if err != nil {
		log.Printf("END   [%d] %s %s\n", id, scriptPath, err.Error())
	} else {
		log.Printf("END   [%d] %s\n", id, scriptPath)
	}
	updateStat(stat, err)

	return stat, err
}

func updateStat(stat *Status, err error) {
	stat.Err = err
	if err != nil {
		stat.ErrorString = err.Error()
	}
	stat.IsRunning = false
	stat.Ended = time.Now()
}

func runFile(id int64) error {
	logFile := strconv.FormatInt(id, 10) + ".log"
	f, _ := os.Create(logFile)
	defer f.Close()

	cmd := exec.Command(scriptPath)

	reader, _ := cmd.StdoutPipe()
	err := cmd.Start()
	if err != nil {
		return err
	}
	bytes, _ := ioutil.ReadAll(reader)
	f.Write(bytes)

	return cmd.Wait()
}
