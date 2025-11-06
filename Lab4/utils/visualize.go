package utils

import (
	"encoding/json"
	"fmt"
	"image/color"
	"os"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

type ClusteringResult struct {
	TaskName         string          `json:"task_name"`
	Config           ClustererConfig `json:"config"`
	SilhouetteScore  float64         `json:"silhouette_score"`
	CalinskiHarabasz float64         `json:"calinski_harabasz"`
	DaviesBouldin    float64         `json:"davies_bouldin"`
	NClusters        int             `json:"n_clusters"`
	ExecutionTime    float64         `json:"execution_time_ms"`
}

type ClustererConfig struct {
	ModelType   string  `json:"model_type"`
	K           int     `json:"k,omitempty"`
	MaxIter     int     `json:"max_iter,omitempty"`
	Eps         float64 `json:"eps,omitempty"`
	MinSamples  int     `json:"min_samples,omitempty"`
	NComponents int     `json:"n_components,omitempty"`
}

type AllResults struct {
	ClusteringResults []ClusteringResult `json:"clustering_results"`
}

func loadResults(filename string) (*AllResults, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var results AllResults
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&results)
	if err != nil {
		return nil, err
	}

	return &results, nil
}

func GenerateSilhouetteComparisonPlot(resultsFile, outputFile string) error {
	results, err := loadResults(resultsFile)
	if err != nil {
		return err
	}

	p := plot.New()
	p.Title.Text = "СРАВНЕНИЕ КАЧЕСТВА КЛАСТЕРИЗАЦИИ\n(Индекс силуэта)"
	p.Title.TextStyle.Font.Size = 14
	p.Title.TextStyle.Color = color.RGBA{R: 0, G: 100, B: 0, A: 255}
	p.Y.Label.Text = "Индекс силуэта"
	p.Y.Label.TextStyle.Font.Size = 14
	p.X.Label.Text = "Алгоритм кластеризации"
	p.X.Label.TextStyle.Font.Size = 14

	modelGroups := make(map[string][]float64)
	modelLabels := make([]string, 0)
	modelShortLabels := make([]string, 0)
	modelMap := make(map[string]bool)

	for _, r := range results.ClusteringResults {
		label := getModelLabel(r.Config)
		if !modelMap[label] {
			modelLabels = append(modelLabels, label)
			modelShortLabels = append(modelShortLabels, getShortModelLabel(r.Config))
			modelMap[label] = true
			modelGroups[label] = make([]float64, 0)
		}
		modelGroups[label] = append(modelGroups[label], r.SilhouetteScore)
	}

	avgValues := make(plotter.Values, len(modelLabels))
	for i, label := range modelLabels {
		sum := 0.0
		for _, v := range modelGroups[label] {
			sum += v
		}
		avgValues[i] = sum / float64(len(modelGroups[label]))
	}

	width := vg.Points(40)
	b, err := plotter.NewBarChart(avgValues, width)
	if err != nil {
		return err
	}

	b.Color = color.RGBA{R: 70, G: 130, B: 180, A: 255}
	b.LineStyle.Color = color.Black
	b.LineStyle.Width = vg.Points(1)

	p.NominalX(modelShortLabels...)
	p.X.Tick.Label.Rotation = 0.785
	p.X.Tick.Label.Font.Size = 10
	p.X.Padding = vg.Points(40)

	p.Add(b)
	p.Legend.Add("Средний индекс силуэта", b)

	if err := p.Save(12*vg.Inch, 6*vg.Inch, outputFile); err != nil {
		return err
	}

	return nil
}

func GenerateCalinskiComparisonPlot(resultsFile, outputFile string) error {
	results, err := loadResults(resultsFile)
	if err != nil {
		return err
	}

	p := plot.New()
	p.Title.Text = "СРАВНЕНИЕ КАЧЕСТВА КЛАСТЕРИЗАЦИИ\n(Индекс Калинского-Харабаша)"
	p.Title.TextStyle.Font.Size = 14
	p.Title.TextStyle.Color = color.RGBA{R: 0, G: 100, B: 0, A: 255}
	p.Y.Label.Text = "Индекс Калинского-Харабаша"
	p.Y.Label.TextStyle.Font.Size = 14
	p.X.Label.Text = "Алгоритм кластеризации"
	p.X.Label.TextStyle.Font.Size = 14

	modelGroups := make(map[string][]float64)
	modelLabels := make([]string, 0)
	modelShortLabels := make([]string, 0)
	modelMap := make(map[string]bool)

	for _, r := range results.ClusteringResults {
		label := getModelLabel(r.Config)
		if !modelMap[label] {
			modelLabels = append(modelLabels, label)
			modelShortLabels = append(modelShortLabels, getShortModelLabel(r.Config))
			modelMap[label] = true
			modelGroups[label] = make([]float64, 0)
		}
		modelGroups[label] = append(modelGroups[label], r.CalinskiHarabasz)
	}

	avgValues := make(plotter.Values, len(modelLabels))
	for i, label := range modelLabels {
		sum := 0.0
		for _, v := range modelGroups[label] {
			sum += v
		}
		avgValues[i] = sum / float64(len(modelGroups[label]))
	}

	width := vg.Points(40)
	b, err := plotter.NewBarChart(avgValues, width)
	if err != nil {
		return err
	}

	b.Color = color.RGBA{R: 50, G: 200, B: 50, A: 255}
	b.LineStyle.Color = color.Black
	b.LineStyle.Width = vg.Points(1)

	p.NominalX(modelShortLabels...)
	p.X.Tick.Label.Rotation = 0.785
	p.X.Tick.Label.Font.Size = 10
	p.X.Padding = vg.Points(40)

	p.Add(b)
	p.Legend.Add("Средний индекс Калинского-Харабаша", b)

	if err := p.Save(12*vg.Inch, 6*vg.Inch, outputFile); err != nil {
		return err
	}

	return nil
}

func GenerateDaviesComparisonPlot(resultsFile, outputFile string) error {
	results, err := loadResults(resultsFile)
	if err != nil {
		return err
	}

	p := plot.New()
	p.Title.Text = "СРАВНЕНИЕ КАЧЕСТВА КЛАСТЕРИЗАЦИИ\n(Индекс Дэвиса-Болдуина, чем меньше - тем лучше)"
	p.Title.TextStyle.Font.Size = 14
	p.Title.TextStyle.Color = color.RGBA{R: 0, G: 100, B: 0, A: 255}
	p.Y.Label.Text = "Индекс Дэвиса-Болдуина"
	p.Y.Label.TextStyle.Font.Size = 14
	p.X.Label.Text = "Алгоритм кластеризации"
	p.X.Label.TextStyle.Font.Size = 14

	modelGroups := make(map[string][]float64)
	modelLabels := make([]string, 0)
	modelShortLabels := make([]string, 0)
	modelMap := make(map[string]bool)

	for _, r := range results.ClusteringResults {
		label := getModelLabel(r.Config)
		if !modelMap[label] {
			modelLabels = append(modelLabels, label)
			modelShortLabels = append(modelShortLabels, getShortModelLabel(r.Config))
			modelMap[label] = true
			modelGroups[label] = make([]float64, 0)
		}
		modelGroups[label] = append(modelGroups[label], r.DaviesBouldin)
	}

	avgValues := make(plotter.Values, len(modelLabels))
	for i, label := range modelLabels {
		sum := 0.0
		for _, v := range modelGroups[label] {
			sum += v
		}
		avgValues[i] = sum / float64(len(modelGroups[label]))
	}

	width := vg.Points(40)
	b, err := plotter.NewBarChart(avgValues, width)
	if err != nil {
		return err
	}

	b.Color = color.RGBA{R: 220, G: 20, B: 60, A: 255}
	b.LineStyle.Color = color.Black
	b.LineStyle.Width = vg.Points(1)

	p.NominalX(modelShortLabels...)
	p.X.Tick.Label.Rotation = 0.785
	p.X.Tick.Label.Font.Size = 10
	p.X.Padding = vg.Points(40)

	p.Add(b)
	p.Legend.Add("Средний индекс Дэвиса-Болдуина", b)

	if err := p.Save(12*vg.Inch, 6*vg.Inch, outputFile); err != nil {
		return err
	}

	return nil
}

func GenerateClusteringComparisonPlot(resultsFile, outputFile string) error {
	results, err := loadResults(resultsFile)
	if err != nil {
		return err
	}

	p := plot.New()
	p.Title.Text = "СРАВНЕНИЕ КАЧЕСТВА КЛАСТЕРИЗАЦИИ\n(Индекс силуэта)"
	p.Title.TextStyle.Font.Size = 14
	p.Title.TextStyle.Color = color.RGBA{R: 0, G: 100, B: 0, A: 255}
	p.Y.Label.Text = "Индекс силуэта"
	p.Y.Label.TextStyle.Font.Size = 14
	p.X.Label.Text = "Алгоритм кластеризации"
	p.X.Label.TextStyle.Font.Size = 14

	modelGroups := make(map[string][]float64)
	modelLabels := make([]string, 0)
	modelShortLabels := make([]string, 0)
	modelMap := make(map[string]bool)

	for _, r := range results.ClusteringResults {
		label := getModelLabel(r.Config)
		if !modelMap[label] {
			modelLabels = append(modelLabels, label)
			modelShortLabels = append(modelShortLabels, getShortModelLabel(r.Config))
			modelMap[label] = true
			modelGroups[label] = make([]float64, 0)
		}
		modelGroups[label] = append(modelGroups[label], r.SilhouetteScore)
	}

	avgValues := make(plotter.Values, len(modelLabels))
	for i, label := range modelLabels {
		sum := 0.0
		for _, v := range modelGroups[label] {
			sum += v
		}
		avgValues[i] = sum / float64(len(modelGroups[label]))
	}

	width := vg.Points(40)
	b, err := plotter.NewBarChart(avgValues, width)
	if err != nil {
		return err
	}

	b.Color = color.RGBA{R: 70, G: 130, B: 180, A: 255}
	b.LineStyle.Color = color.Black
	b.LineStyle.Width = vg.Points(1)

	p.NominalX(modelShortLabels...)
	p.X.Tick.Label.Rotation = 0.785
	p.X.Tick.Label.Font.Size = 10
	p.X.Padding = vg.Points(40)

	p.Add(b)
	p.Legend.Add("Средний индекс силуэта", b)

	if err := p.Save(12*vg.Inch, 6*vg.Inch, outputFile); err != nil {
		return err
	}

	return nil
}

func getModelLabel(config ClustererConfig) string {
	switch config.ModelType {
	case "kmeans":
		return fmt.Sprintf("K-Means (k=%d, max_iter=%d)", config.K, config.MaxIter)
	case "dbscan":
		return fmt.Sprintf("DBSCAN (eps=%.2f, min_samples=%d)", config.Eps, config.MinSamples)
	case "hierarchical":
		return fmt.Sprintf("Иерархическая (k=%d)", config.K)
	case "gaussian_mixture":
		return fmt.Sprintf("Гауссова смесь (компонентов=%d, max_iter=%d)", config.NComponents, config.MaxIter)
	default:
		return config.ModelType
	}
}

func getShortModelLabel(config ClustererConfig) string {
	switch config.ModelType {
	case "kmeans":
		return fmt.Sprintf("K-Means\n(k=%d)", config.K)
	case "dbscan":
		return fmt.Sprintf("DBSCAN\n(eps=%.2f)", config.Eps)
	case "hierarchical":
		return fmt.Sprintf("Иерарх.\n(k=%d)", config.K)
	case "gaussian_mixture":
		return fmt.Sprintf("GMM\n(c=%d)", config.NComponents)
	default:
		return config.ModelType
	}
}
