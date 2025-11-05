package classifiers

import (
	"lab2/data"
	"math"
	"math/rand"
)

type LogisticRegression struct {
	weights   [][]float64
	numClasses int
	numFeatures int
	learningRate float64
	iterations int
}

func NewLogisticRegression(learningRate float64, iterations int) *LogisticRegression {
	return &LogisticRegression{
		learningRate: learningRate,
		iterations:   iterations,
	}
}

func (lr *LogisticRegression) Train(trainSet *data.Dataset) {
	lr.numClasses = trainSet.NumClasses
	lr.numFeatures = trainSet.NumFeatures
	
	lr.weights = make([][]float64, lr.numClasses)
	for i := range lr.weights {
		lr.weights[i] = make([]float64, lr.numFeatures+1)
		for j := range lr.weights[i] {
			lr.weights[i][j] = rand.Float64()*0.01 - 0.005
		}
	}
	
	for iter := 0; iter < lr.iterations; iter++ {
		for _, inst := range trainSet.Instances {
			predictions := lr.predictProbabilities(inst.Features)
			
			for class := 0; class < lr.numClasses; class++ {
				target := 0.0
				if inst.Label == class {
					target = 1.0
				}
				
				error := target - predictions[class]
				
				lr.weights[class][0] += lr.learningRate * error
				for i := 0; i < lr.numFeatures; i++ {
					lr.weights[class][i+1] += lr.learningRate * error * inst.Features[i]
				}
			}
		}
	}
}

func (lr *LogisticRegression) predictProbabilities(features []float64) []float64 {
	scores := make([]float64, lr.numClasses)
	
	for class := 0; class < lr.numClasses; class++ {
		score := lr.weights[class][0]
		for i := 0; i < lr.numFeatures; i++ {
			score += lr.weights[class][i+1] * features[i]
		}
		scores[class] = score
	}
	
	maxScore := math.Inf(-1)
	for _, score := range scores {
		if score > maxScore {
			maxScore = score
		}
	}
	
	sumExp := 0.0
	for i := range scores {
		scores[i] = math.Exp(scores[i] - maxScore)
		sumExp += scores[i]
	}
	
	for i := range scores {
		scores[i] /= sumExp
	}
	
	return scores
}

func (lr *LogisticRegression) Predict(features []float64) int {
	probs := lr.predictProbabilities(features)
	maxProb := 0.0
	predictedClass := 0
	for i, prob := range probs {
		if prob > maxProb {
			maxProb = prob
			predictedClass = i
		}
	}
	return predictedClass
}

func (lr *LogisticRegression) PredictBatch(testSet *data.Dataset) []int {
	predictions := make([]int, len(testSet.Instances))
	for i, inst := range testSet.Instances {
		predictions[i] = lr.Predict(inst.Features)
	}
	return predictions
}

