
type GenEd = number;

export enum GenEds{
    ANY = ~0,
	FSAW = 1 << 0,
	FSPW = 1 << 1,
	FSOC = 1 << 2,
	FSMA = 1 << 3,
	FSAR = 1 << 4,
	DSNS = 1 << 5,
	DSHU = 1 << 6,
	DSSP = 1 << 7,
	SCIS = 1 << 8,
	DVUP = 1 << 9,
	DVCC = 1 << 10,
}

export default GenEd;