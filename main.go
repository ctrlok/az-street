package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/ctrlok/uatranslit/uatranslit"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
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
	out     chan string
	err     chan error
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

var inChan chan Street

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
	http.HandleFunc("/", httpHandler)
	http.HandleFunc("/street.svg", svgHandlerStreet)
	http.HandleFunc("/num.svg", svgHandlerNum)
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
	for true {
		street := <-inChan
		log.WithField("street", street.ID).Info("Start archive generation")
		archive, err := makeArchive(&street)
		if err != nil {
			street.err <- err
			break
		}
		street.out <- archive
	}
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
	if err != nil {
		return err
	}
	return nil
}

func renderSVGnum(dir string, street Street) (err error) {
	file, err := os.Create(fmt.Sprint(dir, "/num.svg"))
	if err != nil {
		return err
	}
	defer file.Close()
	err = numSVG(street, file)
	if err != nil {
		return err
	}
	return nil
}

func renderPNG(dir string) (err error) {
	log.WithField("directory", dir).Debug("Start PNG render")
	err = makePNG("street", dir)
	if err != nil {
		return err
	}
	err = makePNG("num", dir)
	if err != nil {
		return err
	}
	return nil
}

func makePNG(name, dir string) (err error) {
	svgPath := fmt.Sprint(dir, "/", name, ".svg")
	pngPath := fmt.Sprint(dir, "/", name, ".png")
	cmd := exec.Command("inkscape", "-z", "-T", "-e", pngPath, svgPath)
	err = cmd.Run()
	if err != nil {
		return err
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
	if err != nil {
		return err
	}
	return nil
}

func makeEPS(name, dir string) (err error) {
	svgPath := fmt.Sprint(dir, "/", name, ".svg")
	epsPath := fmt.Sprint(dir, "/", name, ".eps")
	cmd := exec.Command("inkscape", "-z", "-T", "-E", epsPath, svgPath)
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func removeSVG(dir string) (err error) {
	cmd := exec.Command("rm", fmt.Sprint(dir, "/street.svg"))
	err = cmd.Run()
	if err != nil {
		return err
	}
	cmd = exec.Command("rm", fmt.Sprint(dir, "/num.svg"))
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func createArchive(dir, archive string) (err error) {
	cmd := exec.Command("zip", "-r", "-j", archive, dir)
	err = cmd.Run()
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

func httpHandler(w http.ResponseWriter, r *http.Request) {
	street, err := decode(r.Body)
	street.createID()
	street.defineStreetTypeUA()
	street.defineStreetName()
	if err != nil {
		http.Redirect(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	var out chan string
	defer close(out)
	var outErr chan error
	defer close(outErr)
	street.out = out
	street.err = outErr

	// Начинаем обработку события
	inChan <- street
	// TODO: add timer and err
	select {
	case path := <-out:
		// TODO: add url
		http.Redirect(w, r, path, http.StatusFound)
	case err := <-outErr:
		http.Redirect(w, r, err.Error(), http.StatusInternalServerError)
	}
}

func decode(r io.Reader) (street Street, err error) {
	decoder := json.NewDecoder(r)
	err = decoder.Decode(&street)
	return street, err
}
