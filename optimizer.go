package main

import (
	"fmt"
	"sort"
)

func optimizeGreedy(camTableSize int) {
	var resultSet []bgpAS
	var negativeResults []bgpAS
	var resultRoutes = 0

	var zero, matched, nospace int
	var nospaceRoutes int
	for i, as := range bgpASlist {
		if bgpASlist[i].routesNumber == 0 {
			zero++
			bgpASlist[i].zero = true
			continue
		}

		if resultRoutes+bgpASlist[i].routesNumber <= camTableSize {
			resultRoutes += as.routesNumber
			bgpASlist[i].picked = true
			matched++
		} else {
			bgpASlist[i].nospace = true
			nospaceRoutes += as.routesNumber
			nospace++
			countrycounter[as.country]++
			negativeResults = append(negativeResults, as)
		}

	}

	/* Not picked */
	//for _, elem := range negativeResults {
	//		fmt.Println("Number", elem.asNumber)
	//}
	fmt.Println(resultRoutes, "in", len(resultSet), "zero", zero, "match", matched, "nospace", nospace, "nospaceRoutes", nospaceRoutes)

	fmt.Println("Routes not installed:")
	fmt.Printf("\t AS not seen in BGP: %d\n", zero)
	fmt.Printf("\t AS not enough CAM: %d\n", nospace)
	fmt.Printf("\t Not installed by country:\n\t")
	for s, i := range countrycounter {
		fmt.Printf("%s:%d ", s, i)
	}
	fmt.Println("")
	generatePrefixList(returnRanges(negativeResults))
}

/* too big, needs to be branched */
func optimizeKnapsack(camTableSize int) {
	var maxRow = len(bgpASlist)
	var maxCol = camTableSize

	/* sort bgplist
	sort.Slice(bgpASlist, func(i, j int) bool {
		return bgpASlist[i].routesNumber > bgpASlist[j].routesNumber
	})
	*/

	/* pointers to the current and compare row */
	var current, last int
	var lastWritten int
	last = 1

	for row := 0; row < maxRow; row++ {
		for col := 1; col < maxCol; col++ {

			// fmt.Println(row, col)
			/* temporary variables to calculate next slot */
			var newValue, newRoutes int

			/* check if weighted fill wit */
			if col-bgpASlist[row].routesNumber >= 0 {
				/* würde auch ohne das IF funktionieren, hängt dann halt nur Leerzeichen mit an */
				newRoutes = bgpASlist[row].routesNumber + matrix[last][col-bgpASlist[row].routesNumber].routesNumber
				newValue = bgpASlist[row].value + matrix[last][col-bgpASlist[row].routesNumber].value
			}

			/* Check whats fitter, this is one-dimensional only :(
			 */
			if matrix[last][col].value >= newValue {
				matrix[current][col].value = matrix[last][col].value
				matrix[current][col].routesNumber = matrix[last][col].routesNumber
				matrix[current][col].asnumbers = matrix[last][col].asnumbers
			} else {
				matrix[current][col].value = newValue
				matrix[current][col].routesNumber = newRoutes
				matrix[current][col].asnumbers = append(matrix[current][col].asnumbers, bgpASlist[row].asNumber)
			}
			lastWritten = current
		}

		last, current = current, last

	}
	fmt.Println(matrix[lastWritten][camTableSize-1].routesNumber)
}

func returnRanges(negativeResults []bgpAS) []string {
	var asResults []string

	/* sort the numbers in ascending order */
	sort.Slice(negativeResults, func(i, j int) bool {
		return negativeResults[i].asNumber < negativeResults[j].asNumber
	})

	if len(negativeResults) == 0 {
		return []string{}
	}

	/* prepare start range */
	var start = negativeResults[0].asNumber
	var cur = start

	/* loop all numbers and try to find submatches */
	for i := 1; i < len(negativeResults); i++ {

		/* are we in an ongoing range? */
		if cur+1 == negativeResults[i].asNumber {
			cur = negativeResults[i].asNumber
			continue
		}

		/* are we closing the range? */
		if start != negativeResults[i-1].asNumber {
			asResults = append(asResults, fmt.Sprintf("%d-%d", start, negativeResults[i-1].asNumber))
		} else {
			asResults = append(asResults, fmt.Sprintf("%d", start))
		}

		/* are we the last element? */
		if i+1 == len(negativeResults) {
			asResults = append(asResults, fmt.Sprintf("%d", negativeResults[i].asNumber))
			break
		}
		/* else, prepare the next loop */
		start = negativeResults[i].asNumber
		cur = start
	}
	return asResults
}
