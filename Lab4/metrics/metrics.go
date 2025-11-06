package metrics

import "math"

type ClusteringMetrics struct {
	SilhouetteScore    float64
	CalinskiHarabasz   float64
	DaviesBouldin      float64
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

func CalculateCalinskiHarabasz(instances []struct {
	Features []float64
	Label    int
}, assignments []int, k int) float64 {
	if len(instances) == 0 || k <= 1 {
		return 0
	}

	n := float64(len(instances))
	overallMean := make([]float64, len(instances[0].Features))
	for _, inst := range instances {
		for j, val := range inst.Features {
			overallMean[j] += val
		}
	}
	for j := range overallMean {
		overallMean[j] /= n
	}

	betweenSS := 0.0
	withinSS := 0.0

	for clusterID := 0; clusterID < k; clusterID++ {
		clusterPoints := make([]int, 0)
		for i, ass := range assignments {
			if ass == clusterID {
				clusterPoints = append(clusterPoints, i)
			}
		}

		if len(clusterPoints) == 0 {
			continue
		}

		clusterMean := make([]float64, len(overallMean))
		for _, idx := range clusterPoints {
			for j, val := range instances[idx].Features {
				clusterMean[j] += val
			}
		}
		for j := range clusterMean {
			clusterMean[j] /= float64(len(clusterPoints))
		}

		clusterSize := float64(len(clusterPoints))
		betweenDist := 0.0
		for j := range overallMean {
			diff := clusterMean[j] - overallMean[j]
			betweenDist += diff * diff
		}
		betweenSS += clusterSize * betweenDist

		for _, idx := range clusterPoints {
			withinDist := 0.0
			for j, val := range instances[idx].Features {
				diff := val - clusterMean[j]
				withinDist += diff * diff
			}
			withinSS += withinDist
		}
	}

	if withinSS < 1e-10 {
		return 0
	}

	return (betweenSS / float64(k-1)) / (withinSS / (n - float64(k)))
}

func CalculateDaviesBouldin(instances []struct {
	Features []float64
	Label    int
}, assignments []int, k int) float64 {
	if len(instances) == 0 || k <= 1 {
		return 1e10
	}

	clusterMeans := make([][]float64, k)
	clusterSizes := make([]int, k)

	for clusterID := 0; clusterID < k; clusterID++ {
		clusterMeans[clusterID] = make([]float64, len(instances[0].Features))
		for i, ass := range assignments {
			if ass == clusterID {
				clusterSizes[clusterID]++
				for j, val := range instances[i].Features {
					clusterMeans[clusterID][j] += val
				}
			}
		}
		if clusterSizes[clusterID] > 0 {
			for j := range clusterMeans[clusterID] {
				clusterMeans[clusterID][j] /= float64(clusterSizes[clusterID])
			}
		}
	}

	avgDistances := make([]float64, k)
	for clusterID := 0; clusterID < k; clusterID++ {
		if clusterSizes[clusterID] == 0 {
			avgDistances[clusterID] = 1e10
			continue
		}

		totalDist := 0.0
		for i, ass := range assignments {
			if ass == clusterID {
				dist := euclideanDistance(instances[i].Features, clusterMeans[clusterID])
				totalDist += dist
			}
		}
		avgDistances[clusterID] = totalDist / float64(clusterSizes[clusterID])
	}

	db := 0.0
	for i := 0; i < k; i++ {
		if clusterSizes[i] == 0 {
			continue
		}

		maxRatio := 0.0
		for j := 0; j < k; j++ {
			if i == j || clusterSizes[j] == 0 {
				continue
			}

			centerDist := euclideanDistance(clusterMeans[i], clusterMeans[j])
			if centerDist < 1e-10 {
				continue
			}

			ratio := (avgDistances[i] + avgDistances[j]) / centerDist
			if ratio > maxRatio {
				maxRatio = ratio
			}
		}
		db += maxRatio
	}

	return db / float64(k)
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
	return math.Sqrt(sum)
}

func CalculateMetrics(instances []struct {
	Features []float64
	Label    int
}, assignments []int, k int) ClusteringMetrics {
	uniqueClusters := make(map[int]bool)
	for _, ass := range assignments {
		if ass >= 0 {
			uniqueClusters[ass] = true
		}
	}
	actualK := len(uniqueClusters)

	if actualK <= 1 {
		return ClusteringMetrics{
			SilhouetteScore:  -1,
			CalinskiHarabasz: 0,
			DaviesBouldin:    1e10,
		}
	}

	return ClusteringMetrics{
		SilhouetteScore:  CalculateSilhouetteScore(instances, assignments, actualK),
		CalinskiHarabasz: CalculateCalinskiHarabasz(instances, assignments, actualK),
		DaviesBouldin:    CalculateDaviesBouldin(instances, assignments, actualK),
	}
}

