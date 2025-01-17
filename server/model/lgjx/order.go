// 自动生成模板Order
package lgjx

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// Order 结构体
type Order struct {
	global.GVA_MODEL
	OrderNo    *string  `json:"orderNo" form:"orderNo" gorm:"type:varchar(64);"`
	ApplyID    *uint    `json:"applyID" form:"applyID"`
	Apply      *Apply   `json:"apply" form:"apply"`
	PayID      *uint    `json:"payID" form:"payID"`
	Pay        *Pay     `json:"pay" form:"pay"`
	LetterID   *uint    `json:"letterID" form:"letterID"`
	Letter     *Letter  `json:"letter" form:"letter"`
	RevokeID   *uint    `json:"revokeID" form:"revokeID"`
	Revoke     *Revoke  `json:"revoke" form:"revoke"`
	DelayID    *uint    `json:"delayID" form:"delayID"`
	Delay      *Delay   `json:"delay" form:"delay"`
	RefundID   *uint    `json:"refundID" form:"refundID"`
	Refund     *Refund  `json:"refund" form:"refund"`
	ClaimID    *uint    `json:"claimID" form:"claimID"`
	Claim      *Claim   `json:"claim" form:"claim"`
	LogoutID   *uint    `json:"logoutID" form:"logoutID"`
	Logout     *Logout  `json:"logout" form:"logout"`
	InvoiceID  *uint    `json:"invoiceID" form:"invoiceID"`
	Invoice    *Invoice `json:"invoice" form:"invoice"`
	ProjectID  *uint    `json:"projectID" form:"projectID"`
	Project    *Project `json:"project" form:"project"`
	IsRepushed *bool    `json:"isRepushed" form:"isRepushed"`
}

// TableName Order 表名
func (Order) TableName() string {
	return "order"
}
