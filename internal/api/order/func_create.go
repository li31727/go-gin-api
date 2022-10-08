package order

import (
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
)

type createRequest struct {
	OrderID string `json:"order_id"`
}

type createResponse struct {
	Code int `json:"code"`
}

// Create 创建订单
// @Summary 创建订单
// @Description 创建订单
// @Tags API.order
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param Request body createRequest true "请求信息"
// @Success 200 {object} createResponse
// @Failure 400 {object} code.Failure
// @Router /api/order/create [post]
func (h *handler) Create() core.HandlerFunc {
	return func(ctx core.Context) {
		res := new(createResponse)
		res.Code = 0
		ctx.Payload(res)
		return

	}
}
