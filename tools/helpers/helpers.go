package helpers

// Исключает массив из массива
func excludeFilter(arr []string, arr2 []string) []string {
	for i := 0; i < len(arr); i++ {
		el := arr[i]
		for _, rem := range arr2 {
			if el == rem {
				arr = append(arr[:i], arr[i+1:]...)
				i-- // Important: decrease index
				break
			}
		}
	}
	return arr
}

// Получить различия между двумя слайсами
func diffData(arr1 []string, arr2 []string) ([]string, []string) {
	var dataAdd []string
	var dataDelete []string

	for i := 0; i < 2; i++ {
		for _, s1 := range arr1 {
			found := false
			for _, s2 := range arr2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			if !found && i == 0 {
				dataAdd = append(dataAdd, s1)
			} else if !found && i == 1 {
				dataDelete = append(dataDelete, s1)
			}
		}
		if i == 0 {
			arr1, arr2 = arr2, arr1
		}
	}
	return dataAdd, dataDelete
}
