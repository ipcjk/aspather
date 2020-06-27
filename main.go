package main

// Copyright 2020 Jörg Kost All rights reserved.
// jk@ip-clear.de
// Use of this source code is governed by Apache 2.0
// license that can be found in the LICENSE.MD file.

import (
	"flag"
	"log"
	"os"
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
var asNumberRoutes map[int]int
var countries = make(map[string]int)
var countrycounter = make(map[string]int)
var personalValue = make(map[int]int)
var routesSeen uint32

/* some global defaults */
var fmtAsPathDefault = "ip as-path access-list %s permit %s$\n"
var fmtASPathName = "savethefib"
var fmtASPathNameFmt *string
var fmtAsPathFmt *string

func main() {
	aslist := flag.String("aslist", "data_asnums", "as number list")
	bestRoutes := flag.String("routes", "bestroutes.slx", "router output, e.g. show ipv6 bgp routes best")
	countryList := flag.String("country", "config_country", "list with country default weight values")
	personalList := flag.String("personal", "", "list with preferred personal as config")
	camSize := flag.Int("camsize", 512000, "size of the routers cam")
	sorttype := flag.Int("sorttype", 1, "type of sort, 0=value,then routesnumber bigger, 1=value,then routesnumber smaller")
	fmtAsPathFmt = flag.String("aspathfmt", fmtAsPathDefault, "default for printing the as-path list")
	fmtASPathNameFmt = flag.String("aspathname", fmtASPathName, "default name for as-path list")
	debug := flag.Bool("debug", false, "if memory or other infos is being printed out")

	flag.Parse()

	if *personalList != "" {
		readPersonalPreference(*personalList)
	}
	readCountryList(*countryList)

	/* open and read routes file*/
	file, err := os.Open(*bestRoutes)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	routesSeen, asNumberRoutes = readRoutes(file)

	readAsList(*aslist)

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

	if *debug {
		log.Printf("Starting optimizer, target cam size %d\n"+
			"input: %d routes in %d bgp autonomous systems\n",
			*camSize, routesSeen, len(bgpASlist))
	}

	optimizeGreedy(*camSize, *debug)

	if *debug {
		PrintMemUsage()
	}

}
