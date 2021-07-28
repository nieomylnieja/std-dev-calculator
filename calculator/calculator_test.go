package calculator

import (
	"context"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// comparing floats yields nondeterministic results,
// a simple and reliable solution should be rounding them with string conversion and ideally ignoring diffs in the last decimal place
// we don't care about perfect decimal precision, we're only interested to see if our algorithm works up to a certain precision
// I do acknowledge though that this is not perfect, I've run the tests for these data sets long enough to be convinced the results don't differ on MY MACHINE (sadly no guarantees on yours)
type stringStdDevResult struct {
	StdDev string
	Data   []int
}

func TestCalculateStdDev(t *testing.T) {

	for name, test := range map[string]struct {
		requests     int
		length       int
		randomResult [][]int
		expected     []stringStdDevResult
	}{
		"single request": {
			requests: 1,
			length:   6,
			randomResult: [][]int{
				{2, 5, 6, 7, 2, 5},
			},
			expected: []stringStdDevResult{
				{StdDev: "1.89297", Data: []int{2, 5, 6, 7, 2, 5}},
				{StdDev: "1.89297", Data: []int{2, 5, 6, 7, 2, 5}}, // sum of all sets
			},
		},
		"single request with a lot of higher order numbers": {
			requests: 1,
			length:   20,
			randomResult: [][]int{
				{36413, 23090, 15194, 50563, 82433, 34147, 74078, 74324, 86159, 11353, 61957, 43721, 37189, 42199, 23000, 68705, 12888, 24538, 79703, 29355},
			},
			expected: []stringStdDevResult{
				{StdDev: "24311.10413", Data: []int{36413, 23090, 15194, 50563, 82433, 34147, 74078, 74324, 86159, 11353, 61957, 43721, 37189, 42199, 23000, 68705, 12888, 24538, 79703, 29355}},
				{StdDev: "24311.10413", Data: []int{36413, 23090, 15194, 50563, 82433, 34147, 74078, 74324, 86159, 11353, 61957, 43721, 37189, 42199, 23000, 68705, 12888, 24538, 79703, 29355}}, // sum of all sets
			},
		},
		"two requests": {
			requests: 2,
			length:   6,
			randomResult: [][]int{
				{9, 7, 4, 4, 2, 9},
				{9, 9, 6, 3, 2, 1},
			},
			expected: []stringStdDevResult{
				{StdDev: "2.67187", Data: []int{9, 7, 4, 4, 2, 9}},
				{StdDev: "3.21455", Data: []int{9, 9, 6, 3, 2, 1}},
				{StdDev: "2.98492", Data: []int{
					9, 7, 4, 4, 2, 9,
					9, 9, 6, 3, 2, 1,
				}}, // sum of all sets
			},
		},
		"ten requests": {
			requests: 10,
			length:   6,
			randomResult: [][]int{
				{8, 9, 7, 1, 8, 2},
				{3, 8, 3, 6, 9, 2},
				{3, 5, 7, 2, 1, 5},
				{4, 1, 8, 3, 4, 5},
				{4, 9, 6, 3, 2, 5},
				{4, 2, 7, 2, 5, 9},
				{1, 5, 6, 5, 4, 1},
				{3, 2, 1, 4, 5, 7},
				{6, 5, 5, 1, 1, 7},
				{8, 9, 1, 8, 9, 3},
			},
			expected: []stringStdDevResult{
				{StdDev: "3.13138", Data: []int{8, 9, 7, 1, 8, 2}},
				{StdDev: "2.67187", Data: []int{3, 8, 3, 6, 9, 2}},
				{StdDev: "2.03443", Data: []int{3, 5, 7, 2, 1, 5}},
				{StdDev: "2.11476", Data: []int{4, 1, 8, 3, 4, 5}},
				{StdDev: "2.26691", Data: []int{4, 9, 6, 3, 2, 5}},
				{StdDev: "2.54406", Data: []int{4, 2, 7, 2, 5, 9}},
				{StdDev: "1.97203", Data: []int{1, 5, 6, 5, 4, 1}},
				{StdDev: "1.97203", Data: []int{3, 2, 1, 4, 5, 7}},
				{StdDev: "2.33928", Data: []int{6, 5, 5, 1, 1, 7}},
				{StdDev: "3.14466", Data: []int{8, 9, 1, 8, 9, 3}},
				{StdDev: "2.60656", Data: []int{
					8, 9, 7, 1, 8, 2,
					3, 8, 3, 6, 9, 2,
					3, 5, 7, 2, 1, 5,
					4, 1, 8, 3, 4, 5,
					4, 9, 6, 3, 2, 5,
					4, 2, 7, 2, 5, 9,
					1, 5, 6, 5, 4, 1,
					3, 2, 1, 4, 5, 7,
					6, 5, 5, 1, 1, 7,
					8, 9, 1, 8, 9, 3,
				}}, // sum of all sets
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			for _, randRes := range test.randomResult {
				if !assert.Len(t, randRes, test.length, "make sure that each random result returned by the mock is of len() equal to test.length") {
					t.FailNow()
				}
			}
			calc := NewCalculator(&mockRandomGetter{result: test.randomResult, mu: new(sync.Mutex)})

			res, err := calc.CalculateStdDev(context.Background(), test.requests, test.length)

			assert.NoError(t, err)
			assert.Len(t, res, test.requests+1)

			convertedResult := convertResultsToString(t, res)
			assert.ElementsMatch(t, test.expected[:len(test.expected)-1], convertedResult[:len(convertedResult)-1])
			// the last list will be randomly ordered as well since the goroutines may finish in different order
			sumOfAllExpected := test.expected[len(test.expected)-1]
			sumOfAllResult := convertedResult[len(convertedResult)-1]
			assert.Equal(t, sumOfAllExpected.StdDev, sumOfAllResult.StdDev)
			assert.ElementsMatch(t, sumOfAllExpected.Data, sumOfAllResult.Data)
		})
	}
}

func convertResultsToString(t *testing.T, results []Result) []stringStdDevResult {
	t.Helper()
	conv := make([]stringStdDevResult, 0, len(results))
	for _, res := range results {
		conv = append(conv, stringStdDevResult{
			StdDev: strconv.FormatFloat(res.StdDev, 'f', 5, 64),
			Data:   res.Data,
		})
	}
	return conv
}

type mockRandomGetter struct {
	result [][]int
	ctr    int
	mu     *sync.Mutex
}

func (m *mockRandomGetter) GetIntegers(_ context.Context, _ int) ([]int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	defer func() { m.ctr++ }()
	return m.result[m.ctr], nil
}
