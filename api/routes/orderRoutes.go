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

	huma.Register(api, huma.Operation{
		OperationID:   "public-order",
		Method:        http.MethodPost,
		Path:          basePath + "/demander/public",
		Summary:       "Public order",
		Tags:          []string{"transaction"},
		DefaultStatus: http.StatusCreated,
	}, handler.PublicOrder)

	huma.Register(api, huma.Operation{
		OperationID:   "supplier-confirm-order",
		Method:        http.MethodPatch,
		Path:          basePath + "/supplier/confirm/{did}",
		Summary:       "Supplier confirm order",
		Tags:          []string{"transaction"},
		DefaultStatus: http.StatusOK,
	}, handler.SupplierConfirmOrder)

	huma.Register(api, huma.Operation{
		OperationID:   "demander-select-order",
		Method:        http.MethodPatch,
		Path:          basePath + "/demander/select/{did}",
		Summary:       "Demander select order",
		Tags:          []string{"transaction"},
		DefaultStatus: http.StatusOK,
	}, handler.DemanderSelectOrder)

	huma.Register(api, huma.Operation{
		OperationID:   "supplier-approve-order",
		Method:        http.MethodPatch,
		Path:          basePath + "/supplier/approve/{did}",
		Summary:       "Supplier approve order",
		Tags:          []string{"transaction"},
		DefaultStatus: http.StatusOK,
	}, handler.SupplierApproveOrder)

	huma.Register(api, huma.Operation{
		OperationID:   "demander-approve-order",
		Method:        http.MethodPatch,
		Path:          basePath + "/demander/approve/{did}",
		Summary:       "Demander approve order",
		Tags:          []string{"transaction"},
		DefaultStatus: http.StatusOK,
	}, handler.DemanderApproveOrder)

	huma.Register(api, huma.Operation{
		OperationID:   "supplier-cancel-order",
		Method:        http.MethodPatch,
		Path:          basePath + "/supplier/cancel/{did}",
		Summary:       "Supplier cancel order",
		Tags:          []string{"transaction"},
		DefaultStatus: http.StatusOK,
	}, handler.SupplierCancelOrder)

	huma.Register(api, huma.Operation{
		OperationID:   "demander-cancel-order",
		Method:        http.MethodPatch,
		Path:          basePath + "/demander/cancel/{did}",
		Summary:       "Demander cancel order",
		Tags:          []string{"transaction"},
		DefaultStatus: http.StatusOK,
	}, handler.DemanderCancelOrder)
}
