package account

import (
	"context"
	"errors"
	"go_casbin/internal/model"
	"go_casbin/pkg/database"

	"gorm.io/gorm"
)

// AccountRepository 账户仓储接口
type AccountRepository interface {
	Create(ctx context.Context, account *model.Account) error
	FindByID(ctx context.Context, id uint) (*model.Account, error)
	FindByName(ctx context.Context, name string) (*model.Account, error)
	FindByEmail(ctx context.Context, email string) (*model.Account, error)
	FindByPhone(ctx context.Context, phone string) (*model.Account, error)
	Update(ctx context.Context, account *model.Account) error
	Delete(ctx context.Context, id uint) error
	FindByCondition(ctx context.Context, condition string, args ...interface{}) ([]*model.Account, error)
	Count(ctx context.Context, condition string, args ...interface{}) (int64, error)
	FindDeletedAccounts(ctx context.Context) ([]*model.Account, error)
	CountDeletedAccounts(ctx context.Context) (int64, error)
	RestoreAccount(ctx context.Context, id uint) error
	ReplaceAccountRoles(ctx context.Context, account *model.Account, roles []model.Role) error
	AppendAccountRoles(ctx context.Context, account *model.Account, roles []model.Role) error
	FindByRoleID(ctx context.Context, roleID uint) ([]*model.Account, error)
}

// AccountRepositoryImpl 账户仓储实现
type AccountRepositoryImpl struct {
	db *gorm.DB
}

// NewAccountRepository 创建账户仓储
func NewAccountRepository() AccountRepository {
	return &AccountRepositoryImpl{db: database.GetDB()}
}

// Create 创建账户
func (r *AccountRepositoryImpl) Create(ctx context.Context, account *model.Account) error {
	return r.db.WithContext(ctx).Create(account).Error
}

// FindByID 根据ID查找账户
func (r *AccountRepositoryImpl) FindByID(ctx context.Context, id uint) (*model.Account, error) {
	var account model.Account
	err := r.db.WithContext(ctx).Preload("Roles").First(&account, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

// FindByName 根据用户名查找账户
func (r *AccountRepositoryImpl) FindByName(ctx context.Context, name string) (*model.Account, error) {
	var account model.Account
	err := r.db.WithContext(ctx).Preload("Roles").Where("name = ?", name).First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

// FindByEmail 根据邮箱查找账户
func (r *AccountRepositoryImpl) FindByEmail(ctx context.Context, email string) (*model.Account, error) {
	var account model.Account
	err := r.db.WithContext(ctx).Preload("Roles").Where("email = ?", email).First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

// FindByPhone 根据手机号查找账户
func (r *AccountRepositoryImpl) FindByPhone(ctx context.Context, phone string) (*model.Account, error) {
	var account model.Account
	err := r.db.WithContext(ctx).Preload("Roles").Where("phone = ?", phone).First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

//根据角色ID查找账户
func (r *AccountRepositoryImpl) FindByRoleID(ctx context.Context, roleID uint) ([]*model.Account, error) {
	var accounts []*model.Account
	err := r.db.WithContext(ctx).Preload("Roles").Where("roles.id = ?", roleID).Find(&accounts).Error
	return accounts, err
}

// Update 更新账户
func (r *AccountRepositoryImpl) Update(ctx context.Context, account *model.Account) error {
	return r.db.WithContext(ctx).Save(account).Error
}

//替换账户的角色
func (r *AccountRepositoryImpl) ReplaceAccountRoles(ctx context.Context, account *model.Account, roles []model.Role) error {
	return r.db.WithContext(ctx).Model(account).Association("Roles").Replace(roles)
}

//追加账户的角色
func (r *AccountRepositoryImpl) AppendAccountRoles(ctx context.Context, account *model.Account, roles []model.Role) error {
	return r.db.WithContext(ctx).Model(account).Association("Roles").Append(roles)
}

// Delete 删除账户
func (r *AccountRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Account{}, id).Error
}

// FindByCondition 根据条件查找账户
func (r *AccountRepositoryImpl) FindByCondition(ctx context.Context, condition string, args ...interface{}) ([]*model.Account, error) {
	var accounts []*model.Account
	err := r.db.WithContext(ctx).Preload("Roles").Where(condition, args...).Find(&accounts).Error
	return accounts, err
}

// Count 统计账户数量
func (r *AccountRepositoryImpl) Count(ctx context.Context, condition string, args ...interface{}) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Account{}).Preload("Roles").Where(condition, args...).Count(&count).Error
	return count, err
} 

//查询软删除的账户
func (r *AccountRepositoryImpl) FindDeletedAccounts(ctx context.Context) ([]*model.Account, error) {
	var accounts []*model.Account
	err := r.db.WithContext(ctx).Preload("Roles").Unscoped().Where("deleted_at IS NOT NULL").Find(&accounts).Error
	return accounts, err
}

//查询软删除的账户数量
func (r *AccountRepositoryImpl) CountDeletedAccounts(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Account{}).Unscoped().Where("deleted_at IS NOT NULL").Count(&count).Error
	return count, err
}

//恢复账户
func (r *AccountRepositoryImpl) RestoreAccount(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&model.Account{}).Unscoped().Where("id = ?", id).Update("deleted_at", nil).Error
}