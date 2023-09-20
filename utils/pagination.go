package utils

import (
	"fmt"
	"little_mangamee/entity"
)

func PaginateIndex(data []entity.IndexData, pageNumber, pageSize int) ([]entity.IndexData, error) {

	if pageNumber <= 0 || pageSize <= 0 {
		return nil, fmt.Errorf("Invalid page number or page size")
	}

	startIndex := (pageNumber - 1) * pageSize
	endIndex := startIndex + pageSize

	// Check if the start and end indices are within the bounds of the data
	if startIndex >= len(data) {
		return nil, fmt.Errorf("Page not found")
	}
	if endIndex > len(data) {
		endIndex = len(data)
	}

	return data[startIndex:endIndex], nil
}
