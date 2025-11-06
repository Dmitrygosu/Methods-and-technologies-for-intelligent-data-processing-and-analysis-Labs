package clusterers

import (
	"lab4/data"
)

type DBSCAN struct {
	eps         float64
	minSamples  int
	labels      []int
	numFeatures int
}

func NewDBSCAN(eps float64, minSamples int) *DBSCAN {
	return &DBSCAN{
		eps:        eps,
		minSamples: minSamples,
	}
}

func (db *DBSCAN) Train(dataset *data.Dataset) {
	if len(dataset.Instances) == 0 {
		return
	}

	db.numFeatures = dataset.NumFeatures
	db.labels = make([]int, len(dataset.Instances))
	for i := range db.labels {
		db.labels[i] = -2
	}

	clusterID := 0
	for i := range dataset.Instances {
		if db.labels[i] != -2 {
			continue
		}

		neighbors := db.getNeighbors(dataset, i)
		if len(neighbors) < db.minSamples {
			db.labels[i] = -1
			continue
		}

		db.expandCluster(dataset, i, neighbors, clusterID)
		clusterID++
	}
}

func (db *DBSCAN) getNeighbors(dataset *data.Dataset, idx int) []int {
	neighbors := make([]int, 0)
	for i := range dataset.Instances {
		if i == idx {
			continue
		}
		dist := data.EuclideanDistance(dataset.Instances[idx].Features, dataset.Instances[i].Features)
		if dist <= db.eps {
			neighbors = append(neighbors, i)
		}
	}
	return neighbors
}

func (db *DBSCAN) expandCluster(dataset *data.Dataset, pointIdx int, neighbors []int, clusterID int) {
	db.labels[pointIdx] = clusterID

	seedSet := make([]int, len(neighbors))
	copy(seedSet, neighbors)

	i := 0
	for i < len(seedSet) {
		currentPoint := seedSet[i]
		i++

		if db.labels[currentPoint] == -1 {
			db.labels[currentPoint] = clusterID
		}

		if db.labels[currentPoint] != -2 {
			continue
		}

		db.labels[currentPoint] = clusterID
		currentNeighbors := db.getNeighbors(dataset, currentPoint)
		if len(currentNeighbors) >= db.minSamples {
			seedSet = append(seedSet, currentNeighbors...)
		}
	}
}

func (db *DBSCAN) Predict(features []float64) int {
	return -1
}

func (db *DBSCAN) PredictBatch(dataset *data.Dataset) []int {
	return db.labels
}

func (db *DBSCAN) GetLabels() []int {
	return db.labels
}
