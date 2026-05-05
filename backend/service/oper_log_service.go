package service

import (
	"backend/model"
	"backend/pkg/global_vars"
	"context"

	"github.com/shangjundragon/dbw"
)

var OperLogService = new(operLogService)

type operLogService struct{}

func (s *operLogService) Save(ctx context.Context, log *model.SysOperLog) error {
	_, err := dbw.New[model.SysOperLog](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).Insert(log)
	return err
}

func (s *operLogService) List(ctx context.Context, pageNum, pageSize int, title, operName string) ([]model.SysOperLog, int64, error) {
	wrapper := dbw.New[model.SysOperLog](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx))
	if title != "" {
		wrapper = wrapper.Like("title", "%"+title+"%")
	}
	if operName != "" {
		wrapper = wrapper.Like("oper_name", "%"+operName+"%")
	}
	list, total, err := wrapper.OrderByDesc("oper_time").SelectPage(pageNum, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
