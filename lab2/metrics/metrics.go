package metrics

type ClassificationMetrics struct {
	Accuracy   float64
	Precision  []float64
	Recall     []float64
	F1         []float64
	MacroF1    float64
	WeightedF1 float64
}

func CalculateMetrics(yTrue, yPred []int, numClasses int) ClassificationMetrics {
	cm := buildConfusionMatrix(yTrue, yPred, numClasses)

	metrics := ClassificationMetrics{
		Precision: make([]float64, numClasses),
		Recall:    make([]float64, numClasses),
		F1:        make([]float64, numClasses),
	}

	correct := 0
	total := len(yTrue)

	for i := 0; i < numClasses; i++ {
		tp := cm[i][i]
		fp := 0
		fn := 0

		for j := 0; j < numClasses; j++ {
			if j != i {
				fp += cm[j][i]
				fn += cm[i][j]
			}
		}

		correct += tp

		if tp+fp > 0 {
			metrics.Precision[i] = float64(tp) / float64(tp+fp)
		}

		if tp+fn > 0 {
			metrics.Recall[i] = float64(tp) / float64(tp+fn)
		}

		if metrics.Precision[i]+metrics.Recall[i] > 0 {
			metrics.F1[i] = 2 * metrics.Precision[i] * metrics.Recall[i] / (metrics.Precision[i] + metrics.Recall[i])
		}
	}

	metrics.Accuracy = float64(correct) / float64(total)

	sumF1 := 0.0
	for i := 0; i < numClasses; i++ {
		sumF1 += metrics.F1[i]
	}
	metrics.MacroF1 = sumF1 / float64(numClasses)

	classCounts := make([]int, numClasses)
	for _, label := range yTrue {
		classCounts[label]++
	}

	weightedSum := 0.0
	totalWeight := 0.0
	for i := 0; i < numClasses; i++ {
		weight := float64(classCounts[i])
		weightedSum += metrics.F1[i] * weight
		totalWeight += weight
	}

	if totalWeight > 0 {
		metrics.WeightedF1 = weightedSum / totalWeight
	}

	return metrics
}

func buildConfusionMatrix(yTrue, yPred []int, numClasses int) [][]int {
	cm := make([][]int, numClasses)
	for i := range cm {
		cm[i] = make([]int, numClasses)
	}

	for i := 0; i < len(yTrue); i++ {
		cm[yTrue[i]][yPred[i]]++
	}

	return cm
}

func BuildConfusionMatrix(yTrue, yPred []int, numClasses int) [][]int {
	return buildConfusionMatrix(yTrue, yPred, numClasses)
}

func CalculateSilhouetteScore(instances []struct {
	Features []float64
	Label    int
}, assignments []int, k int) float64 {
	if len(instances) == 0 {
		return 0
	}

	total := 0.0
	for i := range instances {
		ai := averageIntraClusterDistance(instances, assignments, i, assignments[i])
		bi := averageNearestClusterDistance(instances, assignments, i, assignments[i], k)

		if ai > bi {
			ai, bi = bi, ai
		}

		if bi > 0 {
			total += (bi - ai) / bi
		}
	}

	return total / float64(len(instances))
}

func averageIntraClusterDistance(instances []struct {
	Features []float64
	Label    int
}, assignments []int, idx, cluster int) float64 {
	sum := 0.0
	count := 0

	for i, inst := range instances {
		if i == idx {
			continue
		}
		if assignments[i] == cluster {
			dist := euclideanDistance(instances[idx].Features, inst.Features)
			sum += dist
			count++
		}
	}

	if count == 0 {
		return 0
	}
	return sum / float64(count)
}

func averageNearestClusterDistance(instances []struct {
	Features []float64
	Label    int
}, assignments []int, idx, currentCluster, k int) float64 {
	minAvg := 1e10

	for cluster := 0; cluster < k; cluster++ {
		if cluster == currentCluster {
			continue
		}

		sum := 0.0
		count := 0

		for i, inst := range instances {
			if i == idx {
				continue
			}
			if assignments[i] == cluster {
				dist := euclideanDistance(instances[idx].Features, inst.Features)
				sum += dist
				count++
			}
		}

		if count > 0 {
			avg := sum / float64(count)
			if avg < minAvg {
				minAvg = avg
			}
		}
	}

	if minAvg == 1e10 {
		return 0
	}
	return minAvg
}

func euclideanDistance(a, b []float64) float64 {
	if len(a) != len(b) {
		return 1e10
	}

	sum := 0.0
	for i := range a {
		diff := a[i] - b[i]
		sum += diff * diff
	}

	return sum
}
