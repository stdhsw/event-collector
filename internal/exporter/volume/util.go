// Package volume provides file sorting and number extraction utilities for the Volume exporter.
// SortByNumericSuffix orders file names by their trailing integer suffix,
// and extractNumber parses that suffix from a given filename.
package volume

import (
	"regexp"
	"sort"
	"strconv"
)

var reNumericSuffix = regexp.MustCompile(`\d+$`) // 파일명 끝의 숫자를 추출하기 위한 정규식 패턴

// SortByNumericSuffix files를 파일명 끝의 숫자 기준으로 오름차순 정렬한다.
func SortByNumericSuffix(files []string) {
	sort.Slice(files, func(i, j int) bool {
		numI := extractNumber(files[i])
		numJ := extractNumber(files[j])
		return numI < numJ
	})
}

// extractNumber filename 끝의 숫자를 추출하여 반환한다. 숫자가 없으면 0을 반환한다.
func extractNumber(filename string) int {
	match := reNumericSuffix.FindString(filename)
	if match == "" {
		return 0
	}
	num, err := strconv.Atoi(match)
	if err != nil {
		return 0
	}
	return num
}
