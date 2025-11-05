package classifiers

import (
	"lab2/data"
	"math"
)

type NaiveBayes struct {
	classProbs   map[int]float64
	featureProbs map[int][]Gaussian
	numClasses   int
	numFeatures  int
}

type Gaussian struct {
	Mean float64
	Std  float64
}

func NewNaiveBayes() *NaiveBayes {
	return &NaiveBayes{
		classProbs:   make(map[int]float64),
		featureProbs: make(map[int][]Gaussian),
	}
}

func (nb *NaiveBayes) Train(trainSet *data.Dataset) {
	nb.numClasses = trainSet.NumClasses
	nb.numFeatures = trainSet.NumFeatures
	
	classCounts := make(map[int]int)
	classInstances := make(map[int][]data.Instance)
	
	for _, inst := range trainSet.Instances {
		classCounts[inst.Label]++
		classInstances[inst.Label] = append(classInstances[inst.Label], inst)
	}
	
	totalInstances := float64(len(trainSet.Instances))
	for class, count := range classCounts {
		nb.classProbs[class] = float64(count) / totalInstances
	}
	
	for class := 0; class < nb.numClasses; class++ {
		instances := classInstances[class]
		if len(instances) == 0 {
			continue
		}
		
		features := make([][]float64, nb.numFeatures)
		for i := range features {
			features[i] = make([]float64, len(instances))
		}
		
		for i, inst := range instances {
			for j := range inst.Features {
				features[j][i] = inst.Features[j]
			}
		}
		
		gaussians := make([]Gaussian, nb.numFeatures)
		for i := range features {
			mean := 0.0
			for _, val := range features[i] {
				mean += val
			}
			mean /= float64(len(features[i]))
			
			variance := 0.0
			for _, val := range features[i] {
				diff := val - mean
				variance += diff * diff
			}
			variance /= float64(len(features[i]))
			std := math.Sqrt(variance)
			
			if std < 1e-10 {
				std = 1e-10
			}
			
			gaussians[i] = Gaussian{Mean: mean, Std: std}
		}
		
		nb.featureProbs[class] = gaussians
	}
}

func (nb *NaiveBayes) Predict(features []float64) int {
	maxProb := math.Inf(-1)
	predictedClass := 0
	
	for class := 0; class < nb.numClasses; class++ {
		classProb := math.Log(nb.classProbs[class] + 1e-10)
		
		gaussians, exists := nb.featureProbs[class]
		if !exists {
			continue
		}
		
		for i, feature := range features {
			if i >= len(gaussians) {
				continue
			}
			g := gaussians[i]
			exponent := -0.5 * math.Pow((feature-g.Mean)/g.Std, 2)
			prob := exponent - math.Log(g.Std*math.Sqrt(2*math.Pi))
			classProb += prob
		}
		
		if classProb > maxProb {
			maxProb = classProb
			predictedClass = class
		}
	}
	
	return predictedClass
}

func (nb *NaiveBayes) PredictBatch(testSet *data.Dataset) []int {
	predictions := make([]int, len(testSet.Instances))
	for i, inst := range testSet.Instances {
		predictions[i] = nb.Predict(inst.Features)
	}
	return predictions
}

