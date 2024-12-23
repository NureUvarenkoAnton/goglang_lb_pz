package statistics

type numbers interface {
	int | int8 | int16 | int32 | int64
}

func AvgWeighted[E numbers](nums []E) float64 {
	if len(nums) == 0 {
		return 0
	}
	freq := make(map[E]E)
	for _, num := range nums {
		freq[num]++
	}

	sum := E(0)
	weightSum := E(0)
	for _, num := range nums {
		sum += num * freq[num]
		weightSum += freq[num]
	}

	return float64(sum) / float64(weightSum)
}
