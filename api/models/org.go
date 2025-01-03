package models

import "github.com/gzttcydxx/did/models"

// Org 组织模型
type Org struct {
	Did                 models.DID `json:"did" required:"false"`
	UUID                string     `json:"uuid" required:"false" example:"0effa72b-3a0a-4ca1-8376-106ccd636584"`
	Name                string     `json:"name" required:"false" example:"苏州东仪核电科技股份有限公司"`
	Class               string     `json:"class" required:"false" example:"ENTERPRISE"`
	JuridicalPerson     string     `json:"juridical_person" required:"false" example:"陆建忠"`
	RegisteredCapital   string     `json:"registered_capital" required:"false" example:"2,700万(元)"`
	DateOfEstablishment string     `json:"date_of_establishment" required:"false" example:"1999-07-13"`
	ManagementState     string     `json:"management_state" required:"false" example:"开业"`
	Province            string     `json:"province" required:"false" example:"江苏省"`
	City                string     `json:"city" required:"false" example:"苏州市"`
	County              string     `json:"county" required:"false" example:"吴中区"`
	EnterpriseType      string     `json:"enterprise_type" required:"false" example:"股份有限公司(非上市)"`
	SocialCreditCode    string     `json:"social_credit_code" required:"false" example:"913205007149896072"`
	TaxNumber           string     `json:"tax_number" required:"false" example:"913205007149896072"`
	CompanyRegistration string     `json:"company_registration" required:"false" example:"320506000017856"`
	OrganizationCode    string     `json:"organization_code" required:"false" example:"71498960-7"`
	PhoneNumber         string     `json:"phone_number" required:"false" example:"0512-66265880"`
	Industry            string     `json:"industry" required:"false" example:"专用设备制造业"`
	Address             string     `json:"address" required:"false" example:"苏州市吴中区木渎镇金桥工业园南区（木东路）"`
	Website             string     `json:"website" required:"false" example:"www.eiec.cc"`
	Email               string     `json:"email" required:"false" example:"eiec_ljz@188.com"`
}
