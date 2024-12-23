package models

import (
	"github.com/gzttcydxx/did/models"
)

type Org struct {
	Did                 models.DID `json:"did"`                   // 组织DID
	UUID                string     `json:"uuid"`                  // 唯一标识符
	Name                string     `json:"name"`                  // 组织名称
	Class               string     `json:"class"`                 // 组织类型（如：ENTERPRISE）
	JuridicalPerson     string     `json:"juridical_person"`      // 法人代表
	RegisteredCapital   string     `json:"registered_capital"`    // 注册资本
	DateOfEstablishment string     `json:"date_of_establishment"` // 成立日期
	ManagementState     string     `json:"management_state"`      // 经营状态
	Province            string     `json:"province"`              // 省份
	City                string     `json:"city"`                  // 城市
	County              string     `json:"county"`                // 区县
	EnterpriseType      string     `json:"enterprise_type"`       // 企业类型
	SocialCreditCode    string     `json:"social_credit_code"`    // 统一社会信用代码
	TaxNumber           string     `json:"tax_number"`            // 纳税人识别号
	CompanyRegistration string     `json:"company_registration"`  // 工商注册号
	OrganizationCode    string     `json:"organization_code"`     // 组织机构代码
	PhoneNumber         string     `json:"phone_number"`          // 联系电话
	Industry            string     `json:"industry"`              // 所属行业
	Address             string     `json:"address"`               // 地址
	Website             string     `json:"website"`               // 网站
	Email               string     `json:"email"`                 // 电子邮箱
}
