package wallet

import (
	"errors"
	"io"
	"strings"

	"log"
	"os"
	"strconv"

	"github.com/Muhammad-21/wallet/pkg/types"
	"github.com/google/uuid"
)


var ErrPhoneRegistered = errors.New("phone alredy registered")
var ErrAmmountMustBePositive = errors.New("ammount must be greater then zero")
var ErrAccountNotFound = errors.New("account not found")
var ErrNotEnoughBalance = errors.New("not enough balance ")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrFavoriteNotFound = errors.New("not found")

type Service struct{
	nextAccountID int64
	accounts 	[]*types.Account
	payments 	[]*types.Payment
	favorites   []*types.Favorite
}

func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
		return nil, ErrPhoneRegistered
		}
	}
	s.nextAccountID++
	account := &types.Account{
		ID:		s.nextAccountID,
		Phone: 	phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)
	
	return account, nil
}

func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error){
	payment, err :=s.FindPaymentByID(paymentID)
	if err != nil {
		return nil,err
	}
	favoriteID := uuid.New().String()
	favorite := &types.Favorite{
		ID: 		favoriteID,
		AccountID: 	payment.AccountID,
		Name: 		name,
		Amount: 	payment.Amount,
		Category: 	payment.Category,
	}
	s.favorites = append(s.favorites, favorite)
	return favorite,nil
}

func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error){
	var favorite *types.Favorite
	for _, fav := range s.favorites{
		if fav.ID == favoriteID {
			favorite=fav
			break
		}
	}
	if favorite==nil {
		return nil, ErrFavoriteNotFound				
	}
	new_paymentID := uuid.New().String()
	new_payment := &types.Payment{
		ID: 		new_paymentID,
		AccountID: 	favorite.AccountID,
		Amount: 	favorite.Amount,
		Category: 	favorite.Category,
		Status: 	types.PaymentStatusInProgress,
	}
	account, account_err := s.FindAccountByID(new_payment.AccountID)
	if account_err != nil {
		return nil, account_err
	}
	account.Balance-=new_payment.Amount
	s.payments = append(s.payments, new_payment)
	return new_payment, nil
	
}

func (s *Service) Deposit(accountID int64, ammount types.Money) error {
	if ammount <= 0 {
		return ErrAmmountMustBePositive
	}
	
	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
		account = acc
		break
	}
	}
	
	if account == nil {
		return ErrAccountNotFound
	}
	
	account.Balance += ammount
	return nil
}

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmmountMustBePositive
	}
	
	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
		account = acc
		break
	}
	}
	
	if account == nil {
		return nil, ErrAccountNotFound
	}
	
	if account.Balance < amount {
		return nil, ErrNotEnoughBalance
	}
	
	account.Balance -= amount
	paymentID := uuid.New().String()
	payment := &types.Payment{
		ID: 		paymentID,
		AccountID: 	accountID,
		Amount: 	amount,
		Category: 	category,
		Status: 	types.PaymentStatusInProgress,
	}
	
	s.payments = append(s.payments, payment)
	return payment, nil
	
}
	

func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	for _, account := range s.accounts{
		if account.ID == accountID {
			return account, nil
		}
	}
	return nil, ErrAccountNotFound
}


func (s *Service) Reject(paymentID string) error  {
	var payment_err *types.Payment
	for _, payment:=range s.payments{
		if payment.ID == paymentID{
			payment_err = payment
		}
	}
		if payment_err == nil {
			return ErrPaymentNotFound
		}
			payment_err.Status = types.PaymentStatusFail
			account, err := s.FindAccountByID(payment_err.AccountID)
			if err != nil{
				return nil
			}
			account.Balance+=payment_err.Amount
			return nil
}

func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	for _, payment := range s.payments{
		if payment.ID == paymentID {
			return payment, nil
		}
	}
	return nil, ErrPaymentNotFound
}

func (s *Service) Repeat(paymentID string) (*types.Payment, error)	{
	payment, err :=s.FindPaymentByID(paymentID)
	if err != nil {
		return nil,err
	}
	new_paymentID := uuid.New().String()
	new_payment := &types.Payment{
		ID: 		new_paymentID,
		AccountID: 	payment.AccountID,
		Amount: 	payment.Amount,
		Category: 	payment.Category,
		Status: 	payment.Status,
	}
	account, account_err := s.FindAccountByID(new_payment.AccountID)
	if account_err != nil {
		return nil, account_err
	}
	account.Balance-=new_payment.Amount
	s.payments = append(s.payments, new_payment)
	return new_payment, nil
}

func(s *Service) ExportToFile(path string) error{
	file, err := os.Create(path)
	if err != nil {
		log.Print(err)
	}
	defer func ()  {
		err := file.Close()
		if err != nil {
			log.Print(err)
			return
		}
	}()
	result := ""
	for _, account:= range s.accounts {
		collect := strconv.FormatInt(account.ID, 10) + ";" + string(account.Phone) + ";" +strconv.FormatInt(int64(account.Balance),10) + "|"
		result = result + collect
	}
	_, errr := file.Write([]byte(result))
	if err != nil {
		log.Print(errr)
	}
	return errr

}

func (s *Service) ImportFromFile(path string) error{
	file,err := os.Open(path)
	if err != nil{
		log.Print(err)
		return err
	}
	defer func(){
		err := file.Close()
		if err != nil {
			log.Print(err)
			return
		}
	} ()

	content := make([]byte,0)
	buf := make([]byte,4)
	for {
		read, err := file.Read(buf)
		if err == io.EOF {
			content = append(content, buf[:read]...)
			break
		}
		if err != nil {
			log.Print(err)
			return err
		}
		content = append(content, buf[:read]...)
	}
	data := strings.Split(string(content), "|")
	for _, accounts := range data {
		if len(accounts)>1{
		account := strings.Split(accounts, ";")
		id, err := strconv.ParseInt(account[0],10,64)
		if err != nil {
			log.Print(err)
		}
		balance,err := strconv.ParseInt(account[2],10,64)
		if err != nil {
			log.Print(err)
		}
		accountt := &types.Account{
			ID: id,
			Phone: types.Phone(account[1]),
			Balance: types.Money(balance),
		}
		s.accounts = append(s.accounts, accountt)
	}
}
	return err
}