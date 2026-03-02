package experiment

import (
	"encoding/json"
	"fmt"
	"lab3/data"
	"lab3/metrics"
	"lab3/regressors"
	"math/rand"
	"os"
	"time"
)

type ParamGrid struct {
	RidgeAlpha        []float64
	LassoAlpha        []float64
	DTMaxDepth        []int
	RFNumTrees        []int
	RFMaxDepth        []int
	RFMaxFeatures     []int
	GBNumTrees        []int
	GBMaxDepth        []int
	GBLearningRate    []float64
}

type RegressorConfig struct {
	ModelType    string  `json:"model_type"`
	Alpha        float64 `json:"alpha,omitempty"`
	MaxDepth     int     `json:"max_depth,omitempty"`
	NumTrees     int     `json:"num_trees,omitempty"`
	MaxFeatures  int     `json:"max_features,omitempty"`
	LearningRate float64 `json:"learning_rate,omitempty"`
}

type RegressionResult struct {
	TaskName      string          `json:"task_name"`
	Config        RegressorConfig `json:"config"`
	R2            float64         `json:"r2"`
	RMSE          float64         `json:"rmse"`
	MAE           float64         `json:"mae"`
	MAPE          float64         `json:"mape"`
	ExecutionTime float64         `json:"execution_time_ms"`
}

type AllResults struct {
	RegressionResults []RegressionResult `json:"regression_results"`
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
	dataset   *data.RegressionDataset
	trainSet  *data.RegressionDataset
	testSet   *data.RegressionDataset
}

func NewExperimentRunner(paramGrid ParamGrid) *ExperimentRunner {
	return &ExperimentRunner{
		paramGrid: paramGrid,
	}
}

func (er *ExperimentRunner) RunAllExperiments() (*AllResults, error) {
	results := &AllResults{
		RegressionResults: make([]RegressionResult, 0),
	}

	fmt.Println("Загрузка датасета недвижимости...")
	er.dataset = data.NewHousingDataset()
	fmt.Printf("Загружено %d экземпляров, %d признаков\n",
		len(er.dataset.Instances), er.dataset.NumFeatures)

	fmt.Println("\nРазделение на обучающую и тестовую выборки (80/20)...")
	er.trainSet, er.testSet = er.dataset.Split(0.8, 42)

	trainCopy := &data.RegressionDataset{
		Instances:   make([]data.RegressionInstance, len(er.trainSet.Instances)),
		NumFeatures: er.trainSet.NumFeatures,
	}
	copy(trainCopy.Instances, er.trainSet.Instances)
	trainCopy.Normalize()

	testCopy := &data.RegressionDataset{
		Instances:   make([]data.RegressionInstance, len(er.testSet.Instances)),
		NumFeatures: er.testSet.NumFeatures,
	}
	copy(testCopy.Instances, er.testSet.Instances)
	testCopy.Normalize()

	fmt.Printf("Обучающая выборка: %d экземпляров\n", len(trainCopy.Instances))
	fmt.Printf("Тестовая выборка: %d экземпляров\n", len(testCopy.Instances))

	fmt.Println("\n--- Задача: Регрессия (прогнозирование цены) ---")
	regResults := er.runRegressionExperiments(trainCopy, testCopy)
	results.RegressionResults = append(results.RegressionResults, regResults...)
	fmt.Printf("Выполнено %d конфигураций регрессоров\n", len(regResults))

	return results, nil
}

func (er *ExperimentRunner) runRegressionExperiments(trainSet, testSet *data.RegressionDataset) []RegressionResult {
	results := make([]RegressionResult, 0)

	yTrue := make([]float64, len(testSet.Instances))
	for i, inst := range testSet.Instances {
		yTrue[i] = inst.Target
	}

	configs := er.generateRegressorConfigs()
	configNum := 0
	totalConfigs := len(configs)

	for _, config := range configs {
		configNum++
		if configNum%5 == 0 {
			fmt.Printf("Прогресс регрессии: %d/%d конфигураций\n", configNum, totalConfigs)
		}

		runs := 5
		var totalTime time.Duration
		var r2Sum, rmseSum, maeSum, mapeSum float64

		for run := 0; run < runs; run++ {
			rand.Seed(int64(time.Now().UnixNano() + int64(run)))

			var regressor interface {
				Train(*data.RegressionDataset)
				PredictBatch(*data.RegressionDataset) []float64
			}

			switch config.ModelType {
			case "linear":
				regressor = regressors.NewLinearRegression()
			case "ridge":
				regressor = regressors.NewRidgeRegression(config.Alpha)
			case "lasso":
				regressor = regressors.NewLassoRegression(config.Alpha)
			case "decision_tree":
				regressor = regressors.NewDecisionTreeRegressor(config.MaxDepth)
			case "random_forest":
				regressor = regressors.NewRandomForestRegressor(config.NumTrees, config.MaxDepth, config.MaxFeatures)
			case "gradient_boosting":
				regressor = regressors.NewGradientBoostingRegressor(config.NumTrees, config.MaxDepth, config.LearningRate)
			default:
				continue
			}

			start := time.Now()
			regressor.Train(trainSet)
			yPred := regressor.PredictBatch(testSet)
			elapsed := time.Since(start)

			totalTime += elapsed

			m := metrics.CalculateMetrics(yTrue, yPred)

			r2Sum += m.R2
			rmseSum += m.RMSE
			maeSum += m.MAE
			mapeSum += m.MAPE
		}

		result := RegressionResult{
			TaskName:      "regression",
			Config:        config,
			R2:            r2Sum / float64(runs),
			RMSE:          rmseSum / float64(runs),
			MAE:           maeSum / float64(runs),
			MAPE:          mapeSum / float64(runs),
			ExecutionTime: float64(totalTime.Milliseconds()) / float64(runs),
		}

		results = append(results, result)
	}

	return results
}

func (er *ExperimentRunner) generateRegressorConfigs() []RegressorConfig {
	configs := make([]RegressorConfig, 0)

	configs = append(configs, RegressorConfig{
		ModelType: "linear",
	})

	for _, alpha := range er.paramGrid.RidgeAlpha {
		configs = append(configs, RegressorConfig{
			ModelType: "ridge",
			Alpha:     alpha,
		})
	}

	for _, alpha := range er.paramGrid.LassoAlpha {
		configs = append(configs, RegressorConfig{
			ModelType: "lasso",
			Alpha:     alpha,
		})
	}

	for _, maxDepth := range er.paramGrid.DTMaxDepth {
		configs = append(configs, RegressorConfig{
			ModelType: "decision_tree",
			MaxDepth:  maxDepth,
		})
	}

	for _, numTrees := range er.paramGrid.RFNumTrees {
		for _, maxDepth := range er.paramGrid.RFMaxDepth {
			for _, maxFeatures := range er.paramGrid.RFMaxFeatures {
				configs = append(configs, RegressorConfig{
					ModelType:   "random_forest",
					NumTrees:    numTrees,
					MaxDepth:    maxDepth,
					MaxFeatures: maxFeatures,
				})
			}
		}
	}

	for _, numTrees := range er.paramGrid.GBNumTrees {
		for _, maxDepth := range er.paramGrid.GBMaxDepth {
			for _, learningRate := range er.paramGrid.GBLearningRate {
				configs = append(configs, RegressorConfig{
					ModelType:    "gradient_boosting",
					NumTrees:     numTrees,
					MaxDepth:     maxDepth,
					LearningRate: learningRate,
				})
			}
		}
	}

	return configs
}

