package main

import (
	"fmt"
	"github.com/ajstarks/svgo/float"
	"os"
	"regexp"
)

const textBigFixed string = "font-size:380.8px;font-family:ProbaNav2-Medium, Proba Nav2;letter-spacing:0.01em;fill:#fff;text-anchor:start;"
const textBig1 string = "font-size:220px;font-family:ProbaNav2-Medium, Proba Nav2;letter-spacing:0.01em;fill:#fff;text-anchor:start;"
const textLittle string = "font-size:90px;font-family:ProbaNav2-Regular, Proba Nav2;letter-spacing:0.01em;fill:#fff;text-anchor:start;"
const line string = "fill:none;stroke:#fff;stroke-miterlimit:10;stroke-width:4px;"

var textLittleEnd = fmt.Sprint(textLittle, "text-anchor:end;")
var textBig = fmt.Sprint(textBigFixed, "text-anchor:middle;")

var arrow1x = []float64{130.83, 133.56, 119.11, 133.56, 130.83, 113.23, 130.83}
var arrowY = []float64{447.91, 444.98, 431.56, 418.13, 415.2, 431.56, 447.91}
var arrow2x = []float64{478.62, 475.89, 490.34, 475.89, 478.62, 496.21, 478.62}

func renderSVGnum(street Street, dir string) (err error) {
	file, err := os.Create(fmt.Sprint(dir, "/num.svg"))
	if err != nil {
		return err
	}
	defer file.Close()
	canvas := svg.New(file)
	draw1svg(street.Num, "1", "7", canvas)
	return nil
}

func draw1svg(unaprsedBaseNum, leftNum, rightNum string, canvas *svg.SVG) {
	width := 609.45
	height := 609.45
	canvas.Start(width, height)
	canvas.Gtransform("translate(609.45 0) rotate(90)")
	canvas.Roundrect(0.0, 0.0, height, width, 14.17, 14.17, "fill:#262741")
	canvas.Gend()

	// Little text
	canvas.Text(106.36, 540.94, leftNum, textLittle)
	canvas.Text((width - 106.36), 540.94, rightNum, textLittleEnd)

	// Big text
	re := regexp.MustCompile("([0-9])([^0-9])")
	parsed := re.FindStringSubmatch(unaprsedBaseNum)
	if parsed == nil {
		canvas.Text(width/2, 350.75, unaprsedBaseNum, textBig)
	} else {
		canvas.Text(106.36, 350.75, parsed[1], textBigFixed)
		canvas.Text(130.73+176.65, 350.75, parsed[2], textBig1)
	}

	//Line
	canvas.Line(117.57, 431.55, 295.56, 431.55, line)
	canvas.Polygon(arrow1x, arrowY, "fill:#fff;")

	canvas.Line(313.89, 431.55, 491.88, 431.55, line)
	canvas.Polygon(arrow2x, arrowY, "fill:#fff;")
	canvas.End()
}

func draw2svg(unaprsedBaseNum, leftNum, rightNum string, canvas *svg.SVG) {
	width := 779.53
	height := 609.45
	canvas.Start(width, height)
	canvas.Gtransform("translate(779.53 0) rotate(90)")
	canvas.Roundrect(0.0, 0.0, height, width, 14.17, 14.17, "fill:#262741")
	canvas.Gend()

	// Little text
	canvas.Text(106.36, 540.94, leftNum, textLittle)
	canvas.Text((width - 106.36), 540.94, rightNum, textLittleEnd)
	// Big text
	re := regexp.MustCompile("([0-9]+)([^0-9])")
	numAndLetter := re.FindStringSubmatch(unaprsedBaseNum)
	switch {
	case numAndLetter != nil:
		canvas.Gtransform("translate(92.36 350.75)")
		canvas.Text(0, 0, numAndLetter[1], textBigFixed)
		canvas.Text(409.15, 0, numAndLetter[2], textBig1)
		canvas.Gend()

	// case numAndNum != nil:
	// case numAndNumAndLetter != nil:
	default:
		canvas.Text(width/2, 350.75, unaprsedBaseNum, textBig)
	}

	//Line
	canvas.Line(117.57, 431.55, 382.28, 431.55, line)
	canvas.Polygon(arrow1x, arrowY, "fill:#fff;")

	canvas.Line(397.42, 431.55, 662.13, 431.55, line)
	x := []float64{}
	for _, v := range arrow2x {
		x = append(x, v+170.26)
	}
	canvas.Polygon(x, arrowY, "fill:#fff;")
	canvas.Line(106.36, 0, 106.36, 990, line)
	canvas.End()
}
