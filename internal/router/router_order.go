package router

import orderHandler "github.com/xinliangnote/go-gin-api/internal/api/order"

func setOrderRouter(r *resource) {

	order := r.mux.Group("/api/order", r.interceptors.CheckSignature())
	{
		// order 控制器
		handler := orderHandler.New(r.logger, r.db, r.cache)
		order.POST("/create", handler.Create())
		order.POST("/cancel", handler.Cancel())
		order.GET("/:id", handler.Detail())
	}
}
