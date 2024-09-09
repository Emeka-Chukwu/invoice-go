package bankinfo_repository

import (
	"go-invoice/domain"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestCreateBankInformation(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	bankRepo := NewBankInfoRepository(db)

	req := domain.CreateBankInformationDTO{
		AccountName:   "John Doe",
		AccountNumber: "1234567890",
		ACHRoutingNo:  "111000025",
		BankName:      "First Bank",
		BankAddress:   "123 Bank St.",
		UserID:        101,
	}

	query := `
    INSERT INTO bank_information \(account_name, account_number, ach_routing_no, bank_name, bank_address,user_id, created_at, updated_at\) 
    VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP\) 
    RETURNING id`
	mock.ExpectQuery(query).
		WithArgs(req.AccountName, req.AccountNumber, req.ACHRoutingNo, req.BankName, req.BankAddress, req.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	id, err := bankRepo.CreateBankInformation(req)

	require.NoError(t, err)
	require.Equal(t, int64(1), id)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFetchBankInformation(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	bankRepo := NewBankInfoRepository(db)

	query := `SELECT id, account_name, account_number, ach_routing_no, bank_name, bank_address,user_id, created_at, updated_at 
              FROM bank_information WHERE user_id = \$1`
	rows := sqlmock.NewRows([]string{"id", "account_name", "account_number", "ach_routing_no", "bank_name", "bank_address", "user_id", "created_at", "updated_at"}).
		AddRow(1, "John Doe", "1234567890", "111000025", "First Bank", "123 Bank St.", 101, time.Now(), time.Now())

	mock.ExpectQuery(query).
		WithArgs(101).
		WillReturnRows(rows)

	bankInfo, err := bankRepo.FetchBankInformation(int64(101))
	require.NoError(t, err)
	require.Equal(t, int64(1), bankInfo.ID)
	require.Equal(t, "John Doe", bankInfo.AccountName)
	require.Equal(t, "1234567890", bankInfo.AccountNumber)
	require.Equal(t, "111000025", bankInfo.ACHRoutingNo)
	require.Equal(t, "First Bank", bankInfo.BankName)
	require.Equal(t, "123 Bank St.", bankInfo.BankAddress)
	require.Equal(t, 101, bankInfo.UserID)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
