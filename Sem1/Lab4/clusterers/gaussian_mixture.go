package clusterers

import (
	"lab4/data"
	"math"
	"math/rand"
)

type GaussianMixture struct {
	nComponents int
	means       [][]float64
	covariances [][]float64
	weights     []float64
	labels      []int
	numFeatures int
}

func NewGaussianMixture(nComponents int) *GaussianMixture {
	return &GaussianMixture{
		nComponents: nComponents,
	}
}

func (gm *GaussianMixture) Train(dataset *data.Dataset, maxIterations int) {
	if len(dataset.Instances) == 0 {
		return
	}

	gm.numFeatures = dataset.NumFeatures
	n := len(dataset.Instances)

	gm.means = make([][]float64, gm.nComponents)
	gm.covariances = make([][]float64, gm.nComponents)
	gm.weights = make([]float64, gm.nComponents)

	for i := 0; i < gm.nComponents; i++ {
		gm.means[i] = make([]float64, gm.numFeatures)
		randomIdx := rand.Intn(n)
		copy(gm.means[i], dataset.Instances[randomIdx].Features)

		gm.covariances[i] = make([]float64, gm.numFeatures)
		for j := range gm.covariances[i] {
			gm.covariances[i][j] = 1.0
		}

		gm.weights[i] = 1.0 / float64(gm.nComponents)
	}

	for iter := 0; iter < maxIterations; iter++ {
		responsibilities := make([][]float64, n)
		for i := range responsibilities {
			responsibilities[i] = make([]float64, gm.nComponents)
		}

		for i, inst := range dataset.Instances {
			total := 0.0
			for j := 0; j < gm.nComponents; j++ {
				prob := gm.gaussianPDF(inst.Features, gm.means[j], gm.covariances[j])
				responsibilities[i][j] = gm.weights[j] * prob
				total += responsibilities[i][j]
			}

			if total > 1e-10 {
				for j := 0; j < gm.nComponents; j++ {
					responsibilities[i][j] /= total
				}
			}
		}

		for j := 0; j < gm.nComponents; j++ {
			sumResp := 0.0
			for i := range responsibilities {
				sumResp += responsibilities[i][j]
			}

			if sumResp < 1e-10 {
				continue
			}

			gm.weights[j] = sumResp / float64(n)

			for f := 0; f < gm.numFeatures; f++ {
				meanSum := 0.0
				for i, inst := range dataset.Instances {
					meanSum += responsibilities[i][j] * inst.Features[f]
				}
				gm.means[j][f] = meanSum / sumResp

				varSum := 0.0
				for i, inst := range dataset.Instances {
					diff := inst.Features[f] - gm.means[j][f]
					varSum += responsibilities[i][j] * diff * diff
				}
				gm.covariances[j][f] = math.Max(varSum/sumResp, 1e-6)
			}
		}
	}

	gm.labels = make([]int, n)
	for i, inst := range dataset.Instances {
		bestCluster := 0
		bestProb := 0.0
		for j := 0; j < gm.nComponents; j++ {
			prob := gm.gaussianPDF(inst.Features, gm.means[j], gm.covariances[j])
			weightedProb := gm.weights[j] * prob
			if weightedProb > bestProb {
				bestProb = weightedProb
				bestCluster = j
			}
		}
		gm.labels[i] = bestCluster
	}
}

func (gm *GaussianMixture) gaussianPDF(x []float64, mean []float64, variance []float64) float64 {
	prob := 1.0
	for i := range x {
		diff := x[i] - mean[i]
		prob *= math.Exp(-0.5*diff*diff/variance[i]) / math.Sqrt(2*math.Pi*variance[i])
	}
	return prob
}

func (gm *GaussianMixture) Predict(features []float64) int {
	bestCluster := 0
	bestProb := 0.0
	for j := 0; j < gm.nComponents; j++ {
		prob := gm.gaussianPDF(features, gm.means[j], gm.covariances[j])
		weightedProb := gm.weights[j] * prob
		if weightedProb > bestProb {
			bestProb = weightedProb
			bestCluster = j
		}
	}
	return bestCluster
}

func (gm *GaussianMixture) PredictBatch(dataset *data.Dataset) []int {
	return gm.labels
}

func (gm *GaussianMixture) GetLabels() []int {
	return gm.labels
}
