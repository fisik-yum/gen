package main

import "math"

var phrase_split = []rune(TARGET_PHRASE)

func fitnessMeanRuneDistance(ind individual) float64 {
	var dists [GENE_SIZE]float64
	for i := 0; i < int(GENE_SIZE); i++ {
		dists[i] = (math.Abs(float64(ind[i] - phrase_split[i])))
	}
	var sum float64
	for i := 0; i < int(GENE_SIZE); i++ {
		sum += dists[i]
	}
	return sum / float64(GENE_SIZE)

}

func fitnessHammingDistance(ind individual) float64 {
	var total float64
	for i := 0; i < int(GENE_SIZE); i++ {
		if phrase_split[i] != ind[i] {
			total++
		}
	}
	return total
}
