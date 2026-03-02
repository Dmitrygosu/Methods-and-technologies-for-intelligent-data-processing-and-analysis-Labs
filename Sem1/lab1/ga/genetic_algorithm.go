package ga

import (
	"math"
	"math/rand"
	"sort"
)

type Individual struct {
	Genes   []byte
	Fitness float64
}

type Config struct {
	PopulationSize int
	MaxGenerations int
	CrossoverProb  float64
	MutationProb   float64
	CrossoverType  string
	ElitismCount   int
	BitsPerGene    int
	FitnessFunc    func([]byte) float64
	Seed           int64
}

type GeneticAlgorithm struct {
	config      Config
	population  []Individual
	bestFitness []float64
	rng         *rand.Rand
}

func NewGeneticAlgorithm(config Config) *GeneticAlgorithm {
	return &GeneticAlgorithm{
		config:      config,
		bestFitness: make([]float64, 0),
		rng:         rand.New(rand.NewSource(config.Seed)),
	}
}

func (ga *GeneticAlgorithm) Initialize() {
	ga.population = make([]Individual, ga.config.PopulationSize)
	for i := 0; i < ga.config.PopulationSize; i++ {
		genes := make([]byte, ga.config.BitsPerGene)
		for j := 0; j < ga.config.BitsPerGene; j++ {
			if ga.rng.Float64() < 0.5 {
				genes[j] = 1
			} else {
				genes[j] = 0
			}
		}
		ga.population[i] = Individual{
			Genes:   genes,
			Fitness: ga.config.FitnessFunc(genes),
		}
	}
}

func (ga *GeneticAlgorithm) Run() (Individual, []float64) {
	ga.Initialize()

	for generation := 0; generation < ga.config.MaxGenerations; generation++ {
		sort.Slice(ga.population, func(i, j int) bool {
			return ga.population[i].Fitness > ga.population[j].Fitness
		})

		ga.bestFitness = append(ga.bestFitness, ga.population[0].Fitness)

		newPopulation := make([]Individual, 0, ga.config.PopulationSize)

		for i := 0; i < ga.config.ElitismCount && i < len(ga.population); i++ {
			newPopulation = append(newPopulation, ga.population[i])
		}

		for len(newPopulation) < ga.config.PopulationSize {
			parent1 := ga.tournamentSelection()
			parent2 := ga.tournamentSelection()

			var child1, child2 Individual
			if ga.rng.Float64() < ga.config.CrossoverProb {
				child1, child2 = ga.crossover(parent1, parent2)
			} else {
				child1 = parent1
				child2 = parent2
			}

			ga.mutate(&child1)
			ga.mutate(&child2)

			child1.Fitness = ga.config.FitnessFunc(child1.Genes)
			child2.Fitness = ga.config.FitnessFunc(child2.Genes)

			newPopulation = append(newPopulation, child1)
			if len(newPopulation) < ga.config.PopulationSize {
				newPopulation = append(newPopulation, child2)
			}
		}

		ga.population = newPopulation
	}

	sort.Slice(ga.population, func(i, j int) bool {
		return ga.population[i].Fitness > ga.population[j].Fitness
	})

	return ga.population[0], ga.bestFitness
}

func (ga *GeneticAlgorithm) tournamentSelection() Individual {
	tournamentSize := 3
	best := ga.population[ga.rng.Intn(len(ga.population))]

	for i := 1; i < tournamentSize; i++ {
		candidate := ga.population[ga.rng.Intn(len(ga.population))]
		if candidate.Fitness > best.Fitness {
			best = candidate
		}
	}

	return best
}

func (ga *GeneticAlgorithm) crossover(parent1, parent2 Individual) (Individual, Individual) {
	if ga.config.CrossoverType == "onepoint" {
		return ga.onepointCrossover(parent1, parent2)
	}
	return ga.uniformCrossover(parent1, parent2)
}

func (ga *GeneticAlgorithm) onepointCrossover(parent1, parent2 Individual) (Individual, Individual) {
	point := ga.rng.Intn(len(parent1.Genes))

	child1Genes := make([]byte, len(parent1.Genes))
	child2Genes := make([]byte, len(parent2.Genes))

	copy(child1Genes[:point], parent1.Genes[:point])
	copy(child1Genes[point:], parent2.Genes[point:])

	copy(child2Genes[:point], parent2.Genes[:point])
	copy(child2Genes[point:], parent1.Genes[point:])

	return Individual{Genes: child1Genes}, Individual{Genes: child2Genes}
}

func (ga *GeneticAlgorithm) uniformCrossover(parent1, parent2 Individual) (Individual, Individual) {
	child1Genes := make([]byte, len(parent1.Genes))
	child2Genes := make([]byte, len(parent2.Genes))

	for i := 0; i < len(parent1.Genes); i++ {
		if ga.rng.Float64() < 0.5 {
			child1Genes[i] = parent1.Genes[i]
			child2Genes[i] = parent2.Genes[i]
		} else {
			child1Genes[i] = parent2.Genes[i]
			child2Genes[i] = parent1.Genes[i]
		}
	}

	return Individual{Genes: child1Genes}, Individual{Genes: child2Genes}
}

func (ga *GeneticAlgorithm) mutate(individual *Individual) {
	for i := 0; i < len(individual.Genes); i++ {
		if ga.rng.Float64() < ga.config.MutationProb {
			if individual.Genes[i] == 0 {
				individual.Genes[i] = 1
			} else {
				individual.Genes[i] = 0
			}
		}
	}
}

func BytesToInt(genes []byte) int {
	result := 0
	for i := 0; i < len(genes); i++ {
		if genes[i] == 1 {
			result |= (1 << i)
		}
	}
	return result
}

func BytesToFloat(genes []byte, min, max float64) float64 {
	intVal := BytesToInt(genes)
	maxInt := (1 << len(genes)) - 1
	normalized := float64(intVal) / float64(maxInt)
	return min + normalized*(max-min)
}

func (ga *GeneticAlgorithm) GetBestFitnessHistory() []float64 {
	return ga.bestFitness
}

func StdDev(values []float64, mean float64) float64 {
	sum := 0.0
	for _, v := range values {
		sum += math.Pow(v-mean, 2)
	}
	return math.Sqrt(sum / float64(len(values)))
}
