package main

import (
	"log"
	"math"
	"math/rand"
	"os"

	"./io"
	"./matrix"
)

//Called to determine if annealing should take place.
func Anneal(d, temperature float64) bool {

	if temperature < 1.0E-4 {
		if d > 0.0 {
			return true
		} else {
			return false
		}
	}
	if rand.Float64() < math.Exp(d/temperature) {
		return true
	} else {
		return false
	}
}

func Mod(i, len int) int {
	return i % len
}

//initialOrder is just a simple shuffle of the array
func InitialOrder(s []string) []string {
	for i := range s {
		j := rand.Intn(i + 1)
		s[i], s[j] = s[j], s[i]
	}
	return s
}
func Length(cities []string, m matrix.Matrix) int {
	l := 0
	for i := 0; i < (len(cities)); i++ {
		d, _ := m.RoadLengthBetween(cities[i], cities[i+1])
		l += d
	}
	return l
}

func main() {

	var cities, mincities []string
	var temperature float64
	var length, minlength int
	var minOrder, order []string
	var cycle, sameCount int
	temperature = 800.00
	sameCount = 0

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

	copy(mincities[:], cities)
	cities = InitialOrder(cities)
	mincities = InitialOrder(mincities)
	length = Length(cities, adj)
	minlength = length
	for sameCount < len(cities) {
		for j2 := 0; j2 < len(cities)*len(cities); j2++ {
			i1 := int(math.Floor(float64(len(cities) * rand.Int())))
			j1 := int(math.Floor(float64(len(cities) * rand.Int())))

			first, _ := adj.RoadLengthBetween(cities[i1], cities[i1+1])
			second, _ := adj.RoadLengthBetween(cities[j1], cities[j1+1])
			third, _ := adj.RoadLengthBetween(cities[i1], cities[j1])
			fourth, _ := adj.RoadLengthBetween(cities[i1+1], cities[j1+1])
			d := first + second - third - fourth
			if Anneal(float64(d), temperature) {
				if j1 < i1 {
					k1 := i1
					i1 = j1
					j1 = k1
				}
				for ; j1 > i1; j1-- {
					i2 := cities[i1+1]
					cities[i1+1] = cities[j1]
					cities[j1] = i2
					i1++
				}
			}
		}
		//see if we found improvements
		length = Length(cities, adj)
		if length < minlength {
			minlength = length
		}
	}
}
