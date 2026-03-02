package experiment

import (
	"encoding/json"
	"fmt"
	"lab2/classifiers"
	"lab2/clusterers"
	"lab2/data"
	"lab2/metrics"
	"math/rand"
	"os"
	"time"
)

type ParamGrid struct {
	KNNK           []int
	DTMaxDepth     []int
	LRLearningRate []float64
	LRIterations   []int
	KMeansK        []int
	KMeansMaxIter  []int
}

type ClassifierConfig struct {
	ModelType    string  `json:"model_type"`
	K            int     `json:"k,omitempty"`
	MaxDepth     int     `json:"max_depth,omitempty"`
	LearningRate float64 `json:"learning_rate,omitempty"`
	Iterations   int     `json:"iterations,omitempty"`
}

type ClustererConfig struct {
	K       int `json:"k"`
	MaxIter int `json:"max_iter"`
}

type ClassificationResult struct {
	TaskName        string           `json:"task_name"`
	Config          ClassifierConfig `json:"config"`
	Accuracy        float64          `json:"accuracy"`
	MacroF1         float64          `json:"macro_f1"`
	WeightedF1      float64          `json:"weighted_f1"`
	Precision       []float64        `json:"precision"`
	Recall          []float64        `json:"recall"`
	ConfusionMatrix [][]int          `json:"confusion_matrix"`
	ExecutionTime   float64          `json:"execution_time_ms"`
}

type ClusteringResult struct {
	TaskName        string          `json:"task_name"`
	Config          ClustererConfig `json:"config"`
	SilhouetteScore float64         `json:"silhouette_score"`
	ExecutionTime   float64         `json:"execution_time_ms"`
}

type AllResults struct {
	ClassificationResults []ClassificationResult `json:"classification_results"`
	ClusteringResults     []ClusteringResult     `json:"clustering_results"`
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
	dataset   *data.Dataset
	trainSet  *data.Dataset
	testSet   *data.Dataset
}

func NewExperimentRunner(paramGrid ParamGrid) *ExperimentRunner {
	return &ExperimentRunner{
		paramGrid: paramGrid,
	}
}

func (er *ExperimentRunner) RunAllExperiments() (*AllResults, error) {
	results := &AllResults{
		ClassificationResults: make([]ClassificationResult, 0),
		ClusteringResults:     make([]ClusteringResult, 0),
	}

	fmt.Println("Загрузка датасета Iris...")
	er.dataset = data.NewIrisDataset()
	fmt.Printf("Загружено %d экземпляров, %d признаков, %d классов\n",
		len(er.dataset.Instances), er.dataset.NumFeatures, er.dataset.NumClasses)

	fmt.Println("\nРазделение на обучающую и тестовую выборки (80/20)...")
	er.trainSet, er.testSet = er.dataset.Split(0.8, 42)

	trainCopy := &data.Dataset{
		Instances:   make([]data.Instance, len(er.trainSet.Instances)),
		NumFeatures: er.trainSet.NumFeatures,
		NumClasses:  er.trainSet.NumClasses,
	}
	copy(trainCopy.Instances, er.trainSet.Instances)
	trainCopy.Normalize()

	testCopy := &data.Dataset{
		Instances:   make([]data.Instance, len(er.testSet.Instances)),
		NumFeatures: er.testSet.NumFeatures,
		NumClasses:  er.testSet.NumClasses,
	}
	copy(testCopy.Instances, er.testSet.Instances)
	testCopy.Normalize()

	fmt.Printf("Обучающая выборка: %d экземпляров\n", len(trainCopy.Instances))
	fmt.Printf("Тестовая выборка: %d экземпляров\n", len(testCopy.Instances))

	fmt.Println("\n--- Задача 1: Классификация ---")
	classResults := er.runClassificationExperiments(trainCopy, testCopy)
	results.ClassificationResults = append(results.ClassificationResults, classResults...)
	fmt.Printf("Выполнено %d конфигураций классификаторов\n", len(classResults))

	fmt.Println("\n--- Задача 2: Кластеризация ---")
	clusterResults := er.runClusteringExperiments(trainCopy)
	results.ClusteringResults = append(results.ClusteringResults, clusterResults...)
	fmt.Printf("Выполнено %d конфигураций кластеризации\n", len(clusterResults))

	return results, nil
}

func (er *ExperimentRunner) runClassificationExperiments(trainSet, testSet *data.Dataset) []ClassificationResult {
	results := make([]ClassificationResult, 0)

	yTrue := make([]int, len(testSet.Instances))
	for i, inst := range testSet.Instances {
		yTrue[i] = inst.Label
	}

	configs := er.generateClassifierConfigs()
	configNum := 0
	totalConfigs := len(configs)

	for _, config := range configs {
		configNum++
		if configNum%5 == 0 {
			fmt.Printf("Прогресс классификации: %d/%d конфигураций\n", configNum, totalConfigs)
		}

		runs := 5
		var totalTime time.Duration
		var accuracySum, macroF1Sum, weightedF1Sum float64
		var precisionSum, recallSum []float64
		var confusionMatrix [][]int

		for run := 0; run < runs; run++ {
			rand.Seed(int64(time.Now().UnixNano() + int64(run)))

			var classifier interface {
				Train(*data.Dataset)
				PredictBatch(*data.Dataset) []int
			}

			switch config.ModelType {
			case "knn":
				classifier = classifiers.NewKNN(config.K)
			case "naive_bayes":
				classifier = classifiers.NewNaiveBayes()
			case "decision_tree":
				classifier = classifiers.NewDecisionTree(config.MaxDepth)
			case "logistic_regression":
				classifier = classifiers.NewLogisticRegression(config.LearningRate, config.Iterations)
			default:
				continue
			}

			start := time.Now()
			classifier.Train(trainSet)
			yPred := classifier.PredictBatch(testSet)
			elapsed := time.Since(start)

			totalTime += elapsed

			m := metrics.CalculateMetrics(yTrue, yPred, trainSet.NumClasses)

			if run == 0 {
				precisionSum = make([]float64, len(m.Precision))
				recallSum = make([]float64, len(m.Recall))
				confusionMatrix = metrics.BuildConfusionMatrix(yTrue, yPred, trainSet.NumClasses)
			}

			accuracySum += m.Accuracy
			macroF1Sum += m.MacroF1
			weightedF1Sum += m.WeightedF1

			for i := range m.Precision {
				precisionSum[i] += m.Precision[i]
				recallSum[i] += m.Recall[i]
			}
		}

		avgPrecision := make([]float64, len(precisionSum))
		avgRecall := make([]float64, len(recallSum))
		for i := range precisionSum {
			avgPrecision[i] = precisionSum[i] / float64(runs)
			avgRecall[i] = recallSum[i] / float64(runs)
		}

		result := ClassificationResult{
			TaskName:        "classification",
			Config:          config,
			Accuracy:        accuracySum / float64(runs),
			MacroF1:         macroF1Sum / float64(runs),
			WeightedF1:      weightedF1Sum / float64(runs),
			Precision:       avgPrecision,
			Recall:          avgRecall,
			ConfusionMatrix: confusionMatrix,
			ExecutionTime:   float64(totalTime.Milliseconds()) / float64(runs),
		}

		results = append(results, result)
	}

	return results
}

func (er *ExperimentRunner) runClusteringExperiments(trainSet *data.Dataset) []ClusteringResult {
	results := make([]ClusteringResult, 0)

	configs := er.generateClustererConfigs()
	configNum := 0
	totalConfigs := len(configs)

	for _, config := range configs {
		configNum++
		fmt.Printf("Прогресс кластеризации: %d/%d конфигураций\n", configNum, totalConfigs)

		runs := 5
		var totalTime time.Duration
		var silhouetteSum float64

		for run := 0; run < runs; run++ {
			rand.Seed(int64(time.Now().UnixNano() + int64(run)))

			kmeans := clusterers.NewKMeans(config.K)

			start := time.Now()
			kmeans.Train(trainSet, config.MaxIter)
			assignments := kmeans.PredictBatch(trainSet)
			elapsed := time.Since(start)

			totalTime += elapsed

			instancesForSil := make([]struct {
				Features []float64
				Label    int
			}, len(trainSet.Instances))
			for i, inst := range trainSet.Instances {
				instancesForSil[i] = struct {
					Features []float64
					Label    int
				}{
					Features: inst.Features,
					Label:    inst.Label,
				}
			}

			silhouette := metrics.CalculateSilhouetteScore(instancesForSil, assignments, config.K)
			silhouetteSum += silhouette
		}

		result := ClusteringResult{
			TaskName:        "clustering",
			Config:          config,
			SilhouetteScore: silhouetteSum / float64(runs),
			ExecutionTime:   float64(totalTime.Milliseconds()) / float64(runs),
		}

		results = append(results, result)
	}

	return results
}

func (er *ExperimentRunner) generateClassifierConfigs() []ClassifierConfig {
	configs := make([]ClassifierConfig, 0)

	for _, k := range er.paramGrid.KNNK {
		configs = append(configs, ClassifierConfig{
			ModelType: "knn",
			K:         k,
		})
	}

	configs = append(configs, ClassifierConfig{
		ModelType: "naive_bayes",
	})

	for _, maxDepth := range er.paramGrid.DTMaxDepth {
		configs = append(configs, ClassifierConfig{
			ModelType: "decision_tree",
			MaxDepth:  maxDepth,
		})
	}

	for _, lr := range er.paramGrid.LRLearningRate {
		for _, iters := range er.paramGrid.LRIterations {
			configs = append(configs, ClassifierConfig{
				ModelType:    "logistic_regression",
				LearningRate: lr,
				Iterations:   iters,
			})
		}
	}

	return configs
}

func (er *ExperimentRunner) generateClustererConfigs() []ClustererConfig {
	configs := make([]ClustererConfig, 0)

	for _, k := range er.paramGrid.KMeansK {
		for _, maxIter := range er.paramGrid.KMeansMaxIter {
			configs = append(configs, ClustererConfig{
				K:       k,
				MaxIter: maxIter,
			})
		}
	}

	return configs
}
