package service

import (
	"backend/model"
	"backend/pkg/global_vars"
	"errors"

	"github.com/shangjundragon/dbw"
)

var RoleService = new(roleService)

type roleService struct{}

func (s *roleService) List(pageNum, pageSize int, roleName string) ([]model.SysRole, int64, error) {
	wrapper := dbw.New[model.SysRole](dbw.WithConfig(global_vars.DbConfig))
	if roleName != "" {
		wrapper = wrapper.Like("role_name", "%"+roleName+"%")
	}
	list, total, err := wrapper.OrderByDesc("create_time").SelectPage(pageNum, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (s *roleService) Add(role *model.SysRole) error {
	_, err := dbw.New[model.SysRole](dbw.WithConfig(global_vars.DbConfig)).Insert(role)
	return err
}

func (s *roleService) Update(role *model.SysRole) error {
	_, err := dbw.New[model.SysRole](dbw.WithConfig(global_vars.DbConfig)).UpdateById(role)
	return err
}

func (s *roleService) Delete(id int64) error {
	userRoles, _ := dbw.New[model.SysUserRole](dbw.WithConfig(global_vars.DbConfig)).
		Eq("role_id", id).
		SelectList()
	if len(userRoles) > 0 {
		return errors.New("角色下存在用户，无法删除")
	}

	_, err := dbw.New[model.SysRoleMenu](dbw.WithConfig(global_vars.DbConfig)).
		Eq("role_id", id).
		Delete()
	if err != nil {
		return err
	}

	_, err = dbw.New[model.SysRole](dbw.WithConfig(global_vars.DbConfig)).DeleteById(id)
	return err
}

func (s *roleService) GetById(id int64) (*model.SysRole, error) {
	return dbw.New[model.SysRole](dbw.WithConfig(global_vars.DbConfig)).SelectById(id)
}

func (s *roleService) GetMenuIds(roleId int64) ([]int64, error) {
	roleMenus, err := dbw.New[model.SysRoleMenu](dbw.WithConfig(global_vars.DbConfig)).
		Eq("role_id", roleId).
		SelectList()
	if err != nil {
		return nil, err
	}

	menuIds := make([]int64, len(roleMenus))
	for i, rm := range roleMenus {
		menuIds[i] = rm.MenuId
	}
	return menuIds, nil
}

func (s *roleService) AssignMenu(roleId int64, menuIds []int64) error {
	_, err := dbw.New[model.SysRoleMenu](dbw.WithConfig(global_vars.DbConfig)).
		Eq("role_id", roleId).
		Delete()
	if err != nil {
		return err
	}

	for _, menuId := range menuIds {
		_, err = dbw.New[model.SysRoleMenu](dbw.WithConfig(global_vars.DbConfig)).Insert(
			&model.SysRoleMenu{
				RoleId: roleId,
				MenuId: menuId,
			})
		if err != nil {
			return err
		}
	}
	return nil
}
