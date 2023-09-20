package utils

import (
	"little_mangamee/entity"
	"strings"
)

func SearchIndex(data []entity.IndexData, title string) ([]entity.IndexData, error) {

	title = strings.ToLower(title)
	var matches []entity.IndexData

	for _, v := range data {
		if strings.Contains(strings.ToLower(v.Title), title) {
			matches = append(matches, v)
		}
	}

	return matches, nil
}
