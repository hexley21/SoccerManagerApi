package globe

import (
	"net/http"

	"github.com/hexley21/soccer-manager/internal/common"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/service"
	"github.com/labstack/echo/v4"
)

type handler struct {
	globeService service.GlobeService
}

func NewHandler(globeService service.GlobeService) *handler {
	return &handler{
		globeService: globeService,
	}
}

// @Summary List locales
// @Description Get all available locales
// @Tags globe
// @Produce json
// @Success 200 {object} common.apiResponse{data=[]domain.LocaleCode} "OK"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/globe/locales [get]
func (h *handler) Locales(c echo.Context) error {
	locales, err := h.globeService.ListLocales(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(
		http.StatusOK,
		common.NewApiResponse(locales),
	)
}

// @Summary List Countries
// @Description Get all available countries
// @Tags globe
// @Produce json
// @Success 200 {object} common.apiResponse{data=[]domain.CountryCode} "OK"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/globe/countries [get]
func (h *handler) Countries(c echo.Context) error {
	countries, err := h.globeService.ListCountries(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(
		http.StatusOK,
		common.NewApiResponse(countries),
	)
}
