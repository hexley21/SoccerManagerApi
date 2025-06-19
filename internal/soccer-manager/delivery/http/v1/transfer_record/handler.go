package transfer_record

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/hexley21/soccer-manager/internal/common"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/service"
	"github.com/labstack/echo/v4"
)

type handler struct {
	transferRecordService service.TransferRecordService
	pageSize              int32
	pageLimit             int32
}

func newHandler(
	transferRecordService service.TransferRecordService,
	pageSize int32,
	pageLimit int32,
) *handler {
	return &handler{
		transferRecordService: transferRecordService,
		pageSize:              pageSize,
		pageLimit:             pageLimit,
	}
}

// @Summary Get all transfer records
// @Description Returns a list of all transfer records (paginated or filtered as needed)
// @Tags transfer-records
// @Produce json
// @Param cursor query int false "Pagination cursor"
// @Param page_size query int false "Number of items per page"
// @Success 200 {object} common.apiResponse{data=[]transferRecordResponseDTO} "OK"
// @Failure 400 {object} echo.HTTPError "Bad Request"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/transfer-records [get]
func (h *handler) GetTransferRecords(c echo.Context) error {
	pagination, err := common.ParsePagination(c, h.pageSize, h.pageLimit)
	if err != nil {
		return err
	}

	records, err := h.transferRecordService.ListTransferRecords(
		c.Request().Context(),
		pagination.Cursor,
		h.pageLimit,
	)
	if err != nil {
		return err
	}

	res := make([]transferRecordResponseDTO, len(records))
	for i, r := range records {
		res[i] = transferRecordResponseAdapter(r)
	}

	return c.JSON(http.StatusOK, common.NewApiResponse(res))
}

// @Summary Get transfer record by ID
// @Description Returns a single transfer record by its ID
// @Tags transfer-records
// @Produce json
// @Param record_id path int true "Record ID"
// @Success 200 {object} common.apiResponse{data=transferRecordResponseDTO} "OK"
// @Failure 400 {object} echo.HTTPError "Bad Request"
// @Failure 404 {object} echo.HTTPError "Not Found"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/transfer-records/{record_id} [get]
func (h *handler) GetTransferRecordById(c echo.Context) error {
	recordId, err := strconv.ParseInt(c.Param("record_id"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	record, err := h.transferRecordService.GetTransferRecordByID(
		c.Request().Context(),
		recordId,
	)
	if err != nil {
		if errors.Is(err, service.ErrTransferRecordNotFound) {
			return echo.ErrNotFound.WithInternal(err)
		}

		return err
	}

	return c.JSON(http.StatusOK, common.NewApiResponse(transferRecordResponseAdapter(record)))
}
