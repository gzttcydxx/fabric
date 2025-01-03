package routes

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gzttcydxx/newapi/handlers"
	"github.com/gzttcydxx/newapi/models"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// RegisterCRUD 注册CRUD路由
func RegisterCRUD[T any](api huma.API, contract *client.Contract, basePath string, resourceName string, methods models.CRUDMethods) {
	handler := handlers.NewCRUDHandler[T](resourceName, contract, methods)

	huma.Register(api, huma.Operation{
		OperationID:   "create-" + resourceName,
		Method:        http.MethodPost,
		Path:          basePath,
		Summary:       "Create " + resourceName,
		Tags:          []string{resourceName},
		DefaultStatus: http.StatusCreated,
	}, handler.Create)

	huma.Register(api, huma.Operation{
		OperationID:   "get-" + resourceName,
		Method:        http.MethodGet,
		Path:          basePath + "/{did}",
		Summary:       "Get " + resourceName,
		Tags:          []string{resourceName},
		DefaultStatus: http.StatusOK,
	}, handler.Get)

	huma.Register(api, huma.Operation{
		OperationID:   "query-" + resourceName,
		Method:        http.MethodPost,
		Path:          basePath + "/query",
		Summary:       "Query " + resourceName,
		Tags:          []string{resourceName},
		DefaultStatus: http.StatusOK,
	}, handler.Query)

	huma.Register(api, huma.Operation{
		OperationID:   "update-" + resourceName,
		Method:        http.MethodPatch,
		Path:          basePath,
		Summary:       "Update " + resourceName,
		Tags:          []string{resourceName},
		DefaultStatus: http.StatusOK,
	}, handler.Update)

	huma.Register(api, huma.Operation{
		OperationID:   "delete-" + resourceName,
		Method:        http.MethodDelete,
		Path:          basePath + "/{did}",
		Summary:       "Delete " + resourceName,
		Tags:          []string{resourceName},
		DefaultStatus: http.StatusOK,
	}, handler.Delete)
}
