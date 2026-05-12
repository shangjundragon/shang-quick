package controller

import (
	"backend/model"
	"backend/pkg/constants"
	"backend/pkg/global_vars"
	"backend/pkg/password"
	"backend/pkg/req_util"
	"backend/pkg/res_util"
	"backend/pkg/utils"
	"backend/service"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shangjundragon/dbw"
	"go.uber.org/zap"
)

// UserList 分页查询用户列表，同时批量查询部门名称和角色名称
func UserList(c *gin.Context) {
	traceLogger, _ := req_util.GetTraceLogger(c)
	type ListReq struct {
		PageNum  int    `form:"pageNum" binding:"required"`
		PageSize int    `form:"pageSize" binding:"required"`
		Username string `form:"username"`
		Phone    string `form:"phone"`
		Status   *int   `form:"status"`
		DeptId   *int64 `form:"deptId"`
	}

	req, err := req_util.BindQuery[ListReq](c)
	if err != nil {
		res_util.Fail(c, res_util.WithMsg("参数错误"))
		return
	}

	list, total, err := service.UserService.List(c, req.PageNum, req.PageSize, req.Username, req.Phone, req.Status, req.DeptId)
	if err != nil {
		traceLogger.Error("查询失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("查询失败"))
		return
	}

	type UserVO struct {
		model.SysUser
		DeptName  string   `json:"deptName"`
		RoleNames []string `json:"roleNames"`
	}

	// 批量查询部门名称和角色名称（避免 N+1 问题）
	deptIds := make([]any, 0, len(list))
	for _, user := range list {
		if user.DeptId > 0 {
			deptIds = append(deptIds, user.DeptId)
		}
	}
	deptMap := make(map[int64]string)
	if len(deptIds) > 0 {
		depts, _ := dbw.New[model.SysDept](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(c)).
			In("id", deptIds...).
			SelectList()
		for _, d := range depts {
			deptMap[d.Id] = d.DeptName
		}
	}

	// 批量查询角色名称
	userIds := make([]any, 0, len(list))
	for _, user := range list {
		userIds = append(userIds, user.Id)
	}
	userRoleMap := make(map[int64][]string)
	if len(userIds) > 0 {
		userRoles, _ := dbw.New[model.SysUserRole](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(c)).
			In("user_id", userIds...).
			SelectList()
		if len(userRoles) > 0 {
			roleIdSet := make([]any, 0, len(userRoles))
			for _, ur := range userRoles {
				roleIdSet = append(roleIdSet, ur.RoleId)
			}
			roles, _ := dbw.New[model.SysRole](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(c)).
				In("id", roleIdSet...).
				SelectList()
			roleMap := make(map[int64]string)
			for _, r := range roles {
				roleMap[r.Id] = r.RoleName
			}
			for _, ur := range userRoles {
				if name, ok := roleMap[ur.RoleId]; ok {
					userRoleMap[ur.UserId] = append(userRoleMap[ur.UserId], name)
				}
			}
		}
	}

	voList := make([]UserVO, len(list))
	for i, user := range list {
		voList[i] = UserVO{
			SysUser:   user,
			DeptName:  deptMap[user.DeptId],
			RoleNames: userRoleMap[user.Id],
		}
	}

	res_util.Success(c, res_util.WithData(gin.H{
		"list":  voList,
		"total": total,
	}))
}

func parseRoleIds(ids []string) []int64 {
	result := make([]int64, 0, len(ids))
	for _, s := range ids {
		id, err := strconv.ParseInt(s, 10, 64)
		if err == nil {
			result = append(result, id)
		}
	}
	return result
}

// UserAdd 新增用户，校验密码强度、用户名唯一性，含角色分配
func UserAdd(c *gin.Context) {
	traceLogger, _ := req_util.GetTraceLogger(c)
	type AddReq struct {
		Username string   `json:"username" binding:"required"`
		Password string   `json:"password" binding:"required"`
		Nickname string   `json:"nickname" binding:"required"`
		Phone    string   `json:"phone"`
		Email    string   `json:"email"`
		DeptId   int64    `json:"deptId,string"`
		RoleIds  []string `json:"roleIds"`
		Status   int      `json:"status"`
	}

	req, err := req_util.BindJson[AddReq](c)
	if err != nil {
		traceLogger.Error("参数错误", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("参数错误"))
		return
	}

	if err := password.ValidatePasswordStrong(req.Password); err != nil {
		traceLogger.Warn("校验密码健壮性失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg(err.Error()))
		return
	}

	exists, err := service.UserService.CheckUsernameExists(c, req.Username)
	if err != nil {
		traceLogger.Error("检查用户名失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("操作失败"))
		return
	}
	if exists {
		res_util.Fail(c, res_util.WithMsg("用户名已存在"))
		return
	}

	hashedPwd, err := password.Hash(req.Password)
	if err != nil {
		traceLogger.Error("密码处理失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("密码处理失败"))
		return
	}
	user := &model.SysUser{
		Username: req.Username,
		Password: hashedPwd,
		Nickname: utils.StrPtr(req.Nickname),
		Phone:    utils.StrPtr(req.Phone),
		Email:    utils.StrPtr(req.Email),
		DeptId:   req.DeptId,
		Status:   req.Status,
	}

	err = service.UserService.AddWithRoles(c, user, parseRoleIds(req.RoleIds))
	if err != nil {
		traceLogger.Error("新增失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("新增失败"))
		return
	}

	res_util.Success(c)
}

// UserEdit 编辑用户，修改用户名时检查唯一性
func UserEdit(c *gin.Context) {
	traceLogger, _ := req_util.GetTraceLogger(c)
	type EditReq struct {
		Id       int64    `json:"id,string" binding:"required"`
		Username string   `json:"username"`
		Nickname string   `json:"nickname" binding:"required"`
		Phone    string   `json:"phone"`
		Email    string   `json:"email"`
		DeptId   int64    `json:"deptId,string"`
		RoleIds  []string `json:"roleIds"`
		Status   int      `json:"status"`
	}

	req, err := req_util.BindJson[EditReq](c)
	if err != nil {
		traceLogger.Warn("参数错误", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("参数错误"))
		return
	}

	if req.Username != "" {
		existingUser, err := service.UserService.GetById(c, req.Id)
		if err != nil {
			traceLogger.Error("查询用户失败", zap.Error(err))
			res_util.Fail(c, res_util.WithMsg("操作失败"))
			return
		}
		if existingUser != nil && existingUser.Username != req.Username {
			exists, err := service.UserService.CheckUsernameExists(c, req.Username)
			if err != nil {
				traceLogger.Error("检查用户名失败", zap.Error(err))
				res_util.Fail(c, res_util.WithMsg("操作失败"))
				return
			}
			if exists {
				res_util.Fail(c, res_util.WithMsg("用户名已存在"))
				return
			}
		}
	}

	user := &model.SysUser{
		Id:       req.Id,
		Nickname: utils.StrPtr(req.Nickname),
		Phone:    utils.StrPtr(req.Phone),
		Email:    utils.StrPtr(req.Email),
		DeptId:   req.DeptId,
		Status:   req.Status,
	}
	if req.Username != "" {
		user.Username = req.Username
	}

	updateBy, _ := c.Get(constants.ContextUserIDKey)
	user.UpdateBy = updateBy.(int64)

	err = service.UserService.UpdateWithRoles(c, user, parseRoleIds(req.RoleIds))
	if err != nil {
		traceLogger.Error("编辑失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("编辑失败"))
		return
	}

	res_util.Success(c)
}

// UserChangeStatus 启用/禁用用户，禁止操作超级管理员
func UserChangeStatus(c *gin.Context) {
	traceLogger, _ := req_util.GetTraceLogger(c)
	type StatusReq struct {
		Id     int64 `json:"id,string" binding:"required"`
		Status int   `json:"status" binding:"required"`
	}

	req, err := req_util.BindJson[StatusReq](c)
	if err != nil {
		traceLogger.Warn("参数错误", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("参数错误"))
		return
	}

	isAdmin, err := service.UserService.IsAdmin(c, req.Id)
	if err != nil {
		traceLogger.Error("检查用户角色失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("操作失败"))
		return
	}
	if isAdmin {
		res_util.Fail(c, res_util.WithMsg("不能修改超级管理员状态"))
		return
	}

	user := &model.SysUser{
		Id:     req.Id,
		Status: req.Status,
	}
	updateBy, _ := c.Get(constants.ContextUserIDKey)
	user.UpdateBy = updateBy.(int64)

	err = service.UserService.Update(c, user)
	if err != nil {
		traceLogger.Error("操作失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("操作失败"))
		return
	}

	res_util.Success(c)
}

// UserResetPwd 重置用户密码为随机密码，禁止重置超级管理员
func UserResetPwd(c *gin.Context) {
	traceLogger, _ := req_util.GetTraceLogger(c)
	type ResetReq struct {
		Id int64 `json:"id,string" binding:"required"`
	}

	req, err := req_util.BindJson[ResetReq](c)
	if err != nil {
		traceLogger.Warn("参数错误", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("参数错误"))
		return
	}

	isAdmin, err := service.UserService.IsAdmin(c, req.Id)
	if err != nil {
		traceLogger.Error("检查用户角色失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("操作失败"))
		return
	}
	if isAdmin {
		res_util.Fail(c, res_util.WithMsg("不能重置超级管理员密码"))
		return
	}

	newPwd := password.GenerateRandomPassword(12)
	hashedPwd, err := password.Hash(newPwd)
	if err != nil {
		traceLogger.Error("密码处理失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("密码处理失败"))
		return
	}
	user := &model.SysUser{
		Id:       req.Id,
		Password: hashedPwd,
	}

	err = service.UserService.Update(c, user)
	if err != nil {
		traceLogger.Error("重置失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("重置失败"))
		return
	}

	res_util.Success(c, res_util.WithMsg("密码已重置"))
}

// UserDelete 删除用户（逻辑删除），禁止删除超级管理员
func UserDelete(c *gin.Context) {
	traceLogger, _ := req_util.GetTraceLogger(c)
	type DeleteReq struct {
		Id int64 `json:"id,string" binding:"required"`
	}

	req, err := req_util.BindJson[DeleteReq](c)
	if err != nil {
		traceLogger.Warn("参数错误", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("参数错误"))
		return
	}

	isAdmin, err := service.UserService.IsAdmin(c, req.Id)
	if err != nil {
		traceLogger.Error("检查用户角色失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("操作失败"))
		return
	}
	if isAdmin {
		res_util.Fail(c, res_util.WithMsg("不能删除超级管理员"))
		return
	}

	err = service.UserService.Delete(c, req.Id)
	if err != nil {
		traceLogger.Error("删除失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("删除失败"))
		return
	}

	res_util.Success(c)
}

func UserRoleIds(c *gin.Context) {
	traceLogger, _ := req_util.GetTraceLogger(c)
	userIdStr := c.Query("userId")
	if userIdStr == "" {
		res_util.Fail(c, res_util.WithMsg("参数错误"))
		return
	}

	var userId int64
	fmt.Sscanf(userIdStr, "%d", &userId)

	roles, err := service.UserService.GetUserRoles(c, userId)
	if err != nil {
		traceLogger.Error("查询失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("查询失败"))
		return
	}

	roleIds := make([]string, len(roles))
	for i, role := range roles {
		roleIds[i] = fmt.Sprintf("%d", role.Id)
	}

	res_util.Success(c, res_util.WithData(roleIds))
}
