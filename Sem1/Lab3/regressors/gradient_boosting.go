package regressors

import (
	"lab3/data"
)

type GradientBoostingRegressor struct {
	trees       []*DecisionTreeRegressor
	numTrees    int
	maxDepth    int
	learningRate float64
	initialPred float64
}

func NewGradientBoostingRegressor(numTrees, maxDepth int, learningRate float64) *GradientBoostingRegressor {
	return &GradientBoostingRegressor{
		numTrees:     numTrees,
		maxDepth:     maxDepth,
		learningRate: learningRate,
	}
}

func (gb *GradientBoostingRegressor) Train(dataset *data.RegressionDataset) {
	if len(dataset.Instances) == 0 {
		return
	}

	gb.initialPred = gb.meanTarget(dataset.Instances)
	gb.trees = make([]*DecisionTreeRegressor, gb.numTrees)

	residuals := make([]float64, len(dataset.Instances))
	for i, inst := range dataset.Instances {
		residuals[i] = inst.Target - gb.initialPred
	}

	for i := 0; i < gb.numTrees; i++ {
		treeDataset := gb.createResidualDataset(dataset, residuals)
		tree := NewDecisionTreeRegressor(gb.maxDepth)
		tree.Train(treeDataset)

		predictions := tree.PredictBatch(dataset)
		for j := range residuals {
			residuals[j] -= gb.learningRate * predictions[j]
		}

		gb.trees[i] = tree
	}
}

func (gb *GradientBoostingRegressor) meanTarget(instances []data.RegressionInstance) float64 {
	if len(instances) == 0 {
		return 0.0
	}

	sum := 0.0
	for _, inst := range instances {
		sum += inst.Target
	}
	return sum / float64(len(instances))
}

func (gb *GradientBoostingRegressor) createResidualDataset(original *data.RegressionDataset, residuals []float64) *data.RegressionDataset {
	instances := make([]data.RegressionInstance, len(original.Instances))
	for i, inst := range original.Instances {
		instances[i] = data.RegressionInstance{
			Features: inst.Features,
			Target:   residuals[i],
		}
	}

	return &data.RegressionDataset{
		Instances:   instances,
		NumFeatures: original.NumFeatures,
	}
}

func (gb *GradientBoostingRegressor) Predict(features []float64) float64 {
	prediction := gb.initialPred

	for _, tree := range gb.trees {
		prediction += gb.learningRate * tree.Predict(features)
	}

	return prediction
}

func (gb *GradientBoostingRegressor) PredictBatch(dataset *data.RegressionDataset) []float64 {
	predictions := make([]float64, len(dataset.Instances))
	for i, inst := range dataset.Instances {
		predictions[i] = gb.Predict(inst.Features)
	}
	return predictions
}

