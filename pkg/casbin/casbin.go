package casbin

import (
	"go_casbin/internal/config"
	"go_casbin/internal/logger"
	"go_casbin/pkg/path"
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter"
)

var (
	CasbinService *CasbinEnforcer
	once          sync.Once
)

// CasbinEnforcer Casbin执行器封装
type CasbinEnforcer struct {
	enforcer *casbin.Enforcer
	adapter  *gormadapter.Adapter
	model    *model.Model
}
type CasbinOptions struct {
	Driver string
	DataSource string
	ModelPath string
}

// InitCasbin 初始化casbin服务
func InitCasbin(options CasbinOptions) (initErr error) {
	// once.Do(func() {
	var enforcer *casbin.Enforcer
	var adapter *gormadapter.Adapter
	var model *model.Model
	var err error

	// 根据driver类型选择不同的初始化方式
	if options.Driver == "file" {
		// 文件模式 - 使用项目根目录的绝对路径
		modelPath, err := path.GetAbsolutePath(options.ModelPath)
		if err != nil {
			logger.ErrorWithErr("获取项目根目录失败", err)
			initErr = err
			return
		}

		adapterPath, err := path.GetAbsolutePath(options.DataSource)
		if err != nil {
			logger.ErrorWithErr("获取项目根目录失败", err)
			initErr = err
			return
		}

		enforcer, err = casbin.NewEnforcer(modelPath, adapterPath)
	} else {
		// 数据库模式
		adapter := gormadapter.NewAdapter(options.Driver, options.DataSource)
		enforcer, err = casbin.NewEnforcer(options.ModelPath, adapter)
	}

	if err != nil {
		logger.ErrorWithErr("初始化CasbinService失败", err, logger.String("modelPath", config.ViperConfig.Casbin.ModelPath))
		initErr = err
		return
	}
	CasbinService = &CasbinEnforcer{
		enforcer: enforcer,
		adapter:  adapter,
		model:    model,
	}
	logger.Info("CasbinService初始化成功")
	return
}

func GetCasbinInstance() *CasbinEnforcer {
	return CasbinService
}

// Enforce 权限判断
func (c *CasbinEnforcer) Enforce(sub, obj, act string) (bool, error) {
	ok, err := c.enforcer.Enforce(sub, obj, act)
	if err != nil {
		logger.ErrorWithErr("Casbin权限校验失败", err, logger.String("sub", sub), logger.String("obj", obj), logger.String("act", act))
	}
	return ok, err
}

// AddPolicy 添加策略
func (c *CasbinEnforcer) AddPolicy(params ...interface{}) (bool, error) {
	ok, err := c.enforcer.AddPolicy(params...)
	if err != nil {
		logger.ErrorWithErr("添加Casbin策略失败", err, logger.Field("params", params))
	}
	return ok, err
}

// RemovePolicy 删除策略
func (c *CasbinEnforcer) RemovePolicy(params ...interface{}) (bool, error) {
	ok, err := c.enforcer.RemovePolicy(params...)
	if err != nil {
		logger.ErrorWithErr("删除Casbin策略失败", err, logger.Field("params", params))
	}
	return ok, err
}

// RemoveFilteredPolicy 按条件删除策略
func (c *CasbinEnforcer) RemoveFilteredPolicy(fieldIndex int, fieldValues ...string) (bool, error) {
	ok, err := c.enforcer.RemoveFilteredPolicy(fieldIndex, fieldValues...)
	if err != nil {
		logger.ErrorWithErr("按条件删除Casbin策略失败", err, logger.Int("fieldIndex", fieldIndex), logger.Field("fieldValues", fieldValues))
	}
	return ok, err
}

// AddPolicies 批量添加策略
func (c *CasbinEnforcer) AddPolicies(rules [][]string) (bool, error) {
	ok, err := c.enforcer.AddPolicies(rules)
	if err != nil {
		logger.ErrorWithErr("批量添加Casbin策略失败", err, logger.Field("rules", rules))
	}
	return ok, err
}

// RemovePolicies 批量删除策略
func (c *CasbinEnforcer) RemovePolicies(rules [][]string) (bool, error) {
	ok, err := c.enforcer.RemovePolicies(rules)
	if err != nil {
		logger.ErrorWithErr("批量删除Casbin策略失败", err, logger.Field("rules", rules))
	}
	return ok, err
}

// GetPolicy 获取所有策略
func (c *CasbinEnforcer) GetPolicy() [][]string {
	policies, err := c.enforcer.GetPolicy()
	if err != nil {
		logger.ErrorWithErr("获取策略失败", err)
		return nil
	}
	return policies
}

// SavePolicy 持久化策略
func (c *CasbinEnforcer) SavePolicy() error {
	err := c.enforcer.SavePolicy()
	if err != nil {
		logger.ErrorWithErr("保存Casbin策略失败", err)
	}
	return err
}

// LoadPolicy 重新加载策略
func (c *CasbinEnforcer) LoadPolicy() error {
	err := c.enforcer.LoadPolicy()
	if err != nil {
		logger.ErrorWithErr("加载Casbin策略失败", err)
	}
	return err
}

// GetRolesForUser 获取用户的所有角色
func (c *CasbinEnforcer) GetRolesForUser(user string) []string {
	roles, err := c.enforcer.GetRolesForUser(user)
	if err != nil {
		logger.ErrorWithErr("获取用户角色失败", err, logger.String("user", user))
		return nil
	}
	return roles
}

// GetUsersForRole 获取角色下所有用户
func (c *CasbinEnforcer) GetUsersForRole(role string) []string {
	users, err := c.enforcer.GetUsersForRole(role)
	if err != nil {
		logger.ErrorWithErr("获取角色用户失败", err, logger.String("role", role))
		return nil
	}
	return users
}

// AddRoleForUser 给用户添加角色
func (c *CasbinEnforcer) AddRoleForUser(user, role string) (bool, error) {
	ok, err := c.enforcer.AddRoleForUser(user, role)
	if err != nil {
		logger.ErrorWithErr("添加用户角色失败", err, logger.String("user", user), logger.String("role", role))
	}
	return ok, err
}

// DeleteRoleForUser 移除用户的角色
func (c *CasbinEnforcer) DeleteRoleForUser(user, role string) (bool, error) {
	ok, err := c.enforcer.DeleteRoleForUser(user, role)
	if err != nil {
		logger.ErrorWithErr("移除用户角色失败", err, logger.String("user", user), logger.String("role", role))
	}
	return ok, err
}

// GetAllSubjects 获取所有subject
func (c *CasbinEnforcer) GetAllSubjects() []string {
	subjects, err := c.enforcer.GetAllSubjects()
	if err != nil {
		logger.ErrorWithErr("获取subject失败", err)
		return nil
	}
	return subjects
}

// GetAllObjects 获取所有object
func (c *CasbinEnforcer) GetAllObjects() []string {
	objects, err := c.enforcer.GetAllObjects()
	if err != nil {
		logger.ErrorWithErr("获取object失败", err)
		return nil
	}
	return objects
}

// GetAllActions 获取所有action
func (c *CasbinEnforcer) GetAllActions() []string {
	actions, err := c.enforcer.GetAllActions()
	if err != nil {
		logger.ErrorWithErr("获取action失败", err)
		return nil
	}
	return actions
}

// GetAllRoles 获取所有角色
func (c *CasbinEnforcer) GetAllRoles() []string {
	roles, err := c.enforcer.GetAllRoles()
	if err != nil {
		logger.ErrorWithErr("获取角色失败", err)
		return nil
	}
	return roles
}
