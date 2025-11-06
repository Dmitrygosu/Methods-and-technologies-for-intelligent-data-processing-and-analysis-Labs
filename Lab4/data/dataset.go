package data

import (
	"math"
	"math/rand"
)

type Instance struct {
	Features []float64
}

type Dataset struct {
	Instances   []Instance
	NumFeatures int
}

func NewCustomerDataset() *Dataset {
	instances := []Instance{
		// Возраст, средний чек (тыс. руб.), частота покупок (раз/мес), время на сайте (мин/нед)
		{[]float64{25, 3.5, 2, 45}},
		{[]float64{28, 4.2, 3, 52}},
		{[]float64{32, 5.8, 4, 68}},
		{[]float64{35, 6.5, 5, 75}},
		{[]float64{22, 2.8, 1, 35}},
		{[]float64{45, 8.5, 6, 92}},
		{[]float64{38, 7.2, 5, 85}},
		{[]float64{42, 9.1, 7, 98}},
		{[]float64{29, 4.5, 3, 58}},
		{[]float64{31, 5.2, 4, 65}},
		{[]float64{27, 3.8, 2, 48}},
		{[]float64{33, 6.0, 4, 72}},
		{[]float64{40, 7.8, 6, 88}},
		{[]float64{36, 6.8, 5, 78}},
		{[]float64{44, 9.5, 7, 105}},
		{[]float64{26, 3.2, 2, 42}},
		{[]float64{30, 4.8, 3, 62}},
		{[]float64{34, 6.2, 5, 70}},
		{[]float64{39, 7.5, 6, 90}},
		{[]float64{43, 8.8, 7, 95}},
		{[]float64{24, 3.0, 1, 38}},
		{[]float64{28, 4.0, 3, 55}},
		{[]float64{32, 5.5, 4, 66}},
		{[]float64{37, 7.0, 5, 82}},
		{[]float64{41, 8.2, 6, 93}},
		{[]float64{23, 2.5, 1, 32}},
		{[]float64{29, 4.3, 3, 56}},
		{[]float64{33, 5.9, 4, 69}},
		{[]float64{38, 7.3, 5, 86}},
		{[]float64{42, 8.9, 7, 96}},
		{[]float64{25, 3.3, 2, 40}},
		{[]float64{30, 4.7, 3, 60}},
		{[]float64{35, 6.3, 5, 73}},
		{[]float64{40, 7.6, 6, 87}},
		{[]float64{45, 9.2, 7, 102}},
		{[]float64{26, 3.4, 2, 43}},
		{[]float64{31, 5.0, 4, 64}},
		{[]float64{36, 6.6, 5, 76}},
		{[]float64{41, 8.0, 6, 91}},
		{[]float64{44, 9.3, 7, 100}},
		{[]float64{27, 3.6, 2, 46}},
		{[]float64{32, 5.6, 4, 67}},
		{[]float64{37, 7.1, 5, 83}},
		{[]float64{42, 8.5, 7, 94}},
		{[]float64{24, 2.9, 1, 36}},
		{[]float64{28, 4.1, 3, 53}},
		{[]float64{33, 5.7, 4, 71}},
		{[]float64{39, 7.4, 6, 89}},
		{[]float64{43, 8.7, 7, 97}},
		{[]float64{25, 3.7, 2, 47}},
		{[]float64{30, 4.6, 3, 59}},
		{[]float64{35, 6.4, 5, 74}},
		{[]float64{40, 7.7, 6, 88}},
		{[]float64{45, 9.4, 7, 103}},
		{[]float64{26, 3.9, 2, 49}},
		{[]float64{31, 5.1, 4, 63}},
		{[]float64{36, 6.7, 5, 77}},
		{[]float64{41, 8.1, 6, 92}},
		{[]float64{44, 9.6, 7, 104}},
		{[]float64{27, 3.1, 2, 41}},
		{[]float64{32, 5.4, 4, 68}},
		{[]float64{37, 6.9, 5, 81}},
		{[]float64{42, 8.3, 7, 95}},
		{[]float64{23, 2.7, 1, 34}},
		{[]float64{28, 4.4, 3, 57}},
		{[]float64{33, 5.8, 4, 70}},
		{[]float64{39, 7.6, 6, 90}},
		{[]float64{43, 8.6, 7, 99}},
		{[]float64{25, 3.8, 2, 50}},
		{[]float64{30, 4.9, 3, 61}},
		{[]float64{35, 6.1, 5, 72}},
		{[]float64{40, 7.9, 6, 89}},
		{[]float64{45, 9.7, 7, 106}},
		{[]float64{26, 3.2, 2, 44}},
		{[]float64{31, 5.3, 4, 65}},
		{[]float64{36, 6.5, 5, 75}},
		{[]float64{41, 8.4, 6, 93}},
		{[]float64{44, 9.8, 7, 107}},
		{[]float64{27, 3.5, 2, 48}},
		{[]float64{32, 5.7, 4, 69}},
		{[]float64{37, 7.2, 5, 84}},
		{[]float64{42, 8.7, 7, 96}},
		{[]float64{24, 2.6, 1, 33}},
		{[]float64{28, 4.2, 3, 54}},
		{[]float64{33, 5.9, 4, 70}},
		{[]float64{39, 7.3, 6, 87}},
		{[]float64{43, 8.9, 7, 98}},
		{[]float64{25, 3.4, 2, 45}},
		{[]float64{30, 4.8, 3, 60}},
		{[]float64{35, 6.2, 5, 73}},
		{[]float64{40, 7.8, 6, 86}},
		{[]float64{45, 9.5, 7, 101}},
		{[]float64{26, 3.6, 2, 46}},
		{[]float64{31, 5.2, 4, 64}},
		{[]float64{36, 6.8, 5, 78}},
		{[]float64{41, 8.2, 6, 91}},
		{[]float64{44, 9.9, 7, 105}},
		{[]float64{27, 3.3, 2, 42}},
		{[]float64{32, 5.5, 4, 67}},
		{[]float64{37, 7.0, 5, 82}},
		{[]float64{42, 8.6, 7, 94}},
		{[]float64{23, 2.8, 1, 37}},
		{[]float64{28, 4.3, 3, 56}},
		{[]float64{33, 5.6, 4, 71}},
		{[]float64{39, 7.5, 6, 88}},
		{[]float64{43, 8.8, 7, 97}},
		{[]float64{25, 3.9, 2, 51}},
		{[]float64{30, 4.7, 3, 62}},
		{[]float64{35, 6.3, 5, 74}},
		{[]float64{40, 7.7, 6, 90}},
		{[]float64{45, 9.6, 7, 104}},
		{[]float64{26, 3.1, 2, 43}},
		{[]float64{31, 5.1, 4, 66}},
		{[]float64{36, 6.6, 5, 76}},
		{[]float64{41, 8.3, 6, 92}},
		{[]float64{44, 9.7, 7, 103}},
		{[]float64{27, 3.7, 2, 49}},
		{[]float64{32, 5.8, 4, 68}},
		{[]float64{37, 7.1, 5, 85}},
		{[]float64{42, 8.5, 7, 95}},
		{[]float64{24, 2.9, 1, 39}},
		{[]float64{28, 4.4, 3, 55}},
		{[]float64{33, 5.7, 4, 72}},
		{[]float64{39, 7.4, 6, 89}},
		{[]float64{43, 8.9, 7, 99}},
		{[]float64{25, 3.2, 2, 44}},
		{[]float64{30, 4.6, 3, 61}},
		{[]float64{35, 6.4, 5, 75}},
		{[]float64{40, 7.6, 6, 87}},
		{[]float64{45, 9.8, 7, 102}},
		{[]float64{26, 3.8, 2, 50}},
		{[]float64{31, 5.0, 4, 65}},
		{[]float64{36, 6.7, 5, 79}},
		{[]float64{41, 8.1, 6, 93}},
		{[]float64{44, 9.9, 7, 106}},
		{[]float64{27, 3.4, 2, 47}},
		{[]float64{32, 5.6, 4, 69}},
		{[]float64{37, 6.9, 5, 83}},
		{[]float64{42, 8.4, 7, 96}},
		{[]float64{23, 2.7, 1, 35}},
		{[]float64{28, 4.1, 3, 57}},
		{[]float64{33, 5.9, 4, 70}},
		{[]float64{39, 7.2, 6, 86}},
		{[]float64{43, 8.7, 7, 98}},
		{[]float64{25, 3.5, 2, 48}},
		{[]float64{30, 4.9, 3, 63}},
		{[]float64{35, 6.1, 5, 73}},
		{[]float64{40, 7.9, 6, 91}},
		{[]float64{45, 9.4, 7, 100}},
		{[]float64{26, 3.3, 2, 45}},
		{[]float64{31, 5.4, 4, 67}},
		{[]float64{36, 6.5, 5, 77}},
		{[]float64{41, 8.0, 6, 90}},
		{[]float64{44, 9.5, 7, 101}},
		{[]float64{27, 3.6, 2, 46}},
		{[]float64{32, 5.7, 4, 68}},
		{[]float64{37, 7.0, 5, 84}},
		{[]float64{42, 8.8, 7, 97}},
		{[]float64{24, 3.0, 1, 40}},
		{[]float64{28, 4.5, 3, 58}},
		{[]float64{33, 5.8, 4, 71}},
		{[]float64{39, 7.7, 6, 88}},
		{[]float64{43, 8.6, 7, 99}},
		{[]float64{25, 3.1, 2, 41}},
		{[]float64{30, 4.7, 3, 60}},
		{[]float64{35, 6.2, 5, 74}},
		{[]float64{40, 7.8, 6, 89}},
		{[]float64{45, 9.3, 7, 103}},
		{[]float64{26, 3.9, 2, 51}},
		{[]float64{31, 5.3, 4, 66}},
		{[]float64{36, 6.8, 5, 78}},
		{[]float64{41, 8.2, 6, 92}},
		{[]float64{44, 9.6, 7, 104}},
		{[]float64{27, 3.2, 2, 42}},
		{[]float64{32, 5.5, 4, 69}},
		{[]float64{37, 7.3, 5, 85}},
		{[]float64{42, 8.9, 7, 98}},
	}

	return &Dataset{
		Instances:   instances,
		NumFeatures: 4,
	}
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

func (d *Dataset) BootstrapSample(seed int64) *Dataset {
	rng := rand.New(rand.NewSource(seed))
	sample := make([]Instance, len(d.Instances))

	for i := 0; i < len(d.Instances); i++ {
		idx := rng.Intn(len(d.Instances))
		sample[i] = d.Instances[idx]
	}

	return &Dataset{
		Instances:   sample,
		NumFeatures: d.NumFeatures,
	}
}

