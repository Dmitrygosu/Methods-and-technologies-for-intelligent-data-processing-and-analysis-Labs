package experiment

import (
	"encoding/json"
	"fmt"
	"lab4/clusterers"
	"lab4/data"
	"lab4/metrics"
	"math/rand"
	"os"
	"time"
)

type ParamGrid struct {
	KMeansK          []int
	KMeansMaxIter    []int
	DBSCANEps        []float64
	DBSCANMinSamples []int
	HierarchicalK    []int
	GMMComponents    []int
	GMMMaxIter       []int
}

type ClustererConfig struct {
	ModelType   string  `json:"model_type"`
	K           int     `json:"k,omitempty"`
	MaxIter     int     `json:"max_iter,omitempty"`
	Eps         float64 `json:"eps,omitempty"`
	MinSamples  int     `json:"min_samples,omitempty"`
	NComponents int     `json:"n_components,omitempty"`
}

type ClusteringResult struct {
	TaskName         string          `json:"task_name"`
	Config           ClustererConfig `json:"config"`
	SilhouetteScore  float64         `json:"silhouette_score"`
	CalinskiHarabasz float64         `json:"calinski_harabasz"`
	DaviesBouldin    float64         `json:"davies_bouldin"`
	NClusters        int             `json:"n_clusters"`
	ExecutionTime    float64         `json:"execution_time_ms"`
}

type AllResults struct {
	ClusteringResults []ClusteringResult `json:"clustering_results"`
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
}

func NewExperimentRunner(paramGrid ParamGrid) *ExperimentRunner {
	return &ExperimentRunner{
		paramGrid: paramGrid,
	}
}

func (er *ExperimentRunner) RunAllExperiments() (*AllResults, error) {
	results := &AllResults{
		ClusteringResults: make([]ClusteringResult, 0),
	}

	fmt.Println("Загрузка датасета о клиентах интернет-магазина...")
	er.dataset = data.NewCustomerDataset()
	fmt.Printf("Загружено %d экземпляров, %d признаков\n",
		len(er.dataset.Instances), er.dataset.NumFeatures)

	fmt.Println("\nНормализация признаков...")
	er.dataset.Normalize()

	fmt.Println("\n--- Задача: Кластеризация клиентов ---")

	clusteringResults := er.runClusteringExperiments(er.dataset)
	results.ClusteringResults = clusteringResults

	fmt.Printf("\nВыполнено %d конфигураций кластеризаторов\n", len(clusteringResults))

	return results, nil
}

func (er *ExperimentRunner) runClusteringExperiments(dataset *data.Dataset) []ClusteringResult {
	results := make([]ClusteringResult, 0)

	configs := er.generateClustererConfigs()
	configNum := 0
	totalConfigs := len(configs)

	for _, config := range configs {
		configNum++
		fmt.Printf("Прогресс кластеризации: %d/%d конфигураций\n", configNum, totalConfigs)

		runs := 5
		var totalTime time.Duration
		var silhouetteSum, calinskiSum, daviesSum float64
		var nClustersSum int

		for run := 0; run < runs; run++ {
			rand.Seed(int64(time.Now().UnixNano() + int64(run) + int64(configNum*1000)))

			var assignments []int
			start := time.Now()

			switch config.ModelType {
			case "kmeans":
				kmeans := clusterers.NewKMeans(config.K)
				kmeans.Train(dataset, config.MaxIter)
				assignments = kmeans.PredictBatch(dataset)

			case "dbscan":
				dbscan := clusterers.NewDBSCAN(config.Eps, config.MinSamples)
				dbscan.Train(dataset)
				assignments = dbscan.PredictBatch(dataset)

			case "hierarchical":
				hierarchical := clusterers.NewHierarchical(config.K)
				hierarchical.Train(dataset)
				assignments = hierarchical.PredictBatch(dataset)

			case "gaussian_mixture":
				gmm := clusterers.NewGaussianMixture(config.NComponents)
				gmm.Train(dataset, config.MaxIter)
				assignments = gmm.PredictBatch(dataset)
			}

			elapsed := time.Since(start)
			totalTime += elapsed

			instancesForMetrics := make([]struct {
				Features []float64
				Label    int
			}, len(dataset.Instances))
			for i, inst := range dataset.Instances {
				instancesForMetrics[i] = struct {
					Features []float64
					Label    int
				}{
					Features: inst.Features,
					Label:    0,
				}
			}

			uniqueClusters := make(map[int]bool)
			for _, ass := range assignments {
				if ass >= 0 {
					uniqueClusters[ass] = true
				}
			}
			nClusters := len(uniqueClusters)
			nClustersSum += nClusters

			if nClusters > 1 {
				clusteringMetrics := metrics.CalculateMetrics(instancesForMetrics, assignments, nClusters)
				silhouetteSum += clusteringMetrics.SilhouetteScore
				calinskiSum += clusteringMetrics.CalinskiHarabasz
				daviesSum += clusteringMetrics.DaviesBouldin
			}
		}

		result := ClusteringResult{
			TaskName:         "clustering",
			Config:           config,
			SilhouetteScore:  silhouetteSum / float64(runs),
			CalinskiHarabasz: calinskiSum / float64(runs),
			DaviesBouldin:    daviesSum / float64(runs),
			NClusters:        nClustersSum / runs,
			ExecutionTime:    float64(totalTime.Milliseconds()) / float64(runs),
		}

		results = append(results, result)
	}

	return results
}

func (er *ExperimentRunner) generateClustererConfigs() []ClustererConfig {
	configs := make([]ClustererConfig, 0)

	for _, k := range er.paramGrid.KMeansK {
		for _, maxIter := range er.paramGrid.KMeansMaxIter {
			configs = append(configs, ClustererConfig{
				ModelType: "kmeans",
				K:         k,
				MaxIter:   maxIter,
			})
		}
	}

	for _, eps := range er.paramGrid.DBSCANEps {
		for _, minSamples := range er.paramGrid.DBSCANMinSamples {
			configs = append(configs, ClustererConfig{
				ModelType:  "dbscan",
				Eps:        eps,
				MinSamples: minSamples,
			})
		}
	}

	for _, k := range er.paramGrid.HierarchicalK {
		configs = append(configs, ClustererConfig{
			ModelType: "hierarchical",
			K:         k,
		})
	}

	for _, nComponents := range er.paramGrid.GMMComponents {
		for _, maxIter := range er.paramGrid.GMMMaxIter {
			configs = append(configs, ClustererConfig{
				ModelType:   "gaussian_mixture",
				NComponents: nComponents,
				MaxIter:     maxIter,
			})
		}
	}

	return configs
}
