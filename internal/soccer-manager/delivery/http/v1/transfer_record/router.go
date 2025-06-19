package transfer_record

import (
	"github.com/hexley21/soccer-manager/internal/soccer-manager/delivery"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group, c *delivery.Components) {
	h := newHandler(c.Services.TransferRecordService, c.Cfg.Pagination.M, c.Cfg.Pagination.L)

	g.GET("/transfer-records", h.GetTransferRecords)
	g.GET("/transfer-records/:record_id", h.GetTransferRecordById)
}
