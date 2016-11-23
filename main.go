package main

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/ctrlok/uatranslit/uatranslit"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"
	"unicode/utf8"
)

// Street это струтктура, в которую мы распарсим входящий json
type Street struct {
	StreetNameUA        string `json:"street_name_ua"`
	StreetType          string `json:"street_type"`
	StreetNum           string `json:"street_num"`
	StreetPositionFirst bool   `json:"street_position_first"`
	StreetPositionLast  bool   `json:"street_position_last"`
	StreetNameEng       string
	StreetTypeUA        string
	ID                  string
	out                 chan string
	err                 chan error
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
	dir, err := ioutil.TempDir(tmpDirPath, "archive")
	if err != nil {
		return
	}
	defer os.RemoveAll(dir)
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
	err = renderPDF(dir)
	if err != nil {
		return
	}
	log.WithField("street", street.ID).WithField("directory", dir).WithField("gen_time", t.getDiff()).Info("pdf_generated")
	err = removeSVG(dir)
	if err != nil {
		return
	}
	log.WithField("street", street.ID).WithField("directory", dir).WithField("gen_time", t.getDiff()).Info("svg_removed")
	archive, err = createArchive(dir, street.ID)
	if err != nil {
		return
	}
	log.WithField("street", street.ID).WithField("directory", dir).WithField("gen_time", t.getDiff()).Info("archive_created")
	return
}

// renderSVG will create svg files in directory dir
func renderSVG(street Street, dir string) (err error) {
	log.WithField("street", street.ID).WithField("directory", dir).Debug("Start renderSVG")
	err = renderSVGstreet(street, dir)
	if err != nil {
		log.WithField("street", street.ID).WithField("directory", dir).Errorf("Street template render problem: %s", err)
		return err
	}
	return nil
}

// renderSVGstreet will create street.svg file
func renderSVGstreet(street Street, dir string) (err error) {
	log.WithField("street", street.ID).WithField("directory", dir).Debug("Start renderSVG for street")
	defineStreetTypeUA(&street)
	defineStreetName(&street)
	log.WithField("street", street.ID).WithField("directory", dir).Debug("Start SVG template generation for street")
	templ, err := selectStreetSVGtemplate(street)
	if err != nil {
		return err
	}
	log.WithField("street", street.ID).WithField("directory", dir).Debug("Start render SVG template for street")
	err = renderStreetSVGtemplate(dir, street, templ)
	if err != nil {
		return err
	}
	return nil
}

func defineStreetTypeUA(street *Street) {
	street.StreetTypeUA = streetType[street.StreetType]
	if street.StreetTypeUA == "" {
		log.WithField("street", street.ID).Warnf("No UA definition for type '%s'", street.StreetType)
	}
}

func defineStreetName(street *Street) {
	street.StreetNameUA = strings.Title(street.StreetNameUA)
	street.StreetNameEng = string(uatranslit.ReplaceUARunes([]rune(street.StreetNameUA)))
}

func selectStreetSVGtemplate(street Street) (t *template.Template, err error) {
	var templ string
	streetLen := utf8.RuneCountInString(street.StreetNameUA)
	log.WithField("street", street.ID).Debugf("Street has %v runes", streetLen)
	switch {
	case streetLen <= 8:
		templ = templStreet1
	case streetLen > 8 && streetLen <= 12:
		templ = templStreet2
	case streetLen > 12 && streetLen <= 17:
		templ = templStreet3
	case streetLen > 17 && streetLen <= 21:
		templ = templStreet4
	default:
		return t, errors.New("can't parse street bigger than 21 symbols")
	}
	t, err = template.New("").Parse(templ)
	if err != nil {
		return t, err
	}
	return t, nil
}

func renderStreetSVGtemplate(dir string, street Street, t *template.Template) (err error) {
	file, err := os.Create(fmt.Sprint(dir, "/street.svg"))
	defer file.Close()
	if err != nil {
		return err
	}
	err = t.Execute(file, street)
	if err != nil {
		return err
	}
	return nil
}

func renderPNG(dir string) (err error) {
	log.WithField("directory", dir).Debug("Start PNG render")
	err = renderPNGstreet(dir)
	if err != nil {
		return err
	}
	return nil
}

func renderPNGstreet(dir string) (err error) {
	svgPath := fmt.Sprint(dir, "/street.svg")
	pngPath := fmt.Sprint(dir, "/street.png")
	cmd := exec.Command("inkscape", "-z", "-T", "-e", pngPath, svgPath)
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func renderEPS(dir string) (err error) {
	log.WithField("directory", dir).Debug("Start EPS render")
	err = renderEPSstreet(dir)
	if err != nil {
		return err
	}
	return nil
}

func renderEPSstreet(dir string) (err error) {
	svgPath := fmt.Sprint(dir, "/street.svg")
	epsPath := fmt.Sprint(dir, "/street.eps")
	cmd := exec.Command("inkscape", "-z", "-T", "-E", epsPath, svgPath)
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func renderPDF(dir string) (err error) {
	log.WithField("directory", dir).Debug("Start PDF render")
	err = renderPDFstreet(dir)
	if err != nil {
		return err
	}
	return nil
}

func renderPDFstreet(dir string) (err error) {
	svgPath := fmt.Sprint(dir, "/street.svg")
	pdfPath := fmt.Sprint(dir, "/street.pdf")
	cmd := exec.Command("inkscape", "-z", "-T", "-A", pdfPath, svgPath)
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
	return nil
}

func createArchive(dir string, id string) (archive string, err error) {
	archive = fmt.Sprint(archiveDir, "/", id, ".zip")
	cmd := exec.Command("zip", "-r", archive, dir)
	err = cmd.Run()
	if err != nil {
		return archive, err
	}
	return archive, nil
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	street, err := decode(r.Body)
	street.createID()
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
