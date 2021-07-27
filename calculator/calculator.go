package calculator

import (
	"context"
	"math"
)

type randomGetter interface {
	GetIntegers(ctx context.Context, length int) ([]int, error)
}

func NewCalculator(randomGetter randomGetter) *calculator {
	return &calculator{randomGetter: randomGetter}
}

type calculator struct {
	randomGetter randomGetter
}

func (c calculator) CalculateStdDev(ctx context.Context, requests, length int) ([]Result, error) {
	resultsChan := make(chan *singleCalculationResult, requests)
	defer close(resultsChan)
	for i := 0; i < requests; i++ {
		go c.calculateSingleStdDev(ctx, length, resultsChan)
	}
	results := make([]Result, 0, requests)
	allSets := make([]int, 0, requests)
	var err error
	for i := 0; i < requests; i++ {
		res := <-resultsChan
		if res.err != nil {
			err = res.err
		}
		allSets = append(allSets, res.result.Data...)
		results = append(results, res.result)
	}
	if err != nil {
		return nil, err
	}
	results = append(results, Result{StdDev: c.stdDev(allSets), Data: allSets})
	return results, nil
}

type singleCalculationResult struct {
	result Result
	err    error
}

func (c calculator) calculateSingleStdDev(ctx context.Context, length int, resultsChan chan<- *singleCalculationResult) {
	integers, err := c.randomGetter.GetIntegers(ctx, length)
	if err != nil {
		resultsChan <- &singleCalculationResult{err: err}
		return
	}
	resultsChan <- &singleCalculationResult{result: Result{StdDev: c.stdDev(integers), Data: integers}}
}

func (c calculator) stdDev(integers []int) float64 {
	sum := 0
	for _, i := range integers {
		sum += i
	}
	mean := sum / len(integers)
	deviationsSum := 0
	for _, i := range integers {
		diff := i - mean
		deviationsSum += diff * diff
	}
	variance := deviationsSum / len(integers)
	return math.Sqrt(float64(variance))
}
