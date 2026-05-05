package service

import (
	"backend/model"
	"backend/pkg/global_vars"
	"context"

	"github.com/shangjundragon/dbw"
)

var MenuService = new(menuService)

type menuService struct{}

func (s *menuService) List(ctx context.Context) ([]model.SysMenu, error) {
	return dbw.New[model.SysMenu](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).
		OrderBy("order_num").
		SelectList()
}

func (s *menuService) Add(ctx context.Context, menu *model.SysMenu) error {
	_, err := dbw.New[model.SysMenu](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).Insert(menu)
	return err
}

func (s *menuService) Update(ctx context.Context, menu *model.SysMenu) error {
	_, err := dbw.New[model.SysMenu](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).UpdateById(menu)
	return err
}

func (s *menuService) Delete(ctx context.Context, id int64) error {
	children, _ := dbw.New[model.SysMenu](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).
		Eq("parent_id", id).
		SelectList()

	for _, child := range children {
		s.Delete(ctx, child.Id)
	}

	_, err := dbw.New[model.SysMenu](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).DeleteById(id)
	return err
}

func (s *menuService) GetById(ctx context.Context, id int64) (*model.SysMenu, error) {
	return dbw.New[model.SysMenu](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).SelectById(id)
}
