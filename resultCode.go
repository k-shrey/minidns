package main

type ResultCode uint8

const (
	NOERROR ResultCode = iota
	FORMERR
	SERVFAIL
	NXDOMAIN
	NOTIMP
	REFUSED
)

func (r ResultCode) String() string {
	return [...]string{"NOERROR", "FORMERROR", "SERVFAIL", "NXDOMAIN", "NOTIMP", "REFUSED"}[r]
}

func ResultCodeFromNumber(num uint8) ResultCode {

	switch num {
	case 0:
		return NOERROR
	case 1:
		return SERVFAIL
	case 2:
		return NXDOMAIN
	case 3:
		return NOTIMP
	case 4:
		return REFUSED
	}

	return FORMERR
}
