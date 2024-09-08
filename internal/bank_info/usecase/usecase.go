package bankinfo_usecase

import (
	"go-invoice/domain"
	bankinfo_repository "go-invoice/internal/bank_info/repository"
	"net/http"
)

type BankInfoUsecase interface {
	CreateBankInformation(req domain.CreateBankInformationDTO) (int, int64, error)
	FetchBankInformation(userId int) (int, domain.BankInformation, error)
}

type bankInfoUsecase struct {
	Repo bankinfo_repository.BankInfoRepository
}

func NewBankInfoUsecase(Repo bankinfo_repository.BankInfoRepository) BankInfoUsecase {
	return bankInfoUsecase{Repo: Repo}
}

// FetchActivitieslog implements ActivityUsecase.
func (a bankInfoUsecase) CreateBankInformation(req domain.CreateBankInformationDTO) (int, int64, error) {
	resp, err := a.Repo.CreateBankInformation(req)
	if err != nil {
		return http.StatusInternalServerError, 0, err
	}
	return http.StatusOK, resp, nil
}

// FetchActivityByID implements ActivityUsecase.
func (a bankInfoUsecase) FetchBankInformation(userId int) (int, domain.BankInformation, error) {
	resp, err := a.Repo.FetchBankInformation(int64(userId))
	if err != nil {
		return http.StatusInternalServerError, domain.BankInformation{}, err
	}
	return http.StatusOK, resp, nil
}
