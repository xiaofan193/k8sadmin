package handler

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/go-dev-frame/sponge/pkg/copier"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"github.com/go-dev-frame/sponge/pkg/gin/response"
	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/utils"

	"github.com/xiaofan193/k8sadmin/internal/cache"
	"github.com/xiaofan193/k8sadmin/internal/dao"
	"github.com/xiaofan193/k8sadmin/internal/database"
	"github.com/xiaofan193/k8sadmin/internal/ecode"
	"github.com/xiaofan193/k8sadmin/internal/model"
	"github.com/xiaofan193/k8sadmin/internal/types"
)

var _ UserHandler = (*userHandler)(nil)

// UserHandler defining the handler interface
type UserHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type userHandler struct {
	iDao dao.UserDao
}

// NewUserHandler creating the handler interface
func NewUserHandler() UserHandler {
	return &userHandler{
		iDao: dao.NewUserDao(
			database.GetDB(), // db driver is mysql
			cache.NewUserCache(database.GetCacheType()),
		),
	}
}

// Create a new user
// @Summary Create a new user
// @Description Creates a new user entity using the provided data in the request body.
// @Tags user
// @Accept json
// @Produce json
// @Param data body types.CreateUserRequest true "user information"
// @Success 200 {object} types.CreateUserReply{}
// @Router /api/v1/user [post]
// @Security BearerAuth
func (h *userHandler) Create(c *gin.Context) {
	form := &types.CreateUserRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	user := &model.User{}
	err = copier.Copy(user, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateUser)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, user)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": user.ID})
}

// DeleteByID delete a user by id
// @Summary Delete a user by id
// @Description Deletes a existing user identified by the given id in the path.
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteUserByIDReply{}
// @Router /api/v1/user/{id} [delete]
// @Security BearerAuth
func (h *userHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getUserIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	err := h.iDao.DeleteByID(ctx, id)
	if err != nil {
		logger.Error("DeleteByID error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// UpdateByID update a user by id
// @Summary Update a user by id
// @Description Updates the specified user by given id in the path, support partial update.
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateUserByIDRequest true "user information"
// @Success 200 {object} types.UpdateUserByIDReply{}
// @Router /api/v1/user/{id} [put]
// @Security BearerAuth
func (h *userHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getUserIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateUserByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	user := &model.User{}
	err = copier.Copy(user, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDUser)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, user)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a user by id
// @Summary Get a user by id
// @Description Gets detailed information of a user specified by the given id in the path.
// @Tags user
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetUserByIDReply{}
// @Router /api/v1/user/{id} [get]
// @Security BearerAuth
func (h *userHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getUserIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	user, err := h.iDao.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			logger.Warn("GetByID not found", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.NotFound)
		} else {
			logger.Error("GetByID error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
			response.Output(c, ecode.InternalServerError.ToHTTPCode())
		}
		return
	}

	data := &types.UserObjDetail{}
	err = copier.Copy(data, user)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDUser)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"user": data})
}

// List get a paginated list of users by custom conditions
// @Summary Get a paginated list of users by custom conditions
// @Description Returns a paginated list of user based on query filters, including page number and size.
// @Tags user
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListUsersReply{}
// @Router /api/v1/user/list [post]
// @Security BearerAuth
func (h *userHandler) List(c *gin.Context) {
	form := &types.ListUsersRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	users, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertUsers(users)
	if err != nil {
		response.Error(c, ecode.ErrListUser)
		return
	}

	response.Success(c, gin.H{
		"users": data,
		"total": total,
	})
}

func getUserIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertUser(user *model.User) (*types.UserObjDetail, error) {
	data := &types.UserObjDetail{}
	err := copier.Copy(data, user)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertUsers(fromValues []*model.User) ([]*types.UserObjDetail, error) {
	toValues := []*types.UserObjDetail{}
	for _, v := range fromValues {
		data, err := convertUser(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
