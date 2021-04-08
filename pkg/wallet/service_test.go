package wallet

import (
	"fmt"
	"testing"
)


func TestService_FindAccountByID_possitive(t *testing.T) {
	svc := &Service{}
	account,err := svc.RegisterAccount("+79888888888")
	if err != nil {
		fmt.Println(err)
		return
	}
	accounts, err := svc.FindAccountByID(account.ID)
	if err != nil{
		if account != accounts {
			t.Error(err)
		}
	}
}

func TestService_FindAccountByID_negative(t *testing.T)  {
	svc := &Service{}
	account,err := svc.RegisterAccount("+79888888888")
	if err != nil {
		fmt.Println(err)
		return
	}
	accounts, err := svc.FindAccountByID(account.ID+1)
	if err != nil{
		if err != ErrAccountNotFound{
			t.Error(accounts)
		}
	}
}

func TestService_Reject_found(t *testing.T) {
	svc := &Service{}
	account,err := svc.RegisterAccount("+79888888888")
	if err != nil {
		fmt.Println(err)
		return
	}
	account.Balance=100
	payment, er := svc.Pay(account.ID, 10, "aa")
	err = svc.Reject(payment.ID)
	if err != nil {
		fmt.Println(er)
	}
}


func TestService_Reject_notfound(t *testing.T) {
	svc := &Service{}
	
	err:= svc.Reject("1")
	if err == nil {
		t.Error(ErrPaymentNotFound)
		return 
	}

}