package choquet

import "math"

func PowerSet(original []string) [][]string {
	powerSetSize := int(math.Pow(2, float64(len(original))))
	result := make([][]string, 0, powerSetSize)

	var index int
	for index < powerSetSize {
		var subSet []string

		for j, elem := range original {
			if index&(1<<uint(j)) > 0 {
				subSet = append(subSet, elem)
			}
		}
		if len(subSet) > 0 {
			result = append(result, subSet)
		}
		index++
	}
	return result
}
