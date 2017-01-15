package main

import (
	"crypto/md5"
	"fmt"
	"strings"

	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/ctrlok/uatranslit/uatranslit"
)

type Street struct {
	Name     string `json:"name"`
	Left     string `json:"left"`
	Right    string `json:"right"`
	Class    string `json:"class"`
	Num      string `json:"num"`
	ID       string
	NameEng  string
	ClassEng string
}

func (s *Street) defineStreetClassEng() (err error) {
	s.ClassEng = streetType[s.Class]
	if s.ClassEng == "" {
		log.WithField("street", s.ID).Errorf("No UA definition for type '%s'", s.Class)
		return errors.New("Problem with defining street type")
	}
	log.WithField("street", s.ID).Debug("Class finded")
	return
}

func (s *Street) defineStreetName() {
	s.Name = strings.Title(s.Name)
	s.NameEng = string(uatranslit.ReplaceUARunes([]rune(s.Name)))
	log.WithField("street", s.ID).Debug("Create eng name for street")
}

func (s *Street) createID() {
	s.ID = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprint(s))))
	log.WithField("street", s.ID).Debug("ID generated")
}

func (s *Street) fill() (err error) {
	s.createID()
	s.defineStreetName()
	err = s.defineStreetClassEng()
	return
}

var streetType = map[string]string{
	"вулиця":    "vulitsya",
	"провулок":  "provulok",
	"набережна": "naberezhna",
	"проспект":  "prospekt",
	"бульвар":   "bulvar",
	"узвіз":     "uzviz",
	"тупик":     "tupyk",
}
