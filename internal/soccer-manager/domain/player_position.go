package domain

type PlayerPositionCode string // @name PlayerPositionCode

const (
	PlayerPositionCodeGLK PlayerPositionCode = "GLK"
	PlayerPositionCodeDEF PlayerPositionCode = "DEF"
	PlayerPositionCodeMID PlayerPositionCode = "MID"
	PlayerPositionCodeATK PlayerPositionCode = "ATK"
)

func (e PlayerPositionCode) Valid() bool {
	switch e {
	case PlayerPositionCodeGLK,
		PlayerPositionCodeDEF,
		PlayerPositionCodeMID,
		PlayerPositionCodeATK:
		return true
	}
	return false
}

type PlayerPosition struct {
	Code  PlayerPositionCode
	Label string
}

type PlayerPositionWithLocale struct {
	PlayerPosition
	Locale LocaleCode
}
