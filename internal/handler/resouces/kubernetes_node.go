package resouces

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"github.com/go-dev-frame/sponge/pkg/gin/response"
	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/xiaofan193/k8sadmin/internal/controller"
	"github.com/xiaofan193/k8sadmin/internal/ecode"
	"github.com/xiaofan193/k8sadmin/internal/types"
)

var _ NodeHandler = (*nodeHandler)(nil)

// ResoucesHandler defining the handler interface
type NodeHandler interface {
	GetNodeDetail(c *gin.Context)
	GetNodeList(c *gin.Context)
	UpdateNodeLabel(c *gin.Context)
	UpdateNodeTaint(c *gin.Context)
}

type nodeHandler struct {
}

func NewNodeHandler() NodeHandler {
	return &nodeHandler{}
}

// GetNodeDetail 获取node列表
// @Summary GetNodeDetail 获取node详情
// @Description get node详情
// @Tags node
// @Accept json
// @Produce json
// @Param nodeName query string true "node名称"
// @Success 200 {object} types.GetNodeDetailReply{}
// @Router /api/v1/k8s/node [get]
// @Security BearerAuth
func (n *nodeHandler) GetNodeDetail(c *gin.Context) {
	nodeName := c.Param("name")

	if nodeName == "" {
		response.Error(c, ecode.InvalidParams, fmt.Errorf("nodeName  不能为空"))
	}
	reqParam := &types.NodeDetailRequest{
		NodeName: nodeName,
	}
	detail, err := controller.NewNodeController().GetNodeDetail(c.Request.Context(), reqParam)

	if err != nil {
		logger.Error("GetNodeDetail error", logger.Err(err), logger.Any("parmm", reqParam), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode(), err)
		return
	}
	res := types.GetNodeDetailReply{
		Data: struct {
			Node *types.Node `json:"node"`
		}{Node: detail},
	}
	response.Success(c, res)
}

// GetNodeList 获取node列表
// @Summary GetPodList 获取node列表
// @Description get pod列表
// @Tags node
// @Accept json
// @Produce json
// @Param keyword query string true "关键字"
// @Success 200 {object} types.ListNodeReply{}
// @Router /api/v1/k8s/node [get]
// @Security BearerAuth
func (n *nodeHandler) GetNodeList(c *gin.Context) {
	keyWord := c.Query("keyword")
	reqParam := &types.NodeListRequest{}
	reqParam.KeyWord = keyWord
	list, err := controller.NewNodeController().GetNodeList(c.Request.Context(), reqParam)
	if err != nil {
		logger.Error("GetNodeList error", logger.Err(err), logger.Any("parmm", reqParam), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode(), err)
		return
	}
	res := types.ListNodeReply{
		Data: struct {
			Nodes []*types.Node `json:"nodes"`
		}{Nodes: list},
	}
	response.Success(c, res)
}

// UpdateNodeLabel 打标签
// @Summary UpdateNodeLabel 打标签
// @Description 打标签
// @Tags node
// @Accept json
// @Produce json
// @Param data body types.UpdatedLabelRequest true "请求参数"
// @Success 200 {object} types.UpdateNodeLabelReply{}
// @Router /api/v1/k8s/node [get]
// @Security BearerAuth
func (n *nodeHandler) UpdateNodeLabel(c *gin.Context) {

	reqParam := &types.UpdatedLabelRequest{}

	err := c.ShouldBind(reqParam)
	if err != nil {
		logger.Warn("UpdateNodeLabel error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	if len(reqParam.Labels) == 0 {
		response.Error(c, ecode.InvalidParams, "labels 不能为空")
		return
	}
	err = controller.NewNodeController().UpdateNodeLabel(c.Request.Context(), reqParam)
	if err != nil {
		logger.Error("UpdateNodeLabel error", logger.Err(err), logger.Any("reqParam", reqParam), logger.String("msg", err.Error()), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode(), err)
		return
	}
	response.Success(c)
}

// UpdateNodeTaint 设置node 污点
// @Summary UpdateNodeTaint 设置node 污点
// @Description 设置污点
// @Tags node
// @Accept json
// @Produce json
// @Param data body types.UpdatedTaintRequest true "请求参数"
// @Success 200 {object} types.UpdatedTaintReply{}
// @Router /api/v1/k8s/node [get]
// @Security BearerAuth
func (n *nodeHandler) UpdateNodeTaint(c *gin.Context) {

	reqParam := &types.UpdatedTaintRequest{}

	err := c.ShouldBind(reqParam)
	if err != nil {
		logger.Warn("UpdateNodeTaint error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams, err)
		return
	}

	err = controller.NewNodeController().UpdateNodeTaint(c.Request.Context(), reqParam)
	if err != nil {
		logger.Error("UpdateNodeTaint error", logger.Err(err), logger.Any("reqParam", reqParam), logger.String("msg", err.Error()), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode(), err)
		return
	}
	response.Success(c)
}
