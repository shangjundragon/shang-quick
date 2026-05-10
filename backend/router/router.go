package router

import (
	"backend/controller"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORS())
	r.Use(middleware.TraceLogger())

	r.GET("/uploads/*filepath", controller.FileAccess)

	api := r.Group("/api/v1")
	{
		api.POST("/auth/login", controller.AuthLogin)

		authorized := api.Group("", middleware.JWTAuth())
		{
			authorized.GET("/auth/info", controller.AuthInfo)
		}
		{
			authorized.GET("/user/list", middleware.Permission("user:list"), controller.UserList)
			authorized.POST("/user/add", middleware.OperLog("用户管理", 1), middleware.Permission("user:add"), controller.UserAdd)
			authorized.POST("/user/edit", middleware.OperLog("用户管理", 2), middleware.Permission("user:edit"), controller.UserEdit)
			authorized.POST("/user/changeStatus", middleware.OperLog("用户管理", 2), middleware.Permission("user:edit"), controller.UserChangeStatus)
			authorized.POST("/user/resetPwd", middleware.OperLog("用户管理", 2), middleware.Permission("user:resetPwd"), controller.UserResetPwd)
			authorized.POST("/user/delete", middleware.OperLog("用户管理", 3), middleware.Permission("user:delete"), controller.UserDelete)
			authorized.GET("/user/roleIds", middleware.OperLog("用户管理", 4), middleware.Permission("user:list"), controller.UserRoleIds)

			authorized.GET("/dept/list", middleware.Permission("dept:list"), controller.DeptList)
			authorized.POST("/dept/add", middleware.OperLog("部门管理", 1), middleware.Permission("dept:add"), controller.DeptAdd)
			authorized.POST("/dept/edit", middleware.OperLog("部门管理", 2), middleware.Permission("dept:edit"), controller.DeptEdit)
			authorized.POST("/dept/delete", middleware.OperLog("部门管理", 3), middleware.Permission("dept:delete"), controller.DeptDelete)

			authorized.GET("/menu/list", middleware.Permission("menu:list"), controller.MenuList)
			authorized.POST("/menu/add", middleware.OperLog("菜单管理", 1), middleware.Permission("menu:add"), controller.MenuAdd)
			authorized.POST("/menu/edit", middleware.OperLog("菜单管理", 2), middleware.Permission("menu:edit"), controller.MenuEdit)
			authorized.POST("/menu/delete", middleware.OperLog("菜单管理", 3), middleware.Permission("menu:delete"), controller.MenuDelete)

			authorized.GET("/role/list", middleware.Permission("role:list"), controller.RoleList)
			authorized.POST("/role/add", middleware.OperLog("角色管理", 1), middleware.Permission("role:add"), controller.RoleAdd)
			authorized.POST("/role/edit", middleware.OperLog("角色管理", 2), middleware.Permission("role:edit"), controller.RoleEdit)
			authorized.POST("/role/delete", middleware.OperLog("角色管理", 3), middleware.Permission("role:delete"), controller.RoleDelete)
			authorized.GET("/role/menuIds", middleware.Permission("role:list"), controller.RoleMenuIds)
			authorized.POST("/role/assignMenu", middleware.OperLog("角色管理", 2), middleware.Permission("role:assign"), controller.RoleAssignMenu)

			authorized.GET("/profile", controller.ProfileGet)
			authorized.POST("/profile/update", controller.ProfileUpdate)
			authorized.POST("/profile/updatePwd", controller.ProfileUpdatePwd)

			authorized.GET("/operLog/list", controller.OperLogList)

			authorized.GET("/file/list", middleware.Permission("file:list"), controller.FileList)
			authorized.POST("/file/upload", middleware.Permission("file:upload"), controller.FileUpload)
			authorized.POST("/file/delete", middleware.OperLog("文件管理", 3), middleware.Permission("file:delete"), controller.FileDelete)
			authorized.GET("/file/config", controller.FileConfig)

			authorized.GET("/onlineUser/list", middleware.Permission("onlineUser:list"), controller.OnlineUserList)
			authorized.POST("/onlineUser/kick", middleware.OperLog("在线用户", 3), middleware.Permission("onlineUser:kick"), controller.OnlineUserKick)
		}
	}

	return r
}
