package service

import (
	"backend/model"
	"backend/pkg/global_vars"
	"errors"

	"github.com/shangjundragon/dbw"
)

var UserService = new(userService)

type userService struct{}

func (s *userService) CheckPermission(userID int64, perm string) (bool, error) {
	userRoles, err := dbw.New[model.SysUserRole](dbw.WithConfig(global_vars.DbConfig)).
		Eq("user_id", userID).
		SelectList()
	if err != nil {
		return false, err
	}

	if len(userRoles) == 0 {
		return false, nil
	}

	roleIds := make([]any, len(userRoles))
	for i, ur := range userRoles {
		roleIds[i] = ur.RoleId
	}

	roleMenus, err := dbw.New[model.SysRoleMenu](dbw.WithConfig(global_vars.DbConfig)).
		In("role_id", roleIds...).
		SelectList()
	if err != nil {
		return false, err
	}

	if len(roleMenus) == 0 {
		return false, nil
	}

	menuIds := make([]any, len(roleMenus))
	for i, rm := range roleMenus {
		menuIds[i] = rm.MenuId
	}

	menus, err := dbw.New[model.SysMenu](dbw.WithConfig(global_vars.DbConfig)).
		In("id", menuIds...).
		Eq("menu_type", 2).
		SelectList()
	if err != nil {
		return false, err
	}

	for _, menu := range menus {
		if menu.Perm != nil && *menu.Perm == perm {
			return true, nil
		}
	}

	return false, nil
}

func (s *userService) GetUserRoles(userID int64) ([]model.SysRole, error) {
	userRoles, err := dbw.New[model.SysUserRole](dbw.WithConfig(global_vars.DbConfig)).
		Eq("user_id", userID).
		SelectList()
	if err != nil {
		return nil, err
	}

	if len(userRoles) == 0 {
		return []model.SysRole{}, nil
	}

	roleIds := make([]any, len(userRoles))
	for i, ur := range userRoles {
		roleIds[i] = ur.RoleId
	}

	roles, err := dbw.New[model.SysRole](dbw.WithConfig(global_vars.DbConfig)).
		In("id", roleIds...).
		SelectList()
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (s *userService) GetUserPermissions(userID int64) ([]string, error) {
	roles, err := s.GetUserRoles(userID)
	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		if role.RoleCode == "admin" {
			menus, _ := dbw.New[model.SysMenu](dbw.WithConfig(global_vars.DbConfig)).
				Eq("menu_type", 2).
				SelectList()
			perms := make([]string, 0, len(menus))
			for _, menu := range menus {
				if menu.Perm != nil && *menu.Perm != "" {
					perms = append(perms, *menu.Perm)
				}
			}
			return perms, nil
		}
	}

	userRoles, _ := dbw.New[model.SysUserRole](dbw.WithConfig(global_vars.DbConfig)).
		Eq("user_id", userID).
		SelectList()
	if len(userRoles) == 0 {
		return []string{}, nil
	}

	roleIds := make([]any, len(userRoles))
	for i, ur := range userRoles {
		roleIds[i] = ur.RoleId
	}

	roleMenus, _ := dbw.New[model.SysRoleMenu](dbw.WithConfig(global_vars.DbConfig)).
		In("role_id", roleIds...).
		SelectList()
	if len(roleMenus) == 0 {
		return []string{}, nil
	}

	menuIds := make([]any, len(roleMenus))
	for i, rm := range roleMenus {
		menuIds[i] = rm.MenuId
	}

	menus, _ := dbw.New[model.SysMenu](dbw.WithConfig(global_vars.DbConfig)).
		In("id", menuIds...).
		Eq("menu_type", 2).
		SelectList()

	perms := make([]string, 0, len(menus))
	for _, menu := range menus {
		if menu.Perm != nil && *menu.Perm != "" {
			perms = append(perms, *menu.Perm)
		}
	}

	return perms, nil
}

func (s *userService) GetUserMenus(userID int64) ([]model.SysMenu, error) {
	roles, err := s.GetUserRoles(userID)
	if err != nil {
		return nil, err
	}

	isAdmin := false
	for _, role := range roles {
		if role.RoleCode == "admin" {
			isAdmin = true
			break
		}
	}

	var menus []model.SysMenu
	if isAdmin {
		menus, err = dbw.New[model.SysMenu](dbw.WithConfig(global_vars.DbConfig)).
			SelectList()
	} else {
		userRoles, _ := dbw.New[model.SysUserRole](dbw.WithConfig(global_vars.DbConfig)).
			Eq("user_id", userID).
			SelectList()
		if len(userRoles) == 0 {
			return []model.SysMenu{}, nil
		}

		roleIds := make([]any, len(userRoles))
		for i, ur := range userRoles {
			roleIds[i] = ur.RoleId
		}

		roleMenus, _ := dbw.New[model.SysRoleMenu](dbw.WithConfig(global_vars.DbConfig)).
			In("role_id", roleIds...).
			SelectList()
		if len(roleMenus) == 0 {
			return []model.SysMenu{}, nil
		}

		menuIds := make([]any, len(roleMenus))
		for i, rm := range roleMenus {
			menuIds[i] = rm.MenuId
		}

		menus, _ = dbw.New[model.SysMenu](dbw.WithConfig(global_vars.DbConfig)).
			In("id", menuIds...).
			SelectList()
	}

	return menus, nil
}

func (s *userService) GetByUsername(username string) (*model.SysUser, error) {
	user, err := dbw.New[model.SysUser](dbw.WithConfig(global_vars.DbConfig)).
		Eq("username", username).
		FindOne()
	if err != nil {
		if errors.Is(err, dbw.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) List(pageNum, pageSize int, username, phone string, status *int, deptId *int64) ([]model.SysUser, int64, error) {
	wrapper := dbw.New[model.SysUser](dbw.WithConfig(global_vars.DbConfig))

	if username != "" {
		wrapper = wrapper.Like("username", "%"+username+"%")
	}
	if phone != "" {
		wrapper = wrapper.Like("phone", "%"+phone+"%")
	}
	if status != nil {
		wrapper = wrapper.Eq("status", *status)
	}
	if deptId != nil {
		wrapper = wrapper.Eq("dept_id", *deptId)
	}

	list, total, err := wrapper.OrderByDesc("create_time").SelectPage(pageNum, pageSize)
	if err != nil {
		
		return nil, 0, err
	}
	return list, total, nil
}

func (s *userService) Add(user *model.SysUser) error {
	_, err := dbw.New[model.SysUser](dbw.WithConfig(global_vars.DbConfig)).Insert(user)
	return err
}

func (s *userService) Update(user *model.SysUser) error {
	_, err := dbw.New[model.SysUser](dbw.WithConfig(global_vars.DbConfig)).UpdateById(user)
	return err
}

func (s *userService) Delete(id int64) error {
	_, err := dbw.New[model.SysUser](dbw.WithConfig(global_vars.DbConfig)).DeleteById(id)
	return err
}

func (s *userService) GetById(id int64) (*model.SysUser, error) {
	return dbw.New[model.SysUser](dbw.WithConfig(global_vars.DbConfig)).SelectById(id)
}
