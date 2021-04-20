package wallet

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/khiki1995/wallet/pkg/types"
)

var defaultTestAccount = testAccount{
	phone:   "+992902020102",
	balance: 10_000_00,
	payments: []struct {
		amount   types.Money
		category types.PaymentCategory
	}{
		{
			100_00,
			"ps4",
		},
	},
}

func TestService_FindPaymentByID_success(t *testing.T) {
	s := &Service{}
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	for _, payment := range payments {
		got, err := s.FindPaymentByID(payment.ID)
		if err != nil {
			t.Errorf("FindPaymentByID: cant find payment, error = %v", err)
			return
		}

		if !reflect.DeepEqual(payment, got) {
			t.Errorf("FindPaymentByID: got <> want, got = %v , want = %v", got, payment)
			return
		}
	}
}

func TestService_FindPaymentByID_fail(t *testing.T) {
	s := &Service{}
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment, err := s.FindPaymentByID(uuid.New().String())
	if err == nil {
		t.Errorf("FindPaymentByID: cant find payment, error = %v", err)
		return
	}

	if err != ErrPaymentNotFound {
		t.Errorf("FindPaymentByID: must return ErrPaymentNotFound, returned = %v", payment)
		return
	}
}

func TestSerivce_FindAccountByID(t *testing.T) {
	s := &Service{}
	account, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = s.FindAccountByID(account.ID)
	if err != nil {
		t.Error(err)
	}

	// fails
	account, err = s.FindAccountByID(int64(uuid.New().ID()))
	if err == nil {
		t.Errorf("FindAccountByID(): here should be error, but comes = %v", account)
	}
}

func TestService_Reject(t *testing.T) {
	s := &Service{}
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	for _, payment := range payments {
		err = s.Reject(payment.ID)
		if err != nil {
			t.Error(err)
		}

		pay, err := s.FindPaymentByID(payment.ID)
		if err != nil {
			t.Errorf("Reject(): can't find payment, error = %v", err)
			return
		}
		if pay.Status != types.PaymentStatusFail {
			t.Errorf("Reject(): status didn't changed, payment = %v", pay)
			return
		}

		account, err := s.FindAccountByID(payment.AccountID)
		if err != nil {
			t.Errorf("Reject(): can't find account by id, error = %v", err)
			return
		}
		if account.Balance != defaultTestAccount.balance {
			t.Errorf("Reject(): balance didn't changed, account = %v", account)
			return
		}
	}
}

func TestService_Repeat(t *testing.T) {
	s := &Service{}
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	for _, payment := range payments {
		err = s.Reject(payment.ID)
		if err != nil {
			t.Errorf("Repeat(): can't repeat payment, error = %v", err)
			return
		}
	}
	// fails
	err = s.Reject(uuid.NewString())
	if err == nil {
		t.Errorf("Repeat(): here should be error, but comes = %v", err)
		return
	}
}

func TestService_FavoritePayment(t *testing.T) {
	s := &Service{}
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	for _, payment := range payments {
		favorite, err := s.FavoritePayment(payment.ID, "test")
		if err != nil {
			t.Errorf("FavoritePayment(): can't add payment to favorite, error = %v", err)
			return
		}

		_, err = s.PayFromFavorite(favorite.ID)
		if err != nil {
			t.Errorf("PayFromFavorite(): can't create payment by favorite, error = %v", err)
			return
		}
	}
	// fails
	favor, err := s.FavoritePayment(uuid.NewString(), "test")
	if err == nil {
		t.Errorf("FavoritePayment(): here should be error, but comes = %v", favor)
		return
	}
	pay, err := s.PayFromFavorite(uuid.NewString())
	if err == nil {
		t.Errorf("PayFromFavorite(): here should be error, but comes = %v", pay)
		return
	}
}
