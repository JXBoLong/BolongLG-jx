package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/lgjx"
	"time"
)

type OrderSearch struct {
	lgjx.Order
	ApplyNo         *string     `json:"applyNo" form:"applyNo"`
	ProjectName     *string     `json:"projectName" form:"projectName"`
	InsureName      *string     `json:"insureName" form:"insureName"`
	ElogTemplateId  *string     `json:"elogTemplateId" form:"elogTemplateId"`
	ElogNo          *string     `json:"elogNo" form:"elogNo"`
	OrderStatus     *string     `json:"orderStatus" form:"orderStatus"`
	AuditStatus     *int        `json:"auditStatus" form:"auditStatus"`
	OpenBeginDate   []time.Time `json:"openBeginDate" form:"openBeginDate[]"`
	ApplyCreatedAt  []time.Time `json:"applyCreatedAt" form:"applyCreatedAt[]"`
	LetterCreatedAt []time.Time `json:"letterCreatedAt" form:"letterCreatedAt[]"`
	InsureDay       *int        `json:"insureDay" form:"insureDay"`
	AuditDelay      *bool       `json:"auditDelay" form:"auditDelay"`
	AuditRefund     *bool       `json:"auditRefund" form:"auditRefund"`
	AuditClaim      *bool       `json:"auditClaim" form:"auditClaim"`
	IsPayed         *bool       `json:"isPayed" form:"isPayed"`
	request.PageInfo
}
