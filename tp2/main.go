package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

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

	rand.Seed(time.Now().UTC().UnixNano())
	for i := range s {
		j := rand.Intn(i + 1)
		s[i], s[j] = s[j], s[i]
	}
	return s
}
func Length(cities []string, m matrix.Matrix) int {
	l := 0
	for i := 1; i < len(cities); i++ {
		d, _ := m.RoadLengthBetween(cities[i], cities[i-1])
		l += d
	}
	return l
}

func main() {

	var cities []string
	var temperature, delta, temp_at_min, temp_at_first, temp_at_last, temp_at_best float64
	var orderlength, min_order_length int
	var minOrder, order, maxOrder, first, last []string
	var cycle, sameCount, all_it, it_first, it_last int
	temperature = 900.00
	sameCount = 0
	delta = 0.99
	cycle = 1
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
	copy(minOrder[:], cities)
	order = InitialOrder(cities)
	minOrder = InitialOrder(cities)
	orderlength = Length(order, adj)
	min_order_length = Length(minOrder, adj)
	for sameCount < len(cities) {
		for j2 := 0; j2 < len(cities)*len(cities); j2++ {
			rand.Seed(time.Now().UTC().UnixNano())
			i1 := int(math.Floor(float64(len(cities)-1) * rand.Float64()))
			j1 := int(math.Floor(float64(len(cities)-1) * rand.Float64()))
			for j1 == i1 {
				j1 = int(math.Floor(float64(len(cities)-1) * rand.Float64()))

			}
			if i1 >= 20 || j1 >= 20 {

			}

			first, _ := adj.RoadLengthBetween(order[i1], order[i1+1])
			second, _ := adj.RoadLengthBetween(order[j1], order[j1+1])
			third, _ := adj.RoadLengthBetween(order[i1], order[j1])
			fourth, _ := adj.RoadLengthBetween(order[i1+1], order[j1+1])
			d := first + second - third - fourth
			if Anneal(float64(d), temperature) {
				if j1 < i1 {
					k1 := i1
					i1 = j1
					j1 = k1
				}
				for ; j1 > i1; j1-- {
					i2 := order[i1+1]
					order[i1+1] = order[j1]
					order[j1] = i2
					i1++
				}
			}
		}
		//see if we found improvements
		orderlength = Length(order, adj)
		if orderlength < min_order_length {
			min_order_length = orderlength
			for k2 := 0; k2 < len(cities); k2++ {
				minOrder[k2] = order[k2]
			}
			temp_at_min = temperature
			sameCount = 0
		} else {
			sameCount++
		}
		temperature *= 1 - delta
		fmt.Println(temperature)
		cycle++
	}

	fmt.Println("percurso: %+v\ncusto: %+v\ntemperature:%+v\n", minOrder, min_order_length, temp_at_min)

}
