package workflow

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// 工作流步骤 模版
type WorkFlowStep struct {
	gorm.Model
	Name      string         `json:"name"`        //步骤名称
	Approvers datatypes.JSON `json:"description"` //审批人
	Status    int            `json:"status"`      //0:草稿 1:发布 2:禁用
}

// 工作流模版
type WorkFlow struct {
	gorm.Model
	Name        string         `json:"name"`        //工作流名称
	Version     string         `json:"version"`     //版本号
	Description string         `json:"description"` //描述
	Steps       []WorkFlowStep `json:"steps"`       //步骤
	Status      int            `json:"status"`      //0:草稿 1:发布 2:禁用
}

// 流程实例
type WorkflowInstance struct {
	ID         int64    `gorm:"primaryKey;autoIncrement" json:"id"` // 使用 bigint
	WorkflowID uint     `json:"workflow_id"`
	Workflow   WorkFlow `gorm:"foreignKey:WorkflowID"` // 关联定义

	Current int            `json:"current"` // 当前步骤索引
	Steps   []WorkflowStepInstance `json:"steps"`   // 存储 WorkflowStepInstance 列表
	Done    bool           `json:"done"`    // 是否完成
	Status  int            `json:"status"`  // 状态 0:Pending, 1:Approved, 2:Rejected
}

// 步骤实例结构体（不直接建表，序列化进 JSON）
type WorkflowStepInstance struct {
	StepID    string            `json:"step_id"`   //步骤ID
	Approvals map[string]bool   `json:"approvals"` // 审批记录：userID -> 是否通过
	Reason    map[string]string `json:"reason"`    // 审批意见
	Finished  bool              `json:"finished"`  //是否完成
}
