package testudo

import "strings"

type GenEd uint32

const (
	ANY  GenEd = ^GenEd(0)
	FSAW GenEd = 1 << iota
	FSPW
	FSOC
	FSMA
	FSAR
	DSNS
	DSHU
	DSSP
	SCIS
	DVUP
	DVCC
)

var genedMap = map[string]GenEd{
	"FSAW": FSAW,
	"FSPW": FSPW,
	"FSOC": FSOC,
	"FSMA": FSMA,
	"FSAR": FSAR,
	"DSNS": DSNS,
	"DSHU": DSHU,
	"DSSP": DSSP,
	"SCIS": SCIS,
	"DVUP": DVUP,
	"DVCC": DVCC,
}

func parseGenEd(geneds []string) GenEd {
	var ge GenEd
	for _, g := range geneds {
		ge |= genedMap[strings.ToUpper(g)]
	}
	return ge
}
