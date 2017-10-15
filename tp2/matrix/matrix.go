package matrix

import (
	"strconv"
	"strings"
)

type Matrix [][]string

func CreateAdj(rows []string) Matrix {

	var adj Matrix

	adj = make(Matrix, len(rows)+1)

	for i := 0; i < len(rows); i++ {
		adj[i] = make([]string, len(rows)+1)
	}
	for k, v := range rows {
		adj[k] = strings.Split(v, " ")
	}
	return adj
}

//we have Matrix which has the distances between cities
// as they are ordered by name the first letter of citie is the index we are looking for on the matrix. But for this to work, the order must be bigger ascii corresponds to x and smaller ascii corresponds to y.
//Ascii table has letter k which isn't included on the cities of the matrix. Thats why we have the if and elses.
func (m Matrix) RoadLengthBetween(c1, c2 string) (dist int, e error) {
	if int(c1[0]) < int(c2[0]) {
		temp := c1
		c1 = c2
		c2 = temp
	}

	if c1 == "Lourel" || c2 == "Lourel" {
		if int(c1[0])-67 <= 9 {
			return strconv.Atoi(m[int(c1[0])-67][int(c2[0])-65])
		} else {
			return strconv.Atoi(m[int(c1[0])-67][int(c2[0])-66])
		}
	}
	if int(c1[0]) < 77 {
		if int(c1[0])-66 < 9 {
			return strconv.Atoi(m[int(c1[0])-66][int(c2[0])-65])

		} else {

			return strconv.Atoi(m[int(c1[0])-67][int(c2[0])-65])
		}
	} else {
		if int(c2[0]) < 77 {
			return strconv.Atoi(m[int(c1[0])-67][int(c2[0])-65])
		} else {

			return strconv.Atoi(m[int(c1[0])-67][int(c2[0])-66])
		}
	}
}
