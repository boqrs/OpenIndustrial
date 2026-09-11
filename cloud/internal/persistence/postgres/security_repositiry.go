package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/kernel/security"
	"github.com/boqrs/nexus/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ============================================================================
// CredentialRepository
// ============================================================================

// credentialRepository implements security.CredentialRepository.
type credentialRepository struct {
	db *database.DBProvider
}

// NewCredentialRepository creates a repository for resource credentials.
func NewCredentialRepository(db *database.DBProvider) *credentialRepository {
	return &credentialRepository{
		db: db,
	}
}

var _ security.CredentialRepository = (*credentialRepository)(nil)

// Create creates a resource credential.
func (r *credentialRepository) Create(
	ctx context.Context,
	credential *model.ResourceCredential,
) error {
	if credential == nil {
		return errors.New("credential is nil")
	}

	return dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Create(credential).
		Error
}

// GetActive returns an active credential for a resource.
func (r *credentialRepository) GetActive(
	ctx context.Context,
	resourceID uint,
	credentialType model.CredentialType,
) (*model.ResourceCredential, error) {
	var credential model.ResourceCredential

	err := dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Where(
			"resource_id = ? AND type = ? AND status = ?",
			resourceID,
			credentialType,
			model.CredentialStatusActive,
		).
		First(&credential).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, security.ErrCredentialNotFound
		}

		return nil, err
	}

	return &credential, nil
}

// GetByID returns a credential by primary key.
func (r *credentialRepository) GetByID(
	ctx context.Context,
	id uint,
) (*model.ResourceCredential, error) {
	var credential model.ResourceCredential

	err := dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Where("id = ?", id).
		First(&credential).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, security.ErrCredentialNotFound
		}

		return nil, err
	}

	return &credential, nil
}

// GetForUpdate returns a credential with a pessimistic row lock.
//
// This must be called inside UnitOfWork.Execute() when the caller
// needs the SELECT ... FOR UPDATE semantics.
func (r *credentialRepository) GetForUpdate(
	ctx context.Context,
	id uint,
) (*model.ResourceCredential, error) {
	var credential model.ResourceCredential

	err := dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Clauses(clause.Locking{
			Strength: "UPDATE",
		}).
		Where("id = ?", id).
		First(&credential).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, security.ErrCredentialNotFound
		}

		return nil, err
	}

	return &credential, nil
}

// Consume marks a credential as consumed.
func (r *credentialRepository) Consume(
	ctx context.Context,
	id uint,
	consumedAt time.Time,
) error {
	return dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Model(&model.ResourceCredential{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      model.CredentialStatusConsumed,
			"consumed_at": consumedAt,
		}).
		Error
}

// Revoke revokes a credential.
func (r *credentialRepository) Revoke(
	ctx context.Context,
	id uint,
) error {
	now := time.Now().UTC()

	return dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Model(&model.ResourceCredential{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     model.CredentialStatusRevoked,
			"revoked_at": now,
		}).
		Error
}

// Update updates a resource credential.
func (r *credentialRepository) Update(
	ctx context.Context,
	credential *model.ResourceCredential,
) error {
	if credential == nil {
		return errors.New("credential is nil")
	}

	return dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Save(credential).
		Error
}

// ============================================================================
// IdentityRepository
// ============================================================================

// identityRepository implements security.IdentityRepository.
type identityRepository struct {
	db *database.DBProvider
}

// NewIdentityRepository creates a repository for resource identities.
func NewIdentityRepository(db *database.DBProvider) *identityRepository {
	return &identityRepository{
		db: db,
	}
}

var _ security.IdentityRepository = (*identityRepository)(nil)

// GetByResourceID returns the canonical identity of a resource.
func (r *identityRepository) GetByResourceID(
	ctx context.Context,
	resourceID uint,
) (*model.ResourceIdentity, error) {
	var identity model.ResourceIdentity

	err := dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Where("resource_id = ?", resourceID).
		First(&identity).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, security.ErrIdentityNotFound
		}

		return nil, err
	}

	return &identity, nil
}

// Create creates a resource identity.
func (r *identityRepository) Create(
	ctx context.Context,
	identity *model.ResourceIdentity,
) error {
	if identity == nil {
		return errors.New("identity is nil")
	}

	return dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Create(identity).
		Error
}

// CreateOrUpdate creates the identity when it does not exist,
// otherwise updates the existing identity.
func (r *identityRepository) CreateOrUpdate(
	ctx context.Context,
	identity *model.ResourceIdentity,
) error {
	if identity == nil {
		return errors.New("identity is nil")
	}

	db := dbFromContext(ctx, r.db.Get()).WithContext(ctx)

	var existing model.ResourceIdentity

	err := db.
		Where("resource_id = ?", identity.ResourceID).
		First(&existing).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return db.Create(identity).Error
		}

		return err
	}

	identity.ID = existing.ID

	return db.
		Model(&existing).
		Updates(map[string]interface{}{
			"hardware_id":   identity.HardwareID,
			"serial_number": identity.SerialNumber,
			"updated_at":    identity.UpdatedAt,
		}).
		Error
}

// HardwareIDExists checks whether a hardware ID already belongs
// to another resource.
func (r *identityRepository) HardwareIDExists(
	ctx context.Context,
	hardwareID string,
	excludeResourceID *uint,
) (bool, error) {
	db := dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Model(&model.ResourceIdentity{}).
		Where("hardware_id = ?", hardwareID)

	if excludeResourceID != nil {
		db = db.Where("resource_id <> ?", *excludeResourceID)
	}

	var count int64

	if err := db.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// SerialNumberExists checks whether a serial number already belongs
// to another resource.
//
// ResourceIdentity currently has no tenant_id field, so tenantID is
// intentionally not used in the SQL condition.
func (r *identityRepository) SerialNumberExists(
	ctx context.Context,
	tenantID uuid.UUID,
	serialNumber string,
	excludeResourceID *uint,
) (bool, error) {
	_ = tenantID

	db := dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Model(&model.ResourceIdentity{}).
		Where("serial_number = ?", serialNumber)

	if excludeResourceID != nil {
		db = db.Where("resource_id <> ?", *excludeResourceID)
	}

	var count int64

	if err := db.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// ============================================================================
// CertificateRepository
// ============================================================================

// certificateRepository implements security.CertificateRepository.
type certificateRepository struct {
	db *database.DBProvider
}

// NewCertificateRepository creates a repository for resource certificates.
func NewCertificateRepository(db *database.DBProvider) *certificateRepository {
	return &certificateRepository{
		db: db,
	}
}

var _ security.CertificateRepository = (*certificateRepository)(nil)

// Create creates a resource certificate.
func (r *certificateRepository) Create(
	ctx context.Context,
	certificate *model.ResourceCertificate,
) error {
	if certificate == nil {
		return errors.New("certificate is nil")
	}

	return dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Create(certificate).
		Error
}

// GetActiveByResourceID returns the active certificate of a resource.
func (r *certificateRepository) GetActiveByResourceID(
	ctx context.Context,
	resourceID uint,
) (*model.ResourceCertificate, error) {
	var certificate model.ResourceCertificate

	err := dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Where(
			"resource_id = ? AND status = ?",
			resourceID,
			model.CertificateActive,
		).
		Order("id DESC").
		First(&certificate).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, security.ErrCertificateNotFound
		}

		return nil, err
	}

	return &certificate, nil
}

// GetByCertificateID returns a certificate by external CA certificate ID.
func (r *certificateRepository) GetByCertificateID(
	ctx context.Context,
	certificateID uint,
) (*model.ResourceCertificate, error) {
	var certificate model.ResourceCertificate

	err := dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Where("certificate_id = ?", certificateID).
		First(&certificate).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, security.ErrCertificateNotFound
		}

		return nil, err
	}

	return &certificate, nil
}

// ListByResourceID returns all certificates belonging to a resource.
func (r *certificateRepository) ListByResourceID(
	ctx context.Context,
	resourceID uint,
) ([]model.ResourceCertificate, error) {
	var certificates []model.ResourceCertificate

	err := dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Where("resource_id = ?", resourceID).
		Order("id DESC").
		Find(&certificates).
		Error

	if err != nil {
		return nil, err
	}

	return certificates, nil
}

// Activate activates a certificate.
func (r *certificateRepository) Activate(
	ctx context.Context,
	id uint,
	activatedAt time.Time,
) error {
	return dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Model(&model.ResourceCertificate{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       model.CertificateActive,
			"activated_at": activatedAt,
		}).
		Error
}

// Revoke revokes a certificate.
func (r *certificateRepository) Revoke(
	ctx context.Context,
	id uint,
	revokedAt time.Time,
) error {
	return dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Model(&model.ResourceCertificate{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     model.CertificateRevoked,
			"revoked_at": revokedAt,
		}).
		Error
}

// GetByFingerprint returns a certificate by fingerprint.
func (r *certificateRepository) GetByFingerprint(
	ctx context.Context,
	fingerprint string,
) (*model.ResourceCertificate, error) {
	var certificate model.ResourceCertificate

	err := dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Where("fingerprint = ?", fingerprint).
		First(&certificate).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, security.ErrCertificateNotFound
		}

		return nil, err
	}

	return &certificate, nil
}

// Update updates a resource certificate.
func (r *certificateRepository) Update(
	ctx context.Context,
	certificate *model.ResourceCertificate,
) error {
	if certificate == nil {
		return errors.New("certificate is nil")
	}

	return dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Save(certificate).
		Error
}
