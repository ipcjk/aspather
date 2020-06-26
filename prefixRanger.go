package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func generatePrefixList(aspaths []string) {
	var prefixLists []string
	for _, v := range aspaths {
		if strings.Contains(v, "-") {
			rangeSplit := strings.Split(v, "-")
			start, err := strconv.Atoi(rangeSplit[0])
			if err != nil {
				panic(err)
			}
			end, err := strconv.Atoi(rangeSplit[1])
			if err != nil {
				panic(err)
			}
			if start == end {
				prefixLists = append(prefixLists, fmt.Sprintf(*fmtAsPathFmt, *fmtASPathNameFmt,
					"_"+strconv.Itoa(start)))
			} else {
				prefixLists = append(prefixLists, fmt.Sprintf(*fmtAsPathFmt, *fmtASPathNameFmt,
					GetRegex(start, end)))
			}
		} else {
			prefixLists = append(prefixLists, fmt.Sprintf(*fmtAsPathFmt, *fmtASPathNameFmt, "_"+v))
		}

	}
	sort.Slice(prefixLists,
		func(i, j int) bool {
			return prefixLists[i] < prefixLists[j]
		})

	for _, v := range prefixLists {
		fmt.Fprint(os.Stdout, v)
	}
}
