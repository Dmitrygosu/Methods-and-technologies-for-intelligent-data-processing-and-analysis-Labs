package experiment

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"lab1/ga"
)

type ParamGrid struct {
	PopulationSizes []int
	MaxGenerations  []int
	CrossoverProbs  []float64
	MutationProbs   []float64
	CrossoverTypes  []string
	ElitismCounts   []int
}

type ExperimentConfig struct {
	PopulationSize int     `json:"population_size"`
	MaxGenerations int     `json:"max_generations"`
	CrossoverProb  float64 `json:"crossover_prob"`
	MutationProb   float64 `json:"mutation_prob"`
	CrossoverType  string  `json:"crossover_type"`
	ElitismCount   int     `json:"elitism_count"`
}

type ExperimentResult struct {
	TaskName      string           `json:"task_name"`
	Config        ExperimentConfig `json:"config"`
	BestFitness   float64          `json:"best_fitness"`
	MeanFitness   float64          `json:"mean_fitness"`
	StdDevFitness float64          `json:"std_dev_fitness"`
	ExecutionTime float64          `json:"execution_time_ms"`
	AbsoluteError float64          `json:"absolute_error"`
	RelativeError float64          `json:"relative_error"`
	Convergence   []float64        `json:"convergence"`
}

type LinearSearchResult struct {
	TaskName      string  `json:"task_name"`
	BestValue     float64 `json:"best_value"`
	ExecutionTime float64 `json:"execution_time_ms"`
}

type AllResults struct {
	LinearSearchResults []LinearSearchResult `json:"linear_search_results"`
	GAResults           []ExperimentResult   `json:"ga_results"`
}

func (ar *AllResults) SaveToJSON(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(ar)
}

type ExperimentRunner struct {
	paramGrid ParamGrid
	arrayData []float64
}

func NewExperimentRunner(paramGrid ParamGrid) *ExperimentRunner {
	return &ExperimentRunner{
		paramGrid: paramGrid,
	}
}

func (er *ExperimentRunner) RunAllExperiments() (*AllResults, error) {
	results := &AllResults{
		LinearSearchResults: make([]LinearSearchResult, 0),
		GAResults:           make([]ExperimentResult, 0),
	}

	fmt.Println("Генерация массива с гауссовским распределением (1,000,000 элементов)...")
	er.arrayData = er.generateGaussianArray(1000000, 0.0, 100.0)

	fmt.Println("\n--- Задача 1: Поиск максимума в массиве ---")
	linearResult1 := er.runLinearSearchArray()
	results.LinearSearchResults = append(results.LinearSearchResults, linearResult1)
	fmt.Printf("Линейный поиск: значение=%.6f, время=%.2f мс\n",
		linearResult1.BestValue, linearResult1.ExecutionTime)

	fmt.Println("Запуск генетического алгоритма с различными конфигурациями...")
	gaResults1 := er.runGAForArray(linearResult1.BestValue)
	results.GAResults = append(results.GAResults, gaResults1...)
	fmt.Printf("Выполнено %d конфигураций для задачи 1\n", len(gaResults1))

	fmt.Println("\n--- Задача 2: Оптимизация математической функции ---")
	linearResult2 := er.runLinearSearchFunction()
	results.LinearSearchResults = append(results.LinearSearchResults, linearResult2)
	fmt.Printf("Линейный поиск: значение=%.6f, время=%.2f мс\n",
		linearResult2.BestValue, linearResult2.ExecutionTime)

	fmt.Println("Запуск генетического алгоритма с различными конфигурациями...")
	gaResults2 := er.runGAForFunction(linearResult2.BestValue)
	results.GAResults = append(results.GAResults, gaResults2...)
	fmt.Printf("Выполнено %d конфигураций для задачи 2\n", len(gaResults2))

	return results, nil
}

func (er *ExperimentRunner) generateGaussianArray(size int, mean, stddev float64) []float64 {
	rng := rand.New(rand.NewSource(42))
	arr := make([]float64, size)
	for i := 0; i < size; i++ {
		arr[i] = rng.NormFloat64()*stddev + mean
	}
	return arr
}

func (er *ExperimentRunner) runLinearSearchArray() LinearSearchResult {
	start := time.Now()

	maxVal := er.arrayData[0]
	for _, val := range er.arrayData {
		if val > maxVal {
			maxVal = val
		}
	}

	elapsed := time.Since(start)

	return LinearSearchResult{
		TaskName:      "array_search",
		BestValue:     maxVal,
		ExecutionTime: float64(elapsed.Milliseconds()),
	}
}

func (er *ExperimentRunner) runLinearSearchFunction() LinearSearchResult {
	start := time.Now()

	min, max := 2.7, 7.5
	steps := 1000000
	stepSize := (max - min) / float64(steps)

	maxVal := -math.MaxFloat64
	for i := 0; i <= steps; i++ {
		x := min + float64(i)*stepSize
		val := er.targetFunction(x)
		if val > maxVal {
			maxVal = val
		}
	}

	elapsed := time.Since(start)

	return LinearSearchResult{
		TaskName:      "function_optimization",
		BestValue:     maxVal,
		ExecutionTime: float64(elapsed.Milliseconds()),
	}
}

func (er *ExperimentRunner) targetFunction(x float64) float64 {
	return math.Sin(x) + math.Sin(10.0/3.0*x)
}

func (er *ExperimentRunner) runGAForArray(linearBest float64) []ExperimentResult {
	results := make([]ExperimentResult, 0)
	configs := er.generateConfigs()

	configNum := 0
	totalConfigs := len(configs)

	for _, config := range configs {
		configNum++
		if configNum%10 == 0 {
			fmt.Printf("Прогресс: %d/%d конфигураций\n", configNum, totalConfigs)
		}

		runs := 5
		fitnessValues := make([]float64, runs)
		var totalTime time.Duration
		var convergence []float64

		for run := 0; run < runs; run++ {
			gaConfig := ga.Config{
				PopulationSize: config.PopulationSize,
				MaxGenerations: config.MaxGenerations,
				CrossoverProb:  config.CrossoverProb,
				MutationProb:   config.MutationProb,
				CrossoverType:  config.CrossoverType,
				ElitismCount:   config.ElitismCount,
				BitsPerGene:    20,
				FitnessFunc:    er.arrayFitnessFunc(),
				Seed:           int64(time.Now().UnixNano() + int64(run)),
			}

			algorithm := ga.NewGeneticAlgorithm(gaConfig)

			start := time.Now()
			best, conv := algorithm.Run()
			elapsed := time.Since(start)

			fitnessValues[run] = best.Fitness
			totalTime += elapsed
			if run == 0 {
				convergence = conv
			}
		}

		meanFitness := 0.0
		for _, f := range fitnessValues {
			meanFitness += f
		}
		meanFitness /= float64(runs)

		stdDev := ga.StdDev(fitnessValues, meanFitness)

		bestFitness := fitnessValues[0]
		for _, f := range fitnessValues {
			if f > bestFitness {
				bestFitness = f
			}
		}

		absoluteError := linearBest - bestFitness
		relativeError := absoluteError / linearBest

		result := ExperimentResult{
			TaskName:      "array_search",
			Config:        config,
			BestFitness:   bestFitness,
			MeanFitness:   meanFitness,
			StdDevFitness: stdDev,
			ExecutionTime: float64(totalTime.Milliseconds()) / float64(runs),
			AbsoluteError: absoluteError,
			RelativeError: relativeError,
			Convergence:   convergence,
		}

		results = append(results, result)
	}

	return results
}

func (er *ExperimentRunner) runGAForFunction(linearBest float64) []ExperimentResult {
	results := make([]ExperimentResult, 0)

	configs := er.generateConfigs()

	configNum := 0
	totalConfigs := len(configs)

	for _, config := range configs {
		configNum++
		if configNum%10 == 0 {
			fmt.Printf("Прогресс: %d/%d конфигураций\n", configNum, totalConfigs)
		}

		runs := 5
		fitnessValues := make([]float64, runs)
		var totalTime time.Duration
		var convergence []float64

		for run := 0; run < runs; run++ {
			gaConfig := ga.Config{
				PopulationSize: config.PopulationSize,
				MaxGenerations: config.MaxGenerations,
				CrossoverProb:  config.CrossoverProb,
				MutationProb:   config.MutationProb,
				CrossoverType:  config.CrossoverType,
				ElitismCount:   config.ElitismCount,
				BitsPerGene:    16,
				FitnessFunc:    er.functionFitnessFunc(),
				Seed:           int64(time.Now().UnixNano() + int64(run)),
			}

			algorithm := ga.NewGeneticAlgorithm(gaConfig)

			start := time.Now()
			best, conv := algorithm.Run()
			elapsed := time.Since(start)

			fitnessValues[run] = best.Fitness
			totalTime += elapsed
			if run == 0 {
				convergence = conv
			}
		}

		meanFitness := 0.0
		for _, f := range fitnessValues {
			meanFitness += f
		}
		meanFitness /= float64(runs)

		stdDev := ga.StdDev(fitnessValues, meanFitness)

		bestFitness := fitnessValues[0]
		for _, f := range fitnessValues {
			if f > bestFitness {
				bestFitness = f
			}
		}

		absoluteError := linearBest - bestFitness
		relativeError := absoluteError / linearBest

		result := ExperimentResult{
			TaskName:      "function_optimization",
			Config:        config,
			BestFitness:   bestFitness,
			MeanFitness:   meanFitness,
			StdDevFitness: stdDev,
			ExecutionTime: float64(totalTime.Milliseconds()) / float64(runs),
			AbsoluteError: absoluteError,
			RelativeError: relativeError,
			Convergence:   convergence,
		}

		results = append(results, result)
	}

	return results
}

func (er *ExperimentRunner) arrayFitnessFunc() func([]byte) float64 {
	return func(genes []byte) float64 {
		index := ga.BytesToInt(genes) % len(er.arrayData)
		return er.arrayData[index]
	}
}

func (er *ExperimentRunner) functionFitnessFunc() func([]byte) float64 {
	return func(genes []byte) float64 {
		x := ga.BytesToFloat(genes, 2.7, 7.5)
		return er.targetFunction(x)
	}
}

func (er *ExperimentRunner) generateConfigs() []ExperimentConfig {
	configs := make([]ExperimentConfig, 0)

	for _, popSize := range er.paramGrid.PopulationSizes {
		for _, maxGen := range er.paramGrid.MaxGenerations {
			for _, crossProb := range er.paramGrid.CrossoverProbs {
				for _, mutProb := range er.paramGrid.MutationProbs {
					for _, crossType := range er.paramGrid.CrossoverTypes {
						for _, elitism := range er.paramGrid.ElitismCounts {
							configs = append(configs, ExperimentConfig{
								PopulationSize: popSize,
								MaxGenerations: maxGen,
								CrossoverProb:  crossProb,
								MutationProb:   mutProb,
								CrossoverType:  crossType,
								ElitismCount:   elitism,
							})
						}
					}
				}
			}
		}
	}

	return configs
}
