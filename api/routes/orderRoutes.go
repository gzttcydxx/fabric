package routes

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gzttcydxx/newapi/handlers"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// RegisterOrder 注册订单路由
func RegisterOrder(api huma.API, contract *client.Contract, basePath string) {
	handler := handlers.NewOrderHandler(contract)

	// 买方创建订单并选择卖家
	huma.Register(api, huma.Operation{
		OperationID:   "buyer-create-order-with-seller",
		Method:        http.MethodPost,
		Path:          basePath + "/buyer/create",
		Summary:       "Buyer create order with seller",
		Tags:          []string{"transaction"},
		DefaultStatus: http.StatusCreated,
	}, handler.BuyerCreateOrderWithSeller)

	// 卖家确认订单
	huma.Register(api, huma.Operation{
		OperationID:   "seller-confirm-order",
		Method:        http.MethodPatch,
		Path:          basePath + "/seller/confirm/{did}",
		Summary:       "Seller confirm order",
		Tags:          []string{"transaction"},
		DefaultStatus: http.StatusOK,
	}, handler.SellerConfirmOrder)

	// 买家确认订单
	huma.Register(api, huma.Operation{
		OperationID:   "buyer-confirm-order",
		Method:        http.MethodPatch,
		Path:          basePath + "/buyer/confirm/{did}",
		Summary:       "Buyer confirm order",
		Tags:          []string{"transaction"},
		DefaultStatus: http.StatusOK,
	}, handler.BuyerConfirmOrder)

	// 卖家取消订单
	huma.Register(api, huma.Operation{
		OperationID:   "seller-cancel-order",
		Method:        http.MethodPatch,
		Path:          basePath + "/seller/cancel/{did}",
		Summary:       "Seller cancel order",
		Tags:          []string{"transaction"},
		DefaultStatus: http.StatusOK,
	}, handler.SellerCancelOrder)

	// 买家取消订单
	huma.Register(api, huma.Operation{
		OperationID:   "buyer-cancel-order",
		Method:        http.MethodPatch,
		Path:          basePath + "/buyer/cancel/{did}",
		Summary:       "Buyer cancel order",
		Tags:          []string{"transaction"},
		DefaultStatus: http.StatusOK,
	}, handler.BuyerCancelOrder)
}
