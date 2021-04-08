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