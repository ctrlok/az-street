package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Sirupsen/logrus"
)

func generateArchiveAndPng(log *logrus.Entry, street Street) (filenames []string, err error) {
	dir, err := ioutil.TempDir(tmpDir, "")
	if err != nil {
		log.Errorf("Problem with creating temporary directory in %s", tmpDir)
		return
	}
	log = log.WithField("tmp_dir", dir)
	defer func() {
		if os.RemoveAll(dir) != nil {
			log.Error("Problem with remove temporary directory")
		}
	}()
	log.Debugf("Temparay directory created")
	err = generateSVGfiles(log, dir, street)
	if err != nil {
		return
	}
	err = generatePNGfiles(log, dir)
	if err != nil {
		return
	}
	err = generateEPSfiles(log, dir)
	if err != nil {
		return
	}
	err = removeSVGfiles(log, dir)
	if err != nil {
		return
	}
	archive, err := makeArchive(log, dir)
	if err != nil {
		return
	}
	convertedFiles, err := convertPNGfiles(log, dir)
	if err != nil {
		return
	}
	filenames, err = moveFiles(log, convertedFiles, archive)
	return
}

func generateSVGfiles(log *logrus.Entry, dir string, street Street) (err error) {
	log = log.WithField("func", "generateSVGfiles")
	defer func() { log = log.WithField("func", "") }()
	t := startTimer()
	streetFile, err := os.Create(fmt.Sprint(dir, "/street.svg"))
	if err != nil {
		log.Errorf("Problem creating SVG street file: %s", err)
		return err
	}
	defer streetFile.Close()
	err = streetSVG(street, streetFile)
	log.WithField("operation_time", t.diff()).Debug("Street svg file generated")
	if err != nil {
		log.Error("Problem with generating street SVG file")
		return
	}
	t = startTimer()
	numFile, err := os.Create(fmt.Sprint(dir, "/num.svg"))
	if err != nil {
		log.Errorf("Problem creating SVG num file: %s", err)
		return err
	}
	defer numFile.Close()
	err = numSVG(street, numFile)
	log.WithField("operation_time", t.diff()).Debug("Num svg file generated")
	if err != nil {
		log.Error("Problem with generating num SVG file")
	}
	return
}

func generatePNGfiles(log *logrus.Entry, dir string) (err error) {
	log = log.WithField("func", "generatePNGfiles")
	defer func() { log = log.WithField("func", "") }()
	t := startTimer()
	for _, name := range []string{"street", "num"} {
		svgPath := fmt.Sprint(dir, "/", name, ".svg")
		pngPath := fmt.Sprint(dir, "/", name, ".png")
		command := exec.Command("inkscape", "-z", "-T", "-e", pngPath, svgPath)
		command.Stderr = os.Stderr
		err = command.Run()
		if err != nil {
			log.Errorf("Problem with generating PNG files: %s", err)
			return err
		}
	}
	log.WithField("operation_time", t.diff()).Debug("PNG files generated")
	return
}

func generateEPSfiles(log *logrus.Entry, dir string) (err error) {
	log = log.WithField("func", "generateEPSfiles")
	defer func() { log = log.WithField("func", "") }()
	t := startTimer()
	for _, name := range []string{"street", "num"} {
		svgPath := fmt.Sprint(dir, "/", name, ".svg")
		epsPath := fmt.Sprint(dir, "/", name, ".eps")
		command := exec.Command("inkscape", "-z", "-T", "-E", epsPath, svgPath)
		command.Stderr = os.Stderr
		err = command.Run()
		if err != nil {
			log.Errorf("Problem with generating EPS files: %s", err)
			return err
		}
	}
	log.WithField("operation_time", t.diff()).Debug("PNG files generated")
	return
}

func removeSVGfiles(log *logrus.Entry, dir string) (err error) {
	log = log.WithField("func", "removeSVGfiles")
	defer func() { log = log.WithField("func", "") }()
	for _, name := range []string{"street", "num"} {
		svgPath := fmt.Sprint(dir, "/", name, ".svg")
		command := exec.Command("rm", svgPath)
		command.Stderr = os.Stderr
		err = command.Run()
		if err != nil {
			log.Errorf("Problem with removing SVG files: %s", err)
			return err
		}
	}
	return
}

func makeArchive(log *logrus.Entry, dir string) (archive string, err error) {
	log = log.WithField("func", "makeArchive")
	defer func() { log = log.WithField("func", "") }()
	t := startTimer()
	archive = fmt.Sprint(dir, ".zip")
	command := exec.Command("zip", "-r", "-j", archive, dir)
	command.Stderr = os.Stderr
	err = command.Run()
	if err != nil {
		log.Errorf("Problem with generating Archive: %s", err)
		return "", err
	}
	log.WithField("operation_time", t.diff()).Debug("Archive created")
	return
}

func convertPNGfiles(log *logrus.Entry, dir string) (generatedFiles []string, err error) {
	log = log.WithField("func", "convertPNGfiles")
	defer func() { log = log.WithField("func", "") }()
	t := startTimer()
	for _, name := range []string{"street", "num"} {
		pngPath := fmt.Sprint(dir, "/", name, ".png")
		for _, size := range []string{"80", "160", "240"} {
			resizedName := fmt.Sprint(dir, "_", name, "_", size, ".png")
			command := exec.Command("convert", pngPath, "-resize", size, resizedName)
			command.Stderr = os.Stderr
			err = command.Run()
			if err != nil {
				log.Errorf("Problem with converting PNG files: %s", err)
				return nil, err
			}
			generatedFiles = append(generatedFiles, resizedName)
		}
	}
	log.WithField("operation_time", t.diff()).Debug("PNG files generated")
	return
}

func moveFiles(log *logrus.Entry, files []string, archive string) (returnFiles []string, err error) {
	log = log.WithField("func", "generatePNGfiles")
	defer func() { log = log.WithField("func", "") }()
	fromFiles := []string{archive}
	for _, file := range append(fromFiles, files...) {
		command := exec.Command("mv", file, archiveDir)
		command.Stderr = os.Stderr
		err = command.Run()
		if err != nil {
			log.Errorf("Problem move files to final directory: %s", err)
			return nil, err
		}
		returnFiles = append(returnFiles, fmt.Sprint(archiveDir, "/", filepath.Base(file)))
	}
	return
}
