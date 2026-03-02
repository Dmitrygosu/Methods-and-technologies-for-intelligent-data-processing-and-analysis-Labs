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

type RegressionResult struct {
	TaskName      string          `json:"task_name"`
	Config        RegressorConfig `json:"config"`
	R2            float64         `json:"r2"`
	RMSE          float64         `json:"rmse"`
	MAE           float64         `json:"mae"`
	MAPE          float64         `json:"mape"`
	ExecutionTime float64         `json:"execution_time_ms"`
}

type RegressorConfig struct {
	ModelType    string  `json:"model_type"`
	Alpha        float64 `json:"alpha,omitempty"`
	MaxDepth     int     `json:"max_depth,omitempty"`
	NumTrees     int     `json:"num_trees,omitempty"`
	MaxFeatures  int     `json:"max_features,omitempty"`
	LearningRate float64 `json:"learning_rate,omitempty"`
}

type AllResults struct {
	RegressionResults []RegressionResult `json:"regression_results"`
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

func GenerateR2ComparisonPlot(resultsFile, outputFile string) error {
	results, err := loadResults(resultsFile)
	if err != nil {
		return err
	}

	p := plot.New()
	p.Title.Text = "СРАВНЕНИЕ КАЧЕСТВА РЕГРЕССОРОВ\n(Коэффициент детерминации R²)"
	p.Title.TextStyle.Font.Size = 14
	p.Title.TextStyle.Color = color.RGBA{R: 0, G: 100, B: 0, A: 255}
	p.Y.Label.Text = "Коэффициент детерминации R²"
	p.Y.Label.TextStyle.Font.Size = 14
	p.X.Label.Text = "Модель регрессии"
	p.X.Label.TextStyle.Font.Size = 14

	modelGroups := make(map[string][]float64)
	modelLabels := make([]string, 0)
	modelShortLabels := make([]string, 0)
	modelMap := make(map[string]bool)

	for _, r := range results.RegressionResults {
		label := getModelLabel(r.Config)
		if !modelMap[label] {
			modelLabels = append(modelLabels, label)
			modelShortLabels = append(modelShortLabels, getShortModelLabel(r.Config))
			modelMap[label] = true
			modelGroups[label] = make([]float64, 0)
		}
		modelGroups[label] = append(modelGroups[label], r.R2)
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
	p.Legend.Add("Средний R²", b)

	if err := p.Save(12*vg.Inch, 6*vg.Inch, outputFile); err != nil {
		return err
	}

	return nil
}

func GenerateRMSEComparisonPlot(resultsFile, outputFile string) error {
	results, err := loadResults(resultsFile)
	if err != nil {
		return err
	}

	p := plot.New()
	p.Title.Text = "СРАВНЕНИЕ ОШИБОК РЕГРЕССОРОВ\n(Среднеквадратичная ошибка RMSE)"
	p.Title.TextStyle.Font.Size = 14
	p.Title.TextStyle.Color = color.RGBA{R: 0, G: 100, B: 0, A: 255}
	p.Y.Label.Text = "RMSE (тыс. руб.)"
	p.Y.Label.TextStyle.Font.Size = 14
	p.X.Label.Text = "Модель регрессии"
	p.X.Label.TextStyle.Font.Size = 14

	modelGroups := make(map[string][]float64)
	modelLabels := make([]string, 0)
	modelShortLabels := make([]string, 0)
	modelMap := make(map[string]bool)

	for _, r := range results.RegressionResults {
		label := getModelLabel(r.Config)
		if !modelMap[label] {
			modelLabels = append(modelLabels, label)
			modelShortLabels = append(modelShortLabels, getShortModelLabel(r.Config))
			modelMap[label] = true
			modelGroups[label] = make([]float64, 0)
		}
		modelGroups[label] = append(modelGroups[label], r.RMSE)
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
	p.Legend.Add("Средний RMSE", b)

	if err := p.Save(12*vg.Inch, 6*vg.Inch, outputFile); err != nil {
		return err
	}

	return nil
}

func GenerateR2VsTimePlot(resultsFile, outputFile string) error {
	results, err := loadResults(resultsFile)
	if err != nil {
		return err
	}

	p := plot.New()
	p.Title.Text = "КОМПРОМИСС КАЧЕСТВО/ВРЕМЯ\n(Сравнение регрессоров)"
	p.Title.TextStyle.Font.Size = 14
	p.Title.TextStyle.Color = color.RGBA{R: 0, G: 100, B: 0, A: 255}
	p.X.Label.Text = "Время выполнения (миллисекунды)"
	p.X.Label.TextStyle.Font.Size = 14
	p.Y.Label.Text = "Коэффициент детерминации R²"
	p.Y.Label.TextStyle.Font.Size = 14

	points := make(plotter.XYs, len(results.RegressionResults))
	for i, r := range results.RegressionResults {
		points[i].X = r.ExecutionTime
		points[i].Y = r.R2
	}

	scatter, err := plotter.NewScatter(points)
	if err != nil {
		return err
	}
	scatter.GlyphStyle.Color = color.RGBA{R: 70, G: 130, B: 180, A: 255}
	scatter.GlyphStyle.Radius = vg.Points(3)

	p.Add(scatter)

	if err := p.Save(8*vg.Inch, 5*vg.Inch, outputFile); err != nil {
		return err
	}

	return nil
}

func getModelLabel(config RegressorConfig) string {
	switch config.ModelType {
	case "linear":
		return "Линейная регрессия"
	case "ridge":
		return fmt.Sprintf("Ридж-регрессия (α=%.2f)", config.Alpha)
	case "lasso":
		return fmt.Sprintf("Лассо-регрессия (α=%.2f)", config.Alpha)
	case "decision_tree":
		return fmt.Sprintf("Дерево решений (глубина=%d)", config.MaxDepth)
	case "random_forest":
		return fmt.Sprintf("Случайный лес (деревьев=%d, глубина=%d)", config.NumTrees, config.MaxDepth)
	case "gradient_boosting":
		return fmt.Sprintf("Градиентный бустинг (деревьев=%d, скорость=%.2f)", config.NumTrees, config.LearningRate)
	default:
		return config.ModelType
	}
}

func getShortModelLabel(config RegressorConfig) string {
	switch config.ModelType {
	case "linear":
		return "Линейная"
	case "ridge":
		return fmt.Sprintf("Ридж\n(α=%.2f)", config.Alpha)
	case "lasso":
		return fmt.Sprintf("Лассо\n(α=%.2f)", config.Alpha)
	case "decision_tree":
		return fmt.Sprintf("Дерево\n(d=%d)", config.MaxDepth)
	case "random_forest":
		return fmt.Sprintf("СЛ\n(t=%d)", config.NumTrees)
	case "gradient_boosting":
		return fmt.Sprintf("ГБ\n(t=%d)", config.NumTrees)
	default:
		return config.ModelType
	}
}
