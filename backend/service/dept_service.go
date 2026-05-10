package service

import (
	"backend/model"
	"backend/pkg/constants"
	"backend/pkg/global_vars"
	"context"
	"errors"

	"github.com/shangjundragon/dbw"
)

var DeptService = new(deptService)

type deptService struct{}

func (s *deptService) List(ctx context.Context) ([]model.SysDept, error) {
	return dbw.New[model.SysDept](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).
		OrderBy("order_num").
		SelectList()
}

func (s *deptService) Add(ctx context.Context, dept *model.SysDept) error {
	_, err := dbw.New[model.SysDept](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).Insert(dept)
	return err
}

func (s *deptService) Update(ctx context.Context, dept *model.SysDept) error {
	_, err := dbw.New[model.SysDept](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).UpdateById(dept)
	return err
}

func (s *deptService) Delete(ctx context.Context, id int64) error {
	allDepts, _ := dbw.New[model.SysDept](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).
		Eq("del_flag", constants.N).
		SelectList()
	descendantIds := s.collectDescendantIds(allDepts, id)
	if len(descendantIds) > 0 {
		return errors.New("存在子部门，无法删除")
	}

	users, _ := dbw.New[model.SysUser](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).
		Eq("dept_id", id).
		SelectList()
	if len(users) > 0 {
		return errors.New("部门下存在用户，无法删除")
	}

	_, err := dbw.New[model.SysDept](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).DeleteById(id)
	return err
}

func (s *deptService) collectDescendantIds(allDepts []model.SysDept, parentId int64) []int64 {
	var ids []int64
	for _, dept := range allDepts {
		if dept.ParentId == parentId {
			ids = append(ids, dept.Id)
			ids = append(ids, s.collectDescendantIds(allDepts, dept.Id)...)
		}
	}
	return ids
}

func (s *deptService) GetById(ctx context.Context, id int64) (*model.SysDept, error) {
	return dbw.New[model.SysDept](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).SelectById(id)
}
