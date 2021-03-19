package gozapread

import (
	"encoding/json"
	"fmt"
	"strconv"
)

//GetDepositInvoice generates a LN invoice to your account
func (c *ZapClient) GetDepositInvoice(amount uint) (string, error) {
	jsonStr := fmt.Sprintf(`{"amount":"%d","memo":"ZapRead.com deposit","anon":"0","use":"userDeposit","useId":-1,"useAction":-1}`, amount)
	if resp, err := c.postJSON("Lightning/GetDepositInvoice/", jsonStr, false); err == nil {
		var invoice invoice
		if json.Unmarshal(resp, &invoice) == nil {
			return invoice.Invoice, nil
		} else {
			return "", fmt.Errorf("couldn't unmarshal the json: %w", err)
		}
	} else {
		return "", err
	}
}

//CheckPayment implements Lightning/CheckPayment
func (c *ZapClient) CheckPayment(req string) (bool, error) {
	invoice := invoice{
		Invoice:   req,
		IsDeposit: true,
	}
	if jsonSlc, err := json.Marshal(invoice); err != nil {
		return false, err
	} else {
		if resp, err := c.postJSON("Lightning/CheckPayment/", string(jsonSlc), false); err == nil {
			var check paymentCheck
			if err := json.Unmarshal(resp, &check); err == nil {
				if check.Success {
					return check.Result, nil
				} else {
					return false, fmt.Errorf("checking the payment wasn't successful.")
				}
			} else {
				return false, err
			}
		} else {
			return false, err
		}
	}
}

//SubmitPaymentRequest implements Lightning/SubmitPaymentRequest
func (c *ZapClient) SubmitPaymentRequest(request string) (uint, error) {
	//do basic checks on request
	jsonStr := fmt.Sprintf(`{"request":"%s"}`, request)
	if resp, err := c.postJSON("Lightning/SubmitPaymentRequest", jsonStr, true); err == nil {
		var payment paymentResp
		if json.Unmarshal(resp, &payment) == nil {
			return payment.Fees, nil
		} else {
			return 0, err
		}
	} else {
		return 0, err
	}
}

//ValidatePaymentRequest implements Lightning/ValidatePaymentRequest
func (c *ZapClient) ValidatePaymentRequest(request string) (uint, error) {
	//do basic checks on request
	jsonStr := fmt.Sprintf(`{"request":"%s"}`, request)
	if resp, err := c.postJSON("Lightning/ValidatePaymentRequest", jsonStr, true); err == nil {
		fmt.Println(string(resp))
		var req paymentReq
		if err := json.Unmarshal(resp, &req); err != nil {
			return 0, err
		} else {
			if amount, err := strconv.ParseUint(req.NumSatoshis, 10, 32); err == nil {
				return uint(amount), nil
			} else {
				return 0, err
			}
		}
	} else {
		return 0, err
	}
}
