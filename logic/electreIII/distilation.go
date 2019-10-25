package electreIII

import (
	"../../utils"
	"fmt"
)

var DefaultDistillationFunc = LinearFunctionParameters{A: -.15, B: .3}

type CompareFunction = func(old int, new int) bool

func RankAscending(matrix *AlternativesMatrix, distillationFun *LinearFunctionParameters) *[]int {
	return rank(matrix, distillationFun, greater)
}

func RankDescending(matrix *AlternativesMatrix, distillationFun *LinearFunctionParameters) *[]int {
	ranking := rank(matrix, distillationFun, lower)
	maxPosition := Max(ranking)
	minusValuesFrom(ranking, maxPosition+1)
	return ranking
}

func Max(values *[]int) int {
	if len(*values) == 0 {
		panic(fmt.Errorf("slice is empty"))
	}
	best := (*values)[0]
	for _, v := range *values {
		if v > best {
			best = v
		}
	}
	return best
}

func minusValuesFrom(values *[]int, value int) {
	for i, v := range *values {
		(*values)[i] = value - v
	}
}

func rank(matrix *AlternativesMatrix, distillationFun *LinearFunctionParameters, evaluateFunction CompareFunction) *[]int {
	position := 1
	withoutD := removeDiagonal(matrix)
	maxCred := withoutD.Max()
	positions := make([]int, len(*matrix.Alternatives))
	indices := make([]int, len(positions))
	for i, _ := range indices {
		indices[i] = i
	}
	return distillate(maxCred, position, withoutD, distillationFun, evaluateFunction, false)
}

func samePositions(size, value int) *[]int {
	pos := make([]int, size)
	for i := range pos {
		pos[i] = value
	}
	return &pos
}

func distillate(
	maxCred float64, position int,
	matrix *Matrix,
	distillationFun *LinearFunctionParameters,
	evaluateFunction CompareFunction,
	isInner bool,
) *[]int {
	if maxCred == 0 {
		return samePositions(matrix.Size, position)
	}
	minCred, valuesToConsider := getDistillateMatrix(distillationFun, maxCred, matrix)
	quality := computeQuality(valuesToConsider)
	_, bestIndices := findBestMatch(quality, evaluateFunction)
	positions := samePositions(matrix.Size, 0)
	updatePositions(position, minCred, matrix, bestIndices, positions, distillationFun, evaluateFunction)
	indicesLeftToUpdate := updatedPositions(bestIndices, positions)
	if len(*indicesLeftToUpdate) == matrix.Size || isInner {
		return positions
	}
	position++
	nextIterationMatrix := matrix.Without(indicesLeftToUpdate)
	furtherPositions := distillate(nextIterationMatrix.Max(), position, nextIterationMatrix, distillationFun, evaluateFunction, false)
	writePositionsSequentially(furtherPositions, positions)
	return positions
}

func writePositionsSequentially(positionsToWrite, positions *[]int) {
	toWriteIndex := 0
	for i, p := range *positions {
		if p == 0 {
			if toWriteIndex >= len(*positionsToWrite) {
				panic(fmt.Errorf(
					"position %d is out of scope for possible possitions %v and all positions %v",
					i, *positionsToWrite, *positions,
				))
			}
			(*positions)[i] = (*positionsToWrite)[toWriteIndex]
			toWriteIndex++
		}
	}
}

func updatedPositions(indices, positions *[]int) *[]int {
	newValues := make([]int, 0)
	for _, p := range *indices {
		if (*positions)[p] != 0 {
			newValues = append(newValues, p)
		}
	}
	return &newValues
}

func updatePositions(
	position int, minCred float64,
	valuesToConsider *Matrix, bestIndices, positions *[]int,
	distillationFun *LinearFunctionParameters, evaluateFunction CompareFunction,
) {
	bestIndicesNum := len(*bestIndices)
	if bestIndicesNum > 1 && minCred > 0 {
		nextToFilter := valuesToConsider.Slice(bestIndices)
		subPositions := distillate(minCred, position, nextToFilter, distillationFun, evaluateFunction, true)
		updateValues(bestIndices, positions, subPositions)
	} else if bestIndicesNum > 0 {
		updateValues(bestIndices, positions, samePositions(bestIndicesNum, position))
	}
}

func getDistillateMatrix(distillationFun *LinearFunctionParameters, maxCred float64, matrix *Matrix) (float64, *Matrix) {
	v, _ := distillationFun.evaluate(maxCred)
	minCredThreshold := maxCred - v
	minCred := matrix.FindBest(func(old, new float64) bool {
		// ok because the lowest value is 0, on diagonal for sure.
		return new < minCredThreshold && new > old
	})
	valuesToConsider := matrix.Filter(func(row, col int, v float64) bool {
		if v <= minCred {
			return false
		}
		funcValueForThisField, _ := distillationFun.evaluate(v)
		value := matrix.At(col, row) + funcValueForThisField
		return v > value
	})
	return minCred, valuesToConsider
}

func updateValues(indicesToUpdate, original, new *[]int) {
	for i, v := range *indicesToUpdate {
		(*original)[v] = (*new)[i]
	}
}

func removeDiagonal(matrix *AlternativesMatrix) *Matrix {
	return matrix.Values.Filter(func(row, col int, v float64) bool {
		return row != col
	})
}

func computeQuality(matrix *Matrix) *[]int {
	strength := matrix.MatchesInRow(utils.IsPositive)
	weakness := matrix.MatchesInColumn(utils.IsPositive)
	return calcQuality(&strength, &weakness)
}

func calcQuality(strength, weakness *[]int) *[]int {
	quality := make([]int, len(*strength))
	for i, s := range *strength {
		quality[i] = s - (*weakness)[i]
	}
	return &quality
}

func greater(old, new int) bool {
	return old < new
}

func lower(old, new int) bool {
	return old > new
}

func findBestMatch(values *[]int, isBetter CompareFunction) (value int, indices *[]int) {
	bestValue := (*values)[0]
	bestIndices := make([]int, 0)
	for i, v := range *values {
		if isBetter(bestValue, v) {
			bestValue = v
			bestIndices = []int{i}
		} else if v == bestValue {
			bestIndices = append(bestIndices, i)
		}
	}
	return bestValue, &bestIndices
}
