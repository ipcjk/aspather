package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

func readPersonalPreference(fileName string) {

}

func readBestRoutes(fileName string) (routesSeen uint32) {
	var seenPrefixBefore = make(map[string]bool)
	/* Regular expression helper */
	var aspathReg = regexp.MustCompile(`(\d+)[^\d]*$`)
	var ip4AddressReg = regexp.MustCompile(`(\d{1,3}.\d{1,3}.\d{1,3}.\d{1,3}\/\d{1,2})`)

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	var ignoreNextAS_PATH = false

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if ignoreNextAS_PATH {
			ignoreNextAS_PATH = false
			continue
		}

		ip4 := ip4AddressReg.FindStringSubmatch(scanner.Text())
		if len(ip4) == 2 {
			routesSeen++
			if seenPrefixBefore[ip4[0]] {
				ignoreNextAS_PATH = true
				continue
			}
			seenPrefixBefore[ip4[0]] = true
			continue
		}

		if strings.Contains(scanner.Text(), "AS_PATH:") {
			subMatches := aspathReg.FindStringSubmatch(scanner.Text())
			if len(subMatches) != 2 {
				continue
			}
			asn, err := strconv.Atoi(subMatches[1])
			if err != nil {
				continue
			}
			asNumberRoutes[asn]++
		}

	}
	return routesSeen
}

func readCountryList(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		c := strings.Fields(scanner.Text())
		value, err := strconv.Atoi(c[1])
		if err != nil {
			continue
		}
		countries[c[0]] = value
	}
}

func readAsList(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		as := strings.Fields(scanner.Text())
		/* cut out some things */
		asId := strings.Trim(as[0], "AS")
		id, err := strconv.Atoi(asId)
		if err != nil {
			continue
		}
		as = as[1:]
		country := as[len(as)-1]
		as = as[:len(as)-1]
		asInfo := strings.Join(as, "")
		asNumber[id] = asInfo

		var value int
		if val, ok := countries[country]; !ok {
			value = 10
		} else {
			value += val
		}

		bgpASlist = append(bgpASlist, bgpAS{
			value:        value,
			routesNumber: asNumberRoutes[id],
			asNumber:     id,
			country:      country,
			name:         asId,
		})
	}

}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	log.Printf("\tAlloc = %v MiB", bToMb(m.Alloc))
	log.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	log.Printf("\tSys = %v MiB", bToMb(m.Sys))
	log.Printf("\tNumGC = %v\n", m.NumGC)
}
