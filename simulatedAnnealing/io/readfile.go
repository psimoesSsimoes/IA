package io

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

type rows []string

var CitiesNames = []string{"Atroeira", "Belmar", "Cerdeira", "Douro", "Encosta", "Freira", "Gonta", "Horta", "Infantado", "Jardim", "Lourel", "Monte", "Nelas", "Oura", "Pinhal", "Quebrada", "Roseiral", "Serra", "Teixoso", "Ulgueira", "Vilar"}

func ReadFile(name string) (rows, error) {
	content, err := ioutil.ReadFile(name)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	s := fmt.Sprintf("%s", content)
	file_rows := strings.Split(s, "\n")
	return file_rows, nil
}
