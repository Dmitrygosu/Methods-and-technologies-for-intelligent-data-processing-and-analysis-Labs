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

type ExperimentConfig struct {
	PopulationSize int     `json:"population_size"`
	MaxGenerations int     `json:"max_generations"`
	CrossoverProb  float64 `json:"crossover_prob"`
	MutationProb   float64 `json:"mutation_prob"`
	CrossoverType  string  `json:"crossover_type"`
	ElitismCount   int     `json:"elitism_count"`
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

func GenerateTimeComparisonPlot(resultsFile, outputFile string) error {
	results, err := loadResults(resultsFile)
	if err != nil {
		return err
	}

	p := plot.New()
	p.Title.Text = " СРАВНЕНИЕ ВРЕМЕНИ ВЫПОЛНЕНИЯ "
	p.Title.TextStyle.Font.Size = 16
	p.Title.TextStyle.Color = color.RGBA{R: 0, G: 100, B: 0, A: 255}
	p.Y.Label.Text = "Время выполнения (миллисекунды)"
	p.Y.Label.TextStyle.Font.Size = 14
	p.X.Label.Text = "Тип алгоритма и задача"
	p.X.Label.TextStyle.Font.Size = 14

	var arrayGATimes []float64
	for _, r := range results.GAResults {
		if r.TaskName == "array_search" {
			arrayGATimes = append(arrayGATimes, r.ExecutionTime)
		}
	}

	var funcGATimes []float64
	for _, r := range results.GAResults {
		if r.TaskName == "function_optimization" {
			funcGATimes = append(funcGATimes, r.ExecutionTime)
		}
	}

	var arrayLinearTime, funcLinearTime float64
	for _, r := range results.LinearSearchResults {
		switch r.TaskName {
		case "array_search":
			arrayLinearTime = r.ExecutionTime
		case "function_optimization":
			funcLinearTime = r.ExecutionTime
		}
	}

	avgArrayGA := average(arrayGATimes)
	avgFuncGA := average(funcGATimes)

	w := vg.Points(30)

	values := plotter.Values{avgArrayGA, arrayLinearTime, avgFuncGA, funcLinearTime}

	bars, err := plotter.NewBarChart(values, w)
	if err != nil {
		return err
	}

	colors := []color.RGBA{
		{R: 34, G: 139, B: 34, A: 255},
		{R: 220, G: 20, B: 60, A: 255},
		{R: 0, G: 191, B: 255, A: 255},
		{R: 255, G: 69, B: 0, A: 255},
	}

	for i := 0; i < len(values); i++ {
		bars.Color = colors[i]
	}

	bars.LineStyle.Width = vg.Length(2)
	bars.LineStyle.Color = color.RGBA{R: 0, G: 0, B: 0, A: 255}

	p.Add(bars)

	if avgArrayGA < 0.001 {
		avgArrayGA = 0.001
	}
	if avgFuncGA < 0.001 {
		avgFuncGA = 0.001
	}

	acceleration1 := arrayLinearTime / avgArrayGA
	acceleration2 := funcLinearTime / avgFuncGA

	var interpretation1, interpretation2 string
	if acceleration1 > 1 {
		interpretation1 = fmt.Sprintf("быстрее в %.1fx раз", acceleration1)
	} else if acceleration1 > 0.1 {
		interpretation1 = fmt.Sprintf("медленнее в %.1fx раз", 1/acceleration1)
	} else {
		interpretation1 = "значительно медленнее"
	}
	if acceleration2 > 1 {
		interpretation2 = fmt.Sprintf("быстрее в %.1fx раз", acceleration2)
	} else if acceleration2 > 0.1 {
		interpretation2 = fmt.Sprintf("медленнее в %.1fx раз", 1/acceleration2)
	} else {
		interpretation2 = "значительно медленнее"
	}

	p.Title.Text = fmt.Sprintf("СРАВНЕНИЕ ВРЕМЕНИ ВЫПОЛНЕНИЯ\nГА %s для массива, %s для функции\nРезультат зависит от параметров: малая популяция=быстро, большая=медленно", interpretation1, interpretation2)

	p.NominalX("Генетический\nалгоритм\n(поиск в массиве)",
		"Линейный поиск\n(поиск в массиве)",
		"Генетический\nалгоритм\n(оптимизация функции)",
		"Линейный поиск\n(оптимизация функции)")

	p.Add(plotter.NewGrid())

	if err := p.Save(12*vg.Inch, 8*vg.Inch, outputFile); err != nil {
		return err
	}

	return nil
}

func GenerateConvergencePlot(resultsFile, outputFile string) error {
	results, err := loadResults(resultsFile)
	if err != nil {
		return err
	}

	p := plot.New()
	p.Title.Text = "📈 СХОДИМОСТЬ ГЕНЕТИЧЕСКОГО АЛГОРИТМА 📈"
	p.Title.TextStyle.Font.Size = 16
	p.Title.TextStyle.Color = color.RGBA{R: 0, G: 0, B: 139, A: 255}
	p.X.Label.Text = "Номер поколения (эволюция)"
	p.X.Label.TextStyle.Font.Size = 14
	p.Y.Label.Text = "Лучшая приспособленность (качество решения)"
	p.Y.Label.TextStyle.Font.Size = 14
	p.Legend.Top = false
	p.Legend.TextStyle.Font.Size = 10

	configsToShow := 0
	maxConfigs := 6
	colors := []color.RGBA{
		{R: 255, G: 0, B: 0, A: 255},   // Красный
		{R: 0, G: 128, B: 0, A: 255},   // Зеленый
		{R: 0, G: 0, B: 255, A: 255},   // Синий
		{R: 255, G: 165, B: 0, A: 255}, // Оранжевый
		{R: 128, G: 0, B: 128, A: 255}, // Фиолетовый
		{R: 0, G: 191, B: 255, A: 255}, // Голубой
	}

	for _, r := range results.GAResults {
		if r.TaskName != "array_search" || len(r.Convergence) == 0 {
			continue
		}

		if configsToShow >= maxConfigs {
			break
		}

		pts := make(plotter.XYs, len(r.Convergence))
		for j, val := range r.Convergence {
			pts[j].X = float64(j)
			pts[j].Y = val
		}

		line, err := plotter.NewLine(pts)
		if err != nil {
			return err
		}
		line.Color = colors[configsToShow%len(colors)]
		line.Width = vg.Points(3)

		scatter, err := plotter.NewScatter(pts)
		if err == nil {
			scatter.GlyphStyle.Color = colors[configsToShow%len(colors)]
			scatter.GlyphStyle.Radius = vg.Points(2)
			scatter.GlyphStyle.Shape = draw.CircleGlyph{}
			p.Add(scatter)
		}

		mutationDesc := "низкая"
		if r.Config.MutationProb >= 0.05 {
			mutationDesc = "высокая"
		}

		crossoverDesc := "одноточечное"
		if r.Config.CrossoverType == "uniform" {
			crossoverDesc = "униформное"
		}

		label := fmt.Sprintf("%s | %.2f мутация | %s скрещивание | популяция=%d",
			mutationDesc, r.Config.MutationProb, crossoverDesc, r.Config.PopulationSize)

		p.Add(line)
		p.Legend.Add(label, line)

		configsToShow++
	}

	p.Add(plotter.NewGrid())

	p.Title.Text = "СХОДИМОСТЬ ГЕНЕТИЧЕСКОГО АЛГОРИТМА\nВысокая мутация → больше исследования пространства | Низкая мутация → быстрая сходимость к локальному оптимуму"

	if err := p.Save(14*vg.Inch, 10*vg.Inch, outputFile); err != nil {
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
	p.Title.Text = " КОМПРОМИСС ТОЧНОСТЬ/ВРЕМЯ "
	p.Title.TextStyle.Font.Size = 16
	p.Title.TextStyle.Color = color.RGBA{R: 139, G: 0, B: 139, A: 255}
	p.X.Label.Text = "Время выполнения (миллисекунды)"
	p.X.Label.TextStyle.Font.Size = 14
	p.Y.Label.Text = "Относительная ошибка (процент от оптимального)"
	p.Y.Label.TextStyle.Font.Size = 14

	arrayPts := make(plotter.XYs, 0)
	for _, r := range results.GAResults {
		if r.TaskName == "array_search" {
			arrayPts = append(arrayPts, plotter.XY{
				X: r.ExecutionTime,
				Y: r.RelativeError * 100,
			})
		}
	}

	funcPts := make(plotter.XYs, 0)
	for _, r := range results.GAResults {
		if r.TaskName == "function_optimization" {
			funcPts = append(funcPts, plotter.XY{
				X: r.ExecutionTime,
				Y: r.RelativeError * 100,
			})
		}
	}

	excellentZone := plotter.XYs{
		{X: 0, Y: 0}, {X: 100, Y: 0}, {X: 100, Y: 2}, {X: 0, Y: 2}, {X: 0, Y: 0},
	}
	goodZone := plotter.XYs{
		{X: 0, Y: 2}, {X: 100, Y: 2}, {X: 100, Y: 5}, {X: 0, Y: 5}, {X: 0, Y: 2},
	}

	excellentPoly, err := plotter.NewPolygon(excellentZone)
	if err == nil {
		excellentPoly.Color = color.RGBA{R: 0, G: 255, B: 0, A: 50}
		excellentPoly.LineStyle.Color = color.RGBA{R: 0, G: 200, B: 0, A: 255}
		excellentPoly.LineStyle.Width = vg.Points(1)
		p.Add(excellentPoly)
	}

	goodPoly, err := plotter.NewPolygon(goodZone)
	if err == nil {
		goodPoly.Color = color.RGBA{R: 255, G: 255, B: 0, A: 30}
		goodPoly.LineStyle.Color = color.RGBA{R: 200, G: 200, B: 0, A: 255}
		goodPoly.LineStyle.Width = vg.Points(1)
		p.Add(goodPoly)
	}

	if len(arrayPts) > 0 {
		arrayScatter, err := plotter.NewScatter(arrayPts)
		if err != nil {
			return err
		}
		arrayScatter.GlyphStyle.Color = color.RGBA{R: 255, G: 0, B: 0, A: 180}
		arrayScatter.GlyphStyle.Radius = vg.Points(4)
		arrayScatter.GlyphStyle.Shape = draw.CircleGlyph{}

		p.Add(arrayScatter)
		p.Legend.Add("Поиск в массиве (случайные данные)", arrayScatter)
	}

	if len(funcPts) > 0 {
		funcScatter, err := plotter.NewScatter(funcPts)
		if err != nil {
			return err
		}
		funcScatter.GlyphStyle.Color = color.RGBA{R: 0, G: 0, B: 255, A: 180}
		funcScatter.GlyphStyle.Radius = vg.Points(4)
		funcScatter.GlyphStyle.Shape = draw.TriangleGlyph{}

		p.Add(funcScatter)
		p.Legend.Add("Оптимизация функции (математическая)", funcScatter)
	}

	excellentLegend, err := plotter.NewPolygon(plotter.XYs{{X: 0, Y: 0}})
	if err == nil {
		excellentLegend.Color = color.RGBA{R: 0, G: 255, B: 0, A: 50}
		p.Legend.Add("Отличное качество (ошибка 0-2%)", excellentLegend)
	}

	goodLegend, err := plotter.NewPolygon(plotter.XYs{{X: 0, Y: 0}})
	if err == nil {
		goodLegend.Color = color.RGBA{R: 255, G: 255, B: 0, A: 30}
		p.Legend.Add("Хорошее качество (ошибка 2-5%)", goodLegend)
	}

	p.Legend.Top = true
	p.Legend.TextStyle.Font.Size = 10

	p.Add(plotter.NewGrid())

	p.Title.Text = "КОМПРОМИСС ТОЧНОСТЬ/ВРЕМЯ\nБыстро+точно (идеал) | Быстро+приблизительно (практично) | Медленно+точно (эталон)"

	if err := p.Save(14*vg.Inch, 10*vg.Inch, outputFile); err != nil {
		return err
	}

	return nil
}

func GenerateEfficiencyComparisonPlot(resultsFile, outputFile string) error {
	results, err := loadResults(resultsFile)
	if err != nil {
		return err
	}

	p := plot.New()
	p.Title.Text = " СРАВНЕНИЕ ЭФФЕКТИВНОСТИ АЛГОРИТМОВ "
	p.Title.TextStyle.Font.Size = 18
	p.Title.TextStyle.Color = color.RGBA{R: 0, G: 0, B: 139, A: 255}
	p.X.Label.Text = "Тип алгоритма и задача"
	p.X.Label.TextStyle.Font.Size = 14
	p.Y.Label.Text = "Индекс эффективности (баллы)"
	p.Y.Label.TextStyle.Font.Size = 14

	arrayGAEff := calculateEfficiency(results, "array_search", true)
	arrayLinearEff := calculateEfficiency(results, "array_search", false)
	funcGAEff := calculateEfficiency(results, "function_optimization", true)
	funcLinearEff := calculateEfficiency(results, "function_optimization", false)

	values := plotter.Values{arrayGAEff, arrayLinearEff, funcGAEff, funcLinearEff}

	w := vg.Points(40)
	bars, err := plotter.NewBarChart(values, w)
	if err != nil {
		return err
	}

	colors := []color.RGBA{
		{R: 0, G: 255, B: 0, A: 255},
		{R: 255, G: 0, B: 0, A: 255},
		{R: 0, G: 0, B: 255, A: 255},
		{R: 255, G: 165, B: 0, A: 255},
	}

	for i := 0; i < len(values); i++ {
		bars.Color = colors[i]
	}

	bars.LineStyle.Width = vg.Length(3)
	bars.LineStyle.Color = color.RGBA{R: 0, G: 0, B: 0, A: 255}

	p.Add(bars)

	p.Title.Text = fmt.Sprintf("СРАВНЕНИЕ ЭФФЕКТИВНОСТИ АЛГОРИТМОВ\nФормула эффективности: (100 - ошибка%%) / время_мс × 1000\nГА(массив): %.1f баллов | Линейный(массив): %.1f баллов | ГА(функция): %.1f баллов | Линейный(функция): %.1f баллов\nЧем выше балл, тем лучше соотношение точности и скорости",
		arrayGAEff, arrayLinearEff, funcGAEff, funcLinearEff)

	p.NominalX("Генетический\nалгоритм\n(поиск в массиве)",
		"Линейный поиск\n(поиск в массиве)",
		"Генетический\nалгоритм\n(оптимизация функции)",
		"Линейный поиск\n(оптимизация функции)")

	p.Add(plotter.NewGrid())

	if err := p.Save(14*vg.Inch, 10*vg.Inch, outputFile); err != nil {
		return err
	}

	return nil
}

func calculateEfficiency(results *AllResults, taskName string, isGA bool) float64 {
	if isGA {
		var totalTime, totalError float64
		count := 0
		for _, r := range results.GAResults {
			if r.TaskName == taskName {
				totalTime += r.ExecutionTime
				totalError += r.RelativeError * 100
				count++
			}
		}
		if count == 0 || totalTime == 0 {
			return 0
		}
		avgTime := totalTime / float64(count)
		avgError := totalError / float64(count)
		if avgTime < 0.001 {
			avgTime = 0.001
		}
		if avgError > 99 {
			avgError = 99
		}
		if avgError < 0.1 {
			avgError = 0.1
		}
		efficiency := (100 - avgError) / avgTime * 1000
		if efficiency > 100000 {
			efficiency = 100000
		}
		return efficiency
	} else {
		for _, r := range results.LinearSearchResults {
			if r.TaskName == taskName {
				if r.ExecutionTime == 0 {
					return 0
				}
				return 100 / r.ExecutionTime * 1000
			}
		}
	}
	return 0
}

func average(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}
