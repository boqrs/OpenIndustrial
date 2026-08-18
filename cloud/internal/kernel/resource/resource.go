package resource

import (
	"context"
	"crypto"
	"crypto/x509/pkix" // 修正: 导入 pkix
	"net"
	"time"

	"github.com/google/uuid"

	"github.com/OpenIndustrial/cloud/internal/persistence/model"

)

// ResourceRepository defines the persistence interface for Resource entities.
type ResourceRepository interface {
	CreateResource(ctx context.Context, resource *model.Resource) error
	GetResourceByID(ctx context.Context, tenantID, resourceID uuid.UUID) (*model.Resource, error)
	UpdateResource(ctx context.Context, resource *model.Resource) error
	DeleteResource(ctx context.Context, tenantID, resourceID uuid.UUID) error
	ListResources(ctx context.Context, tenantID uuid.UUID, resourceType string, limit, offset int) ([]*model.Resource, error)
	CheckUserInSameGroupAsResource(ctx context.Context, userID, resourceID uuid.UUID) (bool, error)
	BatchCreateResources(ctx context.Context, resources []*model.Resource) error // New
	FindResourceByNameAndType(ctx context.Context, tenantID uuid.UUID, name, resourceType string) (*model.Resource, error)
	UpdateParent(ctx context.Context, tenantID, resourceID, newParentID uuid.UUID) error
	Exists(ctx context.Context, id uuid.UUID) (bool, error) 
	FindByParentID(ctx context.Context, tenantID, parentID uuid.UUID) ([]*model.Resource, error)
}

// AttributeDefinitionRepository defines the persistence interface for AttributeDefinition entities.
type AttributeDefinitionRepository interface {
	CreateAttributeDefinition(ctx context.Context, def *model.AttributeDefinition) error
	BatchCreateAttributeDefinition(ctx context.Context, attrs []*model.AttributeDefinition)error
	GetAttributeDefinitionByID(ctx context.Context, tenantID, defID uuid.UUID) (*model.AttributeDefinition, error)
	FindByIDs(ctx context.Context, tenantID uuid.UUID, ids []uuid.UUID) ([]*model.AttributeDefinition, error)
	GetAttributeDefinitionByKey(ctx context.Context, tenantID uuid.UUID, key string) (*model.AttributeDefinition, error)
	ListAttributeDefinitions(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*model.AttributeDefinition, error)
	UpdateAttributeDefinition(ctx context.Context, def *model.AttributeDefinition) error
	DeleteAttributeDefinition(ctx context.Context, tenantID, defID uuid.UUID) error
	FindByName(ctx context.Context, tenantID uuid.UUID, name string) (*model.AttributeDefinition, error)
	FindAttributeDefinitionByResourceID(ctx context.Context, resourceID uuid.UUID)([]*model.AttributeDefinition, error) 
}

// ResourceAttributeRepository defines the persistence interface for ResourceAttribute values.
type ResourceAttributeRepository interface {
	SetAttribute(ctx context.Context, attr *model.ResourceAttribute) error
	SetAttributes(ctx context.Context, attrs []*model.ResourceAttribute) error
	GetAttribute(ctx context.Context, resourceID, attributeID uuid.UUID) (*model.ResourceAttribute, error)
	GetAttributesByResourceID(ctx context.Context, resourceID uuid.UUID) ([]*model.ResourceAttribute, error)
	DeleteAttribute(ctx context.Context, resourceID, attributeID uuid.UUID) error
	GetForResource(ctx context.Context, resourceID uuid.UUID) (map[string]interface{}, error)
	UpsertForResource(ctx context.Context, tenantID, resourceID uuid.UUID, attributes map[string]interface{}) error
	BatchCreateResourceAttributes(ctx context.Context, attr []*model.ResourceAttribute) error
    GetAttributesForResource(ctx context.Context, tenantID, resourceID uuid.UUID) (map[string]interface{}, error)

}

type ResourceConnectionsRepository interface {
	CreateConnection(ctx context.Context, conn *model.ResourceConnection) error
	GetConnectionByID(ctx context.Context, conID uint)(*model.ResourceConnection, error)
	DeleteConnection(ctx context.Context, connectionID uint) error
	ListConnectionsByResourceID(ctx context.Context, resourceID uuid.UUID) ([]*model.ResourceConnection, error)
}

type ParsedCSR struct {
	Subject      pkix.Name // 修正: 类型从 string 变为 pkix.Name
	PublicKey    crypto.PublicKey
	DNSNames     []string
	IPAddresses  []net.IP
	EmailAddress []string
}

// IssueCertificateRequest is the request to issue a new certificate.
type IssueCertificateRequest struct {
	CSR      string
	Validity time.Duration // 修正: 增加 Validity 字段
}

// IssuedCertificate is the result of a successful certificate issuance.
type IssuedCertificate struct {
	CertificatePEM string
	CertificateID  string
}

type CertificateAuthority interface {
	ValidateCSR(
		csrPEM string,
	) (*ParsedCSR, error)

	IssueCertificate(
		ctx context.Context,
		req IssueCertificateRequest,
	) (*IssuedCertificate, error)

	RevokeCertificate(
		ctx context.Context,
		certificateID string,
		reason string,
	) error
}