package domain

// TODO: replace to more appropriate place
const LocaleCtxKey = "loc"

type CountryCode string

func (c CountryCode) Valid() bool {
	return len(c) == 2
}

type LocaleCode string

func (l LocaleCode) Valid() bool {
	return len(l) == 2
}
