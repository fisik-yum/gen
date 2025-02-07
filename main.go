package main

import (
	"fmt"
	"gonum.org/v1/gonum/stat"
	"math/rand"
	"sort"
	//"github.com/ka-weihe/fast-levenshtein"
)

// Problem setup
const TARGET_PHRASE = "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum."

// const TARGET_PHRASE="Hello, World!"
const GENE_SIZE = int32(len(TARGET_PHRASE))
const MAX_GENERATIONS = 2500

// Population parameters
const POP_SIZE = 2048
const MUTATION_CHANCE = 0.01
const POP_HOLD = 10

var FITNESS_FUNC = fitnessMeanRuneDistance

// [min,max] for rune values
const UTF8_MIN = int32(' ')
const UTF8_MAX = int32(126)

type individual [GENE_SIZE]rune

// Population
var population []individual

func (i individual) String() string {
	var ret string
	for _, c := range i {
		ret += string(c)
	}
	return ret
}

// Perform a two-point crossover between two individuals
func cross(a individual, b individual) (individual, individual) {
	point1 := rand.Int31n(GENE_SIZE)
	point2 := point1 + rand.Int31n(GENE_SIZE-point1) //ensure second point is >= first point
	//fmt.Printf("Cross points: %v %v\n", point1, point2)
	crossA := splice(a[:point1], b[point1:point2], a[point2:])
	crossB := splice(b[:point1], a[point1:point2], b[point2:])
	return crossA, crossB
}

// Splice two pre-scliced genomes together
func splice(genomes ...[]rune) individual {
	var ret individual
	cpos := 0
	for _, sector := range genomes {
		for _, g := range sector {
			ret[cpos] = g
			cpos++
		}
	}
	return ret
}
func genRand() individual {
	var ret individual
	for i := 0; i < int(GENE_SIZE); i++ {
		ret[i] = rand.Int31n(100)
	}
	return ret
}
func calcFitness(ind individual) float64 {
	//calculate fitness score of a member
	// mean character diffference
	//return float64(levenshtein.Distance(ind.String(),TARGET_PHRASE))//float64(0)
	return FITNESS_FUNC(ind)
}

// Summarize statistics for a generation. Pass in a sorted slice by fitness
func summarize(pop []individual, gen int) {
	fitnesses := make([]float64, POP_SIZE)
	for i := 0; i < POP_SIZE; i++ {
		fitnesses[i] = calcFitness(pop[i])
	}
	fmt.Printf("Generation: %d\nFitness Range: [%f, %f]\n", gen, fitnesses[0], fitnesses[POP_SIZE-1])
	fmt.Printf("Mean Fitness: %f\nStd.Dev: %f", stat.Mean(fitnesses, nil), stat.StdDev(fitnesses, nil))
	fmt.Printf("\nBest Individual: %v\n\n", pop[0])
}
func genRandGene() rune {
	return UTF8_MIN + rand.Int31n(UTF8_MAX-UTF8_MIN+1)
}
func genRandIndividual() (ret [GENE_SIZE]rune) {
	for i := 0; i < int(GENE_SIZE); i++ {
		ret[i] = genRandGene()
	}
	return
}
func mutate(ind individual, rand_chance float32) (ret individual) {
	for i := 0; i < int(GENE_SIZE); i++ {
		if rand.Float32() < rand_chance {
			ret[i] = genRandGene()
		} else {
			ret[i] = ind[i]
		}
	}
	return
}
func simulateGeneration(pop []individual) ([]individual, bool) {
	if calcFitness(pop[0]) == 0 {
		return pop, true
	} else {
		//create new slice for new population
		newpop := make([]individual, 0)
		newpop = append(newpop, pop[0:POP_HOLD]...)
		for i := 0; i < POP_SIZE-POP_HOLD; i += 2 {
			a, b := cross(newpop[0], newpop[1])
			newpop = append(newpop, a, b)
		}
		for i := 0; i < POP_SIZE; i++ {
			newpop[i] = mutate(newpop[i], MUTATION_CHANCE)
		}
		return newpop, false

	}
}

// initialize population and generation 0
func init() {
	population = make([]individual, POP_SIZE)
	for i := 0; i < POP_SIZE; i++ {
		population[i] = genRandIndividual()
	}

}
func main() {
	state := false
	for i := 0; i < MAX_GENERATIONS && !state; i++ {
		//sort population in descending order by fitness
		sort.Slice(population, func(i int, j int) bool {
			return calcFitness(population[i]) < calcFitness(population[j])
		})
		summarize(population, i)
		population, state = simulateGeneration(population)
	}
	fmt.Println("Completed!")
}
