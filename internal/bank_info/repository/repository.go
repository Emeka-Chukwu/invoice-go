package bankinfo_repository

import (
	"context"
	"database/sql"
	"go-invoice/domain"
	"go-invoice/util"
)

type bankInfoRepository struct {
	Db *sql.DB
}

type BankInfoRepository interface {
	CreateBankInformation(req domain.CreateBankInformationDTO) (int64, error)
	FetchBankInformation(userId int64) (domain.BankInformation, error)
}

func NewBankInfoRepository(db *sql.DB) BankInfoRepository {
	return &bankInfoRepository{Db: db}
}

func (br *bankInfoRepository) CreateBankInformation(req domain.CreateBankInformationDTO) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), util.DbTimeout)
	defer cancel()
	query := `
    INSERT INTO bank_information (account_name, account_number, ach_routing_no, bank_name, bank_address,user_id, created_at, updated_at) 
    VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) 
    RETURNING id`
	var id int64
	err := br.Db.QueryRowContext(ctx, query, req.AccountName, req.AccountNumber, req.ACHRoutingNo, req.BankName, req.BankAddress, req.UserID).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (br *bankInfoRepository) FetchBankInformation(userId int64) (domain.BankInformation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), util.DbTimeout)
	defer cancel()
	query := `SELECT id, account_name, account_number, ach_routing_no, bank_name, bank_address,user_id, created_at, updated_at 
              FROM bank_information WHERE user_id = $1`

	var bankInfo domain.BankInformation
	err := br.Db.QueryRowContext(ctx, query, userId).Scan(
		&bankInfo.ID, &bankInfo.AccountName, &bankInfo.AccountNumber,
		&bankInfo.ACHRoutingNo, &bankInfo.BankName, &bankInfo.BankAddress,
		&bankInfo.UserID, &bankInfo.CreatedAt, &bankInfo.UpdatedAt)
	if err != nil {
		return bankInfo, err
	}
	return bankInfo, nil
}
