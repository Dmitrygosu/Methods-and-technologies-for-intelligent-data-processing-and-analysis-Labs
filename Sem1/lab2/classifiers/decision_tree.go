package classifiers

import (
	"lab2/data"
)

type DecisionTree struct {
	root *TreeNode
	maxDepth int
}

type TreeNode struct {
	IsLeaf     bool
	Label      int
	FeatureIdx int
	Threshold  float64
	Left       *TreeNode
	Right      *TreeNode
}

func NewDecisionTree(maxDepth int) *DecisionTree {
	return &DecisionTree{
		maxDepth: maxDepth,
	}
}

func (dt *DecisionTree) Train(trainSet *data.Dataset) {
	dt.root = dt.buildTree(trainSet.Instances, 0)
}

func (dt *DecisionTree) buildTree(instances []data.Instance, depth int) *TreeNode {
	if len(instances) == 0 {
		return &TreeNode{IsLeaf: true, Label: 0}
	}
	
	if depth >= dt.maxDepth {
		return &TreeNode{IsLeaf: true, Label: dt.majorityLabel(instances)}
	}
	
	if dt.isPure(instances) {
		return &TreeNode{IsLeaf: true, Label: instances[0].Label}
	}
	
	bestFeature, bestThreshold, _ := dt.findBestSplit(instances)
	if bestFeature == -1 {
		return &TreeNode{IsLeaf: true, Label: dt.majorityLabel(instances)}
	}
	
	left, right := dt.split(instances, bestFeature, bestThreshold)
	
	node := &TreeNode{
		IsLeaf:     false,
		FeatureIdx: bestFeature,
		Threshold:  bestThreshold,
		Left:       dt.buildTree(left, depth+1),
		Right:      dt.buildTree(right, depth+1),
	}
	
	return node
}

func (dt *DecisionTree) isPure(instances []data.Instance) bool {
	if len(instances) == 0 {
		return true
	}
	firstLabel := instances[0].Label
	for _, inst := range instances {
		if inst.Label != firstLabel {
			return false
		}
	}
	return true
}

func (dt *DecisionTree) majorityLabel(instances []data.Instance) int {
	if len(instances) == 0 {
		return 0
	}
	counts := make(map[int]int)
	for _, inst := range instances {
		counts[inst.Label]++
	}
	maxCount := 0
	majority := 0
	for label, count := range counts {
		if count > maxCount {
			maxCount = count
			majority = label
		}
	}
	return majority
}

func (dt *DecisionTree) gini(instances []data.Instance) float64 {
	if len(instances) == 0 {
		return 0
	}
	counts := make(map[int]int)
	for _, inst := range instances {
		counts[inst.Label]++
	}
	
	gini := 1.0
	for _, count := range counts {
		prob := float64(count) / float64(len(instances))
		gini -= prob * prob
	}
	return gini
}

func (dt *DecisionTree) findBestSplit(instances []data.Instance) (int, float64, float64) {
	if len(instances) == 0 {
		return -1, 0, 1.0
	}
	
	numFeatures := len(instances[0].Features)
	bestFeature := -1
	bestThreshold := 0.0
	bestGini := 1.0
	
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
			
			leftGini := dt.gini(left)
			rightGini := dt.gini(right)
			
			leftWeight := float64(len(left)) / float64(len(instances))
			rightWeight := float64(len(right)) / float64(len(instances))
			
			weightedGini := leftWeight*leftGini + rightWeight*rightGini
			
			if weightedGini < bestGini {
				bestGini = weightedGini
				bestFeature = featureIdx
				bestThreshold = threshold
			}
		}
	}
	
	return bestFeature, bestThreshold, bestGini
}

func (dt *DecisionTree) split(instances []data.Instance, featureIdx int, threshold float64) ([]data.Instance, []data.Instance) {
	left := make([]data.Instance, 0)
	right := make([]data.Instance, 0)
	
	for _, inst := range instances {
		if inst.Features[featureIdx] <= threshold {
			left = append(left, inst)
		} else {
			right = append(right, inst)
		}
	}
	
	return left, right
}

func (dt *DecisionTree) Predict(features []float64) int {
	node := dt.root
	for node != nil && !node.IsLeaf {
		if features[node.FeatureIdx] <= node.Threshold {
			node = node.Left
		} else {
			node = node.Right
		}
	}
	if node == nil {
		return 0
	}
	return node.Label
}

func (dt *DecisionTree) PredictBatch(testSet *data.Dataset) []int {
	predictions := make([]int, len(testSet.Instances))
	for i, inst := range testSet.Instances {
		predictions[i] = dt.Predict(inst.Features)
	}
	return predictions
}

