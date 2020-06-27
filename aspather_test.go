package main

import (
	"strings"
	"testing"
)

func TestParsingRoutes(t *testing.T) {
	var routesSeen uint32
	var asRoutes = make(map[int]int)
	testinput := strings.NewReader(
		`1      2001::/32          02:0000:ac02::1
                                          10513      90         0      BI    
         AS_PATH: 3356 6939
2      2001::/32          02:0000:ac02::19
                                           none      90         0      I     
         AS_PATH: 1299 25192
3      2001:4:112::/48    02:0000:ac02::1
                                          100        160        0      BI    
         AS_PATH: 6724 112
1815006223.238.16.0/22    0.81.194.250    none      140        0      BEx   
         AS_PATH: 9498 45609
1815007223.238.16.0/22    0.81.194.250    none      140        0      E     
         AS_PATH: 9498 45609
1815008223.238.20.0/22    0.81.194.250    none      140        0      BEx   
         AS_PATH: 9498 45609
1815009223.238.20.0/22    0.81.194.250    none      140        0      E     
         AS_PATH: 9498 45609
749746 109.230.222.0/24   0.81.196.125   0          140        0      BEx   
         AS_PATH: 200738
749747 109.230.222.0/24   0.81.196.125   0          140        0      E     
         AS_PATH: 200738
1815337223.255.254.0/24   6.23.18.19    none      90         0      I     
         AS_PATH: 1299 7473 3758 55415
`)

	routesSeen, asRoutes = readRoutes(testinput)

	if routesSeen != 6 {
		t.Errorf("Not enough routes parsed: %d / 10", routesSeen)
	}

	if asRoutes[6939] != 1 {
		t.Error("Too much / less routes from 6939", asRoutes[6939])
	}

	if asRoutes[45609] != 2 {
		t.Error("Too much / less routes from 45609", asRoutes[45609])
	}

	if asRoutes[25192] != 0 {
		t.Error("Too much / less routes from 25192", asRoutes[25192])
	}

	if asRoutes[200738] != 1 {
		t.Error("Too much / less routes from 200738", asRoutes[200738])
	}

	if asRoutes[55415] != 1 {
		t.Error("Too much / less routes from 200738", asRoutes[55415])
	}
}
