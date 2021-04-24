package wallet

import (
	"fmt"
	"testing"

	"github.com/Muhammad-21/wallet/pkg/types"
)

type testService struct {
	*Service
}
	
type testAccount struct {
	phone types.Phone
	balance types.Money
	payments []struct {
	amount types.Money
	category types.PaymentCategory
}
}
	
var defaultTestAccount=testAccount {
	phone: "+7999999999",
	balance: 100,
	payments: []struct{
	amount types.Money
	category types.PaymentCategory
	}{{100, "auto"},
	},
}
	
func newTestService() *testService {
	return &testService{Service: &Service{}}
}

func (s *testService) addAccount(data testAccount) (*types.Account, []*types.Payment, error) {
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("can`t register account, erro = %v", err)
	}
	
	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("can`t deposit account, error = %v", err)
	}
	
	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("can`t make payment, error = %v", err)
		}
	}
	
	return account, payments, nil
}
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
	
	errr := svc.Deposit(account.ID,100)
	if errr != nil {
		fmt.Println(errr)
		return
	}
	payment, er := svc.Pay(account.ID, 10, "auto")
	if er != nil {
		fmt.Println(er)
	}
	errrr := svc.Reject(payment.ID)
	if errrr != nil {
		fmt.Println(errrr)
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


func TestService_Repeat_found(t *testing.T) {
	svc := &Service{}
	account,err := svc.RegisterAccount("+79888888888")
	if err != nil {
		fmt.Println(err)
		return
	}
	errr := svc.Deposit(account.ID,100)
	if errr != nil {
		fmt.Println(errr)
		return
	}
	payment, er := svc.Pay(account.ID, 10, "auto")
	if er != nil {
		fmt.Println(er)
	}
	_,errrr := svc.Repeat(payment.ID)
	if errrr != nil {
		fmt.Println(errrr)
	}
}

func TestService_FavoritePayment_ok(t *testing.T) {
	s := newTestService()
	_,payments,err:=s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return 
	}
	
	payment:=payments[0]
	_, err=s.FavoritePayment(payment.ID,"auto")
	if err != nil {
		fmt.Println(err)
	return 
	}
}

func TestService_PayFromFavorite_ok(t *testing.T) {
	s := newTestService()
	_,payments,err:=s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
	return 
	}
	
	payment:=payments[0]
	fv, err:=s.FavoritePayment(payment.ID,"auto")
	if err != nil {
		fmt.Println(err)
	return 
	}
	
	_,err=s.PayFromFavorite(fv.ID)
	if err != nil {
		fmt.Println()
	return 
	}
	
	savedAccount, err:=s.FindAccountByID(payment.AccountID)
	if err != nil {
		fmt.Println(err)
	return
	}
	if savedAccount.Balance==defaultTestAccount.balance{
		fmt.Println(savedAccount)
	return
	}
	
}



func TestService_Import_success_user(t *testing.T) {
	var svc Service

	err := svc.ImportFromFile("export.txt")

	if err != nil {
		t.Errorf("method ExportToFile returned not nil error, err => %v", err)
	}

}


func TestService_Export_success(t *testing.T) {
	svc := Service{}

	svc.RegisterAccount("+992000000001")
	svc.RegisterAccount("+992000000002")
	svc.RegisterAccount("+992000000003")
	svc.RegisterAccount("+992000000004")

	err := svc.Export("data")
	if err != nil {
		t.Errorf("method ExportToFile returned not nil error, err => %v", err)
	}

	err = svc.Import("data")
	if err != nil {
		t.Errorf("method ExportToFile returned not nil error, err => %v", err)
	}
}

func TestService_ExportHistory_success_user(t *testing.T) {
	svc := Service{}

	acc, err := svc.RegisterAccount("+992000000001")

	if err != nil {
		t.Errorf("method RegisterAccount returned not nil error, account => %v", acc)
	}

	err = svc.Deposit(acc.ID, 100_00)
	if err != nil {
		t.Errorf("method Deposit returned not nil error, error => %v", err)
	}

	_, err = svc.Pay(acc.ID, 1, "Cafe")
	if err != nil {
		t.Errorf("method Pay returned not nil error, err => %v", err)
	}
	_, err = svc.Pay(acc.ID, 2, "Auto")
	if err != nil {
		t.Errorf("method Pay returned not nil error, err => %v", err)
	}
	_, err = svc.Pay(acc.ID, 3, "MarketShop")
	if err != nil {
		t.Errorf("method Pay returned not nil error, err => %v", err)
	}
	

	payments, err := svc.ExportAccountHistory(acc.ID)
	if err != nil {
		t.Errorf("method ExportAccountHistory returned not nil error, err => %v", err)
	}

	err = svc.HistoryToFiles(payments, "../../data", 2)
	if err != nil {
		t.Errorf("method HistoryToFiles returned not nil error, err => %v", err)
	}
}


func BenchmarkSumPayments(b *testing.B) {
	// svc:=&Service{}
	var svc Service
	account, err := svc.RegisterAccount("+992927777777")
	if err != nil {
		b.Errorf("account => %v",account)
	}
	err = svc.Deposit(account.ID, 100_00)
	if err != nil {
		b.Errorf("error => %v", err)
	}
	want:=types.Money(55)
	for i := types.Money(1); i <= 10; i++ {
		_, err := svc.Pay(account.ID, i, "aa")
		if  err != nil {
			b.Errorf("error => %v", err)
		}
	}
	got:=svc.SumPayments(5)
	if want != got{
		b.Errorf("want => %v got => %v", want, got)
	}
}

// func BenchmarkFilterPayments(b *testing.B) {
// 	// svc:=&Service{}
// 	var svc Service
// 	account, err := svc.RegisterAccount("+992927777777")
// 	if err != nil {
// 		b.Errorf("account => %v",account)
// 	}
// 	err = svc.Deposit(account.ID, 100_00)
// 	if err != nil {
// 		b.Errorf("error => %v", err)
// 	}
// 	//want:=types.Money(55)
// 	for i := types.Money(1); i <= 10; i++ {
// 		_, err := svc.Pay(account.ID, i, "aa")
// 		if  err != nil {
// 			b.Errorf("error => %v", err)
// 		}
// 	}
// 	got,err:=svc.FilterPayments(1,5)
// 	if err != nil{
// 		b.Errorf("got => %v",got)
// 	}
// }