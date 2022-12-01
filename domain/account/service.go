package account

import (
	"context"

	"github.com/google/uuid"
	"github.com/iowanobos/postgres-client/postgres"
	pq "github.com/iowanobos/postgres-client/postgres/query"
	"github.com/sirupsen/logrus"
)

type Service struct {
	txManager  postgres.TransactionManager
	repository Repository
}

func NewService(
	txManager postgres.TransactionManager,
	repository Repository,
) *Service {
	return &Service{
		txManager:  txManager,
		repository: repository,
	}
}

func (s *Service) Test(ctx context.Context) error {
	ctx, err := s.txManager.BeginTx(ctx)
	if err != nil {
		logrus.WithError(err).Error("begin transaction failed")
		return err
	}

	defer func() {
		err := s.txManager.Rollback(ctx)
		if err != nil {
			logrus.WithError(err).Error("rollback transaction failed")
			return
		}

		accounts, err := s.repository.List(context.Background(), pq.NewListQuery(nil))
		if err != nil {
			logrus.WithError(err).Error("get account list failed")
			return
		}

		logrus.WithField("count", accounts.TotalCount).Info("accounts list result")
	}()

	count := 10

	if err := s.repository.BatchCreate(ctx, generateAccounts(count)); err != nil {
		logrus.WithError(err).Error("create account list failed")
		return err
	}

	accounts, err := s.repository.List(ctx, pq.NewListQuery(nil).AddSort("name", true).WithIteration(pq.Pagination{Number: 1, Size: 3}))
	if err != nil {
		logrus.WithError(err).Error("get account list failed")
		return err
	}

	logrus.WithField("count", accounts.TotalCount).Info("accounts list result in transaction")
	return nil
}

func generateAccounts(count int) []Account {
	accounts := make([]Account, count)
	for i := 0; i < count; i++ {
		accounts[i].Name = uuid.New().String()
	}
	return accounts
}
