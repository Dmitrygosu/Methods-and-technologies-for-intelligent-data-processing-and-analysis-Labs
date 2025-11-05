package metrics

import "math"

type RegressionMetrics struct {
	R2    float64
	RMSE  float64
	MAE   float64
	MAPE  float64
}

func CalculateMetrics(yTrue, yPred []float64) RegressionMetrics {
	if len(yTrue) != len(yPred) || len(yTrue) == 0 {
		return RegressionMetrics{}
	}

	var sumSqResiduals, sumSqTotal float64
	var sumAbsError, sumPercError float64
	meanTrue := mean(yTrue)

	for i := 0; i < len(yTrue); i++ {
		residual := yTrue[i] - yPred[i]
		sumSqResiduals += residual * residual
		sumSqTotal += (yTrue[i] - meanTrue) * (yTrue[i] - meanTrue)
		sumAbsError += math.Abs(residual)
		
		if math.Abs(yTrue[i]) > 1e-10 {
			sumPercError += math.Abs(residual / yTrue[i])
		}
	}

	n := float64(len(yTrue))
	r2 := 1.0 - (sumSqResiduals / sumSqTotal)
	if sumSqTotal < 1e-10 {
		r2 = 0.0
	}

	rmse := math.Sqrt(sumSqResiduals / n)
	mae := sumAbsError / n
	mape := (sumPercError / n) * 100.0

	return RegressionMetrics{
		R2:   r2,
		RMSE: rmse,
		MAE:  mae,
		MAPE: mape,
	}
}

func mean(values []float64) float64 {
	if len(values) == 0 {
		return 0.0
	}

	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

