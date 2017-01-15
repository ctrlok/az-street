package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_generateFiles(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip integration test in short mode")
	}
	street := Street{
		Name:  "Коломийський",
		Num:   "22/1а",
		Class: "провулок",
		Left:  "33",
		Right: "31",
	}
	street.fill()
	dir, _ := ioutil.TempDir(tmpDir, "archive")
	defer os.RemoveAll(dir)

	log := logrus.WithField("test", "test")
	// SVG files
	err := generateSVGfiles(log, dir, street)
	assert.NoError(t, err, "SVG files should be generated")
	_, err = os.Stat(fmt.Sprint(dir, "/street.svg"))
	assert.NoError(t, err, "Check svg street file exist")
	_, err = os.Stat(fmt.Sprint(dir, "/num.svg"))
	assert.NoError(t, err, "Check svg num file exist")

	// PNG files
	err = generatePNGfiles(log, dir)
	assert.NoError(t, err, "Files should be proper generated")
	for _, file := range []string{"street.png", "num.png"} {
		_, err = os.Stat(fmt.Sprint(dir, "/", file))
		assert.NoError(t, err, "PNG file should exist")
	}

	// EPS files
	err = generateEPSfiles(log, dir)
	assert.NoError(t, err, "Files should be proper generated")
	for _, file := range []string{"street.eps", "num.eps"} {
		_, err = os.Stat(fmt.Sprint(dir, "/", file))
		assert.NoError(t, err, "EPS file should exist")
	}

	// Remove SVG files from tmp dir
	err = removeSVGfiles(log, dir)
	assert.NoError(t, err, "SVG files should be removed")
	_, err = os.Stat(fmt.Sprint(dir, "/street.svg"))
	assert.Error(t, err, "Check svg street file not exist")
	_, err = os.Stat(fmt.Sprint(dir, "/num.svg"))
	assert.Error(t, err, "Check svg num file not exist")

	// Archive
	archive, err := makeArchive(log, dir)
	assert.NoError(t, err, "Archive should be created")
	_, err = os.Stat(archive)
	assert.NoError(t, err, "Archive file should exist")

	// convert PNG files
	files, err := convertPNGfiles(log, dir)
	assert.NoError(t, err, "PNG files should be converted without error")
	for _, file := range files {
		_, err = os.Stat(file)
		assert.NoError(t, err, "Files should exist")
	}

	// move files to DST
	dst, err := moveFiles(log, files, archive)
	assert.NoError(t, err, "Files should move without error")
	for _, file := range dst {
		_, err = os.Stat(file)
		assert.NoError(t, err, "Files should exist")
	}
	for _, file := range append(files, archive) {
		_, err = os.Stat(file)
		assert.Error(t, err, "Files should not exist")
	}
}
