package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/lgjx"
	"os"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/example"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Gorm 初始化数据库并产生数据库全局变量
// Author SliverHorn
func Gorm() *gorm.DB {
	switch global.GVA_CONFIG.System.DbType {
	case "mysql":
		return GormMysql()
	case "pgsql":
		return GormPgSql()
	case "oracle":
		return GormOracle()
	case "mssql":
		return GormMssql()
	default:
		return GormMysql()
	}
}

// RegisterTables 注册数据库表专用
// Author SliverHorn
func RegisterTables(db *gorm.DB) {
	err := db.AutoMigrate(
		// 系统模块表
		system.SysApi{},
		system.SysUser{},
		system.SysBaseMenu{},
		system.JwtBlacklist{},
		system.SysAuthority{},
		system.SysDictionary{},
		system.SysOperationRecord{},
		system.SysAutoCodeHistory{},
		system.SysDictionaryDetail{},
		system.SysBaseMenuParameter{},
		system.SysBaseMenuBtn{},
		system.SysAuthorityBtn{},
		system.SysAutoCode{},

		// 示例模块表
		example.ExaFile{},
		example.ExaCustomer{},
		example.ExaFileChunk{},
		example.ExaFileUploadAndDownload{},

		// 自动化模块表
		// Code generated by github.com/flipped-aurora/gin-vue-admin/server Begin; DO NOT EDIT.

		// Code generated by github.com/flipped-aurora/gin-vue-admin/server End; DO NOT EDIT.
	)

	if err != nil {
		global.GVA_LOG.Error("register table failed", zap.Error(err))
		os.Exit(0)
	}

	err = global.MustGetGlobalDBByDBName("lg-jx").AutoMigrate(
		lgjx.Order{},
		lgjx.Apply{},
		lgjx.Pay{},
		lgjx.Letter{},
		lgjx.Revoke{},
		lgjx.Delay{},
		lgjx.Refund{},
		lgjx.Claim{},
		lgjx.Logout{},
		lgjx.Invoice{},
		lgjx.InvoiceApply{},
		lgjx.Project{},
		lgjx.Template{},
		lgjx.File{},
	)

	if err != nil {
		global.GVA_LOG.Error("lg-jx register table failed", zap.Error(err))
		os.Exit(0)
	}

	err = global.MustGetGlobalDBByDBName("lg-jx-test").AutoMigrate(
		lgjx.Order{},
		lgjx.Apply{},
		lgjx.Pay{},
		lgjx.Letter{},
		lgjx.Revoke{},
		lgjx.Delay{},
		lgjx.Refund{},
		lgjx.Claim{},
		lgjx.Logout{},
		lgjx.Invoice{},
		lgjx.InvoiceApply{},
		lgjx.Project{},
		lgjx.Template{},
		lgjx.File{},
	)

	if err != nil {
		global.GVA_LOG.Error("lg-jx-test register table failed", zap.Error(err))
		os.Exit(0)
	}
	global.GVA_LOG.Info("register table success")
}
