package class

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

func (ge GenEd) String() string {
	return genedStringMap[ge]
}

var genedStringMap = map[GenEd]string{
	FSAW: "FSAW",
	FSPW: "FSPW",
	FSOC: "FSOC",
	FSMA: "FSMA",
	FSAR: "FSAR",
	DSNS: "DSNS",
	DSHU: "DSHU",
	DSSP: "DSSP",
	SCIS: "SCIS",
	DVUP: "DVUP",
	DVCC: "DVCC",
}
var stringGenedMap = map[string]GenEd{
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

func ParseGenEd(geneds []string) GenEd {
	var ge GenEd
	for _, g := range geneds {
		ge |= stringGenedMap[strings.ToUpper(g)]
	}
	return ge
}
