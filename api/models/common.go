package models

// CRUDMethods 定义 CRUD 方法名
type CRUDMethods struct {
	Create string // 创建方法名
	Read   string // 读取方法名
	Query  string // 查询方法名
	Update string // 更新方法名
	Delete string // 删除方法名
}

// JSONBody 通用请求体
type JSONBody[T any] struct {
	Body T
}

// Status 定义响应体结构
type Status struct {
	Success bool `json:"success" example:"true"`
}

// ListBody 通用列表响应体
type List[T any] struct {
	Items []T `json:"items"`
}

// GetInput 通用获取请求
type GetInput struct {
	Did string `path:"did" doc:"DID"`
}
