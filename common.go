package main

import (
	"log"
	"strconv"
)

func PanicIfError(err error) {
	if err != nil {
		log.Panic(err.Error())
	}
}

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}
