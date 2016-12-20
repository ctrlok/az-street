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
    "name_ua": "Коломийський",
    "type": "provulok",
    "num": "22/1а",
		"prev": "31",
		"next": "33"
}`)
	street, err := decode(body)
	assert.NoError(err, "Decode problem")
	assert.Equal("Коломийський", street.NameUA, "")
	assert.Equal("22/1а", street.Num, "")
	assert.Equal("31", street.Prev)
	assert.Equal("33", street.Next)
	assert.Equal("provulok", street.Type, "")
}

func TestDecodeErr(t *testing.T) {
	assert := assert.New(t)
	// Test not json input
	body := strings.NewReader("{not:json}")
	_, err := decode(body)
	assert.Error(err, "")

	// Test not valid json
	body = strings.NewReader(`{"name_ua":1}`)
	_, err = decode(body)
	assert.Error(err, "")
}

func TestDefineStreetName(t *testing.T) {
	street := Street{NameUA: "Коломийський"}
	street.defineStreetName()
	assert.Equal(t, "Kolomyiskyi", street.NameEng, "")

	street2 := Street{NameUA: "коломийський"}
	street2.defineStreetName()
	assert.Equal(t, "Kolomyiskyi", street2.NameEng, "test capitalize")
}

func TestRenderStepSucc(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip integration test in short mode")
	}
	street := Street{
		NameUA: "Коломийський",
		Num:    "22/1а",
		Type:   "provulok",
		Next:   "33",
		Prev:   "31",
	}
	street.createID()
	dir, _ := ioutil.TempDir(tmpDirPath, "archive")
	defer removeDirs(dir)
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
	err = removeSVG(dir)
	assert.NoError(t, err, "remove SVG files")
	_, err = os.Stat(fmt.Sprint(dir, "/street.svg"))
	assert.NotNil(t, err, "Check png file was removed")
}

func TestRenderStepFail(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip integration test in short mode")
	}
	dir, _ := ioutil.TempDir(tmpDirPath, "archive")
	defer removeDirs(dir)
	err := renderPNG(dir)
	assert.Error(t, err, "render PNG files")
	_, err = os.Stat(fmt.Sprint(dir, "/street.png"))
	assert.NotNil(t, err, "Check png file exist")
	err = renderEPS(dir)
	assert.Error(t, err, "render EPS files")
	_, err = os.Stat(fmt.Sprint(dir, "/street.eps"))
	assert.NotNil(t, err, "Check eps file exist")
}

func TestMakeArchiveSucc(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip integration test in short mode")
	}
	street := Street{
		NameUA: "Варенична",
		Num:    "22",
		Prev:   "20",
		Next:   "24",
		Type:   "vulitsya",
	}
	street.createID()
	archive, err := makeArchive(&street)
	assert.NoError(t, err, "Check archive created")
	_, err = os.Stat(archive)
	assert.Nil(t, err, "Check archive file exist")
	filesInTmp, _ := ioutil.ReadDir(archiveDir)
	assert.Equal(t, 1, len(filesInTmp), "should be only one file in archive dir")
}

func TestMakeAllImages(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip integration test in short mode")
	}
	street := Street{
		NameUA: "Щорса",
		Num:    "2",
		Type:   "provulok",
		Prev:   "31",
		Next:   "33",
	}
	street.createID()
	dir, _ := ioutil.TempDir(tmpDirPath, "archive")
	defer removeDirs(dir)
	err := renderSVG(street, dir)
	assert.NoError(t, err, "")
	err = renderPNG(dir)
	assert.NoError(t, err, "")
	err = renderEPS(dir)
	assert.NoError(t, err, "")

}
