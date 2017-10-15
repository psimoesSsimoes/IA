package main

import (
	"fmt"
	"log"
	"os"

	"./io"
	"./matrix"
)

func main() {

	var cities []string

	params := os.Args[1:]
	if params[0] == "all" {
		cities = []string{"Atroeira", "Belmar", "Cerdeira", "Douro", "Encosta", "Freira", "Gonta", "Horta", "Infantado", "Jardim", "Lourel", "Monte", "Nelas", "Oura", "Pinhal", "Quebrada", "Roseiral", "Serra", "Teixoso", "Ulgueira", "Vilar"}
	} else {
		cities = params
	}
	rows, err := io.ReadFile("/home/psimoes/Github/IA/tp2/data/CitiesDist.txt")
	if err != nil {
		log.Fatal(err)
	}
	adj := matrix.CreateAdj(rows)
	//0/1/2/3/9/10 passed
	for i := 0; i < 21; i++ {
		if i == 7 {

		} else {
			dist, err := adj.RoadLengthBetween(cities[7], cities[i])
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(cities[7], cities[i], dist)
		}
	}

}
