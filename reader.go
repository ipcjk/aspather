package main

import (
	"bufio"
	"io"
	"log"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

func readPersonalPreference(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		as := strings.Fields(scanner.Text())
		if len(as) != 2 {
			continue
		}
		id, err := strconv.Atoi(as[0])
		if err != nil {
			continue
		}
		value, err := strconv.Atoi(as[1])
		if err != nil {
			continue
		}
		personalValue[id] = value
	}

}

func readRoutes(reader io.Reader) (uint32, map[int]int) {
	var seenPrefixBefore = make(map[string]bool)
	var asRoutes = make(map[int]int)
	var routesSeen uint32

	/* Regular expression helper */
	var aspathReg = regexp.MustCompile(`(\d+)[^\d]*$`)
	var ipAddressReg = regexp.MustCompile(`^\d{1,7}\s+([\d\.\/abcdef:]+)|^\d{7}([\d\.\/abcdef:]+)`)
	var skipNextLine = false
	var reset = false

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {

		if skipNextLine {
			skipNextLine = false
			continue
		}

		brocadeIPmatcher := ipAddressReg.FindStringSubmatch(scanner.Text())

		if len(brocadeIPmatcher) == 3 {
			var prefix string
			if brocadeIPmatcher[2] != "" {
				prefix = brocadeIPmatcher[2]
			} else {
				prefix = brocadeIPmatcher[1]
			}

			if seenPrefixBefore[prefix] {
				skipNextLine = true
				continue
			}
			routesSeen++
			seenPrefixBefore[prefix] = true
			reset = false
			continue
		}

		if strings.Contains(scanner.Text(), "AS_PATH:") && reset == false {
			reset = true
			subMatches := aspathReg.FindStringSubmatch(scanner.Text())
			if len(subMatches) != 2 {
				continue
			}
			asn, err := strconv.Atoi(subMatches[1])
			if err != nil {
				continue
			}
			asRoutes[asn]++
			continue
		}

	}
	return routesSeen, asRoutes
}

func readCountryList(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		c := strings.Fields(scanner.Text())
		value, err := strconv.Atoi(c[1])
		if err != nil {
			continue
		}
		countries[c[0]] = value
	}
}

func readAsList(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
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
			value = *bgpDefaultValue
		} else {
			value += val
		}

		if _, ok := personalValue[id]; ok {
			value = personalValue[id]
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
	log.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	log.Printf("TotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	log.Printf("Sys = %v MiB", bToMb(m.Sys))
	log.Printf("NumGC = %v\n", m.NumGC)
}
