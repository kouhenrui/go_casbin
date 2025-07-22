package service

import (
	"context"
	"go_casbin/internal/model"
	"go_casbin/internal/repository/account"
	"go_casbin/pkg/database"

	"gorm.io/gorm"
)

type AccountService interface {
	GetAccountByID(ctx context.Context, id uint) (*model.Account, error)
	AddAccount(ctx context.Context, account *model.Account) error
	UpdateAccount(ctx context.Context, account *model.Account) error
	RemoveAccount(ctx context.Context, id uint) error
	ReplaceAccountRoles(ctx context.Context, account *model.Account, roles []model.Role) error
	// 新增事务方法
	CreateAccountWithRoles(ctx context.Context, account *model.Account, roles []model.Role) error
	UpdateAccountWithRoles(ctx context.Context, account *model.Account, roles []model.Role) error
	DeleteAccountWithCleanup(ctx context.Context, id uint) error
}

type AccountServiceImpl struct {
	accountRepository account.AccountRepository
}

func NewAccountService() *AccountServiceImpl {
	return &AccountServiceImpl{accountRepository: account.NewAccountRepository()}
}

func (s *AccountServiceImpl) GetAccountByID(ctx context.Context, id uint) (*model.Account, error) {
	return s.accountRepository.FindByID(ctx, id)
}

func (s *AccountServiceImpl) AddAccount(ctx context.Context, account *model.Account) error {
	return s.accountRepository.Create(ctx, account)
}

func (s *AccountServiceImpl) UpdateAccount(ctx context.Context, account *model.Account) error {
	return s.accountRepository.Update(ctx, account)
}

func (s *AccountServiceImpl) RemoveAccount(ctx context.Context, id uint) error {
	return s.accountRepository.Delete(ctx, id)
}

func (s *AccountServiceImpl) ReplaceAccountRoles(ctx context.Context, account *model.Account, roles []model.Role) error {
	return s.accountRepository.ReplaceAccountRoles(ctx, account, roles)
}

// CreateAccountWithRoles 创建账户并分配角色（事务）
func (s *AccountServiceImpl) CreateAccountWithRoles(ctx context.Context, account *model.Account, roles []model.Role) error {
	return database.GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 创建账户
		if err := tx.Create(account).Error; err != nil {
			return err
		}
		
		// 2. 分配角色
		if len(roles) > 0 {
			if err := tx.Model(account).Association("Roles").Replace(roles); err != nil {
				return err
			}
		}
		
		return nil
	})
}
// 事务处理操作中,遇到错误会自动回滚
// UpdateAccountWithRoles 更新账户和角色（事务）
func (s *AccountServiceImpl) UpdateAccountWithRoles(ctx context.Context, account *model.Account, roles []model.Role) error {
	return database.GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 更新账户基本信息
		if err := tx.Save(account).Error; err != nil {
			return err
		}
		
		// 2. 更新角色关联
		if err := tx.Model(account).Association("Roles").Replace(roles); err != nil {
			return err
		}
		
		return nil
	})
}

// DeleteAccountWithCleanup 删除账户并清理相关数据（事务）
func (s *AccountServiceImpl) DeleteAccountWithCleanup(ctx context.Context, id uint) error {
	return database.GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 先查找账户
		var account model.Account
		if err := tx.Preload("Roles").First(&account, id).Error; err != nil {
			return err
		}
		
		// 2. 清理角色关联
		if err := tx.Model(&account).Association("Roles").Clear(); err != nil {
			return err
		}
		
		// 3. 删除账户
		if err := tx.Delete(&account).Error; err != nil {
			return err
		}
		
		return nil
	})
}