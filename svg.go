package main

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"unicode/utf8"

	"github.com/ajstarks/svgo"
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

const mainNumBig = "font-family:Proba Nav2 Md;font-size:1590;font-weight:400;fill:#FFFFFF;"
const mainNumLittle = "font-family:Proba Nav2 SmBd;font-size:910;font-weight:500;fill:#FFFFFF;"
const mainNumLitLit = "font-family:Proba Nav2;font-size:600;font-weight:bold;fill:#FFFFFF;"

func numSVG(street Street, w io.Writer) (err error) {
	var height = 2539
	var heightPhys = 215
	var width int
	var widthPhys int
	var yMainText = 1414
	var widthLineShort int
	var widthLineLong int
	var arrowNumLeft = "font-family:Proba Nav2;font-size:372;fill:#FFFFFF;text-anchor:start"
	var arrowNumRight = "font-family:Proba Nav2;font-size:372;fill:#FFFFFF;text-anchor:end"
	arrow := "M22.3614357,78.754834 L89.0954544,12.0208153 L77.7817459,0.707106781 L2.13162821e-14,78.4888527 L0.265981284,78.754834 L2.13162821e-14,79.0208153 L77.7817459,156.802561 L89.0954544,145.488853 L22.3614357,78.754834 Z"
	var arrowTranslate1 int
	var arrowTranslate2 int
	canvas := svg.New(w)
	numLen := utf8.RuneCountInString(street.Num)
	matchedString := regexp.MustCompile("^([-0-9]+)([^0-9]?/?[0-9]*)?([^0-9])?$").FindStringSubmatch(street.Num)
	var bigNums []bigNumStruct
	switch {
	case numLen <= 2:
		width = 2539
		widthPhys = 215
		widthLineShort = 720
		widthLineLong = 1619
		arrowTranslate1 = 1290
		arrowTranslate2 = 920
		bigNums = subMatch(width, matchedString)
	case numLen == 3 || numLen == 4 && len(matchedString[1]) == 1:
		width = 3248
		widthPhys = 275
		widthLineShort = 1049
		widthLineLong = 2293
		arrowTranslate1 = 1776
		arrowTranslate2 = 1244
		bigNums = subMatch(width, matchedString)
	case numLen == 4 || numLen == 5 && len(matchedString[1]) == 2:
		width = 4015
		widthPhys = 340
		widthLineShort = 1428
		widthLineLong = 3060
		arrowTranslate1 = 2349
		arrowTranslate2 = 1623
		bigNums = subMatch(width, matchedString)
	default:
		width = 5196
		widthPhys = 440
		widthLineShort = 2018
		widthLineLong = 4236
		arrowTranslate1 = 3237
		arrowTranslate2 = 2218
		bigNums = subMatch(width, matchedString)
	}
	canvas.Startunit(widthPhys, heightPhys, "mm", fmt.Sprintf(`viewBox="0 0 %v %v"`, width, height))
	canvas.Roundrect(0, 0, width, height, 50, 50, "fill:#262741")
	for _, num := range bigNums {
		if num.num == "/" {
			canvas.Gtransform(fmt.Sprintf("translate(%v, 694.000000)", num.x+50))
			canvas.Gtransform("translate(126.709838, 392.581988) rotate(-75.000000) translate(-126.709838, -392.581988)")
			canvas.Rect(-273, 368, 800, 48, "fill:#FFFFFF")
			canvas.Gend()
			canvas.Gend()
		} else {
			canvas.Text(num.x, yMainText, num.num, num.style)
		}
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

	canvas.Text(470, 2226, street.Prev, arrowNumLeft)
	canvas.Text(490+widthLineLong, 2226, street.Next, arrowNumRight)

	canvas.End()
	return nil
}

func subMatch(width int, parsed []string) (bigNums []bigNumStruct) {
	widthText := (len(parsed[1]) * 878 / 2) + (utf8.RuneCountInString(parsed[2]) * 500 / 2) + (utf8.RuneCountInString(parsed[3]) * 500 / 2)
	if regexp.MustCompile("/").MatchString(parsed[2]) {
		widthText -= 150
	}
	if regexp.MustCompile("-").MatchString(parsed[1]) {
		widthText -= 308
	}
	for _, r := range parsed[1] {
		bigNums = append(bigNums, bigNumStruct{style: mainNumBig, num: string(r), x: width/2 - widthText})
		if r == '-' {
			widthText -= 500
		} else {
			widthText -= 808
		}
	}
	for _, r := range parsed[2] {
		bigNums = append(bigNums, bigNumStruct{style: mainNumLittle, num: string(r), x: width/2 - widthText})
		if r == '/' {
			widthText -= 300
		} else {
			widthText -= 500
		}
	}
	for _, r := range parsed[3] {
		bigNums = append(bigNums, bigNumStruct{style: mainNumLitLit, num: string(r), x: width/2 - widthText})
	}
	return bigNums
}
