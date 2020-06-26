package main

import (
	"flag"
	"log"
	"sort"
)

type bgpAS struct {
	name         string
	value        int
	routesNumber int
	asNumber     int
	country      string
	picked       bool
	nospace      bool
	zero         bool
}

type node struct {
	asnumbers    []int
	value        int
	routesNumber int
}

/* matrix */
var matrix [][]node
var bgpASlist []bgpAS

/* helper maps */
var asNumber = make(map[int]string)
var asNumberRoutes = make(map[int]int)
var countries = make(map[string]int)
var countrycounter = make(map[string]int)
var fmtAsPathDefault = "ip as-path access-list %s permit %s$\n"
var fmtASPathName = "savethefib"
var fmtASPathNameFmt *string
var fmtAsPathFmt *string

func main() {
	aslist := flag.String("aslist", "data_asnums", "as number list")
	bestRoutes := flag.String("routes", "bestroutes.slx", "router output, e.g. show ipv6 bgp routes best")
	countryList := flag.String("country", "config_country", "list with country default weight values")
	asconfig := flag.String("asconfig", "asconfig", "personal as configuration and weights")
	camSize := flag.Int("camsize", 512000, "size of the routers cam")
	sorttype := flag.Int("sorttype", 1, "type of sort, 0=value,then routesnumber bigger, 1=value,then routesnumber smaller")
	fmtAsPathFmt = flag.String("aspathfmt", fmtAsPathDefault, "default for printing the as-path list")
	fmtASPathNameFmt = flag.String("aspathname", fmtASPathName, "default name for as-path list")

	flag.Parse()

	readCountryList(*countryList)
	routesSeen := readBestRoutes(*bestRoutes)
	readAsList(*aslist)
	readPersonalPreference(*asconfig)

	/* sortswitch */
	if *sorttype == 0 {
		sort.Slice(bgpASlist, func(i, j int) bool {
			switch {
			case bgpASlist[i].value != bgpASlist[j].value:
				return bgpASlist[i].value >= bgpASlist[j].value
			default:
				return bgpASlist[i].routesNumber >= bgpASlist[j].routesNumber
			}
		})
	} else {
		sort.Slice(bgpASlist, func(i, j int) bool {
			switch {
			case bgpASlist[i].value != bgpASlist[j].value:
				return bgpASlist[i].value >= bgpASlist[j].value
			default:
				return bgpASlist[i].routesNumber <= bgpASlist[j].routesNumber
			}
		})
	}

	log.Printf("Starting optimizer, target cam size %d\n"+
		"input: %d routes in %d bgp autonomous systems\n",
		*camSize, routesSeen, len(bgpASlist))

	optimizeGreedy(*camSize)

	PrintMemUsage()

}
