package regressors

import (
	"lab3/data"
)

type LinearRegression struct {
	weights []float64
	bias    float64
}

func NewLinearRegression() *LinearRegression {
	return &LinearRegression{}
}

func (lr *LinearRegression) Train(dataset *data.RegressionDataset) {
	numFeatures := dataset.NumFeatures
	numInstances := len(dataset.Instances)

	if numInstances == 0 {
		return
	}

	lr.weights = make([]float64, numFeatures)
	lr.bias = 0.0

	learningRate := 0.01
	iterations := 1000

	for iter := 0; iter < iterations; iter++ {
		totalError := 0.0
		gradWeights := make([]float64, numFeatures)
		gradBias := 0.0

		for _, inst := range dataset.Instances {
			prediction := lr.Predict(inst.Features)
			error := inst.Target - prediction
			totalError += error * error

			for i := 0; i < numFeatures; i++ {
				gradWeights[i] += -2.0 * error * inst.Features[i]
			}
			gradBias += -2.0 * error
		}

		for i := 0; i < numFeatures; i++ {
			lr.weights[i] -= learningRate * gradWeights[i] / float64(numInstances)
		}
		lr.bias -= learningRate * gradBias / float64(numInstances)
	}
}

func (lr *LinearRegression) Predict(features []float64) float64 {
	if len(lr.weights) == 0 {
		return 0.0
	}

	prediction := lr.bias
	for i := 0; i < len(features) && i < len(lr.weights); i++ {
		prediction += lr.weights[i] * features[i]
	}
	return prediction
}

func (lr *LinearRegression) PredictBatch(dataset *data.RegressionDataset) []float64 {
	predictions := make([]float64, len(dataset.Instances))
	for i, inst := range dataset.Instances {
		predictions[i] = lr.Predict(inst.Features)
	}
	return predictions
}

type RidgeRegression struct {
	weights []float64
	bias    float64
	alpha   float64
}

func NewRidgeRegression(alpha float64) *RidgeRegression {
	return &RidgeRegression{
		alpha: alpha,
	}
}

func (rr *RidgeRegression) Train(dataset *data.RegressionDataset) {
	numFeatures := dataset.NumFeatures
	numInstances := len(dataset.Instances)

	if numInstances == 0 {
		return
	}

	rr.weights = make([]float64, numFeatures)
	rr.bias = 0.0

	learningRate := 0.01
	iterations := 1000

	for iter := 0; iter < iterations; iter++ {
		gradWeights := make([]float64, numFeatures)
		gradBias := 0.0

		for _, inst := range dataset.Instances {
			prediction := rr.Predict(inst.Features)
			error := inst.Target - prediction

			for i := 0; i < numFeatures; i++ {
				gradWeights[i] += -2.0*error*inst.Features[i] + 2.0*rr.alpha*rr.weights[i]
			}
			gradBias += -2.0 * error
		}

		for i := 0; i < numFeatures; i++ {
			rr.weights[i] -= learningRate * gradWeights[i] / float64(numInstances)
		}
		rr.bias -= learningRate * gradBias / float64(numInstances)
	}
}

func (rr *RidgeRegression) Predict(features []float64) float64 {
	if len(rr.weights) == 0 {
		return 0.0
	}

	prediction := rr.bias
	for i := 0; i < len(features) && i < len(rr.weights); i++ {
		prediction += rr.weights[i] * features[i]
	}
	return prediction
}

func (rr *RidgeRegression) PredictBatch(dataset *data.RegressionDataset) []float64 {
	predictions := make([]float64, len(dataset.Instances))
	for i, inst := range dataset.Instances {
		predictions[i] = rr.Predict(inst.Features)
	}
	return predictions
}

type LassoRegression struct {
	weights []float64
	bias    float64
	alpha   float64
}

func NewLassoRegression(alpha float64) *LassoRegression {
	return &LassoRegression{
		alpha: alpha,
	}
}

func (las *LassoRegression) Train(dataset *data.RegressionDataset) {
	numFeatures := dataset.NumFeatures
	numInstances := len(dataset.Instances)

	if numInstances == 0 {
		return
	}

	las.weights = make([]float64, numFeatures)
	las.bias = 0.0

	learningRate := 0.01
	iterations := 1000

	for iter := 0; iter < iterations; iter++ {
		gradWeights := make([]float64, numFeatures)
		gradBias := 0.0

		for _, inst := range dataset.Instances {
			prediction := las.Predict(inst.Features)
			error := inst.Target - prediction

			for i := 0; i < numFeatures; i++ {
				sign := 1.0
				if las.weights[i] < 0 {
					sign = -1.0
				}
				gradWeights[i] += -2.0*error*inst.Features[i] + las.alpha*sign
			}
			gradBias += -2.0 * error
		}

		for i := 0; i < numFeatures; i++ {
			las.weights[i] -= learningRate * gradWeights[i] / float64(numInstances)
		}
		las.bias -= learningRate * gradBias / float64(numInstances)
	}
}

func (las *LassoRegression) Predict(features []float64) float64 {
	if len(las.weights) == 0 {
		return 0.0
	}

	prediction := las.bias
	for i := 0; i < len(features) && i < len(las.weights); i++ {
		prediction += las.weights[i] * features[i]
	}
	return prediction
}

func (las *LassoRegression) PredictBatch(dataset *data.RegressionDataset) []float64 {
	predictions := make([]float64, len(dataset.Instances))
	for i, inst := range dataset.Instances {
		predictions[i] = las.Predict(inst.Features)
	}
	return predictions
}

