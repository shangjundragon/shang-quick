package service

import (
	"backend/model"
	"backend/pkg/global_vars"
	"context"
	"database/sql"
	"errors"

	"github.com/shangjundragon/dbw"
)

var RoleService = new(roleService)

type roleService struct{}

func (s *roleService) List(ctx context.Context, pageNum, pageSize int, roleName string) ([]model.SysRole, int64, error) {
	wrapper := dbw.New[model.SysRole](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx))
	if roleName != "" {
		wrapper = wrapper.Like("role_name", "%"+roleName+"%")
	}
	list, total, err := wrapper.OrderByDesc("create_time").SelectPage(pageNum, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (s *roleService) CheckRoleCodeExists(ctx context.Context, roleCode string, excludeId int64) (bool, error) {
	role, err := dbw.New[model.SysRole](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).
		Eq("role_code", roleCode).
		FindOne()
	if err != nil {
		if errors.Is(err, dbw.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	if excludeId > 0 && role.Id == excludeId {
		return false, nil
	}
	return true, nil
}

func (s *roleService) Add(ctx context.Context, role *model.SysRole) error {
	exists, err := s.CheckRoleCodeExists(ctx, role.RoleCode, 0)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("角色编码已存在")
	}
	_, err = dbw.New[model.SysRole](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).Insert(role)
	return err
}

func (s *roleService) Update(ctx context.Context, role *model.SysRole) error {
	exists, err := s.CheckRoleCodeExists(ctx, role.RoleCode, role.Id)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("角色编码已存在")
	}
	_, err = dbw.New[model.SysRole](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).UpdateById(role)
	return err
}

func (s *roleService) Delete(ctx context.Context, id int64) error {
	userRoles, _ := dbw.New[model.SysUserRole](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).
		Eq("role_id", id).
		SelectList()
	if len(userRoles) > 0 {
		return errors.New("角色下存在用户，无法删除")
	}

	err := dbw.ExecuteTx(func(tx *sql.Tx) error {
		_, err := dbw.New[model.SysRoleMenu](dbw.WithConfig(global_vars.DbConfig), dbw.WithTx(tx), dbw.WithContext(ctx)).
			Eq("role_id", id).
			Delete()
		if err != nil {
			return err
		}

		_, err = dbw.New[model.SysRole](dbw.WithConfig(global_vars.DbConfig), dbw.WithTx(tx), dbw.WithContext(ctx)).DeleteById(id)
		return err
	}, global_vars.DbConfig.Db)

	if err == nil {
		UserService.ClearAllPermissionCache(ctx)
	}
	return err
}

func (s *roleService) GetById(ctx context.Context, id int64) (*model.SysRole, error) {
	return dbw.New[model.SysRole](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).SelectById(id)
}

func (s *roleService) GetMenuIds(ctx context.Context, roleId int64) ([]int64, error) {
	roleMenus, err := dbw.New[model.SysRoleMenu](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).
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

func (s *roleService) AssignMenu(ctx context.Context, roleId int64, menuIds []int64) error {
	err := dbw.ExecuteTx(func(tx *sql.Tx) error {
		_, err := dbw.New[model.SysRoleMenu](dbw.WithConfig(global_vars.DbConfig), dbw.WithTx(tx), dbw.WithContext(ctx)).
			Eq("role_id", roleId).
			Delete()
		if err != nil {
			return err
		}

		for _, menuId := range menuIds {
			_, err = dbw.New[model.SysRoleMenu](dbw.WithConfig(global_vars.DbConfig), dbw.WithTx(tx), dbw.WithContext(ctx)).Insert(
				&model.SysRoleMenu{
					RoleId: roleId,
					MenuId: menuId,
				})
			if err != nil {
				return err
			}
		}
		return nil
	}, global_vars.DbConfig.Db)
	if err != nil {
		return err
	}
	UserService.ClearAllPermissionCache(ctx)
	return nil
}
