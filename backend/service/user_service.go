package service

import (
	"backend/model"
	"backend/pkg/cache"
	"backend/pkg/global_vars"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/shangjundragon/dbw"
)

var UserService = new(userService)

type userService struct{}

func (s *userService) CheckPermission(ctx context.Context, userID int64, perm string) (bool, error) {
	userRoles, err := dbw.New[model.SysUserRole](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).
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

	roleMenus, err := dbw.New[model.SysRoleMenu](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).
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

	menus, err := dbw.New[model.SysMenu](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).
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

func (s *userService) GetUserRoles(ctx context.Context, userID int64) ([]model.SysRole, error) {
	userRoles, err := dbw.New[model.SysUserRole](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).
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

	roles, err := dbw.New[model.SysRole](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).
		In("id", roleIds...).
		SelectList()
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (s *userService) GetUserPermissions(ctx context.Context, userID int64) ([]string, error) {
	cacheKey := fmt.Sprintf("perms:%d", userID)
	if cache.GlobalCache != nil {
		cached, err := cache.GlobalCache.Get(ctx, cacheKey)
		if err == nil && cached != "" {
			var perms []string
			if json.Unmarshal([]byte(cached), &perms) == nil {
				return perms, nil
			}
		}
	}

	roles, err := s.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}

	var perms []string
	for _, role := range roles {
		if role.RoleCode == "admin" {
			menus, _ := dbw.New[model.SysMenu](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).
				Eq("menu_type", 2).
				SelectList()
			perms = make([]string, 0, len(menus))
			for _, menu := range menus {
				if menu.Perm != nil && *menu.Perm != "" {
					perms = append(perms, *menu.Perm)
				}
			}
			if cache.GlobalCache != nil {
				if b, e := json.Marshal(perms); e == nil {
					cache.GlobalCache.Set(ctx, cacheKey, string(b), 10*time.Minute)
				}
			}
			return perms, nil
		}
	}

	userRoles, _ := dbw.New[model.SysUserRole](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).
		Eq("user_id", userID).
		SelectList()
	if len(userRoles) == 0 {
		return []string{}, nil
	}

	roleIds := make([]any, len(userRoles))
	for i, ur := range userRoles {
		roleIds[i] = ur.RoleId
	}

	roleMenus, _ := dbw.New[model.SysRoleMenu](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).
		In("role_id", roleIds...).
		SelectList()
	if len(roleMenus) == 0 {
		return []string{}, nil
	}

	menuIds := make([]any, len(roleMenus))
	for i, rm := range roleMenus {
		menuIds[i] = rm.MenuId
	}

	menus, _ := dbw.New[model.SysMenu](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).
		In("id", menuIds...).
		Eq("menu_type", 2).
		SelectList()

	perms = make([]string, 0, len(menus))
	for _, menu := range menus {
		if menu.Perm != nil && *menu.Perm != "" {
			perms = append(perms, *menu.Perm)
		}
	}

	if cache.GlobalCache != nil {
		if b, e := json.Marshal(perms); e == nil {
			cache.GlobalCache.Set(ctx, cacheKey, string(b), 10*time.Minute)
		}
	}

	return perms, nil
}

func (s *userService) GetUserMenus(ctx context.Context, userID int64) ([]model.SysMenu, error) {
	roles, err := s.GetUserRoles(ctx, userID)
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
		menus, err = dbw.New[model.SysMenu](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).
			OrderBy("order_num").
			SelectList()
	} else {
		userRoles, _ := dbw.New[model.SysUserRole](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).
			Eq("user_id", userID).
			SelectList()
		if len(userRoles) == 0 {
			return []model.SysMenu{}, nil
		}

		roleIds := make([]any, len(userRoles))
		for i, ur := range userRoles {
			roleIds[i] = ur.RoleId
		}

		roleMenus, _ := dbw.New[model.SysRoleMenu](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).
			In("role_id", roleIds...).
			SelectList()
		if len(roleMenus) == 0 {
			return []model.SysMenu{}, nil
		}

		menuIds := make([]any, len(roleMenus))
		for i, rm := range roleMenus {
			menuIds[i] = rm.MenuId
		}

		menus, _ = dbw.New[model.SysMenu](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).
			In("id", menuIds...).
			OrderBy("order_num").
			SelectList()
	}

	return menus, nil
}

func (s *userService) GetByUsername(ctx context.Context, username string) (*model.SysUser, error) {
	user, err := dbw.New[model.SysUser](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).
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

func (s *userService) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	user, err := s.GetByUsername(ctx, username)
	if err != nil {
		return false, err
	}
	return user != nil, nil
}

func (s *userService) List(ctx context.Context, pageNum, pageSize int, username, phone string, status *int, deptId *int64) ([]model.SysUser, int64, error) {
	wrapper := dbw.New[model.SysUser](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx))

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

func (s *userService) Add(ctx context.Context, user *model.SysUser) error {
	_, err := dbw.New[model.SysUser](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).Insert(user)
	return err
}

func (s *userService) AddWithRoles(ctx context.Context, user *model.SysUser, roleIds []int64) error {
	return dbw.ExecuteTx(func(tx *sql.Tx) error {
		_, err := dbw.New[model.SysUser](dbw.WithConfig(global_vars.DbConfig), dbw.WithTx(tx), dbw.WithContext(ctx)).Insert(user)
		if err != nil {
			return err
		}

		for _, roleId := range roleIds {
			_, err = dbw.New[model.SysUserRole](dbw.WithConfig(global_vars.DbConfig), dbw.WithTx(tx), dbw.WithContext(ctx)).Insert(
				&model.SysUserRole{UserId: user.Id, RoleId: roleId})
			if err != nil {
				return err
			}
		}
		return nil
	}, global_vars.DbConfig.Db)
}

func (s *userService) Update(ctx context.Context, user *model.SysUser) error {
	_, err := dbw.New[model.SysUser](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).UpdateById(user)
	return err
}

func (s *userService) UpdateWithRoles(ctx context.Context, user *model.SysUser, roleIds []int64) error {
	return dbw.ExecuteTx(func(tx *sql.Tx) error {
		_, err := dbw.New[model.SysUser](dbw.WithConfig(global_vars.DbConfig), dbw.WithTx(tx), dbw.WithContext(ctx)).UpdateById(user)
		if err != nil {
			return err
		}

		_, err = dbw.New[model.SysUserRole](dbw.WithConfig(global_vars.DbConfig), dbw.WithTx(tx), dbw.WithContext(ctx)).
			Eq("user_id", user.Id).
			Delete()
		if err != nil {
			return err
		}

		for _, roleId := range roleIds {
			_, err = dbw.New[model.SysUserRole](dbw.WithConfig(global_vars.DbConfig), dbw.WithTx(tx), dbw.WithContext(ctx)).Insert(
				&model.SysUserRole{UserId: user.Id, RoleId: roleId})
			if err != nil {
				return err
			}
		}
		return nil
	}, global_vars.DbConfig.Db)
}

func (s *userService) Delete(ctx context.Context, id int64) error {
	return dbw.ExecuteTx(func(tx *sql.Tx) error {
		_, err := dbw.New[model.SysUserRole](dbw.WithConfig(global_vars.DbConfig), dbw.WithTx(tx), dbw.WithContext(ctx)).
			Eq("user_id", id).
			Delete()
		if err != nil {
			return err
		}

		_, err = dbw.New[model.SysUser](dbw.WithConfig(global_vars.DbConfig), dbw.WithTx(tx), dbw.WithContext(ctx)).DeleteById(id)
		return err
	}, global_vars.DbConfig.Db)
}

func (s *userService) GetById(ctx context.Context, id int64) (*model.SysUser, error) {
	return dbw.New[model.SysUser](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).SelectById(id)
}
