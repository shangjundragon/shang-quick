package service

import (
	"backend/model"
	"backend/pkg/global_vars"

	"github.com/shangjundragon/dbw"
)

var MenuService = new(menuService)

type menuService struct{}

func (s *menuService) List() ([]model.SysMenu, error) {
	return dbw.New[model.SysMenu](dbw.WithConfig(global_vars.DbConfig)).
		OrderBy("order_num").
		SelectList()
}

func (s *menuService) Add(menu *model.SysMenu) error {
	_, err := dbw.New[model.SysMenu](dbw.WithConfig(global_vars.DbConfig)).Insert(menu)
	return err
}

func (s *menuService) Update(menu *model.SysMenu) error {
	_, err := dbw.New[model.SysMenu](dbw.WithConfig(global_vars.DbConfig)).UpdateById(menu)
	return err
}

func (s *menuService) Delete(id int64) error {
	children, _ := dbw.New[model.SysMenu](dbw.WithConfig(global_vars.DbConfig)).
		Eq("parent_id", id).
		SelectList()

	for _, child := range children {
		s.Delete(child.Id)
	}

	_, err := dbw.New[model.SysMenu](dbw.WithConfig(global_vars.DbConfig)).DeleteById(id)
	return err
}

func (s *menuService) GetById(id int64) (*model.SysMenu, error) {
	return dbw.New[model.SysMenu](dbw.WithConfig(global_vars.DbConfig)).SelectById(id)
}
