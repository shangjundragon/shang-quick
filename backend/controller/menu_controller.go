package controller

import (
	"backend/model"
	"backend/pkg/req_util"
	"backend/pkg/res_util"
	"backend/pkg/utils"
	"backend/service"

	"github.com/gin-gonic/gin"
)

func MenuList(c *gin.Context) {
	list, err := service.MenuService.List()
	if err != nil {
		res_util.Fail(c, res_util.WithMsg("查询失败"))
		return
	}
	res_util.Success(c, res_util.WithData(list))
}

func MenuAdd(c *gin.Context) {
	type AddReq struct {
		ParentId  int64  `json:"parentId"`
		MenuName  string `json:"menuName" binding:"required"`
		MenuType  int    `json:"menuType" binding:"required"`
		Icon      string `json:"icon"`
		Path      string `json:"path"`
		Component string `json:"component"`
		Perm      string `json:"perm"`
		OrderNum  int    `json:"orderNum"`
		IsFrame   int    `json:"isFrame"`
		IsCache   int    `json:"isCache"`
		IsVisible int    `json:"isVisible"`
	}

	req, err := req_util.BindJson[AddReq](c)
	if err != nil {
		res_util.Fail(c, res_util.WithMsg("参数错误"))
		return
	}

	menu := &model.SysMenu{
		ParentId:  req.ParentId,
		MenuName:  req.MenuName,
		MenuType:  req.MenuType,
		Icon:      utils.StrPtr(req.Icon),
		Path:      utils.StrPtr(req.Path),
		Component: utils.StrPtr(req.Component),
		Perm:      utils.StrPtr(req.Perm),
		OrderNum:  req.OrderNum,
		IsFrame:   req.IsFrame,
		IsCache:   req.IsCache,
		IsVisible: req.IsVisible,
		Status:    1,
	}

	err = service.MenuService.Add(menu)
	if err != nil {
		res_util.Fail(c, res_util.WithMsg("新增失败"))
		return
	}

	res_util.Success(c)
}

func MenuEdit(c *gin.Context) {
	type EditReq struct {
		Id        int64  `json:"id" binding:"required"`
		ParentId  int64  `json:"parentId"`
		MenuName  string `json:"menuName" binding:"required"`
		MenuType  int    `json:"menuType" binding:"required"`
		Icon      string `json:"icon"`
		Path      string `json:"path"`
		Component string `json:"component"`
		Perm      string `json:"perm"`
		OrderNum  int    `json:"orderNum"`
		IsFrame   int    `json:"isFrame"`
		IsCache   int    `json:"isCache"`
		IsVisible int    `json:"isVisible"`
	}

	req, err := req_util.BindJson[EditReq](c)
	if err != nil {
		res_util.Fail(c, res_util.WithMsg("参数错误"))
		return
	}

	menu := &model.SysMenu{
		Id:        req.Id,
		ParentId:  req.ParentId,
		MenuName:  req.MenuName,
		MenuType:  req.MenuType,
		Icon:      utils.StrPtr(req.Icon),
		Path:      utils.StrPtr(req.Path),
		Component: utils.StrPtr(req.Component),
		Perm:      utils.StrPtr(req.Perm),
		OrderNum:  req.OrderNum,
		IsFrame:   req.IsFrame,
		IsCache:   req.IsCache,
		IsVisible: req.IsVisible,
	}

	err = service.MenuService.Update(menu)
	if err != nil {
		res_util.Fail(c, res_util.WithMsg("编辑失败"))
		return
	}

	res_util.Success(c)
}

func MenuDelete(c *gin.Context) {
	type DeleteReq struct {
		Id int64 `json:"id" binding:"required"`
	}

	req, err := req_util.BindJson[DeleteReq](c)
	if err != nil {
		res_util.Fail(c, res_util.WithMsg("参数错误"))
		return
	}

	err = service.MenuService.Delete(req.Id)
	if err != nil {
		res_util.Fail(c, res_util.WithMsg("删除失败"))
		return
	}

	res_util.Success(c)
}
