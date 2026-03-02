package clusterers

import (
	"lab4/data"
	"math"
)

type Hierarchical struct {
	k           int
	labels      []int
	numFeatures int
}

func NewHierarchical(k int) *Hierarchical {
	return &Hierarchical{
		k: k,
	}
}

func (h *Hierarchical) Train(dataset *data.Dataset) {
	if len(dataset.Instances) == 0 {
		return
	}

	h.numFeatures = dataset.NumFeatures
	n := len(dataset.Instances)

	clusters := make([][]int, n)
	for i := range clusters {
		clusters[i] = []int{i}
	}

	distances := make([][]float64, n)
	for i := range distances {
		distances[i] = make([]float64, n)
		for j := range distances[i] {
			if i < j {
				distances[i][j] = data.EuclideanDistance(dataset.Instances[i].Features, dataset.Instances[j].Features)
			} else if i > j {
				distances[i][j] = distances[j][i]
			}
		}
	}

	for len(clusters) > h.k {
		minDist := math.Inf(1)
		mergeI, mergeJ := -1, -1

		for i := 0; i < len(clusters); i++ {
			for j := i + 1; j < len(clusters); j++ {
				dist := h.completeLinkage(dataset, clusters[i], clusters[j], distances)
				if dist < minDist {
					minDist = dist
					mergeI, mergeJ = i, j
				}
			}
		}

		if mergeI == -1 || mergeJ == -1 {
			break
		}

		clusters[mergeI] = append(clusters[mergeI], clusters[mergeJ]...)
		clusters = append(clusters[:mergeJ], clusters[mergeJ+1:]...)
	}

	h.labels = make([]int, n)
	for clusterID, cluster := range clusters {
		for _, pointIdx := range cluster {
			h.labels[pointIdx] = clusterID
		}
	}
}

func (h *Hierarchical) completeLinkage(dataset *data.Dataset, cluster1, cluster2 []int, distances [][]float64) float64 {
	maxDist := 0.0
	for _, i := range cluster1 {
		for _, j := range cluster2 {
			dist := distances[i][j]
			if dist > maxDist {
				maxDist = dist
			}
		}
	}
	return maxDist
}

func (h *Hierarchical) Predict(features []float64) int {
	return -1
}

func (h *Hierarchical) PredictBatch(dataset *data.Dataset) []int {
	return h.labels
}

func (h *Hierarchical) GetLabels() []int {
	return h.labels
}
