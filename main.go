package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
)

var concurencyLevel, _ = strconv.Atoi(os.Getenv("CONCURENCY"))
var archiveDir = os.Getenv("DEST")
var tmpDir = os.Getenv("RTMP")

var globalChan = make(chan inEvent)

var inEventClass = struct {
	archiveAndPng int
}{
	archiveAndPng: 0,
}

type inEvent struct {
	street  Street
	class   int
	outChan chan response
}

type response struct {
	files []string
	err   error
}

func init() {
	if concurencyLevel == 0 {
		log.Debug("Can't read env variable CONCURENCY... Will set to 4")
		concurencyLevel = 4
	}
	if tmpDir == "" {
		tmpDir = os.TempDir()
		log.Debugf("Can't read env variable RTMP... Will set to %s", tmpDir)
	}
	if archiveDir == "" {
		archiveDir = os.TempDir()
		log.Debugf("Can't read env variable DEST... Will set to %s", archiveDir)
	}
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
	if os.Getenv("DEBUG") != "" {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	log.Info("Starting server")
	log.Info("Setting up variables")
	log.Infof("CONCURENCY = %s", concurencyLevel)
	go startLoop()
	http.HandleFunc("/", handlerGenerateArchiveAndPng)
	err := http.ListenAndServe(":3001", nil)
	if err != nil {
		log.Panic(err)
	}
}

func startLoop() {
	for i := 0; i < concurencyLevel; i++ {
		startWorker(i)
	}
}

func startWorker(i int) {
	for {
		var filenames []string
		var err error
		log := log.WithField("worker", i)
		log.Debug("Worker wait for recieving event event")
		event := <-globalChan
		log.Debugf("Recieve event: %#s", event)
		log = log.WithField("street", event.street.ID)
		switch event.class {
		case inEventClass.archiveAndPng:
			log = log.WithField("class", "archive_and_png")
			log.Debug("Start processing...")
			t := startTimer()
			filenames, err = generateArchiveAndPng(log, event.street)
			log.WithField("operation_time", t.stop).Info("Finish event processing")
		}
		log.Debug("Send response to web handler")
		event.outChan <- response{files: filenames, err: err}
	}
}

func handlerGenerateArchiveAndPng(w http.ResponseWriter, r *http.Request) {
	log := log.WithField("func", "handlerGenerateArchiveAndPng")
	log.Debug("Recieve http request")
	t := startTimer()
	street := Street{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&street)
	if err != nil {
		log.Errorf("Problem decoding json: %s", err)
		http.Redirect(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	log.WithField("operation_time", t.diff()).Debug("json unmarshall")

	err = street.fill()
	if err != nil {
		log.Errorf("Problem fill string: %s", err)
		http.Redirect(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	log = log.WithField("street", street.ID)

	out := make(chan response)
	defer close(out)
	log.Debug("Sending event to processing...")
	tWaitForProcessing := startTimer()
	globalChan <- inEvent{
		street:  street,
		class:   inEventClass.archiveAndPng,
		outChan: out,
	}
	log.WithField("operation_time", tWaitForProcessing.diff()).Debug("Event was send to processing")
	tProcessing := startTimer()
	result := <-out
	log.WithField("operation_time", tProcessing.diff()).Debug("Finish processing")
	if result.err != nil {
		log.Errorf("Problem processing: %s", result.err)
		http.Redirect(w, r, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	res := `{
		"archive": "%s", 
		"street_80": "%s",
		"num_80": "%s",
		"street_160": "%s",
		"num_160": "%s",
		"street_240": "%s",
		"num_240": "%s",
	}`
	tSendToUser := startTimer()
	w.Write([]byte(fmt.Sprintf(res, result.files[0], result.files[1], result.files[2], result.files[3], result.files[4], result.files[5], result.files[6])))
	log.WithField("operation_time", tSendToUser.diff()).Debug("Response was sended to user")
	log.WithField("operation_time", t.diff()).Debug("Finish all request")

}
