package gozapread

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

func (c *ZapClient) GetDepositInvoice(amount uint) (string, error) {
	jsonStr := fmt.Sprintf(`{"amount":"%d","memo":"ZapRead.com deposit","anon":"0","use":"userDeposit","useId":-1,"useAction":-1}`, amount)
	if resp, err := c.postJSON("Lightning/GetDepositInvoice/", jsonStr, false); err == nil {
		var invoice Invoice
		if json.Unmarshal(resp, &invoice) == nil {
			return invoice.Invoice, nil
		}
	}
	return "", errors.New("GetDepositInvoice failed")
}

func (c *ZapClient) CheckPayment(req string) (bool, error) {
	invoice := InvoiceResp{
		Invoice:   req,
		IsDeposit: true,
	}
	if jsonSlc, err := json.Marshal(invoice); err == nil {
		if resp, err := c.postJSON("Lightning/CheckPayment/", string(jsonSlc), false); err == nil {
			var check PaymentCheck
			if json.Unmarshal(resp, &check) == nil {
				if check.Success {
					return check.Result, nil
				}
			}
		}
	}
	return false, errors.New("CheckPayment failed")
}

func (c *ZapClient) SubmitPaymentRequest(request string) (uint, error) {
	//do basic checks on request
	jsonStr := fmt.Sprintf(`{"request":"%s"}`, request)
	if resp, err := c.postJSON("Lightning/SubmitPaymentRequest", jsonStr, true); err == nil {
		var payment PaymentResp
		if json.Unmarshal(resp, &payment) == nil {
			return payment.Fees, nil
		}
	} else {
		return 0, err
	}
	return 0, errors.New("SubmitPaymentRequest failed")
}

func (c *ZapClient) ValidatePaymentRequest(request string) (uint, error) {
	//do basic checks on request
	jsonStr := fmt.Sprintf(`{"request":"%s"}`, request)
	if resp, err := c.postJSON("Lightning/ValidatePaymentRequest", jsonStr, true); err == nil {
		fmt.Println(string(resp))
		var req PaymentReq
		if json.Unmarshal(resp, &req) == nil {
			if amount, err := strconv.ParseUint(req.NumSatoshis, 10, 32); err == nil {
				return uint(amount), nil
			}
		}
	}
	return 0, errors.New("ValidatePaymentRequest failed")
}
