package main

import "errors"

func parseMap(field interface{}) (map[string]interface{}, error) {
	fieldMap, ok := field.(map[string]interface{})
	if !ok {
		return nil, errors.New("field is not a map")
	}

	return fieldMap, nil
}

func parseMapIntArr(field interface{}) (map[string][]int, error) {
	fieldMap, err := parseMap(field)
	if err != nil {
		return nil, err
	}

	isValid := true
	intArrMap := make(map[string][]int)
	for key, value := range fieldMap {
		// Assert value is []int
		valueSlice, isSlice := value.([]interface{})
		if !isSlice {
			isValid = false
			break
		}

		intArrMap[key] = make([]int, len(valueSlice))
		// Check each element in the slice is an int
		for _, elem := range valueSlice {
			if _, isInt := elem.(float64); !isInt {
				isValid = false
				break
			}
			intArrMap[key] = append(intArrMap[key], int(elem.(float64)))
		}

		if !isValid {
			break
		}
	}

	return intArrMap, nil
}

func parseMapInt(field interface{}) (map[string]int, error) {
	fieldMap, err := parseMap(field)
	if err != nil {
		return nil, err
	}

	isValid := true
	intMap := make(map[string]int)
	for key, value := range fieldMap {
		// Assert value is int
		valueInt, isInt := value.(float64)
		if !isInt {
			isValid = false
			break
		}

		intMap[key] = int(valueInt)
	}

	if !isValid {
		return nil, errors.New("value is not an int")
	}

	return intMap, nil
}
