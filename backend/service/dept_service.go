package service

import (
	"backend/model"
	"backend/pkg/global_vars"
	"errors"

	"github.com/shangjundragon/dbw"
)

var DeptService = new(deptService)

type deptService struct{}

func (s *deptService) List() ([]model.SysDept, error) {
	return dbw.New[model.SysDept](dbw.WithConfig(global_vars.DbConfig)).
		OrderBy("order_num").
		SelectList()
}

func (s *deptService) Add(dept *model.SysDept) error {
	_, err := dbw.New[model.SysDept](dbw.WithConfig(global_vars.DbConfig)).Insert(dept)
	return err
}

func (s *deptService) Update(dept *model.SysDept) error {
	_, err := dbw.New[model.SysDept](dbw.WithConfig(global_vars.DbConfig)).UpdateById(dept)
	return err
}

func (s *deptService) Delete(id int64) error {
	children, _ := dbw.New[model.SysDept](dbw.WithConfig(global_vars.DbConfig)).
		Eq("parent_id", id).
		SelectList()
	if len(children) > 0 {
		return errors.New("存在子部门，无法删除")
	}

	users, _ := dbw.New[model.SysUser](dbw.WithConfig(global_vars.DbConfig)).
		Eq("dept_id", id).
		SelectList()
	if len(users) > 0 {
		return errors.New("部门下存在用户，无法删除")
	}

	_, err := dbw.New[model.SysDept](dbw.WithConfig(global_vars.DbConfig)).DeleteById(id)
	return err
}

func (s *deptService) GetById(id int64) (*model.SysDept, error) {
	return dbw.New[model.SysDept](dbw.WithConfig(global_vars.DbConfig)).SelectById(id)
}
