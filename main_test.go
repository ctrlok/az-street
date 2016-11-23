package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestHttpHandlerErrJson(t *testing.T) {
	assert := assert.New(t)
	body := strings.NewReader("")
	req := httptest.NewRequest("POST", "/", body)
	res := httptest.NewRecorder()
	httpHandler(res, req)
	assert.Equal(http.StatusInternalServerError, res.Result().StatusCode, "")

}

func TestDecodeSucc(t *testing.T) {
	assert := assert.New(t)
	body := strings.NewReader(`{
    "street_name_ua": "Коломийський",
    "street_type": "provulok",
    "street_num": "22/1а",
    "street_position_first": true,
    "street_position_last": true
}`)
	street, err := decode(body)
	assert.NoError(err, "Decode problem")
	assert.Equal("Коломийський", street.StreetNameUA, "")
	assert.Equal("22/1а", street.StreetNum, "")
	assert.True(street.StreetPositionFirst, "")
	assert.True(street.StreetPositionLast, "")
	assert.Equal("provulok", street.StreetType, "")
}

func TestDecodeErr(t *testing.T) {
	assert := assert.New(t)
	// Test not json input
	body := strings.NewReader("{not:json}")
	_, err := decode(body)
	assert.Error(err, "")

	// Test not valid json
	body = strings.NewReader(`{"street_name_ua":1}`)
	_, err = decode(body)
	assert.Error(err, "")
}

func TestDefineStreetName(t *testing.T) {
	street := Street{StreetNameUA: "Коломийський"}
	defineStreetName(&street)
	assert.Equal(t, "Kolomyiskyi", street.StreetNameEng, "")

	street2 := Street{StreetNameUA: "коломийський"}
	defineStreetName(&street2)
	assert.Equal(t, "Kolomyiskyi", street2.StreetNameEng, "test capitalize")
}

func TestRenderStepSucc(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip integration test in short mode")
	}
	street := Street{
		StreetNameUA:        "Коломийський",
		StreetNum:           "22/1а",
		StreetType:          "provulok",
		StreetPositionFirst: true,
		StreetPositionLast:  true,
	}
	street.createID()
	dir, _ := ioutil.TempDir(tmpDirPath, "archive")
	defer os.RemoveAll(dir)
	err := renderSVG(street, dir)
	assert.NoError(t, err, "render SVG files")
	_, err = os.Stat(fmt.Sprint(dir, "/street.svg"))
	assert.Nil(t, err, "Check svg file exist")
	err = renderPNG(dir)
	assert.NoError(t, err, "render PNG files")
	_, err = os.Stat(fmt.Sprint(dir, "/street.png"))
	assert.Nil(t, err, "Check png file exist")
	err = renderEPS(dir)
	assert.NoError(t, err, "render EPS files")
	_, err = os.Stat(fmt.Sprint(dir, "/street.eps"))
	assert.Nil(t, err, "Check eps file exist")
	err = renderPDF(dir)
	assert.NoError(t, err, "render PDF files")
	_, err = os.Stat(fmt.Sprint(dir, "/street.pdf"))
	assert.Nil(t, err, "Check pdf file exist")
	err = removeSVG(dir)
	assert.NoError(t, err, "remove SVG files")
	_, err = os.Stat(fmt.Sprint(dir, "/street.svg"))
	assert.NotNil(t, err, "Check png file was removed")
}

func TestRenderStepFail(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip integration test in short mode")
	}
	street := Street{
		StreetNameUA:        "Коломийський",
		StreetNum:           "22/1а",
		StreetType:          "provulok",
		StreetPositionFirst: true,
		StreetPositionLast:  true,
	}
	street.createID()
	dir, _ := ioutil.TempDir(tmpDirPath, "archive")
	defer os.RemoveAll(dir)
	err := renderPNG(dir)
	assert.Error(t, err, "render PNG files")
	_, err = os.Stat(fmt.Sprint(dir, "/street.png"))
	assert.NotNil(t, err, "Check png file exist")
	err = renderEPS(dir)
	assert.Error(t, err, "render EPS files")
	_, err = os.Stat(fmt.Sprint(dir, "/street.eps"))
	assert.NotNil(t, err, "Check eps file exist")
	err = renderPDF(dir)
	assert.Error(t, err, "render PDF files")
	_, err = os.Stat(fmt.Sprint(dir, "/street.pdf"))
	assert.NotNil(t, err, "Check pdf file exist")
}

func TestMakeArchiveSucc(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip integration test in short mode")
	}
	street := Street{
		StreetNameUA:        "Коломийський",
		StreetNum:           "22/1а",
		StreetType:          "provulok",
		StreetPositionFirst: true,
		StreetPositionLast:  true,
	}
	street.createID()
	archive, err := makeArchive(&street)
	assert.NoError(t, err, "Check archive created")
	_, err = os.Stat(archive)
	assert.Nil(t, err, "Check archive file exist")
	filesInTmp, _ := ioutil.ReadDir(archiveDir)
	assert.Equal(t, 1, len(filesInTmp), "should be only one file in archive dir")
}
