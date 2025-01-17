package lgjx

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/lgjx"
	"github.com/flipped-aurora/gin-vue-admin/server/model/lgjx/jrapi/jrclientrequest"
	"github.com/flipped-aurora/gin-vue-admin/server/model/lgjx/jrapi/jrclientresponse"
	"github.com/flipped-aurora/gin-vue-admin/server/model/lgjx/jrapi/jrrequest"
	"github.com/flipped-aurora/gin-vue-admin/server/model/lgjx/jrapi/jrresponse"
	lgjxReq "github.com/flipped-aurora/gin-vue-admin/server/model/lgjx/request"
	lgjx2 "github.com/flipped-aurora/gin-vue-admin/server/utils/lgjx"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type TestInvoiceApplyService struct {
}

func (testInvoiceApplyService *TestInvoiceApplyService) CreateInvoiceApply(invoiceApply lgjx.InvoiceApply) (err error) {
	err = global.MustGetGlobalDBByDBName("lg-jx-test").Create(&invoiceApply).Error
	return err
}

func (testInvoiceApplyService *TestInvoiceApplyService) DeleteInvoiceApply(invoiceApply lgjx.InvoiceApply) (err error) {
	err = global.MustGetGlobalDBByDBName("lg-jx-test").Delete(&invoiceApply).Error
	return err
}

func (testInvoiceApplyService *TestInvoiceApplyService) DeleteInvoiceApplyByIds(ids request.IdsReq) (err error) {
	err = global.MustGetGlobalDBByDBName("lg-jx-test").Delete(&[]lgjx.InvoiceApply{}, "id in ?", ids.Ids).Error
	return err
}

func (testInvoiceApplyService *TestInvoiceApplyService) UpdateInvoiceApply(invoiceApply lgjx.InvoiceApply) (err error) {
	err = global.MustGetGlobalDBByDBName("lg-jx-test").Save(&invoiceApply).Error
	return err
}

func (testInvoiceApplyService *TestInvoiceApplyService) GetInvoiceApply(id uint) (invoiceApply lgjx.InvoiceApply, err error) {
	err = global.MustGetGlobalDBByDBName("lg-jx-test").Where("id = ?", id).First(&invoiceApply).Error
	return
}

func (testInvoiceApplyService *TestInvoiceApplyService) GetInvoiceApplyInfoList(info lgjxReq.InvoiceApplySearch) (list []lgjx.InvoiceApply, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.MustGetGlobalDBByDBName("lg-jx-test").Model(&lgjx.InvoiceApply{})
	var invoiceApplys []lgjx.InvoiceApply
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Order("created_at desc").Offset(offset).Find(&invoiceApplys).Error
	return invoiceApplys, total, err
}

func (testInvoiceApplyService *TestInvoiceApplyService) ApproveInvoiceApply(invoiceApply lgjx.InvoiceApply) (err error) {
	err = global.MustGetGlobalDBByDBName("lg-jx-test").Transaction(func(tx *gorm.DB) error {
		auditStatus := int64(2)
		auditOpinion := "受理成功"
		auditDate := time.Now().Format("2006-01-02 15:04:05")
		invoiceApply.AuditStatus = &auditStatus
		invoiceApply.AuditOpinion = &auditOpinion
		invoiceApply.AuditDate = &auditDate
		var orderList []jrrequest.Order
		var invoiceList []jrclientrequest.Invoice
		var totalAmount float64
		err = json.Unmarshal([]byte(*invoiceApply.OrderList), &orderList)
		if err != nil {
			return err
		}
		if len(orderList) > 0 {
			totalAmount = 0
			for _, order := range orderList {
				totalAmount += *order.OrderPayAmount
				invoice := jrclientrequest.Invoice{
					InvoiceNo:          "0123456789abcdef",
					InvoiceAmount:      *order.OrderPayAmount,
					InvoiceType:        "A1",
					InvoiceForm:        "B1",
					InvoicePoint:       0.0001,
					InvoiceContent:     "发票内容",
					InvoiceRemark:      "发票备注",
					InvoiceTime:        "2022-12-19 17:00:00",
					InvoiceDownloadUrl: global.GVA_CONFIG.Insurance.APIDoaminTest + "/invoiceFileDownload?param=testInvoiceDownloadUrl",
					OrderList: []jrclientrequest.Order{{
						OrderNo:            *order.OrderNo,
						OrderInvoiceAmount: *order.OrderPayAmount,
					}},
				}
				invoiceList = append(invoiceList, invoice)
			}
		}
		jsonInvoiceList, err := json.Marshal(invoiceList)
		if err != nil {
			return err
		}
		invoiceListString := string(jsonInvoiceList)
		invoiceApply.InvoiceList = &invoiceListString
		err = tx.Save(&invoiceApply).Error
		if err != nil {
			return err
		}

		if global.GVA_CONFIG.Insurance.JRAPIDomainTest != "" {
			apiPath := "/jrapi/lg/invoicePush"
			var invoicePush = jrclientrequest.InvoicePush{
				ApplyNo:            *invoiceApply.ApplyNo,
				AuditStatus:        *invoiceApply.AuditStatus,
				AuditOpinion:       *invoiceApply.AuditOpinion,
				AuditDate:          *invoiceApply.AuditDate,
				InvoiceTotalAmount: totalAmount,
				InvoiceList:        invoiceList,
			}
			req, err := lgjx2.GenJRRequest(invoicePush)
			if err != nil {
				return err
			}
			var res jrresponse.JRResponse
			client := resty.New()
			resp, err := client.R().
				SetBody(&req).
				SetResult(&res).
				Post(global.GVA_CONFIG.Insurance.JRAPIDomainTest + apiPath)
			if err != nil {
				return err
			}
			if resp.StatusCode() == http.StatusOK {
				if res.Code != 0 {
					err := errors.New(res.Msg)
					global.GVA_LOG.Error("调用"+apiPath+"失败", zap.Error(err))
					return err
				} else {
					byteEncryptData, err := base64.StdEncoding.DecodeString(res.Data)
					if err != nil {
						return err
					}
					jsonData, err := lgjx2.Sm4Decrypt(byteEncryptData)
					if err != nil {
						return err
					}
					var resData jrclientresponse.Response
					err = json.Unmarshal([]byte(jsonData), &resData)
					if err != nil {
						return err
					}
					if resData.ReceiveResult != "success" {
						global.GVA_LOG.Error("调用"+apiPath+"结果不为success", zap.Error(err))
						return errors.New("接收结果不为success")
					}
				}
			} else {
				return errors.New("交易中心响应失败")
			}
		}

		return nil
	})
	return err
}

func (testInvoiceApplyService *TestInvoiceApplyService) RejectInvoiceApply(invoiceApply lgjx.InvoiceApply) (err error) {
	err = global.MustGetGlobalDBByDBName("lg-jx-test").Transaction(func(tx *gorm.DB) error {
		auditStatus := int64(3)
		auditOpinion := "受理失败"
		auditDate := time.Now().Format("2006-01-02 15:04:05")
		invoiceApply.AuditStatus = &auditStatus
		invoiceApply.AuditOpinion = &auditOpinion
		invoiceApply.AuditDate = &auditDate
		err := tx.Save(&invoiceApply).Error
		if err != nil {
			return err
		}

		if global.GVA_CONFIG.Insurance.JRAPIDomainTest != "" {
			apiPath := "/jrapi/lg/invoicePush"
			var invoicePush = jrclientrequest.InvoicePush{
				ApplyNo:      *invoiceApply.ApplyNo,
				AuditStatus:  *invoiceApply.AuditStatus,
				AuditOpinion: *invoiceApply.AuditOpinion,
				AuditDate:    *invoiceApply.AuditDate,
			}
			req, err := lgjx2.GenJRRequest(invoicePush)
			if err != nil {
				return err
			}
			var res jrresponse.JRResponse
			client := resty.New()
			resp, err := client.R().
				SetBody(&req).
				SetResult(&res).
				Post(global.GVA_CONFIG.Insurance.JRAPIDomainTest + apiPath)
			if err != nil {
				return err
			}
			if resp.StatusCode() == http.StatusOK {
				if res.Code != 0 {
					err := errors.New(res.Msg)
					global.GVA_LOG.Error("调用"+apiPath+"失败", zap.Error(err))
					return err
				} else {
					byteEncryptData, err := base64.StdEncoding.DecodeString(res.Data)
					if err != nil {
						return err
					}
					jsonData, err := lgjx2.Sm4Decrypt(byteEncryptData)
					if err != nil {
						return err
					}
					var resData jrclientresponse.Response
					err = json.Unmarshal([]byte(jsonData), &resData)
					if err != nil {
						return err
					}
					if resData.ReceiveResult != "success" {
						global.GVA_LOG.Error("调用"+apiPath+"结果不为success", zap.Error(err))
						return errors.New("接收结果不为success")
					}
				}
			} else {
				return errors.New("交易中心响应失败")
			}
		}

		return nil
	})
	return err
}
