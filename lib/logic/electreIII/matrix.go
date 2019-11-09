package electreIII

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"sort"
	"strings"
)

type Matrix struct {
	Size int
	Data []float64
}

func NewMatrix(values *[][]float64) *Matrix {
	size := len(*values)
	data := make([]float64, size*size)
	for i, v := range *values {
		copy(data[i*size:(i+1)*size], v)
	}
	return &Matrix{Size: size, Data: data}
}

func (m *Matrix) At(row, col int) float64 {
	return m.Data[row*m.Size+col]
}

func calcCoords(index, size int) (row, col int) {
	return index / size, index % size
}

func (m *Matrix) Matches(groupsNumber int, groupEvaluator func(row, col int) int, predicate func(value float64) bool) []int {
	groups := make([]int, groupsNumber)
	for i, v := range m.Data {
		groupIndex := groupEvaluator(calcCoords(i, m.Size))
		if predicate(v) {
			groups[groupIndex] += 1
		}
	}
	return groups
}

func (m *Matrix) MatchesInRow(predicate func(value float64) bool) []int {
	return m.Matches(m.Size, func(row, col int) int {
		return row
	}, predicate)
}

func (m *Matrix) MatchesInColumn(predicate func(value float64) bool) []int {
	return m.Matches(m.Size, func(row, col int) int {
		return col
	}, predicate)
}

func (m *Matrix) Filter(filter func(row, col int, v float64) bool) *Matrix {
	newVals := make([]float64, m.Size*m.Size)
	for i, v := range m.Data {
		row, col := calcCoords(i, m.Size)
		if filter(row, col, v) {
			newVals[i] = v
		} else {
			newVals[i] = 0
		}
	}
	return &Matrix{Size: m.Size, Data: newVals}
}

func (m *Matrix) FindBest(isBetter func(old, new float64) bool) float64 {
	if m.Size == 0 {
		panic(fmt.Errorf("matrix is empty"))
	}
	best := m.Data[0]
	for _, v := range m.Data {
		if isBetter(best, v) {
			best = v
		}
	}
	return best
}

func (m *Matrix) Max() float64 {
	return m.FindBest(func(old, new float64) bool {
		return new > old
	})
}

func (m *Matrix) Min() float64 {
	return m.FindBest(func(old, new float64) bool {
		return new < old
	})
}

func (m *Matrix) String() string {
	r := m.Size
	res := make([]string, r)
	for i := 0; i < r; i++ {
		row := make([]string, r)
		for j := 0; j < r; j++ {
			row[j] = fmt.Sprintf("%.2f", m.At(i, j))
		}
		res[i] = strings.Join(row, "\t")
	}
	return strings.Join(res, "\n")
}

func (m *Matrix) Without(indices *[]int) *Matrix {
	size := m.Size
	if len(*indices) == size {
		return m
	}
	data := make([]float64, m.Size*m.Size)
	copy(data, m.Data)
	toRemove := len(*indices)
	sorted := make([]int, toRemove)
	copy(sorted, *indices)
	sort.Sort(sort.Reverse(sort.IntSlice(sorted)))
	for _, v := range sorted {
		data = append(data[0:v*size], data[(v+1)*size:]...)
	}
	size -= toRemove
	resultData := make([]float64, size*size)
	dataIndex := 0
	for i, v := range data {
		rowIndex := i % m.Size
		if !utils.ContainsInts(&sorted, &rowIndex) {
			resultData[dataIndex] = v
			dataIndex++
		}
	}
	return &Matrix{Size: size, Data: resultData}
}

func (m *Matrix) Slice(indices *[]int) *Matrix {
	resultSize := len(*indices)
	if resultSize == m.Size {
		return m
	}
	data := make([]float64, 0)
	sort.Ints(*indices)
	for _, v := range *indices {
		data = append(data, m.Data[v*m.Size:(v+1)*m.Size]...)
	}
	resultData := make([]float64, resultSize*resultSize)
	dataIndex := 0
	for i, v := range data {
		rowIndex := i % m.Size
		if utils.ContainsInts(indices, &rowIndex) {
			resultData[dataIndex] = v
			dataIndex++
		}
	}
	return &Matrix{Size: resultSize, Data: resultData}
}
