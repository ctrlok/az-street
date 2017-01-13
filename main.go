package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/ctrlok/uatranslit/uatranslit"
)

// Street это струтктура, в которую мы распарсим входящий json
type Street struct {
	NameUA  string `json:"name_ua"`
	Type    string `json:"type"`
	Num     string `json:"num"`
	Prev    string `json:"prev"`
	Next    string `json:"next"`
	NameEng string
	TypeUA  string
	ID      string
}

func (s *Street) defineStreetTypeUA() {
	s.TypeUA = streetType[s.Type]
	if s.TypeUA == "" {
		log.WithField("street", s.ID).Warnf("No UA definition for type '%s'", s.Type)
	}
}

func (s *Street) defineStreetName() {
	s.NameUA = strings.Title(s.NameUA)
	s.NameEng = string(uatranslit.ReplaceUARunes([]rune(s.NameUA)))
}

func (s *Street) createID() {
	s.ID = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprint(s))))
	log.WithField("streetLong", *s).WithField("street", s.ID)
}

var streetType = map[string]string{
	"vulitsya":   "вулиця",
	"provulok":   "провулок",
	"naberezhna": "набережна",
	"prospekt":   "проспект",
	"bulvar":     "бульвар",
	"uzviz":      "узвіз",
	"tupyk":      "тупик",
}

type timer struct {
	startTime int64
	time      int64
}

func (t *timer) stopTimer() {
	t.time = t.getDiff()
}

func (t *timer) getDiff() int64 {
	return (time.Now().UnixNano() - t.startTime)
}

func startTimer() *timer {
	t := timer{}
	t.startTime = time.Now().UnixNano()
	return &t
}

// ConcLevel - постоянная, которая определяет сколько одновременных созданий изображений возможно
const ConcLevel int = 4

// tmpDirPath is a path for saving files before it was rendered. Default to os.Tempdir (/tmp at linux)
var tmpDirPath = os.Getenv("RTMP")

// tmpDirPath is a path for saving files before it was rendered. Default to os.Tempdir
var archiveDir = os.Getenv("ARCHIVE_DIR")

const reqTypeArchive = 0
const reqTypePng = 1

type inputEvent struct {
	street   Street
	reqType  int
	response chan response
}

var inChan = make(chan inputEvent)

type response struct {
	filename string
	error    error
}

func init() {
	if archiveDir == "" {
		archiveDir = os.TempDir()
	}
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
	if os.Getenv("DEBUG") != "" {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	startInscapeHandlers()
	http.HandleFunc("/zip", func(w http.ResponseWriter, r *http.Request) {
		log.Debugln("Recieve /zip request")
		archiveHandler(w, r, reqTypeArchive)
	})
	http.HandleFunc("/png", func(w http.ResponseWriter, r *http.Request) {
		log.Debugln("Recieve /png request")
		archiveHandler(w, r, reqTypePng)
	})
	// http.HandleFunc("/street.svg", svgHandlerStreet)
	// http.HandleFunc("/num.svg", svgHandlerNum)
	err := http.ListenAndServe(":3001", nil)
	if err != nil {
		log.Panic(err)
	}
}

func startInscapeHandlers() {
	for i := 0; i < ConcLevel; i++ {
		log.Infof("Start inscape handler: %s", i)
		go inscapeHandler()
	}
}

func inscapeHandler() {
	for {
		log.Debug("Wait until event")
		inEvent := <-inChan
		log.WithField("street", inEvent.street.ID).Info("Start event handling")
		var filename string
		var err error
		switch inEvent.reqType {
		case reqTypeArchive:
			log.WithField("street", inEvent.street.ID).Info("Start archive generation")
			filename, err = makeArchive(&inEvent.street)
			log.WithField("street", inEvent.street.ID).Debugf("Archive created: %v. \nerr: %v", filename, err)
		case reqTypePng:
			log.WithField("street", inEvent.street.ID).Info("Start png generation")
			filename, err = makePng(&inEvent.street)

		}
		inEvent.response <- response{
			filename: filename,
			error:    err,
		}
		log.WithField("street", inEvent.street.ID).Debug("response return from handler")
	}
}

func makePng(street *Street) (archive string, err error) {
	archive = fmt.Sprint(archiveDir, "/", street.ID)
	dir, err := ioutil.TempDir(tmpDirPath, "")
	if err != nil {
		return
	}
	// defer removeDirs(dir)
	log.WithField("street", street.ID).WithField("directory", dir).Debug("Temporary directory was created")
	t := startTimer()
	err = renderSVG(*street, dir)
	if err != nil {
		return
	}
	log.WithField("street", street.ID).WithField("directory", dir).WithField("gen_time", t.getDiff()).Info("svg_generated")
	err = renderPNG(dir)
	if err != nil {
		return
	}

	sizes := []string{"_80", "_160", "_240", ""}
	for _, t := range []string{"street", "num"} {
		for _, size := range sizes {
			name := fmt.Sprint(t, size, ".png")
			command := exec.Command("mv", fmt.Sprint(dir, "/", name), fmt.Sprint(archive, "_", name))
			err = command.Run()
			if err != nil {
				log.WithField("street", street.ID).
					WithField("directory", dir).
					WithField("cmd", command.Args).
					Debugf("Problem with move files")
				return "", err
			}
		}
	}
	return
}

func makeArchive(street *Street) (archive string, err error) {
	archive = fmt.Sprint(archiveDir, "/", street.ID, ".zip")
	_, err = os.Stat(archive)
	if err == nil {
		log.WithField("street", street.ID).Info("Archive already exist, skip...")
		return
	}
	dir, err := ioutil.TempDir(tmpDirPath, "")
	if err != nil {
		return
	}
	defer removeDirs(dir)
	log.WithField("street", street.ID).WithField("directory", dir).Debug("Temporary directory was created")
	t := startTimer()
	err = renderSVG(*street, dir)
	if err != nil {
		return
	}
	log.WithField("street", street.ID).WithField("directory", dir).WithField("gen_time", t.getDiff()).Info("svg_generated")
	err = renderPNG(dir)
	if err != nil {
		return
	}
	log.WithField("street", street.ID).WithField("directory", dir).WithField("gen_time", t.getDiff()).Info("png_generated")
	err = renderEPS(dir)
	if err != nil {
		return
	}
	log.WithField("street", street.ID).WithField("directory", dir).WithField("gen_time", t.getDiff()).Info("eps_generated")
	err = removeSVG(dir)
	if err != nil {
		return
	}
	log.WithField("street", street.ID).WithField("directory", dir).WithField("gen_time", t.getDiff()).Info("svg_removed")
	err = createArchive(dir, archive)
	if err != nil {
		return
	}
	log.WithField("street", street.ID).WithField("directory", dir).WithField("gen_time", t.getDiff()).Info("archive_created")
	return
}

func removeDirs(dir string) {
	err := os.RemoveAll(dir)
	if err != nil {
		log.WithField("directory", dir).Error(err)
	}
}

// renderSVG will create svg files in directory dir
func renderSVG(street Street, dir string) (err error) {
	log.WithField("street", street.ID).WithField("directory", dir).Debug("Start renderSVG")
	street.defineStreetTypeUA()
	street.defineStreetName()
	err = renderSVGstreet(dir, street)
	if err != nil {
		return err
	}
	err = renderSVGnum(dir, street)
	if err != nil {
		return err
	}
	return
}

func renderSVGstreet(dir string, street Street) (err error) {
	file, err := os.Create(fmt.Sprint(dir, "/street.svg"))
	if err != nil {
		return err
	}
	defer file.Close()
	err = streetSVG(street, file)
	return err
}

func renderSVGnum(dir string, street Street) (err error) {
	file, err := os.Create(fmt.Sprint(dir, "/num.svg"))
	if err != nil {
		return err
	}
	defer file.Close()
	err = numSVG(street, file)
	return
}

func renderPNG(dir string) (err error) {
	log.WithField("directory", dir).Debug("Start PNG render")
	err = makePNG("street", dir)
	if err != nil {
		return err
	}
	err = makePNG("num", dir)
	return
}

func makePNG(name, dir string) (err error) {
	svgPath := fmt.Sprint(dir, "/", name, ".svg")
	pngPath := fmt.Sprint(dir, "/", name, ".png")
	command := exec.Command("inkscape", "-z", "-T", "-e", pngPath, svgPath)
	err = command.Run()
	if err != nil {
		log.WithField("directory", dir).Debugf("Problem with PNG render, cmd: %s", command.Args)
		return err
	}

	sizes := []string{"80", "160", "240"}
	for _, size := range sizes {
		command = exec.Command("convert", pngPath, "-resize", size, fmt.Sprint(dir, "/", name, "_", size, ".png"))
		err = command.Run()
		if err != nil {
			log.WithField("directory", dir).Debugf("Problem with PNG render, cmd: %s", command.CombinedOutput)
			return err
		}
	}
	return nil
}

func renderEPS(dir string) (err error) {
	log.WithField("directory", dir).Debug("Start EPS render")
	err = makeEPS("street", dir)
	if err != nil {
		return err
	}
	err = makeEPS("num", dir)
	return err
}

func makeEPS(name, dir string) (err error) {
	svgPath := fmt.Sprint(dir, "/", name, ".svg")
	epsPath := fmt.Sprint(dir, "/", name, ".eps")
	command := exec.Command("inkscape", "-z", "-T", "-E", epsPath, svgPath)
	err = command.Run()
	return err
}

func removeSVG(dir string) (err error) {
	command := exec.Command("rm", fmt.Sprint(dir, "/street.svg"))
	err = command.Run()
	if err != nil {
		return err
	}
	command = exec.Command("rm", fmt.Sprint(dir, "/num.svg"))
	err = command.Run()
	return err
}

func createArchive(dir, archive string) (err error) {
	command := exec.Command("zip", "-r", "-j", archive, dir)
	err = command.Run()
	return err
}

func svgHandlerNum(w http.ResponseWriter, r *http.Request) {
	street, _ := decode(r.Body)
	street.createID()
	street.defineStreetTypeUA()
	street.defineStreetName()
	numSVG(street, w)
}

func svgHandlerStreet(w http.ResponseWriter, r *http.Request) {
	street, _ := decode(r.Body)
	street.createID()
	street.defineStreetTypeUA()
	street.defineStreetName()
	streetSVG(street, w)
}

func archiveHandler(w http.ResponseWriter, r *http.Request, reqType int) {
	log.Debugln("Start request decofing")
	street, err := decode(r.Body)
	if err != nil {
		http.Redirect(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	street.createID()
	log.WithField("street", street.ID).Debugln("Recieve ID")
	street.defineStreetTypeUA()
	log.WithField("street", street.ID).Debugln("Define type")
	street.defineStreetName()
	log.WithField("street", street.ID).Debugln("Define name")
	out := make(chan response)
	defer close(out)
	reqArch := inputEvent{
		street:   street,
		reqType:  0,
		response: out,
	}
	reqPng := inputEvent{
		street:   street,
		reqType:  1,
		response: out,
	}
	log.WithField("street", street.ID).Debugln("Make req")

	// Начинаем обработку события
	inChan <- reqArch
	log.WithField("street", street.ID).Debugln("Req sended")
	respArch := <-out
	log.WithField("street", street.ID).Debugln("Response recieved")
	if respArch.error != nil {
		log.WithField("street", street.ID).Error(respArch.error)
		http.Redirect(w, r, respArch.error.Error(), http.StatusInternalServerError)
	}
	inChan <- reqPng
	log.WithField("street", street.ID).Debugln("Req sended")
	respPng := <-out
	log.WithField("street", street.ID).Debugln("Response recieved")
	if respPng.error != nil {
		log.WithField("street", street.ID).Error(respPng.error)
		http.Redirect(w, r, respPng.error.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	res := `{
		"archive": "%s", 
		"street": "%s_street.png", 
		"num": "%s_num.png",
		"street_80": "%s_street_80.png",
		"num_80": "%s_num_80.png",
		"street_160": "%s_street_160.png",
		"num_160": "%s_num_160.png",
		"street_240": "%s_street_240.png",
		"num_240": "%s_num_240.png",
	}`
	n := respPng.filename
	w.Write([]byte(fmt.Sprintf(res, respArch.filename, n, n, n, n, n, n, n, n)))
}

func decode(r io.Reader) (street Street, err error) {
	decoder := json.NewDecoder(r)
	err = decoder.Decode(&street)
	return street, err
}
