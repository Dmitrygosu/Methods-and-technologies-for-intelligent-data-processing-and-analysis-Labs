package utils

import (
	"encoding/json"
	"fmt"
	"image/color"
	"os"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

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

type ClassifierConfig struct {
	ModelType    string  `json:"model_type"`
	K            int     `json:"k,omitempty"`
	MaxDepth     int     `json:"max_depth,omitempty"`
	LearningRate float64 `json:"learning_rate,omitempty"`
	Iterations   int     `json:"iterations,omitempty"`
}

type ClusteringResult struct {
	TaskName        string          `json:"task_name"`
	Config          ClustererConfig `json:"config"`
	SilhouetteScore float64         `json:"silhouette_score"`
	ExecutionTime   float64         `json:"execution_time_ms"`
}

type ClustererConfig struct {
	K       int `json:"k"`
	MaxIter int `json:"max_iter"`
}

type AllResults struct {
	ClassificationResults []ClassificationResult `json:"classification_results"`
	ClusteringResults     []ClusteringResult     `json:"clustering_results"`
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

func GenerateClassificationComparisonPlot(resultsFile, outputFile string) error {
	results, err := loadResults(resultsFile)
	if err != nil {
		return err
	}

	p := plot.New()
	p.Title.Text = "СРАВНЕНИЕ ТОЧНОСТИ КЛАССИФИКАТОРОВ\n(На датасете Iris: 3 класса цветков ириса)"
	p.Title.TextStyle.Font.Size = 14
	p.Title.TextStyle.Color = color.RGBA{R: 0, G: 100, B: 0, A: 255}
	p.Y.Label.Text = "Точность классификации"
	p.Y.Label.TextStyle.Font.Size = 14
	p.X.Label.Text = "Классификатор (модель)"
	p.X.Label.TextStyle.Font.Size = 14

	modelGroups := make(map[string][]float64)
	modelLabels := make([]string, 0)
	modelShortLabels := make([]string, 0)
	modelMap := make(map[string]bool)

	for _, r := range results.ClassificationResults {
		label := getModelLabel(r.Config)
		if !modelMap[label] {
			modelLabels = append(modelLabels, label)
			modelShortLabels = append(modelShortLabels, getShortModelLabel(r.Config))
			modelMap[label] = true
			modelGroups[label] = make([]float64, 0)
		}
		modelGroups[label] = append(modelGroups[label], r.Accuracy)
	}

	bars := make([]plotter.Values, len(modelLabels))
	for i, label := range modelLabels {
		bars[i] = modelGroups[label]
	}

	width := vg.Points(40)
	b, err := plotter.NewBarChart(plotter.Values(bars[0]), width)
	if err != nil {
		return err
	}

	avgValues := make(plotter.Values, len(modelLabels))
	for i, label := range modelLabels {
		sum := 0.0
		for _, v := range modelGroups[label] {
			sum += v
		}
		avgValues[i] = sum / float64(len(modelGroups[label]))
	}

	b.Values = avgValues
	b.Color = color.RGBA{R: 70, G: 130, B: 180, A: 255}
	b.LineStyle.Color = color.Black
	b.LineStyle.Width = vg.Points(1)

	p.NominalX(modelShortLabels...)
	p.X.Tick.Label.Rotation = 0.785
	p.X.Tick.Label.YAlign = draw.YCenter
	p.X.Tick.Label.XAlign = draw.XCenter
	p.X.Tick.Label.Font.Size = 10
	p.X.Padding = vg.Points(40)

	p.Add(b)
	p.Legend.Add("Средняя точность", b)

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
	p.Title.Text = "СРАВНЕНИЕ КАЧЕСТВА КЛАСТЕРИЗАЦИИ\n(Алгоритм K-Means на датасете Iris)"
	p.Title.TextStyle.Font.Size = 14
	p.Title.TextStyle.Color = color.RGBA{R: 0, G: 100, B: 0, A: 255}
	p.Y.Label.Text = "Силуэтный коэффициент (качество кластеризации)"
	p.Y.Label.TextStyle.Font.Size = 14
	p.X.Label.Text = "Количество кластеров (K)"
	p.X.Label.TextStyle.Font.Size = 14

	kGroups := make(map[int][]float64)
	kValues := make([]int, 0)
	kMap := make(map[int]bool)

	for _, r := range results.ClusteringResults {
		k := r.Config.K
		if !kMap[k] {
			kValues = append(kValues, k)
			kMap[k] = true
			kGroups[k] = make([]float64, 0)
		}
		kGroups[k] = append(kGroups[k], r.SilhouetteScore)
	}

	points := make(plotter.XYs, len(kValues))
	for i, k := range kValues {
		sum := 0.0
		for _, v := range kGroups[k] {
			sum += v
		}
		points[i].X = float64(k)
		points[i].Y = sum / float64(len(kGroups[k]))
	}

	line, err := plotter.NewLine(points)
	if err != nil {
		return err
	}
	line.LineStyle.Width = vg.Points(2)
	line.LineStyle.Color = color.RGBA{R: 220, G: 20, B: 60, A: 255}

	scatter, err := plotter.NewScatter(points)
	if err != nil {
		return err
	}
	scatter.GlyphStyle.Color = color.RGBA{R: 220, G: 20, B: 60, A: 255}
	scatter.GlyphStyle.Radius = vg.Points(4)

	p.Add(line, scatter)
	p.Legend.Add("Силуэтный коэффициент", line)

	if err := p.Save(8*vg.Inch, 5*vg.Inch, outputFile); err != nil {
		return err
	}

	return nil
}

func GenerateAccuracyVsTimePlot(resultsFile, outputFile string) error {
	results, err := loadResults(resultsFile)
	if err != nil {
		return err
	}

	p := plot.New()
	p.Title.Text = "КОМПРОМИСС ТОЧНОСТЬ/ВРЕМЯ\n(Сравнение классификаторов)"
	p.Title.TextStyle.Font.Size = 14
	p.Title.TextStyle.Color = color.RGBA{R: 0, G: 100, B: 0, A: 255}
	p.X.Label.Text = "Время выполнения (миллисекунды)"
	p.X.Label.TextStyle.Font.Size = 14
	p.Y.Label.Text = "Точность классификации"
	p.Y.Label.TextStyle.Font.Size = 14

	points := make(plotter.XYs, len(results.ClassificationResults))
	for i, r := range results.ClassificationResults {
		points[i].X = r.ExecutionTime
		points[i].Y = r.Accuracy
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

func getModelLabel(config ClassifierConfig) string {
	switch config.ModelType {
	case "knn":
		return fmt.Sprintf("Метод k-ближайших соседей (k=%d)", config.K)
	case "naive_bayes":
		return "Наивный байесовский классификатор"
	case "decision_tree":
		return fmt.Sprintf("Дерево решений (глубина=%d)", config.MaxDepth)
	case "logistic_regression":
		return fmt.Sprintf("Логистическая регрессия (скорость обучения=%.2f)", config.LearningRate)
	default:
		return config.ModelType
	}
}

func getShortModelLabel(config ClassifierConfig) string {
	switch config.ModelType {
	case "knn":
		return fmt.Sprintf("KNN\n(k=%d)", config.K)
	case "naive_bayes":
		return "Наивный\nБайес"
	case "decision_tree":
		return fmt.Sprintf("Дерево\n(d=%d)", config.MaxDepth)
	case "logistic_regression":
		return fmt.Sprintf("Лог. регр.\n(lr=%.2f)", config.LearningRate)
	default:
		return config.ModelType
	}
}

func GenerateConfusionMatrixPlot(resultsFile, outputFile, modelType string) error {
	results, err := loadResults(resultsFile)
	if err != nil {
		return err
	}

	var targetResult *ClassificationResult
	for _, r := range results.ClassificationResults {
		if r.Config.ModelType == modelType {
			if targetResult == nil || r.Accuracy > targetResult.Accuracy {
				targetResult = &r
			}
		}
	}

	if targetResult == nil || len(targetResult.ConfusionMatrix) == 0 {
		return fmt.Errorf("матрица ошибок не найдена для модели %s", modelType)
	}

	cm := targetResult.ConfusionMatrix
	numClasses := len(cm)

	p := plot.New()
	p.Title.Text = fmt.Sprintf("МАТРИЦА ОШИБОК: %s\n(Классы: Сетоза, Версиколор, Виргиника - виды ириса)", getModelLabel(targetResult.Config))
	p.Title.TextStyle.Font.Size = 14
	p.Title.TextStyle.Color = color.RGBA{R: 0, G: 100, B: 0, A: 255}

	classNames := []string{"Сетоза", "Версиколор", "Виргиника"}
	if numClasses > 3 {
		classNames = make([]string, numClasses)
		for i := 0; i < numClasses; i++ {
			classNames[i] = fmt.Sprintf("Класс %d", i)
		}
	}

	values := make(plotter.Values, 0)
	xLabels := make([]string, 0)

	for i := 0; i < numClasses; i++ {
		for j := 0; j < numClasses; j++ {
			values = append(values, float64(cm[i][j]))
			if i == j {
				xLabels = append(xLabels, fmt.Sprintf("%s→%s\n(правильно)", classNames[i], classNames[j]))
			} else {
				xLabels = append(xLabels, fmt.Sprintf("%s→%s\n(ошибка)", classNames[i], classNames[j]))
			}
		}
	}

	width := vg.Points(20)
	bars, err := plotter.NewBarChart(values, width)
	if err != nil {
		return err
	}

	bars.Color = color.RGBA{R: 70, G: 130, B: 180, A: 255}
	bars.LineStyle.Color = color.Black
	bars.LineStyle.Width = vg.Points(1)

	p.NominalX(xLabels...)
	p.X.Tick.Label.Rotation = 0.785
	p.X.Tick.Label.YAlign = draw.YCenter
	p.X.Tick.Label.XAlign = draw.XCenter
	p.X.Tick.Label.Font.Size = 9
	p.X.Padding = vg.Points(48)
	p.X.Label.Text = "Истинный класс → Предсказанный класс\n(Сетоза, Версиколор, Виргиника - классы цветков ириса)"
	p.X.Label.TextStyle.Font.Size = 12
	p.Y.Label.Text = "Количество объектов"
	p.Y.Label.TextStyle.Font.Size = 14

	p.Add(bars)

	if err := p.Save(14*vg.Inch, 7*vg.Inch, outputFile); err != nil {
		return err
	}

	return nil
}
