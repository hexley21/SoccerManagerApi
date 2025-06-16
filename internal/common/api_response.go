package common

type apiResponse struct {
	Data any `json:"data"`
} // @name ApiResponse

func NewApiResponse(data any) apiResponse {
	return apiResponse{Data: data}
}
