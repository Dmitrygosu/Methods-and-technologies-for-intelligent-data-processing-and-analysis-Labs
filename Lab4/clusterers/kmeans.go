package clusterers

import (
	"lab4/data"
	"math"
	"math/rand"
)

type KMeans struct {
	k           int
	centroids   [][]float64
	numFeatures int
}

func NewKMeans(k int) *KMeans {
	return &KMeans{
		k: k,
	}
}

func (km *KMeans) Train(dataset *data.Dataset, maxIterations int) {
	if len(dataset.Instances) == 0 {
		return
	}

	km.numFeatures = dataset.NumFeatures
	km.centroids = make([][]float64, km.k)

	for i := range km.centroids {
		km.centroids[i] = make([]float64, km.numFeatures)
		randomIdx := rand.Intn(len(dataset.Instances))
		copy(km.centroids[i], dataset.Instances[randomIdx].Features)
	}

	for iter := 0; iter < maxIterations; iter++ {
		clusters := make([][]data.Instance, km.k)
		for i := range clusters {
			clusters[i] = make([]data.Instance, 0)
		}

		for _, inst := range dataset.Instances {
			clusterIdx := km.findClosestCentroid(inst.Features)
			clusters[clusterIdx] = append(clusters[clusterIdx], inst)
		}

		changed := false
		for i := 0; i < km.k; i++ {
			if len(clusters[i]) == 0 {
				continue
			}

			newCentroid := make([]float64, km.numFeatures)
			for _, inst := range clusters[i] {
				for j := range inst.Features {
					newCentroid[j] += inst.Features[j]
				}
			}

			for j := range newCentroid {
				newCentroid[j] /= float64(len(clusters[i]))
			}

			if !vectorsEqual(km.centroids[i], newCentroid) {
				changed = true
			}
			km.centroids[i] = newCentroid
		}

		if !changed {
			break
		}
	}
}

func (km *KMeans) findClosestCentroid(features []float64) int {
	minDist := math.Inf(1)
	closestIdx := 0

	for i, centroid := range km.centroids {
		dist := data.EuclideanDistance(features, centroid)
		if dist < minDist {
			minDist = dist
			closestIdx = i
		}
	}

	return closestIdx
}

func (km *KMeans) Predict(features []float64) int {
	return km.findClosestCentroid(features)
}

func (km *KMeans) PredictBatch(dataset *data.Dataset) []int {
	predictions := make([]int, len(dataset.Instances))
	for i, inst := range dataset.Instances {
		predictions[i] = km.Predict(inst.Features)
	}
	return predictions
}

func (km *KMeans) GetCentroids() [][]float64 {
	return km.centroids
}

func vectorsEqual(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if math.Abs(a[i]-b[i]) > 1e-10 {
			return false
		}
	}
	return true
}
