package postgres

import (
	"context"

	"github.com/OpenIndustrial/cloud/internal/identity"
	"github.com/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)



// groupRepository implements the identity.GroupRepository interface using GORM.
type groupRepository struct {
	db *gorm.DB
}

// NewGroupRepository creates a new GORM-based GroupRepository.
// It now correctly returns the interface type.
func NewGroupRepository(db *gorm.DB) identity.GroupRepository {
	return &groupRepository{db: db}
}

// CreateGroup creates a new group in the database.
func (r *groupRepository) CreateGroup(ctx context.Context, group *model.Group) error {
	return r.db.WithContext(ctx).Create(group).Error
}

// GetGroupByID retrieves a group by its UUID and tenant ID.
func (r *groupRepository) GetGroupByID(ctx context.Context, tenantID, groupID uuid.UUID) (*model.Group, error) {
	var group model.Group
	err := r.db.WithContext(ctx).
		Where("uuid = ? AND tenant_id = ?", groupID, tenantID).
		First(&group).Error
	if err != nil {
		return nil, err // GORM will return gorm.ErrRecordNotFound if not found.
	}
	return &group, nil
}

// AddUserToGroup adds a user to a group.
// This demonstrates handling a many-to-many relationship with GORM.
func (r *groupRepository) AddUserToGroup(ctx context.Context, tenantID, userID, groupID uuid.UUID) error {
	// First, find the target group and user.
	var group model.Group
	if err := r.db.WithContext(ctx).Where("uuid = ? AND tenant_id = ?", groupID, tenantID).First(&group).Error; err != nil {
		return err // Group not found
	}

	var user model.User
	if err := r.db.WithContext(ctx).Where("uuid = ? AND tenant_id = ?", userID, tenantID).First(&user).Error; err != nil {
		return err // User not found
	}

	// Use GORM's Association to append the user to the group's Users.
	// GORM will handle the `user_groups` join table insertion.
	return r.db.WithContext(ctx).Model(&group).Association("Users").Append(&user)
}

// RemoveUserFromGroup removes a user from a group.
func (r *groupRepository) RemoveUserFromGroup(ctx context.Context, tenantID, userID, groupID uuid.UUID) error {
	var group model.Group
	if err := r.db.WithContext(ctx).Where("uuid = ? AND tenant_id = ?", groupID, tenantID).First(&group).Error; err != nil {
		return err
	}

	var user model.User
	if err := r.db.WithContext(ctx).Where("uuid = ? AND tenant_id = ?", userID, tenantID).First(&user).Error; err != nil {
		return err
	}

	// Use GORM's Association to delete the relationship.
	return r.db.WithContext(ctx).Model(&group).Association("Users").Delete(&user)
}

// ListGroupsByUserID retrieves all groups a specific user is a member of.
func (r *groupRepository) ListGroupsByUserID(ctx context.Context, tenantID, userID uuid.UUID) ([]*model.Group, error) {
	var user model.User
	// Find the user and preload their associated groups.
	err := r.db.WithContext(ctx).
		Preload("Groups", "tenant_id = ?", tenantID). // Preload groups, ensuring they belong to the correct tenant.
		Where("uuid = ? AND tenant_id = ?", userID, tenantID).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	// The preloaded groups are now available in the user object.
	return user.Groups, nil
}