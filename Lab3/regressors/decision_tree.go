package regressors

import (
	"lab3/data"
	"math"
)

type DecisionTreeRegressor struct {
	root     *RegressorTreeNode
	maxDepth int
}

type RegressorTreeNode struct {
	IsLeaf     bool
	Value      float64
	FeatureIdx int
	Threshold  float64
	Left       *RegressorTreeNode
	Right      *RegressorTreeNode
}

func NewDecisionTreeRegressor(maxDepth int) *DecisionTreeRegressor {
	return &DecisionTreeRegressor{
		maxDepth: maxDepth,
	}
}

func (dt *DecisionTreeRegressor) Train(dataset *data.RegressionDataset) {
	dt.root = dt.buildTree(dataset.Instances, 0)
}

func (dt *DecisionTreeRegressor) buildTree(instances []data.RegressionInstance, depth int) *RegressorTreeNode {
	if len(instances) == 0 {
		return &RegressorTreeNode{IsLeaf: true, Value: 0.0}
	}

	if depth >= dt.maxDepth {
		return &RegressorTreeNode{IsLeaf: true, Value: dt.meanValue(instances)}
	}

	if dt.isPure(instances) {
		return &RegressorTreeNode{IsLeaf: true, Value: instances[0].Target}
	}

	bestFeature, bestThreshold, _ := dt.findBestSplit(instances)
	if bestFeature == -1 {
		return &RegressorTreeNode{IsLeaf: true, Value: dt.meanValue(instances)}
	}

	left, right := dt.split(instances, bestFeature, bestThreshold)

	node := &RegressorTreeNode{
		IsLeaf:     false,
		FeatureIdx: bestFeature,
		Threshold:  bestThreshold,
		Left:       dt.buildTree(left, depth+1),
		Right:      dt.buildTree(right, depth+1),
	}

	return node
}

func (dt *DecisionTreeRegressor) isPure(instances []data.RegressionInstance) bool {
	if len(instances) <= 1 {
		return true
	}

	firstValue := instances[0].Target
	epsilon := 1e-6
	for _, inst := range instances {
		if math.Abs(inst.Target-firstValue) > epsilon {
			return false
		}
	}
	return true
}

func (dt *DecisionTreeRegressor) meanValue(instances []data.RegressionInstance) float64 {
	if len(instances) == 0 {
		return 0.0
	}

	sum := 0.0
	for _, inst := range instances {
		sum += inst.Target
	}
	return sum / float64(len(instances))
}

func (dt *DecisionTreeRegressor) variance(instances []data.RegressionInstance) float64 {
	if len(instances) == 0 {
		return 0.0
	}

	mean := dt.meanValue(instances)
	sumSqDiff := 0.0
	for _, inst := range instances {
		diff := inst.Target - mean
		sumSqDiff += diff * diff
	}
	return sumSqDiff / float64(len(instances))
}

func (dt *DecisionTreeRegressor) findBestSplit(instances []data.RegressionInstance) (int, float64, float64) {
	if len(instances) == 0 {
		return -1, 0, 1e10
	}

	numFeatures := len(instances[0].Features)
	bestFeature := -1
	bestThreshold := 0.0
	bestVariance := 1e10

	for featureIdx := 0; featureIdx < numFeatures; featureIdx++ {
		values := make([]float64, len(instances))
		for i, inst := range instances {
			values[i] = inst.Features[featureIdx]
		}

		for _, threshold := range values {
			left, right := dt.split(instances, featureIdx, threshold)
			if len(left) == 0 || len(right) == 0 {
				continue
			}

			leftVariance := dt.variance(left)
			rightVariance := dt.variance(right)

			leftWeight := float64(len(left)) / float64(len(instances))
			rightWeight := float64(len(right)) / float64(len(instances))

			weightedVariance := leftWeight*leftVariance + rightWeight*rightVariance

			if weightedVariance < bestVariance {
				bestVariance = weightedVariance
				bestFeature = featureIdx
				bestThreshold = threshold
			}
		}
	}

	return bestFeature, bestThreshold, bestVariance
}

func (dt *DecisionTreeRegressor) split(instances []data.RegressionInstance, featureIdx int, threshold float64) ([]data.RegressionInstance, []data.RegressionInstance) {
	left := make([]data.RegressionInstance, 0)
	right := make([]data.RegressionInstance, 0)

	for _, inst := range instances {
		if inst.Features[featureIdx] <= threshold {
			left = append(left, inst)
		} else {
			right = append(right, inst)
		}
	}

	return left, right
}

func (dt *DecisionTreeRegressor) Predict(features []float64) float64 {
	node := dt.root
	for node != nil && !node.IsLeaf {
		if features[node.FeatureIdx] <= node.Threshold {
			node = node.Left
		} else {
			node = node.Right
		}
	}
	if node == nil {
		return 0.0
	}
	return node.Value
}

func (dt *DecisionTreeRegressor) PredictBatch(dataset *data.RegressionDataset) []float64 {
	predictions := make([]float64, len(dataset.Instances))
	for i, inst := range dataset.Instances {
		predictions[i] = dt.Predict(inst.Features)
	}
	return predictions
}

func (dt *DecisionTreeRegressor) trainOnSubset(instances []data.RegressionInstance, featureSubset []int) {
	dt.root = dt.buildTreeWithFeatureSubset(instances, 0, featureSubset)
}

func (dt *DecisionTreeRegressor) buildTreeWithFeatureSubset(instances []data.RegressionInstance, depth int, featureSubset []int) *RegressorTreeNode {
	if len(instances) == 0 {
		return &RegressorTreeNode{IsLeaf: true, Value: 0.0}
	}

	if depth >= dt.maxDepth {
		return &RegressorTreeNode{IsLeaf: true, Value: dt.meanValue(instances)}
	}

	if dt.isPure(instances) {
		return &RegressorTreeNode{IsLeaf: true, Value: instances[0].Target}
	}

	bestFeature, bestThreshold, _ := dt.findBestSplitWithSubset(instances, featureSubset)
	if bestFeature == -1 {
		return &RegressorTreeNode{IsLeaf: true, Value: dt.meanValue(instances)}
	}

	left, right := dt.split(instances, bestFeature, bestThreshold)

	node := &RegressorTreeNode{
		IsLeaf:     false,
		FeatureIdx: bestFeature,
		Threshold:  bestThreshold,
		Left:       dt.buildTreeWithFeatureSubset(left, depth+1, featureSubset),
		Right:      dt.buildTreeWithFeatureSubset(right, depth+1, featureSubset),
	}

	return node
}

func (dt *DecisionTreeRegressor) findBestSplitWithSubset(instances []data.RegressionInstance, featureSubset []int) (int, float64, float64) {
	if len(instances) == 0 {
		return -1, 0, 1e10
	}

	bestFeature := -1
	bestThreshold := 0.0
	bestVariance := 1e10

	for _, featureIdx := range featureSubset {
		values := make([]float64, len(instances))
		for i, inst := range instances {
			values[i] = inst.Features[featureIdx]
		}

		for _, threshold := range values {
			left, right := dt.split(instances, featureIdx, threshold)
			if len(left) == 0 || len(right) == 0 {
				continue
			}

			leftVariance := dt.variance(left)
			rightVariance := dt.variance(right)

			leftWeight := float64(len(left)) / float64(len(instances))
			rightWeight := float64(len(right)) / float64(len(instances))

			weightedVariance := leftWeight*leftVariance + rightWeight*rightVariance

			if weightedVariance < bestVariance {
				bestVariance = weightedVariance
				bestFeature = featureIdx
				bestThreshold = threshold
			}
		}
	}

	return bestFeature, bestThreshold, bestVariance
}

