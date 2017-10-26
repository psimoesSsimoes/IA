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

func DoLogStuff() {

	var format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} ▶ %{message}`,
	)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)

	backend2Formatter := logging.NewBackendFormatter(backend2, format)
	logging.SetBackend(backend2Formatter)

}

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

	return float64(LargestsDistances[1] + LargestsDistances[0] - SmallestsDistances[1] - SmallestsDistances[0])

}

//Called to determine if annealing should take place.
func Anneal(d, temperature float64) bool {
	/**If the temperature is below this amount, then the temperature plays no role in determining if annealing will take place or not*/
	if temperature < 1.0E-4 {
		//If the temperature is below this threshold then we simply check to see if the distance is greater than 0
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
	//close cycle
	d, _ := m.RoadLengthBetween(cities[len(cities)-1], cities[0])
	l += d

	return l
}

func Calculate2RandomNumbers(cities []string) (int, int) {
	//seed for random
	rand.Seed(time.Now().UTC().UnixNano())
	//The input values are then randomized according to the temperature. A higher temperature will result in more randomization, a lower temperature will result in less
	//randomization.
	i1 := int(math.Floor(float64(len(cities)-1) * rand.Float64()))
	j1 := int(math.Floor(float64(len(cities)-1) * rand.Float64()))
	//if the 2 indexes are the same, generate again
	for j1 == i1 {
		j1 = int(math.Floor(float64(len(cities)-1) * rand.Float64()))

	}
	return i1, j1
}

func Distance2FedAnnealing(i1, j1 int, adj matrix.Matrix, order []string) float64 {
	one, _ := adj.RoadLengthBetween(order[i1], order[i1+1])
	two, _ := adj.RoadLengthBetween(order[j1], order[j1+1])
	three, _ := adj.RoadLengthBetween(order[i1], order[j1])
	four, _ := adj.RoadLengthBetween(order[i1+1], order[j1+1])
	d := one + two - three - four
	return float64(d)
}

func PrintSolutions(best, worst, first, last Solution, cycle int, d time.Duration) {
	var log = logging.MustGetLogger("example")
	DoLogStuff()
	log.Notice("BEST SOLUTION ORDER:", best.Order)
	log.Notice("BEST SOLUTION LENGTH: ", best.Length)
	log.Notice("BEST SOLUTION FOUND AT TEMP:", best.Temperature)
	log.Notice("BEST SOLUTION FOUND AT ITERATION:", best.Iteration)
	log.Error("WORST SOLUTION ORDER:", worst.Order)
	log.Error("WORST SOLUTION LENGTH:", worst.Length)
	log.Error("WORST SOLUTION FOUND AT TEMP:", worst.Temperature)
	log.Error("WORST SOLUTION FOUND AT ITERATION:", worst.Iteration)
	log.Warning("FIRST SOLUTION ORDER:", first.Order)
	log.Warning("FIRST SOLUTION LENGTH:", first.Length)
	log.Warning("FIRST SOLUTION FOUND AT TEMP:", first.Temperature)
	log.Warning("FIRST SOLUTION FOUND AT ITERATION:", first.Iteration)
	log.Debugf("LAST SOLUTION ORDER: %+v", last.Order)
	log.Debugf("LAST SOLUTION LENGTH: %+v", last.Length)
	log.Debugf("LAST SOLUTION FOUND AT TEMP: %+v", last.Temperature)
	log.Debugf("LAST SOLUTION ORDER AT ITERATION: %+v", last.Iteration)
	log.Critical("TOTAL NUMBER OF ITERATIONS:", cycle)
	log.Critical("TIME OF EXECUTION:", d)

}

type Solution struct {
	Temperature float64
	Length      int
	Order       []string
	Iteration   int
}

func main() {

	var (
		best                          Solution = Solution{}
		worst                         Solution = Solution{}
		first                         Solution = Solution{}
		last                          Solution = Solution{}
		temperature                   float64
		orderlength, cycle, sameCount int
		cities, order                 []string
		isFirst                       bool
	)
	const (
		delta float64 = 0.03
	)
	sameCount = 0
	cycle = 1
	isFirst = true

	params := os.Args[1:]
	//grab params passed to program
	if params[0] == "all" {
		cities = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	} else {
		cities = params
	}
	//read rows of filrows, err := io.ReadFile("/home/psimoes/Github/IA/tp2/data/CitiesDist.txt")
	rows, err := io.ReadFile("/home/psimoes/Github/IA/tp2/data/teste1tp2")
	if err != nil {
	}

	fmt.Println(cities)
	//create matrix from rows read
	adj := matrix.CreateAdj(rows)
	temperature = FindInitialTemperature(cities, adj)
	worst.Temperature = temperature

	//copy(best.Order[:], cities)
	order = InitialOrder(cities)

	worst.Order = order
	best.Order = InitialOrder(cities)
	orderlength = Length(order, adj)
	best.Length = Length(best.Order, adj)
	worst.Length = orderlength

	//start measuring algorithm time
	start := time.Now()

	for sameCount < len(cities) {
		//for each temperature, the simulated annealing algorithm runs through a number of cycles predeterminated by us
		for j2 := 0; j2 < len(cities)*len(cities); j2++ {
			cycle++
			i1, j1 := Calculate2RandomNumbers(cities)
			d := Distance2FedAnnealing(i1, j1, adj, order)
			//copia por valor quando devia fazer copia por referencia
			if orderlength > worst.Length {
				worst.Length = orderlength
				temporary := make([]string, len(order))
				copy(temporary, order)
				worst.Order = temporary
				worst.Temperature = temperature
				worst.Iteration = cycle
			}
			if Anneal(d, temperature) {
				if isFirst {
					isFirst = false
					first.Temperature = temperature
					temporary := make([]string, len(order))
					copy(temporary, order)
					first.Order = temporary
					first.Length = orderlength
					first.Iteration = cycle
				}
				//If it is determined that we should anneal (j1 is less than i1), then we swap the values
				if j1 < i1 {
					k1 := i1
					i1 = j1
					j1 = k1
				}
				//and we loop between i1 and j1 and swap the values as we progress, i.e, reorders the path
				for ; j1 > i1; j1-- {
					i2 := order[i1+1]
					order[i1+1] = order[j1]
					order[j1] = i2
					i1++
				}

			}

		}
		//see if we found a worst solution
		orderlength = Length(order, adj)

		//if we found improvements change best and restart sameCOunt
		if orderlength < best.Length {
			best.Length = orderlength
			temporary := make([]string, len(order))
			copy(temporary, order)

			best.Order = temporary

			best.Temperature = temperature
			best.Iteration = cycle
			sameCount = 0
		} else {
			sameCount++
		}

		last.Order = order
		last.Length = orderlength
		last.Temperature = temperature
		last.Iteration = cycle

		// simply reduce the temperature by a fixed amount through each cycle.
		temperature *= delta
		cycle++
	}
	PrintSolutions(best, worst, first, last, cycle, time.Since(start))
}
