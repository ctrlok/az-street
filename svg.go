package main

import (
	"errors"
	"fmt"
	"github.com/ajstarks/svgo"
	"io"
	"net/http"
	"unicode/utf8"
)

const bigTextStyle string = "font-family:Proba Nav2 SmBd;font-size:910;font-weight:500;fill:#FFFFFF"
const littleTextStyle string = "font-family:Proba Nav2;font-size:372;font-weight:normal;fill:#FFFFFF"
const height int = 2539
const xTextBegin int = 460
const yBigText = 1420

func httpHandlerSVG(w http.ResponseWriter, r *http.Request) {
	street := Street{
		StreetNameUA: "Дзжерельна",
		StreetType:   "vulitsya",
	}
	defineStreetTypeUA(&street)
	defineStreetName(&street)
	streetSVG(street, w)
}

func streetSVG(street Street, w io.Writer) (err error) {
	canvas := svg.New(w)
	var width int
	streetLen := utf8.RuneCountInString(street.StreetNameUA)
	switch {
	case streetLen <= 8:
		width = 5197
	case streetLen > 8 && streetLen <= 12:
		width = 7440
	case streetLen > 12 && streetLen <= 17:
		width = 9321
	case streetLen > 17 && streetLen <= 21:
		width = 11681
	default:
		return errors.New("can't parse street bigger than 21 symbols")
	}
	canvas.Start(width, height, "")
	canvas.Roundrect(0, 0, width, height, 50, 50, "fill:#262741")
	canvas.Text(xTextBegin, yBigText, street.StreetNameUA, bigTextStyle)
	canvas.Text(xTextBegin, 493, street.StreetTypeUA, littleTextStyle)
	canvas.Rect(xTextBegin, 1787, width-920, 16, "fill:#FFFFFF")
	canvas.Text(xTextBegin, 2226, fmt.Sprintf("%s %s", street.StreetNameEng, street.StreetType), littleTextStyle)
	canvas.End()
	return err
}
