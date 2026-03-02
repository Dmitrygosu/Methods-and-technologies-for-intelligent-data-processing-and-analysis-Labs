package data

import (
	"math"
	"math/rand"
)

type RegressionInstance struct {
	Features []float64
	Target   float64
}

type RegressionDataset struct {
	Instances   []RegressionInstance
	NumFeatures int
}

func NewHousingDataset() *RegressionDataset {
	instances := []RegressionInstance{
		// Площадь, количество комнат, этаж, возраст дома -> цена (в тыс. руб.)
		{[]float64{45.0, 1, 5, 2}, 3200},
		{[]float64{55.0, 1, 3, 5}, 3800},
		{[]float64{62.0, 2, 7, 1}, 4200},
		{[]float64{68.0, 2, 2, 8}, 4100},
		{[]float64{75.0, 2, 10, 3}, 4800},
		{[]float64{52.0, 1, 4, 6}, 3500},
		{[]float64{80.0, 3, 5, 2}, 5200},
		{[]float64{85.0, 3, 1, 12}, 5100},
		{[]float64{90.0, 3, 8, 1}, 5800},
		{[]float64{58.0, 2, 6, 4}, 3900},
		{[]float64{95.0, 3, 9, 2}, 6000},
		{[]float64{70.0, 2, 3, 7}, 4500},
		{[]float64{105.0, 4, 6, 3}, 6800},
		{[]float64{110.0, 4, 2, 10}, 6500},
		{[]float64{48.0, 1, 8, 1}, 3300},
		{[]float64{65.0, 2, 4, 9}, 4300},
		{[]float64{78.0, 2, 11, 4}, 4900},
		{[]float64{88.0, 3, 7, 2}, 5600},
		{[]float64{100.0, 4, 5, 5}, 6400},
		{[]float64{115.0, 4, 1, 15}, 7000},
		{[]float64{50.0, 1, 9, 2}, 3400},
		{[]float64{72.0, 2, 12, 3}, 4700},
		{[]float64{82.0, 3, 6, 6}, 5300},
		{[]float64{98.0, 3, 4, 8}, 5900},
		{[]float64{108.0, 4, 3, 11}, 6700},
		{[]float64{120.0, 4, 13, 1}, 7200},
		{[]float64{42.0, 1, 14, 1}, 3000},
		{[]float64{60.0, 2, 15, 5}, 4000},
		{[]float64{76.0, 2, 8, 7}, 5000},
		{[]float64{92.0, 3, 10, 4}, 5700},
		{[]float64{102.0, 4, 9, 6}, 6300},
		{[]float64{125.0, 4, 11, 2}, 7500},
		{[]float64{54.0, 1, 16, 3}, 3600},
		{[]float64{66.0, 2, 5, 10}, 4400},
		{[]float64{84.0, 3, 12, 5}, 5400},
		{[]float64{96.0, 3, 2, 13}, 6100},
		{[]float64{112.0, 4, 4, 9}, 6900},
		{[]float64{130.0, 5, 6, 3}, 7800},
		{[]float64{46.0, 1, 17, 4}, 3100},
		{[]float64{64.0, 2, 13, 8}, 4200},
		{[]float64{74.0, 2, 14, 6}, 4600},
		{[]float64{86.0, 3, 15, 7}, 5500},
		{[]float64{104.0, 4, 7, 12}, 6600},
		{[]float64{118.0, 4, 8, 14}, 7100},
		{[]float64{135.0, 5, 10, 4}, 8000},
		{[]float64{56.0, 2, 18, 2}, 3700},
		{[]float64{69.0, 2, 19, 11}, 4400},
		{[]float64{79.0, 3, 20, 9}, 5100},
		{[]float64{94.0, 3, 1, 16}, 6000},
		{[]float64{106.0, 4, 5, 10}, 6800},
		{[]float64{122.0, 5, 7, 5}, 7300},
		{[]float64{140.0, 5, 9, 6}, 8200},
		{[]float64{51.0, 1, 21, 3}, 3300},
		{[]float64{63.0, 2, 22, 7}, 4100},
		{[]float64{77.0, 2, 23, 8}, 4900},
		{[]float64{87.0, 3, 24, 11}, 5600},
		{[]float64{99.0, 3, 25, 13}, 6200},
		{[]float64{114.0, 4, 3, 15}, 7000},
		{[]float64{128.0, 5, 6, 7}, 7700},
		{[]float64{145.0, 5, 8, 8}, 8500},
		{[]float64{47.0, 1, 26, 5}, 3100},
		{[]float64{61.0, 2, 27, 9}, 3900},
		{[]float64{73.0, 2, 28, 12}, 4700},
		{[]float64{81.0, 3, 29, 14}, 5200},
		{[]float64{97.0, 3, 30, 16}, 6000},
		{[]float64{109.0, 4, 2, 18}, 6700},
		{[]float64{124.0, 5, 4, 9}, 7400},
		{[]float64{142.0, 5, 5, 10}, 8300},
		{[]float64{53.0, 1, 31, 6}, 3500},
		{[]float64{67.0, 2, 32, 10}, 4300},
		{[]float64{71.0, 2, 33, 13}, 4500},
		{[]float64{83.0, 3, 34, 15}, 5400},
		{[]float64{91.0, 3, 35, 17}, 5800},
		{[]float64{107.0, 4, 11, 19}, 6800},
		{[]float64{126.0, 5, 12, 11}, 7600},
		{[]float64{148.0, 5, 13, 12}, 8800},
		{[]float64{49.0, 1, 36, 7}, 3200},
		{[]float64{59.0, 2, 37, 11}, 3800},
		{[]float64{68.0, 2, 38, 14}, 4300},
		{[]float64{89.0, 3, 39, 16}, 5700},
		{[]float64{93.0, 3, 40, 18}, 5900},
		{[]float64{111.0, 4, 14, 20}, 6900},
		{[]float64{127.0, 5, 15, 13}, 7700},
		{[]float64{150.0, 5, 16, 14}, 9000},
		{[]float64{55.0, 1, 1, 8}, 3600},
		{[]float64{57.0, 2, 2, 12}, 3700},
		{[]float64{70.0, 2, 3, 15}, 4400},
		{[]float64{78.0, 3, 4, 17}, 5000},
		{[]float64{101.0, 4, 5, 21}, 6400},
		{[]float64{116.0, 4, 6, 22}, 7000},
		{[]float64{132.0, 5, 7, 15}, 7900},
		{[]float64{155.0, 5, 8, 16}, 9200},
		{[]float64{44.0, 1, 9, 9}, 3000},
		{[]float64{62.0, 2, 10, 13}, 4100},
		{[]float64{75.0, 2, 11, 16}, 4800},
		{[]float64{88.0, 3, 12, 19}, 5600},
		{[]float64{103.0, 4, 13, 23}, 6500},
		{[]float64{119.0, 4, 14, 24}, 7200},
		{[]float64{138.0, 5, 15, 17}, 8100},
		{[]float64{160.0, 5, 16, 18}, 9500},
		{[]float64{52.0, 1, 17, 10}, 3400},
		{[]float64{65.0, 2, 18, 14}, 4200},
		{[]float64{76.0, 2, 19, 17}, 4900},
		{[]float64{85.0, 3, 20, 20}, 5400},
		{[]float64{105.0, 4, 21, 25}, 6600},
		{[]float64{121.0, 5, 22, 26}, 7300},
		{[]float64{143.0, 5, 23, 19}, 8400},
		{[]float64{165.0, 5, 24, 20}, 9800},
		{[]float64{48.0, 1, 25, 11}, 3100},
		{[]float64{58.0, 2, 26, 15}, 3800},
		{[]float64{72.0, 2, 27, 18}, 4600},
		{[]float64{82.0, 3, 28, 21}, 5300},
		{[]float64{108.0, 4, 29, 27}, 6800},
		{[]float64{123.0, 5, 30, 28}, 7400},
		{[]float64{147.0, 5, 1, 21}, 8600},
		{[]float64{170.0, 5, 2, 22}, 10000},
		{[]float64{50.0, 1, 3, 12}, 3200},
		{[]float64{63.0, 2, 4, 16}, 4000},
		{[]float64{74.0, 2, 5, 19}, 4700},
		{[]float64{90.0, 3, 6, 22}, 5800},
		{[]float64{113.0, 4, 7, 29}, 7000},
		{[]float64{129.0, 5, 8, 30}, 7800},
		{[]float64{152.0, 5, 9, 23}, 8900},
		{[]float64{175.0, 5, 10, 24}, 10500},
		{[]float64{54.0, 1, 11, 13}, 3500},
		{[]float64{66.0, 2, 12, 17}, 4300},
		{[]float64{77.0, 2, 13, 20}, 5000},
		{[]float64{95.0, 3, 14, 23}, 6100},
		{[]float64{117.0, 4, 15, 31}, 7200},
		{[]float64{134.0, 5, 16, 32}, 8000},
		{[]float64{158.0, 5, 17, 25}, 9200},
		{[]float64{180.0, 5, 18, 26}, 11000},
		{[]float64{56.0, 1, 19, 14}, 3600},
		{[]float64{68.0, 2, 20, 18}, 4400},
		{[]float64{79.0, 2, 21, 21}, 5100},
		{[]float64{98.0, 3, 22, 24}, 6300},
		{[]float64{120.0, 4, 23, 33}, 7400},
		{[]float64{137.0, 5, 24, 34}, 8200},
		{[]float64{162.0, 5, 25, 27}, 9400},
		{[]float64{185.0, 5, 26, 28}, 11200},
		{[]float64{51.0, 1, 27, 15}, 3300},
		{[]float64{64.0, 2, 28, 19}, 4100},
		{[]float64{80.0, 2, 29, 22}, 5200},
		{[]float64{100.0, 3, 30, 25}, 6400},
		{[]float64{125.0, 4, 1, 35}, 7700},
		{[]float64{141.0, 5, 2, 36}, 8400},
		{[]float64{167.0, 5, 3, 29}, 9700},
		{[]float64{190.0, 5, 4, 30}, 11500},
	}

	return &RegressionDataset{
		Instances:   instances,
		NumFeatures: 4,
	}
}

func (d *RegressionDataset) Split(trainRatio float64, seed int64) (*RegressionDataset, *RegressionDataset) {
	rand.Seed(seed)
	perm := rand.Perm(len(d.Instances))

	trainSize := int(float64(len(d.Instances)) * trainRatio)

	train := &RegressionDataset{
		Instances:   make([]RegressionInstance, trainSize),
		NumFeatures: d.NumFeatures,
	}

	test := &RegressionDataset{
		Instances:   make([]RegressionInstance, len(d.Instances)-trainSize),
		NumFeatures: d.NumFeatures,
	}

	for i := 0; i < trainSize; i++ {
		train.Instances[i] = d.Instances[perm[i]]
	}

	for i := trainSize; i < len(d.Instances); i++ {
		test.Instances[i-trainSize] = d.Instances[perm[i]]
	}

	return train, test
}

func (d *RegressionDataset) Normalize() {
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

