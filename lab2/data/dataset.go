package data

import (
	"math"
	"math/rand"
)

type Instance struct {
	Features []float64
	Label    int
}

type Dataset struct {
	Instances   []Instance
	NumFeatures int
	NumClasses  int
}

func NewIrisDataset() *Dataset {
	// Ирисы Фишера - классический датасет
	// 4 признака: длина чашелистика, ширина чашелистика, длина лепестка, ширина лепестка
	// 3 класса: Setosa (0), Versicolor (1), Virginica (2)
	instances := []Instance{
		// Setosa (класс 0)
		{[]float64{5.1, 3.5, 1.4, 0.2}, 0},
		{[]float64{4.9, 3.0, 1.4, 0.2}, 0},
		{[]float64{4.7, 3.2, 1.3, 0.2}, 0},
		{[]float64{4.6, 3.1, 1.5, 0.2}, 0},
		{[]float64{5.0, 3.6, 1.4, 0.2}, 0},
		{[]float64{5.4, 3.9, 1.7, 0.4}, 0},
		{[]float64{4.6, 3.4, 1.4, 0.3}, 0},
		{[]float64{5.0, 3.4, 1.5, 0.2}, 0},
		{[]float64{4.4, 2.9, 1.4, 0.2}, 0},
		{[]float64{4.9, 3.1, 1.5, 0.1}, 0},
		{[]float64{5.4, 3.7, 1.5, 0.2}, 0},
		{[]float64{4.8, 3.4, 1.6, 0.2}, 0},
		{[]float64{4.8, 3.0, 1.4, 0.1}, 0},
		{[]float64{4.3, 3.0, 1.1, 0.1}, 0},
		{[]float64{5.8, 4.0, 1.2, 0.2}, 0},
		{[]float64{5.7, 4.4, 1.5, 0.4}, 0},
		{[]float64{5.4, 3.9, 1.3, 0.4}, 0},
		{[]float64{5.1, 3.5, 1.4, 0.3}, 0},
		{[]float64{5.7, 3.8, 1.7, 0.3}, 0},
		{[]float64{5.1, 3.8, 1.5, 0.3}, 0},
		{[]float64{5.4, 3.4, 1.7, 0.2}, 0},
		{[]float64{5.1, 3.7, 1.5, 0.4}, 0},
		{[]float64{4.6, 3.6, 1.0, 0.2}, 0},
		{[]float64{5.1, 3.3, 1.7, 0.5}, 0},
		{[]float64{4.8, 3.4, 1.9, 0.2}, 0},
		{[]float64{5.0, 3.0, 1.6, 0.2}, 0},
		{[]float64{5.0, 3.4, 1.6, 0.4}, 0},
		{[]float64{5.2, 3.5, 1.5, 0.2}, 0},
		{[]float64{5.2, 3.4, 1.4, 0.2}, 0},
		{[]float64{4.7, 3.2, 1.6, 0.2}, 0},
		{[]float64{4.8, 3.1, 1.6, 0.2}, 0},
		{[]float64{5.4, 3.4, 1.5, 0.4}, 0},
		{[]float64{5.2, 4.1, 1.5, 0.1}, 0},
		{[]float64{5.5, 4.2, 1.4, 0.2}, 0},
		{[]float64{4.9, 3.1, 1.5, 0.1}, 0},
		{[]float64{5.0, 3.2, 1.2, 0.2}, 0},
		{[]float64{5.5, 3.5, 1.3, 0.2}, 0},
		{[]float64{4.9, 3.1, 1.5, 0.1}, 0},
		{[]float64{4.4, 3.0, 1.3, 0.2}, 0},
		{[]float64{5.1, 3.4, 1.5, 0.2}, 0},
		{[]float64{5.0, 3.5, 1.3, 0.3}, 0},
		{[]float64{4.5, 2.3, 1.3, 0.3}, 0},
		{[]float64{4.4, 3.2, 1.3, 0.2}, 0},
		{[]float64{5.0, 3.5, 1.6, 0.6}, 0},
		{[]float64{5.1, 3.8, 1.9, 0.4}, 0},
		{[]float64{4.8, 3.0, 1.4, 0.3}, 0},
		{[]float64{5.1, 3.8, 1.6, 0.2}, 0},
		{[]float64{4.6, 3.2, 1.4, 0.2}, 0},
		{[]float64{5.3, 3.7, 1.5, 0.2}, 0},
		{[]float64{5.0, 3.3, 1.4, 0.2}, 0},

		// Versicolor (класс 1)
		{[]float64{7.0, 3.2, 4.7, 1.4}, 1},
		{[]float64{6.4, 3.2, 4.5, 1.5}, 1},
		{[]float64{6.9, 3.1, 4.9, 1.5}, 1},
		{[]float64{5.5, 2.3, 4.0, 1.3}, 1},
		{[]float64{6.5, 2.8, 4.6, 1.5}, 1},
		{[]float64{5.7, 2.8, 4.5, 1.3}, 1},
		{[]float64{6.3, 3.3, 4.7, 1.6}, 1},
		{[]float64{4.9, 2.4, 3.3, 1.0}, 1},
		{[]float64{6.6, 2.9, 4.6, 1.3}, 1},
		{[]float64{5.2, 2.7, 3.9, 1.4}, 1},
		{[]float64{5.0, 2.0, 3.5, 1.0}, 1},
		{[]float64{5.9, 3.0, 4.2, 1.5}, 1},
		{[]float64{6.0, 2.2, 4.0, 1.0}, 1},
		{[]float64{6.1, 2.9, 4.7, 1.4}, 1},
		{[]float64{5.6, 2.9, 3.6, 1.3}, 1},
		{[]float64{6.7, 3.1, 4.4, 1.4}, 1},
		{[]float64{5.6, 3.0, 4.5, 1.5}, 1},
		{[]float64{5.8, 2.7, 4.1, 1.0}, 1},
		{[]float64{6.2, 2.2, 4.5, 1.5}, 1},
		{[]float64{5.6, 2.5, 3.9, 1.1}, 1},
		{[]float64{5.9, 3.2, 4.8, 1.8}, 1},
		{[]float64{6.1, 2.8, 4.0, 1.3}, 1},
		{[]float64{6.3, 2.5, 4.9, 1.5}, 1},
		{[]float64{6.1, 2.8, 4.7, 1.2}, 1},
		{[]float64{6.4, 2.9, 4.3, 1.3}, 1},
		{[]float64{6.6, 3.0, 4.4, 1.4}, 1},
		{[]float64{6.8, 2.8, 4.8, 1.4}, 1},
		{[]float64{6.7, 3.0, 5.0, 1.7}, 1},
		{[]float64{6.0, 2.9, 4.5, 1.5}, 1},
		{[]float64{5.7, 2.6, 3.5, 1.0}, 1},
		{[]float64{5.5, 2.4, 3.8, 1.1}, 1},
		{[]float64{5.5, 2.4, 3.7, 1.0}, 1},
		{[]float64{5.8, 2.7, 3.9, 1.2}, 1},
		{[]float64{6.0, 2.7, 5.1, 1.6}, 1},
		{[]float64{5.4, 3.0, 4.5, 1.5}, 1},
		{[]float64{6.0, 3.4, 4.5, 1.6}, 1},
		{[]float64{6.7, 3.1, 4.7, 1.5}, 1},
		{[]float64{6.3, 2.3, 4.4, 1.3}, 1},
		{[]float64{5.6, 3.0, 4.1, 1.3}, 1},
		{[]float64{5.5, 2.5, 4.0, 1.3}, 1},
		{[]float64{5.5, 2.6, 4.4, 1.2}, 1},
		{[]float64{6.1, 3.0, 4.6, 1.4}, 1},
		{[]float64{5.8, 2.6, 4.0, 1.2}, 1},
		{[]float64{5.0, 2.3, 3.3, 1.0}, 1},
		{[]float64{5.6, 2.7, 4.2, 1.3}, 1},
		{[]float64{5.7, 3.0, 4.2, 1.2}, 1},
		{[]float64{5.7, 2.9, 4.2, 1.3}, 1},
		{[]float64{6.2, 2.9, 4.3, 1.3}, 1},
		{[]float64{5.1, 2.5, 3.0, 1.1}, 1},
		{[]float64{5.7, 2.8, 4.1, 1.3}, 1},

		// Virginica (класс 2)
		{[]float64{6.3, 3.3, 6.0, 2.5}, 2},
		{[]float64{5.8, 2.7, 5.1, 1.9}, 2},
		{[]float64{7.1, 3.0, 5.9, 2.1}, 2},
		{[]float64{6.3, 2.9, 5.6, 1.8}, 2},
		{[]float64{6.5, 3.0, 5.8, 2.2}, 2},
		{[]float64{7.6, 3.0, 6.6, 2.1}, 2},
		{[]float64{4.9, 2.5, 4.5, 1.7}, 2},
		{[]float64{7.3, 2.9, 6.3, 1.8}, 2},
		{[]float64{6.7, 2.5, 5.8, 1.8}, 2},
		{[]float64{7.2, 3.6, 6.1, 2.5}, 2},
		{[]float64{6.5, 3.2, 5.1, 2.0}, 2},
		{[]float64{6.4, 2.7, 5.3, 1.9}, 2},
		{[]float64{6.8, 3.0, 5.5, 2.1}, 2},
		{[]float64{5.7, 2.5, 5.0, 2.0}, 2},
		{[]float64{5.8, 2.8, 5.1, 2.4}, 2},
		{[]float64{6.4, 3.2, 5.3, 2.3}, 2},
		{[]float64{6.5, 3.0, 5.5, 1.8}, 2},
		{[]float64{7.7, 3.8, 6.7, 2.2}, 2},
		{[]float64{7.7, 2.6, 6.9, 2.3}, 2},
		{[]float64{6.0, 2.2, 5.0, 1.5}, 2},
		{[]float64{6.9, 3.2, 5.7, 2.3}, 2},
		{[]float64{5.6, 2.8, 4.9, 2.0}, 2},
		{[]float64{7.7, 2.8, 6.7, 2.0}, 2},
		{[]float64{6.3, 2.7, 4.9, 1.8}, 2},
		{[]float64{6.7, 3.3, 5.7, 2.1}, 2},
		{[]float64{7.2, 3.2, 6.0, 1.8}, 2},
		{[]float64{6.2, 2.8, 4.8, 1.8}, 2},
		{[]float64{6.1, 3.0, 4.9, 1.8}, 2},
		{[]float64{6.4, 2.8, 5.6, 2.1}, 2},
		{[]float64{7.2, 3.0, 5.8, 1.6}, 2},
		{[]float64{7.4, 2.8, 6.1, 1.9}, 2},
		{[]float64{7.9, 3.8, 6.4, 2.0}, 2},
		{[]float64{6.4, 2.8, 5.6, 2.2}, 2},
		{[]float64{6.3, 2.8, 5.1, 1.5}, 2},
		{[]float64{6.1, 2.6, 5.6, 1.4}, 2},
		{[]float64{7.7, 3.0, 6.1, 2.3}, 2},
		{[]float64{6.3, 3.4, 5.6, 2.4}, 2},
		{[]float64{6.4, 3.1, 5.5, 1.8}, 2},
		{[]float64{6.0, 3.0, 4.8, 1.8}, 2},
		{[]float64{6.9, 3.1, 5.4, 2.1}, 2},
		{[]float64{6.7, 3.1, 5.6, 2.4}, 2},
		{[]float64{6.9, 3.1, 5.1, 2.3}, 2},
		{[]float64{5.8, 2.7, 5.1, 1.9}, 2},
		{[]float64{6.8, 3.2, 5.9, 2.3}, 2},
		{[]float64{6.7, 3.3, 5.7, 2.5}, 2},
		{[]float64{6.7, 3.0, 5.2, 2.3}, 2},
		{[]float64{6.3, 2.5, 5.0, 1.9}, 2},
		{[]float64{6.5, 3.0, 5.2, 2.0}, 2},
		{[]float64{6.2, 3.4, 5.4, 2.3}, 2},
		{[]float64{5.9, 3.0, 5.1, 1.8}, 2},
	}

	return &Dataset{
		Instances:   instances,
		NumFeatures: 4,
		NumClasses:  3,
	}
}

func (d *Dataset) Split(trainRatio float64, seed int64) (*Dataset, *Dataset) {
	rand.Seed(seed)
	perm := rand.Perm(len(d.Instances))

	trainSize := int(float64(len(d.Instances)) * trainRatio)

	train := &Dataset{
		Instances:   make([]Instance, trainSize),
		NumFeatures: d.NumFeatures,
		NumClasses:  d.NumClasses,
	}

	test := &Dataset{
		Instances:   make([]Instance, len(d.Instances)-trainSize),
		NumFeatures: d.NumFeatures,
		NumClasses:  d.NumClasses,
	}

	for i := 0; i < trainSize; i++ {
		train.Instances[i] = d.Instances[perm[i]]
	}

	for i := trainSize; i < len(d.Instances); i++ {
		test.Instances[i-trainSize] = d.Instances[perm[i]]
	}

	return train, test
}

func (d *Dataset) Normalize() {
	if len(d.Instances) == 0 {
		return
	}

	numFeatures := len(d.Instances[0].Features)
	min := make([]float64, numFeatures)
	max := make([]float64, numFeatures)

	for i := range min {
		min[i] = math.Inf(1)
		max[i] = math.Inf(-1)
	}

	for _, inst := range d.Instances {
		for i, val := range inst.Features {
			if val < min[i] {
				min[i] = val
			}
			if val > max[i] {
				max[i] = val
			}
		}
	}

	for i := range d.Instances {
		for j := range d.Instances[i].Features {
			rangeVal := max[j] - min[j]
			if rangeVal > 1e-10 {
				d.Instances[i].Features[j] = (d.Instances[i].Features[j] - min[j]) / rangeVal
			}
		}
	}
}

func EuclideanDistance(a, b []float64) float64 {
	if len(a) != len(b) {
		return math.Inf(1)
	}

	sum := 0.0
	for i := range a {
		diff := a[i] - b[i]
		sum += diff * diff
	}
	return math.Sqrt(sum)
}
