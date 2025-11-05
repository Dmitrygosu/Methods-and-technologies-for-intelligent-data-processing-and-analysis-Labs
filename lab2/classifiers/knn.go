package classifiers

import (
	"lab2/data"
	"sort"
)

type KNN struct {
	k        int
	trainSet *data.Dataset
}

type Neighbor struct {
	Distance float64
	Label    int
}

func NewKNN(k int) *KNN {
	return &KNN{
		k: k,
	}
}

func (knn *KNN) Train(trainSet *data.Dataset) {
	knn.trainSet = trainSet
}

func (knn *KNN) Predict(features []float64) int {
	if knn.trainSet == nil || len(knn.trainSet.Instances) == 0 {
		return 0
	}
	
	neighbors := make([]Neighbor, 0, len(knn.trainSet.Instances))
	
	for _, inst := range knn.trainSet.Instances {
		dist := data.EuclideanDistance(features, inst.Features)
		neighbors = append(neighbors, Neighbor{
			Distance: dist,
			Label:    inst.Label,
		})
	}
	
	sort.Slice(neighbors, func(i, j int) bool {
		return neighbors[i].Distance < neighbors[j].Distance
	})
	
	labelCounts := make(map[int]int)
	for i := 0; i < knn.k && i < len(neighbors); i++ {
		labelCounts[neighbors[i].Label]++
	}
	
	maxCount := 0
	predictedLabel := 0
	for label, count := range labelCounts {
		if count > maxCount {
			maxCount = count
			predictedLabel = label
		}
	}
	
	return predictedLabel
}

func (knn *KNN) PredictBatch(testSet *data.Dataset) []int {
	predictions := make([]int, len(testSet.Instances))
	for i, inst := range testSet.Instances {
		predictions[i] = knn.Predict(inst.Features)
	}
	return predictions
}

