package postgres

import (
	"context"
	"time"

	"github.com/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

)

// credentialRepository implements the security.CredentialRepository interface.
type CertificateRepository struct {
	db *gorm.DB
}

// NewCredentialRepository creates a new repository for resource credentials.
func NewCredentialRepository(db *gorm.DB) *CredentialRepository {
	return &CredentialRepository{db: db}
}

// Compile-time check to ensure credentialRepository implements the interface.
//var _ security.CredentialRepository = (*credentialRepository)(nil)

func (r *CertificateRepository) Create(ctx context.Context, cred *model.ResourceCertificate) error {
	return r.db.WithContext(ctx).Create(cred).Error
}

func (r *CertificateRepository) FindBySecretHash(ctx context.Context, hash string) (*model.ResourceCredential, error) {
	var cred model.ResourceCredential
	err := r.db.WithContext(ctx).Where("secret_hash = ? AND status = ?", hash, model.CredentialStatusActive).First(&cred).Error
	if err != nil {
		return nil, err
	}
	return &cred, nil
}

func (r *CertificateRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status model.CredentialStatus, consumedAt *time.Time, revokedAt *time.Time) error {
	return r.db.WithContext(ctx).Model(&model.ResourceCredential{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      status,
		"consumed_at": consumedAt,
		"revoked_at":  revokedAt,
	}).Error
}

func (r *CertificateRepository) GetByFingerprint(ctx context.Context, fingerprint string) (*model.ResourceCertificate, error) {
	var cert model.ResourceCertificate
	err := r.db.WithContext(ctx).Where("fingerprint = ?", fingerprint).First(&cert).Error
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

func (r *CertificateRepository) Update(ctx context.Context, cert *model.ResourceCertificate) error {
	return r.db.WithContext(ctx).Save(cert).Error
}
func(r *CertificateRepository) GetActiveByResourceID(ctx context.Context,resourceID uuid.UUID) (*model.ResourceCertificate, error){
	return nil, nil
}

func ( r *CertificateRepository)    GetByCertificateID(ctx context.Context,certificateID string) (*model.ResourceCertificate, error){
	return nil, nil
}

func (CertificateRepository)  ListByResourceID(ctx context.Context,resourceID uuid.UUID) ([]model.ResourceCertificate, error){
	return nil, nil
}


func ( r *CertificateRepository)    Activate(ctx context.Context,id uuid.UUID,activatedAt time.Time) error{
	return nil
}

func (r *CertificateRepository)    Revoke(ctx context.Context,id uuid.UUID,revokedAt time.Time) error{
	return nil
}

// func (r *CertificateRepository)    GetByCertificateID(ctx context.Context,certificateID string) (*model.ResourceCertificate, error){
// 	return nil, nil
// }



// ===================================================================
// IdentityRepository Implementation
// ===================================================================

// identityRepository implements the security.IdentityRepository interface.
type IdentityRepository struct {
	db *gorm.DB
}

// NewIdentityRepository creates a new repository for resource identities.
func NewIdentityRepository(db *gorm.DB) *IdentityRepository {
	return &IdentityRepository{db: db}
}

// Compile-time check to ensure identityRepository implements the interface.

func (r *IdentityRepository) Create(ctx context.Context, identity *model.ResourceIdentity) error {
	return r.db.WithContext(ctx).Create(identity).Error
}

func (r *IdentityRepository) FindByResourceID(ctx context.Context, resourceID uuid.UUID) (*model.ResourceIdentity, error) {
	var identity model.ResourceIdentity
	err := r.db.WithContext(ctx).Where("resource_id = ?", resourceID).First(&identity).Error
	if err != nil {
		return nil, err
	}
	return &identity, nil
}

func (r *IdentityRepository) FindByHardwareID(ctx context.Context, hardwareID string) (*model.ResourceIdentity, error) {
	var identity model.ResourceIdentity
	err := r.db.WithContext(ctx).Where("hardware_id = ?", hardwareID).First(&identity).Error
	if err != nil {
		return nil, err
	}
	return &identity, nil
}

func (r *IdentityRepository) FindBySerialNumber(ctx context.Context, serialNumber string) (*model.ResourceIdentity, error) {
	var identity model.ResourceIdentity
	err := r.db.WithContext(ctx).Where("serial_number = ?", serialNumber).First(&identity).Error
	if err != nil {
		return nil, err
	}
	return &identity, nil
}

func (r *IdentityRepository)    GetByResourceID(ctx context.Context,resourceID uuid.UUID) (*model.ResourceIdentity, error){
	return nil, nil
}

func (r *IdentityRepository)    CreateOrUpdate(ctx context.Context,identity *model.ResourceIdentity) error{
	return nil
}

func (r *IdentityRepository)    HardwareIDExists(ctx context.Context,hardwareID string,excludeResourceID *uuid.UUID) (bool, error){
	return false, nil
}

func (r *IdentityRepository)    SerialNumberExists(ctx context.Context,tenantID uuid.UUID,serialNumber string,excludeResourceID *uuid.UUID) (bool, error){
	return false, nil
}

// ===================================================================
// CertificateRepository Implementation
// ===================================================================

// certificateRepository implements the security.CertificateRepository interface.
type CredentialRepository struct {
	db *gorm.DB
}

// NewCertificateRepository creates a new repository for resource certificates.
func NewCertificateRepository(db *gorm.DB) *CertificateRepository {
	return &CertificateRepository{db: db}
}
func(r *CredentialRepository)    Consume(ctx context.Context,id uuid.UUID,consumedAt time.Time) error{
	return nil
}
// Compile-time check to ensure certificateRepository implements the interface.
//var _ security.CertificateRepository = (*certificateRepository)(nil)
func (r *CredentialRepository) Create(ctx context.Context, cert *model.ResourceCredential) error {
	return r.db.WithContext(ctx).Create(cert).Error
}

func (r *CredentialRepository) ListByResourceID(ctx context.Context, resourceID uuid.UUID) ([]*model.ResourceCertificate, error) {
	var certs []*model.ResourceCertificate
	err := r.db.WithContext(ctx).Where("resource_id = ?", resourceID).Find(&certs).Error
	return certs, err
}

func (r *CredentialRepository) FindByCertificateID(ctx context.Context, certID string) (*model.ResourceCertificate, error) {
	var cert model.ResourceCertificate
	err := r.db.WithContext(ctx).Where("certificate_id = ?", certID).First(&cert).Error
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

func (r *CredentialRepository) FindByFingerprint(ctx context.Context, fingerprint string) (*model.ResourceCertificate, error) {
	var cert model.ResourceCertificate
	err := r.db.WithContext(ctx).Where("fingerprint = ?", fingerprint).First(&cert).Error
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

func (r *CredentialRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status model.CertificateStatus, activatedAt *time.Time, revokedAt *time.Time) error {
	return r.db.WithContext(ctx).Model(&model.ResourceCertificate{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       status,
		"activated_at": activatedAt,
		"revoked_at":   revokedAt,
	}).Error
}

func (r *CredentialRepository) GetForUpdate(ctx context.Context, id uuid.UUID) (*model.ResourceCredential, error) {
	var cred model.ResourceCredential
	// Use pessimistic locking to prevent race conditions on the credential record.
	err := r.db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).First(&cred, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &cred, nil
}

func (r *CredentialRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&model.ResourceCredential{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     model.CredentialStatusRevoked,
		"revoked_at": &now,
	}).Error
}

func (r *CredentialRepository) Update(ctx context.Context, cred *model.ResourceCredential) error {
	return r.db.WithContext(ctx).Save(cred).Error
}

func (r *CredentialRepository)    GetActive(ctx context.Context,resourceID uuid.UUID,credentialType model.CredentialType) (*model.ResourceCredential, error){
	return nil, nil
}

func (r *CredentialRepository)    GetByID(ctx context.Context,id uuid.UUID) (*model.ResourceCredential, error){
	return nil, nil
}