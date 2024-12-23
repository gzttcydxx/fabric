package identity

// import (
// 	"encoding/pem"
// 	"fmt"

// 	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
// 	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
// 	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
// 	"github.com/hyperledger/fabric-sdk-go/pkg/msp"
// )

// // MSPManager 管理MSP相关操作
// type MSPManager struct {
// 	sdk       *fabsdk.FabricSDK
// 	mspClient *mspclient.Client
// }

// // NewMSPManager 创建新的MSP管理器
// func NewMSPManager(configPath string, orgName string) (*MSPManager, error) {
// 	// 初始化SDK
// 	sdk, err := fabsdk.New(config.FromFile(configPath))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create SDK: %v", err)
// 	}

// 	// 创建MSP客户端上下文
// 	ctxProvider := sdk.Context(fabsdk.WithOrg(orgName))

// 	// 创建MSP客户端
// 	mspClient, err := mspclient.New(ctxProvider)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create MSP client: %v", err)
// 	}

// 	return &MSPManager{
// 		sdk:       sdk,
// 		mspClient: mspClient,
// 	}, nil
// }

// // RegisterUser 注册新用户并获取证书
// func (m *MSPManager) RegisterUser(userdid string, orgdid string) (string, error) {
// 	// 注册用户
// 	registerRequest := &mspclient.RegistrationRequest{
// 		Name:        userdid,
// 		Type:        "client",
// 		Affiliation: orgdid,
// 		Attributes: []mspclient.Attribute{
// 			{
// 				Name:  "crosschain",
// 				Value: "true",
// 				ECert: true,
// 			},
// 		},
// 	}

// 	secret, err := m.mspClient.Register(registerRequest)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to register user: %v", err)
// 	}

// 	// 登记用户并获取证书
// 	enrollRequest := &msp.EnrollmentRequest{
// 		Name:   userdid,
// 		Secret: secret,
// 	}

// 	enrollment, err := m.mspClient.Enroll(enrollRequest)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to enroll user: %v", err)
// 	}

// 	// 将证书转换为PEM格式
// 	certPEM := pem.EncodeToMemory(&pem.Block{
// 		Type:  "CERTIFICATE",
// 		Bytes: enrollment.Identity.GetECert().Cert().Raw,
// 	})

// 	return string(certPEM), nil
// }

// // ValidateCertificate 验证证书
// func (m *MSPManager) ValidateCertificate(certPEM string) (bool, error) {
// 	// 解析证书
// 	block, _ := pem.Decode([]byte(certPEM))
// 	if block == nil {
// 		return false, fmt.Errorf("failed to decode certificate PEM")
// 	}

// 	// 验证证书是否在当前MSP中有效
// 	identity, err := m.mspClient.GetSigningIdentity(block.Bytes)
// 	if err != nil {
// 		return false, nil
// 	}

// 	return identity != nil, nil
// }
