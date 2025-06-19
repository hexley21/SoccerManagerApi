package transfer

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/hexley21/soccer-manager/internal/common"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/jwt/access"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/service"
	"github.com/labstack/echo/v4"
)

type handler struct {
	transferService service.TransferService
	pageSize        int32
	pageLimit       int32
}

func newHandler(transferService service.TransferService, pageSize int32, pageLimit int32) *handler {
	return &handler{
		transferService: transferService,
		pageSize:        pageSize,
		pageLimit:       pageLimit,
	}
}

// @Summary List all transfers
// @Description Returns a list of all transfers (paginated)
// @Tags transfers
// @Produce json
// @Param cursor query int false "Pagination cursor"
// @Param page_size query int false "Number of items per page"
// @Success 200 {object} common.apiResponse{data=[]transferResponseDTO} "OK"
// @Failure 400 {object} echo.HTTPError "Bad Request"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/transfers [get]
func (h *handler) GetTransfers(c echo.Context) error {
	pagination, err := common.ParsePagination(c, h.pageSize, h.pageLimit)
	if err != nil {
		return err
	}

	transfers, err := h.transferService.ListTransfers(
		c.Request().Context(),
		pagination.Cursor,
		pagination.PageSize,
	)
	if err != nil {
		return err
	}

	res := make([]transferResponseDTO, len(transfers))
	for i, tr := range transfers {
		res[i] = transferResponseAdapter(tr)
	}

	return c.JSON(http.StatusOK, common.NewApiResponse(res))
}

// @Summary Get transfer with id
// @Description Returns a transfer by id
// @Tags transfers
// @Produce json
// @Param transfer_id path int true "Transfer ID"
// @Success 200 {object} common.apiResponse{data=transferResponseDTO} "OK"
// @Failure 400 {object} echo.HTTPError "Bad Request"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/transfers/{transfer_id} [get]
func (h *handler) GetTransferById(c echo.Context) error {
	transferId, err := strconv.ParseInt(c.Param("transfer_id"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	transfer, err := h.transferService.GetTransferByID(c.Request().Context(), transferId)
	if err != nil {
		if errors.Is(err, service.ErrTransferNotFound) {
			return echo.ErrNotFound.WithInternal(err)
		}

		return err
	}

	return c.JSON(http.StatusOK, common.NewApiResponse(transferResponseAdapter(transfer)))
}

// @Summary List transfers by team ID
// @Description Returns a list of transfers for the provided team ID (paginated)
// @Tags transfers
// @Produce json
// @Param team_id path int true "Team ID"
// @Param cursor query int false "Pagination cursor"
// @Param page_size query int false "Number of items per page"
// @Success 200 {object} common.apiResponse{data=[]transferResponseDTO} "OK"
// @Failure 400 {object} echo.HTTPError "Bad Request"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/teams/{team_id}/transfers [get]
func (h *handler) GetTransfersByTeamId(c echo.Context) error {
	teamId, err := strconv.ParseInt(c.Param("team_id"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	pagination, err := common.ParsePagination(c, h.pageSize, h.pageLimit)
	if err != nil {
		return err
	}

	transfers, err := h.transferService.ListTransfersByTeamId(
		c.Request().Context(),
		teamId,
		pagination.Cursor,
		pagination.PageSize,
	)
	if err != nil {
		return err
	}

	res := make([]transferResponseDTO, len(transfers))
	for i, tr := range transfers {
		res[i] = transferResponseAdapter(tr)
	}

	return c.JSON(http.StatusOK, common.NewApiResponse(res))
}

// @Summary Get transfer of a player
// @Description Returns a transfer record of a player
// @Tags transfers
// @Produce json
// @Param player_id path int true "Player ID"
// @Success 200 {object} common.apiResponse{data=transferResponseDTO} "OK"
// @Failure 400 {object} echo.HTTPError "Bad Request"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/players/{player_id}/transfer [get]
func (h *handler) GetTransferByPlayer(c echo.Context) error {
	playerId, err := strconv.ParseInt(c.Param("player_id"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	transfer, err := h.transferService.GetTransferByPlayerID(c.Request().Context(), playerId)
	if err != nil {
		if errors.Is(err, service.ErrTransferNotFound) {
			return echo.ErrNotFound.WithInternal(err)
		}

		return err
	}

	return c.JSON(http.StatusOK, common.NewApiResponse(transferResponseAdapter(transfer)))
}

// @Summary Create transfer
// @Description Creates a new transfer listing
// @Tags transfers
// @Accept json
// @Produce json
// @Security AccessToken
// @Param request body createTransferRequestDTO true "Transfer details"
// @Success 201 {object} common.apiResponse{data=int} "Created"
// @Failure 400 {object} echo.HTTPError "Bad Request"
// @Failure 404 {object} echo.HTTPError "Not Found"
// @Failure 409 {object} echo.HTTPError "Conflict"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/transfers [post]
func (h *handler) CreateTransfer(c echo.Context) error {
	userData, ok := c.Get(access.CtxKey).(access.Data)
	if !ok {
		return echo.ErrUnauthorized.WithInternal(access.NewInvalidTokenError(userData))
	}

	var req createTransferRequestDTO
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	transferId, err := h.transferService.CreateTransfer(
		c.Request().Context(),
		userData.UserID,
		req.PlayerID,
		req.Price.IntPart(),
	)
	if err != nil {
		if errors.Is(err, service.ErrNonexistentCode) {
			return echo.ErrNotFound.WithInternal(err)
		}
		if errors.Is(err, service.ErrPlayerAlreadyInTransfers) {
			return echo.ErrConflict.WithInternal(err)
		}
		if errors.Is(err, service.ErrInvalidArguments) {
			return echo.ErrBadRequest.WithInternal(err)
		}

		return err
	}

	return c.JSON(
		http.StatusCreated,
		common.NewApiResponse(transferId),
	)
}

// @Summary Delete transfer
// @Description Deletes a transfer by ID
// @Tags transfers
// @Produce json
// @Security AccessToken
// @Param transfer_id path int true "Transfer ID"
// @Success 204 "No Content"
// @Failure 400 {object} echo.HTTPError "Bad Request"
// @Failure 404 {object} echo.HTTPError "Not Found"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/transfers/{transfer_id} [delete]
func (h *handler) DeleteTransfer(c echo.Context) error {
	transferId, err := strconv.ParseInt(c.Param("transfer_id"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	userData, ok := c.Get(access.CtxKey).(access.Data)
	if !ok {
		return echo.ErrUnauthorized.WithInternal(access.NewInvalidTokenError(userData))
	}

	if err := h.transferService.DeleteTransfer(c.Request().Context(), transferId, userData.UserID); err != nil {
		if errors.Is(err, service.ErrTransferNotFound) {
			return echo.ErrNotFound.WithInternal(err)
		}

		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// @Summary Update transfer
// @Description Updates an existing transfer by ID
// @Tags transfers
// @Accept json
// @Produce json
// @Security AccessToken
// @Param transfer_id path int true "Transfer ID"
// @Param request body updateTransferRequestDTO true "Updated transfer info"
// @Success 200 "OK"
// @Failure 400 {object} echo.HTTPError "Bad Request"
// @Failure 404 {object} echo.HTTPError "Not Found"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/transfers/{transfer_id} [put]
func (h *handler) UpdateTransfer(c echo.Context) error {
	transferId, err := strconv.ParseInt(c.Param("transfer_id"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	userData, ok := c.Get(access.CtxKey).(access.Data)
	if !ok {
		return echo.ErrUnauthorized.WithInternal(access.NewInvalidTokenError(userData))
	}

	var req updateTransferRequestDTO
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	if err := h.transferService.UpdateTransferPrice(c.Request().Context(), transferId, userData.UserID, req.Price.IntPart()); err != nil {
		if errors.Is(err, service.ErrTransferNotFound) {
			return echo.ErrNotFound.WithInternal(err)
		}

		return err
	}

	return c.NoContent(http.StatusOK)
}

// @Summary Buy a player
// @Description Purchases a player in a specific transfer listing
// @Tags transfers
// @Produce json
// @Security AccessToken
// @Param transfer_id path int true "Transfer ID"
// @Success 201 "Created"
// @Failure 400 {object} echo.HTTPError "Bad Request"
// @Failure 401 {object} echo.HTTPError "Unauthorized"
// @Failure 404 {object} echo.HTTPError "Not Found"
// @Failure 409 {object} echo.HTTPError "Conflict"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/transfers/{transfer_id}/buy [post]
func (h *handler) BuyPlayer(c echo.Context) error {
	transferId, err := strconv.ParseInt(c.Param("transfer_id"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	userData, ok := c.Get(access.CtxKey).(access.Data)
	if !ok {
		return echo.ErrUnauthorized.WithInternal(access.NewInvalidTokenError(userData))
	}

	if err := h.transferService.BuyPlayer(c.Request().Context(), transferId, userData.UserID); err != nil {
		if errors.Is(err, service.ErrTransferNotFound) {
			return echo.ErrNotFound.WithInternal(err)
		}
		if errors.Is(err, service.ErrCantBuyFromYourself) {
			return echo.ErrConflict.WithInternal(err)
		}
		if errors.Is(err, service.ErrNotEnoughFunds) {
			return echo.ErrPaymentRequired.WithInternal(err)
		}

		return err
	}

	return c.NoContent(http.StatusOK)
}
