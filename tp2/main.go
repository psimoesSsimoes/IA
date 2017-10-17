package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"./io"
	"./matrix"
	"github.com/op/go-logging"
)

func FindInitialTemperature(cities []string, adj matrix.Matrix) float64 {

	LargestsDistances := make([]int, 2)
	SmallestsDistances := make([]int, 2)
	smallest := math.MaxUint32
	largest := 0

	for i := 0; i < len(cities); i++ {
		for j := 0; j < len(cities); j++ {
			p, err := adj.RoadLengthBetween(cities[i], cities[j])
			if err == nil {
				if smallest > p {
					fmt.Println("entered smallest")
					smallest = p
					SmallestsDistances[1] = SmallestsDistances[0]
					SmallestsDistances[0] = p
				}
				if largest < p {
					largest = p
					LargestsDistances[1] = LargestsDistances[0]
					LargestsDistances[0] = p
				}
			}
		}
	}
	fmt.Printf("%+v , %+v , %+v ,%+v", LargestsDistances[1], LargestsDistances[0], SmallestsDistances[1], SmallestsDistances[0])

	return float64(LargestsDistances[1] + LargestsDistances[0] - SmallestsDistances[1] - SmallestsDistances[0])

}

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

	var log = logging.MustGetLogger("example")

	var format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} ▶ %{message}`,
	)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)

	backend2Formatter := logging.NewBackendFormatter(backend2, format)
	logging.SetBackend(backend2Formatter)

	var cities []string
	var temperature, delta, temp_at_min /**, temp_at_first, temp_at_last, temp_at_best*/ float64
	var orderlength, min_order_length int
	var minOrder, order /**, maxOrder, first, last*/ []string
	var cycle, sameCount /**, all_it, it_first, it_last*/ int

	sameCount = 0
	delta = 0.03
	cycle = 1
	params := os.Args[1:]
	if params[0] == "all" {
		cities = []string{"Atroeira", "Belmar", "Cerdeira", "Douro", "Encosta", "Freira", "Gonta", "Horta", "Infantado", "Jardim", "Lourel", "Monte", "Nelas", "Oura", "Pinhal", "Quebrada", "Roseiral", "Serra", "Teixoso", "Ulgueira", "Vilar"}
	} else {
		cities = params
	}
	rows, err := io.ReadFile("/home/psimoes/Github/IA/tp2/data/CitiesDist.txt")
	if err != nil {
	}
	adj := matrix.CreateAdj(rows)
	temperature = FindInitialTemperature(cities, adj)
	fmt.Print("Initial Temperature: ")
	fmt.Println(temperature)
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
			fmt.Println(i1)
			fmt.Println(j1)
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
		temperature *= delta
		fmt.Println(temperature)
		cycle++
	}
	log.Notice("BEST SOLUTION ORDER:", minOrder)
	log.Notice("BEST SOLUTION LENGTH: ", min_order_length)
	log.Notice("BEST SOLUTION FOUND AT TEMP:", temp_at_min)
	log.Notice("BEST SOLUTION FOUND AT ITERATION:", temp_at_min)
	log.Error("WORST SOLUTION ORDER:")
	log.Error("WORST SOLUTION LENGTH:")
	log.Error("WORST SOLUTION FOUND AT TEMP:")
	log.Error("WORST SOLUTION FOUND AT ITERATION:")
	log.Warning("FIRST SOLUTION ORDER:")
	log.Warning("FIRST SOLUTION LENGTH:")
	log.Warning("FIRST SOLUTION FOUND AT TEMP:")
	log.Warning("FIRST SOLUTION FOUND AT ITERATION:")
	log.Debugf("LAST SOLUTION ORDER:")
	log.Debugf("LAST SOLUTION LENGTH:")
	log.Debugf("LAST SOLUTION FOUND AT TEMP:")
	log.Debugf("LAST SOLUTION ORDER AT ITERATION:")
	log.Critical("TOTAL NUMBER OF ITERATIONS:")
	log.Critical("TIME OF EXECUTION:")
}
