package main

import (
	"errors"
	"fmt"
	"github.com/ajstarks/svgo"
	"io"
	"regexp"
	"unicode/utf8"
)

const bigTextStyle string = "font-family:Proba Nav2 SmBd;font-size:910;font-weight:500;fill:#FFFFFF"
const littleTextStyle string = "font-family:Proba Nav2;font-size:372;font-weight:normal;fill:#FFFFFF"

func streetSVG(street Street, w io.Writer) (err error) {
	var height = 2539
	var xTextBegin = 460
	var yBigText = 1420
	var heightPhys = 215
	canvas := svg.New(w)
	var width int
	var widthPhys int
	streetLen := utf8.RuneCountInString(street.NameUA)
	switch {
	case streetLen <= 8:
		width = 5197
		widthPhys = 440
	case streetLen > 8 && streetLen <= 12:
		width = 7440
		widthPhys = 630
	case streetLen > 12 && streetLen <= 17:
		width = 9321
		widthPhys = 790
	case streetLen > 17 && streetLen <= 21:
		width = 11681
		widthPhys = 990
	default:
		return errors.New("can't parse street bigger than 21 symbols")
	}
	canvas.Startunit(widthPhys, heightPhys, "mm", fmt.Sprintf(`viewBox="0 0 %v %v"`, width, height))
	canvas.Roundrect(0, 0, width, height, 50, 50, "fill:#262741")
	canvas.Text(422, yBigText, street.NameUA, bigTextStyle)
	canvas.Text(xTextBegin, 493, street.TypeUA, littleTextStyle)
	canvas.Rect(xTextBegin, 1787, width-920, 16, "fill:#FFFFFF")
	canvas.Text(xTextBegin, 2226, fmt.Sprintf("%s %s", street.NameEng, street.Type), littleTextStyle)
	canvas.End()
	return err
}

type bigNumStruct struct {
	style string
	num   string
	x     int
}

func numSVG(street Street, w io.Writer) (err error) {
	var height = 2539
	var heightPhys = 215
	var width int
	var widthPhys int
	var yMainText = 1414
	var widthLineShort int
	var widthLineLong int
	var mainNumBig = "font-family:ProbaNav2-Medium, Proba Nav2;font-size:1590;font-weight:400;fill:#FFFFFF;"
	var mainNumLittle = "font-family:ProbaNav2-Medium, Proba Nav2;font-size:910;font-weight:500;fill:#FFFFFF;"
	var mainNumBigMiddle = fmt.Sprint(mainNumBig, "text-anchor:middle;")
	arrow := "M22.3614357,78.754834 L89.0954544,12.0208153 L77.7817459,0.707106781 L2.13162821e-14,78.4888527 L0.265981284,78.754834 L2.13162821e-14,79.0208153 L77.7817459,156.802561 L89.0954544,145.488853 L22.3614357,78.754834 Z"
	var arrowTranslate1 int
	var arrowTranslate2 int
	canvas := svg.New(w)
	numLen := utf8.RuneCountInString(street.Num)
	reNum := regexp.MustCompile("^([0-9]+)$")
	reNumAndLetter := regexp.MustCompile("([0-9]+)([^0-9])$")
	reSlash := regexp.MustCompile("/")
	bigNums := []bigNumStruct{}
	switch {
	case numLen <= 2:
		width = 2539
		widthPhys = 215
		widthLineShort = 720
		widthLineLong = 1619
		arrowTranslate1 = 1290
		arrowTranslate2 = 920
		if parsed := reNum.FindStringSubmatch(street.Num); parsed != nil {
			bigNums = append(bigNums, bigNumStruct{style: mainNumBigMiddle, num: parsed[0], x: width / 2})
		} else if parsed := reNumAndLetter.FindStringSubmatch(street.Num); parsed != nil {
			bigNums = append(bigNums, bigNumStruct{style: mainNumBig, num: parsed[0], x: 385})
			bigNums = append(bigNums, bigNumStruct{style: mainNumLittle, num: parsed[1], x: 385 + 400})
		}
	case numLen == 3 || numLen == 4 && reSlash.MatchString(street.Num):
		width = 3248
		widthPhys = 275
		widthLineShort = 1049
		arrowTranslate1 = 1776
		arrowTranslate2 = 1244
	case numLen == 4 || numLen == 5 && reSlash.MatchString(street.Num):
		width = 4015
		widthPhys = 340
		widthLineShort = 1428
		arrowTranslate1 = 2349
		arrowTranslate2 = 1623
	default:
		width = 5196
		widthPhys = 440
		widthLineShort = 2018
		arrowTranslate1 = 3237
		arrowTranslate2 = 2218
	}
	canvas.Startunit(widthPhys, heightPhys, "mm", fmt.Sprintf(`viewBox="0 0 %v %v"`, width, height))
	canvas.Roundrect(0, 0, width, height, 50, 50, "fill:#262741")
	for _, num := range bigNums {
		canvas.Text(num.x, yMainText, num.num, num.style)
	}
	canvas.Gtransform("translate(470, 1716)")
	if street.Prev != "" {
		canvas.Path(arrow, "fill:#FFFFFF;")
	}
	if street.Next != "" {
		canvas.Gtransform(fmt.Sprintf("translate(%v, 78) scale(-1, 1) translate(-%v, -78) translate(%v, 0)", arrowTranslate1, arrowTranslate1, arrowTranslate2))
		canvas.Path(arrow, "fill:#FFFFFF;")
		canvas.Gend()
	}
	if street.Next != "" && street.Prev != "" {
		canvas.Rect(20, 71, widthLineShort, 16, "fill:#FFFFFF;")
		canvas.Gtransform(fmt.Sprintf("translate(%v, 78) scale(-1, 1) translate(-%v, -78) translate(%v, 0)", arrowTranslate1, arrowTranslate1, arrowTranslate2))
		canvas.Rect(20, 71, widthLineShort, 16, "fill:#FFFFFF;")
		canvas.Gend()
	} else {
		canvas.Rect(20, 71, widthLineLong, 16, "fill:#FFFFFF;")
	}
	canvas.Gend()

	canvas.End()
	return nil
}
