package main

import (
	"log"
	"math/rand"
	"strconv"
)

//TODO : more flashy color
func GenerateColor() string {
	var r, g, b int

	r = (rand.Intn(256) + 60) / 2
	g = (rand.Intn(256) + 60) / 2
	b = (rand.Intn(256) + 60) / 2

	return ("rgb(" + strconv.Itoa(r) + ", " + strconv.Itoa(g) + ", " + strconv.Itoa(b) + ")")
}

func Atoi(str string) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		log.Println(err)
	}
	return val
}
