package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hexley21/soccer-manager/internal/soccer-manager/domain"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/jwt/access"
	mock_jwt "github.com/hexley21/soccer-manager/internal/soccer-manager/jwt/mock"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/server/middleware"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_JWTAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManager := mock_jwt.NewMockManagerWithTTL[access.Data](ctrl)
	mw := middleware.JWTAuth("", mockManager)

	t.Run("valid token", func(t *testing.T) {
		mockManager.EXPECT().
			ParseTokenString("validToken").
			Return(access.NewData(123, domain.UserRoleADMIN), nil)

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set("Authorization", "Bearer validToken")
		rec := httptest.NewRecorder()
		ctx := echo.New().NewContext(req, rec)

		_ = mw(func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		})(ctx)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "test", rec.Body.String())
		assert.Equal(t, access.NewData(123, domain.UserRoleADMIN), ctx.Get(access.CtxKey))
	})

	t.Run("invalid token", func(t *testing.T) {
		mockManager.EXPECT().
			ParseTokenString("invalidToken").
			Return(access.Data{}, echo.ErrUnauthorized)

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set("Authorization", "Bearer invalidToken")
		rec := httptest.NewRecorder()
		ctx := echo.New().NewContext(req, rec)

		_ = mw(func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		})(ctx)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("no token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		ctx := echo.New().NewContext(req, rec)

		_ = mw(func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		})(ctx)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("skip allowed paths", func(t *testing.T) {
		t.Run("GET request", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			ctx := echo.New().NewContext(req, rec)

			_ = mw(func(c echo.Context) error {
				return c.String(http.StatusOK, "test")
			})(ctx)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "test", rec.Body.String())
		})

		testCases := []string{"/auth"}
		for _, path := range testCases {
			t.Run(path, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, path, nil)
				rec := httptest.NewRecorder()
				ctx := echo.New().NewContext(req, rec)

				_ = mw(func(c echo.Context) error {
					return c.String(http.StatusOK, "test")
				})(ctx)

				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Equal(t, "test", rec.Body.String())
			})
		}
	})

	t.Run("include GET admin paths", func(t *testing.T) {
		testCases := []string{"/users"}
		for _, path := range testCases {
			t.Run(path, func(t *testing.T) {
				mockManager.EXPECT().
					ParseTokenString("validToken").
					Return(access.NewData(123, domain.UserRoleADMIN), nil)

				req := httptest.NewRequest(http.MethodGet, path, nil)
				req.Header.Set("Authorization", "Bearer validToken")
				rec := httptest.NewRecorder()
				ctx := echo.New().NewContext(req, rec)

				_ = mw(func(c echo.Context) error {
					return c.String(http.StatusOK, "test")
				})(ctx)

				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Equal(t, "test", rec.Body.String())
				assert.Equal(t, access.NewData(123, domain.UserRoleADMIN), ctx.Get(access.CtxKey))
			})
		}
	})
}
