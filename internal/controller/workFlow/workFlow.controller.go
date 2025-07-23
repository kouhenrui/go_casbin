package workFlow

import (
	"go_casbin/internal/middleware/response"
	workInstance "go_casbin/internal/model/workFlow"
	workService "go_casbin/internal/service/workFlow"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type WorkFlowController interface {
	CreateWorkFlow(c *gin.Context)
	GetWorkFlow(c *gin.Context)
	GetWorkFlowList(c *gin.Context)
	UpdateWorkFlow(c *gin.Context)

	CreateWorkFlowInstance(c *gin.Context)
	GetWorkFlowInstance(c *gin.Context)
	GetWorkFlowInstanceList(c *gin.Context)
	// UpdateWorkFlowInstance(c *gin.Context)
	DeleteWorkFlowInstance(c *gin.Context)
	ApproveWorkFlowInstance(c *gin.Context)
	
}

type WorkFlowControllerImpl struct {
	workFlowService workService.WorkFlowService
}

func NewWorkFlowController() WorkFlowController {
	return &WorkFlowControllerImpl{
		workFlowService: workService.NewWorkFlowServiceImpl(),
	}
}

//创建工作流模版
func (w *WorkFlowControllerImpl) CreateWorkFlow(c *gin.Context) {
	err:=w.workFlowService.CreateWorkFlow(c.Request.Context(), &workInstance.WorkFlow{
		Name: "测试模版",
		Version: "1.0.0",
		Description: "测试模版描述",
		Steps: []workInstance.WorkFlowStep{
			{
				Name: "测试步骤1",
				Approvers: datatypes.JSON([]byte(`[{"name":"张三","id":"1"},{"name":"李四","id":"2"}]`)),
				Status: 0,
			},
			{
				Name: "测试步骤2",
				Approvers: datatypes.JSON([]byte(`[{"name":"张三","id":"1"},{"name":"李四","id":"2"}]`)),
				Status: 0,
			},
			{
				Name: "测试步骤3",
				Approvers: datatypes.JSON([]byte(`[{"name":"张三","id":"1"},{"name":"李四","id":"2"}]`)),
				Status: 0,
			},
		},
		Status: 0,
	})
	if err!=nil{
		response.LogicError(c, err.Error())
		return
	}
	response.Success(c, "创建成功")
}

//获取工作流模版
func (w *WorkFlowControllerImpl) GetWorkFlow(c *gin.Context) {
	workFlow,err:=w.workFlowService.GetWorkFlow(c.Request.Context(), 1)
	if err!=nil{
		response.LogicError(c, err.Error())
		return
	}
	response.Success(c, workFlow)
}

//获取工作流模版列表
func (w *WorkFlowControllerImpl) GetWorkFlowList(c *gin.Context) {
	workFlows,err:=w.workFlowService.ListWorkFlows(c.Request.Context(), 10, 0)
	if err!=nil{
		response.LogicError(c, err.Error())
		return
	}
	response.Success(c, workFlows)
}

//更新工作流模版
func (w *WorkFlowControllerImpl) UpdateWorkFlow(c *gin.Context) {

}

//删除工作流模版
func (w *WorkFlowControllerImpl) DeleteWorkFlow(c *gin.Context) {

}

//创建工作流实例
func (w *WorkFlowControllerImpl) CreateWorkFlowInstance(c *gin.Context) {
	err:=w.workFlowService.CreateWorkFlowInstance(c.Request.Context(), &workInstance.WorkflowInstance{
		WorkflowID: 1,
		Current: 0,
		Steps: []workInstance.WorkflowStepInstance{},
		Done: false,
	})
	if err!=nil{
		response.LogicError(c, err.Error())
		return
	}
	response.Success(c, "创建成功")
}

//获取工作流实例
func (w *WorkFlowControllerImpl) GetWorkFlowInstance(c *gin.Context) {
	workInstance,err:=w.workFlowService.GetWorkFlowInstance(c.Request.Context(), "0")
	if err!=nil{
		response.LogicError(c, err.Error())
		return
	}
	response.Success(c, workInstance)
}

//更新工作流实例
func (w *WorkFlowControllerImpl) UpdateWorkFlowInstance(c *gin.Context) {
	w.workFlowService.UpdateWorkFlowInstance(c.Request.Context(), "1", &workInstance.WorkflowInstance{
		WorkflowID: 1,
		Current: 0,
		Steps: []workInstance.WorkflowStepInstance{
			{
				StepID: "1",
				Approvals: map[string]bool{"1": false, "2": false},
				Reason: map[string]string{"1": "不合适", "2": "不合适"},
				Finished: false,
			},
		},
		Done: false,
	})
}

//删除工作流实例
func (w *WorkFlowControllerImpl) DeleteWorkFlowInstance(c *gin.Context) {
	err:=w.workFlowService.DeleteWorkFlowInstance(c.Request.Context(), "1")	
	if err!=nil{
		response.LogicError(c, err.Error())
		return
	}
	response.Success(c, "删除成功")
}

//获取工作流实例列表
func (w *WorkFlowControllerImpl) GetWorkFlowInstanceList(c *gin.Context) {
	workInstances,err:=w.workFlowService.ListWorkFlowInstances(c.Request.Context())
	if err!=nil{
		response.LogicError(c, err.Error())
		return
	}
	response.Success(c, workInstances)
}

//审批工作流实例
func (w *WorkFlowControllerImpl) ApproveWorkFlowInstance(c *gin.Context) {
	err:=w.workFlowService.ApproveWorkFlowInstance(c.Request.Context(), "0", &workInstance.WorkflowStepInstance{
		StepID: "1",
		Approvals: map[string]bool{"1": false, "2": false},
		Reason: map[string]string{"1": "不合适", "2": "不合适"},
		Finished: true,
	})
	if err!=nil{
		response.LogicError(c, err.Error())
		return
	}
	response.Success(c, "审批成功")
}

//获取工作流实例的当前步骤
func (w *WorkFlowControllerImpl) GetCurrentStep(c *gin.Context) {
	step,err:=w.workFlowService.GetCurrentStep(c.Request.Context(), "1")
	if err!=nil{
		response.LogicError(c, err.Error())
		return
	}
	response.Success(c, step)
}

//获取工作流实例的下一步
func (w *WorkFlowControllerImpl) GetNextStep(c *gin.Context) {
	step,err:=w.workFlowService.GetNextStep(c.Request.Context(), "1")
	if err!=nil{
		response.LogicError(c, err.Error())
		return
	}
	response.Success(c, step)
}