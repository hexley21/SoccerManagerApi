package common

import (
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Params struct {
	Cursor   int64
	PageSize int32
}

// ParsePagination parses inputs from the request and returns processed cursor and page size params
// the error is always instance of *echo.HttpError
func ParsePagination(c echo.Context, defaultPageSize int32, maxPageSize int32) (Params, error) {
	var cursor int64
	var pageSize int32

	cursorQuery := c.QueryParam("cursor")
	if cursorQuery == "" {
		cursor = 0
	} else {
		cTmp, err := strconv.ParseInt(cursorQuery, 10, 64)
		if err != nil {
			return Params{}, echo.ErrBadRequest.WithInternal(fmt.Errorf("invalid cursor: %w", err))
		}
		cursor = cTmp
	}

	pageSizeQuery := c.QueryParam("page_size")
	if pageSizeQuery == "" || pageSizeQuery == "0" {
		pageSize = defaultPageSize
	} else {
		psTmp, err := strconv.ParseInt(pageSizeQuery, 10, 32)
		if err != nil {
			return Params{}, echo.ErrBadRequest.WithInternal(fmt.Errorf("invalid page_size: %w", err))
		}
		pageSize = int32(psTmp)
	}

	if (pageSize > maxPageSize) || (pageSize < 0) {
		return Params{}, echo.ErrBadRequest.WithInternal(fmt.Errorf("invalid page_size: %v", pageSize))
	}

	return Params{
		Cursor:   cursor,
		PageSize: pageSize,
	}, nil
}
