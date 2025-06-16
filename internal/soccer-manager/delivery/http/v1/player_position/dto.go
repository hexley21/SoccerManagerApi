package player_position

import "github.com/hexley21/soccer-manager/internal/soccer-manager/domain"

type playerPositionsResponseDTO map[domain.PlayerPositionCode]string                                  // @name PlayerPositionsResponse
type playerPositionsWithLocalesResponseDTO map[domain.LocaleCode]map[domain.PlayerPositionCode]string // @name PlayerPositionsWithLocalesResponse

type createPlayerPositionTranslationRequestDTO struct {
	Code   domain.PlayerPositionCode `json:"code"   validate:"required,playerpos"`
	Locale domain.LocaleCode         `json:"locale" validate:"required,localecode"`
	Label  string                    `json:"label"  validate:"required,alphaunicode,min=2,max=30"`
} // @name CreatePlayerPositionTranslationRequest

type updatePlayerPositionTranslationRequestDTO struct {
	Code   domain.PlayerPositionCode `json:"code"   validate:"required,playerpos"`
	Locale domain.LocaleCode         `json:"locale" validate:"required,localecode"`
	Label  string                    `json:"label"  validate:"required,alphaunicode,min=2,max=30"`
} // @name UpdatePlayerPositionTranslationRequest

type deletePlayerPositionTranslationRequestDTO struct {
	Code   domain.PlayerPositionCode `json:"code"   validate:"required,playerpos"`
	Locale domain.LocaleCode         `json:"locale" validate:"required,localecode"`
} // @name DeletelayerPositionTranslationRequest
