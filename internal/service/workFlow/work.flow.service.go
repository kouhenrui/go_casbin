package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	workflow "go_casbin/internal/model/workFlow"
	"sort"
	"strconv"
)

type WorkFlowService interface {
	// 创建工作流模板
	CreateWorkFlow(ctx context.Context, workflow *workflow.WorkFlow) error

	// 根据ID获取工作流模板，id建议使用uint，因为gorm.Model中ID是uint
	GetWorkFlow(ctx context.Context, id uint) (*workflow.WorkFlow, error)

	// 更新工作流模板
	UpdateWorkFlow(ctx context.Context, id uint, workflow *workflow.WorkFlow) error

	// 删除工作流模板（软删除）
	DeleteWorkFlow(ctx context.Context, id uint) error

	// 列表获取所有工作流模板（分页可选）
	ListWorkFlows(ctx context.Context, limit, offset int) ([]*workflow.WorkFlow, error)

	// 创建工作流实例
	CreateWorkFlowInstance(ctx context.Context, instance *workflow.WorkflowInstance) error

	// 根据ID获取工作流实例
	GetWorkFlowInstance(ctx context.Context, id string) (*workflow.WorkflowInstance, error)

	// 更新工作流实例
	UpdateWorkFlowInstance(ctx context.Context, id string, instance *workflow.WorkflowInstance) error

	// 删除工作流实例
	DeleteWorkFlowInstance(ctx context.Context, id string) error

	// 列表获取所有工作流实例（分页可选）
	ListWorkFlowInstances(ctx context.Context) ([]*workflow.WorkflowInstance, error)

	// 审批工作流实例
	ApproveWorkFlowInstance(ctx context.Context, id string, approval *workflow.WorkflowStepInstance) error

	// 获取工作流实例的当前步骤
	GetCurrentStep(ctx context.Context, id string) (*workflow.WorkflowStepInstance, error)

	// 获取工作流实例的下一步
	GetNextStep(ctx context.Context, id string) (*workflow.WorkflowStepInstance, error)
	
}

type MockWorkFlowService struct {
	workflows      map[uint]*workflow.WorkFlow
	instances      map[string]*workflow.WorkflowInstance
	nextWorkflowID uint64
}

func NewWorkFlowServiceImpl() *MockWorkFlowService {
	return &MockWorkFlowService{
		workflows:      make(map[uint]*workflow.WorkFlow),
		instances:      make(map[string]*workflow.WorkflowInstance),
		nextWorkflowID: 1,
	}
}

func (m *MockWorkFlowService) CreateWorkFlow(ctx context.Context, w *workflow.WorkFlow) error {
	w.ID = uint(m.nextWorkflowID)
	m.workflows[w.ID] = w
	m.nextWorkflowID++
	return nil
}

func (m *MockWorkFlowService) GetWorkFlow(ctx context.Context, id uint) (*workflow.WorkFlow, error) {
	if wf, ok := m.workflows[id]; ok {
		return wf, nil
	}
	return nil, fmt.Errorf("workflow not found")
}

func (m *MockWorkFlowService) UpdateWorkFlow(ctx context.Context, id uint, w *workflow.WorkFlow) error {
	if _, ok := m.workflows[id]; !ok {
		return fmt.Errorf("workflow not found")
	}
	w.ID = id
	m.workflows[id] = w
	return nil
}

func (m *MockWorkFlowService) DeleteWorkFlow(ctx context.Context, id uint) error {
	if _, ok := m.workflows[id]; !ok {
		return fmt.Errorf("workflow not found")
	}
	delete(m.workflows, id)
	return nil
}

func (m *MockWorkFlowService) ListWorkFlows(ctx context.Context, limit, offset int) ([]*workflow.WorkFlow, error) {
	values := make([]*workflow.WorkFlow, 0, len(m.workflows))
	for _, v := range m.workflows {
		values = append(values, v)
	}
	sort.Slice(values, func(i, j int) bool { return values[i].ID < values[j].ID })

	end := offset + limit
	if end > len(values) {
		end = len(values)
	}
	if offset > len(values) {
		return []*workflow.WorkFlow{}, nil
	}
	return values[offset:end], nil
}

func (m *MockWorkFlowService) CreateWorkFlowInstance(ctx context.Context, inst *workflow.WorkflowInstance) error {
	m.instances[strconv.FormatInt(inst.ID, 10)] = inst	
	return nil
}

func (m *MockWorkFlowService) GetWorkFlowInstance(ctx context.Context, id string) (*workflow.WorkflowInstance, error) {
	if inst, ok := m.instances[id]; ok {
		return inst, nil
	}
	return nil, fmt.Errorf("instance not found")
}

func (m *MockWorkFlowService) UpdateWorkFlowInstance(ctx context.Context, id string, inst *workflow.WorkflowInstance) error {
	if _, ok := m.instances[id]; !ok {
		return fmt.Errorf("instance not found")
	}
	inst.ID,_ = strconv.ParseInt(id, 10, 64)
	m.instances[id] = inst
	return nil
}

func (m *MockWorkFlowService) DeleteWorkFlowInstance(ctx context.Context, id string) error {
	if _, ok := m.instances[id]; !ok {
		return fmt.Errorf("instance not found")
	}
	delete(m.instances, id)
	return nil
}

func (m *MockWorkFlowService) ListWorkFlowInstances(ctx context.Context) ([]*workflow.WorkflowInstance, error) {
	values := make([]*workflow.WorkflowInstance, 0, len(m.instances))
	for _, v := range m.instances {
		values = append(values, v)
	}
	sort.Slice(values, func(i, j int) bool { return values[i].ID < values[j].ID })

	// end := offset + limit
	// if end > len(values) {
	// 	end = len(values)
	// }
	// if offset > len(values) {
	// 	return []*workflow.WorkflowInstance{}, nil
	// }
	return values, nil
}

func (m *MockWorkFlowService) ApproveWorkFlowInstance(ctx context.Context, id string, approval *workflow.WorkflowStepInstance) error {
	if _, ok := m.instances[id]; !ok {
		return fmt.Errorf("instance not found")
	}
	
	m.instances[id].Steps = append(m.instances[id].Steps, *approval)
	total:=len(approval.Approvals)
	rejectCount:=0
	for _, v := range approval.Approvals {
		if !v {
			rejectCount++
		}
	}
	if rejectCount > total/2{
		m.instances[id].Status=3
		m.instances[id].Steps[m.instances[id].Current].Finished=true
		m.instances[id].Done=true
	}else{
		m.instances[id].Status=2
	}
	return nil
}

func (m *MockWorkFlowService) GetCurrentStep(ctx context.Context, id string) (*workflow.WorkflowStepInstance, error) {
	if _, ok := m.instances[id]; !ok {
		return nil, fmt.Errorf("instance not found")
	}
	jsonBytes, _ := json.Marshal(m.instances[id].Steps[m.instances[id].Current])
	var step workflow.WorkflowStepInstance
	json.Unmarshal(jsonBytes, &step)
	return &step, nil
}

func (m *MockWorkFlowService) GetNextStep(ctx context.Context, id string) (*workflow.WorkflowStepInstance, error) {
	if _, ok := m.instances[id]; !ok {
		return nil, fmt.Errorf("instance not found")
	}
	jsonBytes, _ := json.Marshal(m.instances[id].Steps[m.instances[id].Current+1])
	var step workflow.WorkflowStepInstance
	json.Unmarshal(jsonBytes, &step)
	return &step, nil
}