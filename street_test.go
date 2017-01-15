package main

import "testing"
import "github.com/stretchr/testify/assert"

func TestStreet_defineStreetClassEng(t *testing.T) {
	s := Street{Class: "вулиця"}
	err := s.defineStreetClassEng()
	assert.NoError(t, err, "Proper class")
	assert.Equal(t, "vulitsya", s.ClassEng)

	s = Street{Class: "BAD CLASS"}
	err = s.defineStreetClassEng()
	assert.Error(t, err, "Bad class should error")
}

func TestStreet_defineStreetName(t *testing.T) {
	s := Street{Name: "вознякова"}
	s.defineStreetName()
	assert.Equal(t, "Vozniakova", s.NameEng)
	assert.Equal(t, "Вознякова", s.Name, "street Name should auto capitalize")
}

func TestStreet_createID(t *testing.T) {
	s := Street{}
	s.createID()
	assert.NotEqual(t, "", s.ID)
}
