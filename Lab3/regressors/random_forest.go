package regressors

import (
	"lab3/data"
	"math/rand"
)

type RandomForestRegressor struct {
	trees       []*DecisionTreeRegressor
	numTrees    int
	maxDepth    int
	maxFeatures int
}

func NewRandomForestRegressor(numTrees, maxDepth, maxFeatures int) *RandomForestRegressor {
	return &RandomForestRegressor{
		numTrees:    numTrees,
		maxDepth:    maxDepth,
		maxFeatures: maxFeatures,
	}
}

func (rf *RandomForestRegressor) Train(dataset *data.RegressionDataset) {
	rf.trees = make([]*DecisionTreeRegressor, rf.numTrees)

	for i := 0; i < rf.numTrees; i++ {
		bootstrap := rf.bootstrapSample(dataset, int64(i))
		featureSubset := rf.selectFeatureSubset(dataset.NumFeatures, int64(i))

		tree := NewDecisionTreeRegressor(rf.maxDepth)
		tree.trainOnSubset(bootstrap, featureSubset)
		rf.trees[i] = tree
	}
}

func (rf *RandomForestRegressor) bootstrapSample(dataset *data.RegressionDataset, seed int64) []data.RegressionInstance {
	rng := rand.New(rand.NewSource(seed))
	sample := make([]data.RegressionInstance, len(dataset.Instances))

	for i := 0; i < len(dataset.Instances); i++ {
		idx := rng.Intn(len(dataset.Instances))
		sample[i] = dataset.Instances[idx]
	}

	return sample
}

func (rf *RandomForestRegressor) selectFeatureSubset(numFeatures int, seed int64) []int {
	rng := rand.New(rand.NewSource(seed))
	numSelected := rf.maxFeatures
	if numSelected > numFeatures {
		numSelected = numFeatures
	}

	perm := rng.Perm(numFeatures)
	return perm[:numSelected]
}

func (rf *RandomForestRegressor) Predict(features []float64) float64 {
	if len(rf.trees) == 0 {
		return 0.0
	}

	sum := 0.0
	for _, tree := range rf.trees {
		sum += tree.Predict(features)
	}
	return sum / float64(len(rf.trees))
}

func (rf *RandomForestRegressor) PredictBatch(dataset *data.RegressionDataset) []float64 {
	predictions := make([]float64, len(dataset.Instances))
	for i, inst := range dataset.Instances {
		predictions[i] = rf.Predict(inst.Features)
	}
	return predictions
}

