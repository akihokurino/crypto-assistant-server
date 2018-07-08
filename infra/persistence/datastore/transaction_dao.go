package datastore

import (
	"github.com/akihokurino/crypto-assistant-server/domain/models"
	"github.com/mjibson/goon"
	"github.com/pkg/errors"
	"github.com/akihokurino/crypto-assistant-server/domain/repositories"
	"context"
	"google.golang.org/appengine/datastore"
)

const kindTransaction = "Transaction"

type TransactionDAO struct {
	_kind     string `goon:"kind,Transaction"`
	Id        string `datastore:"-" goon:"id"`
	UserId    string
	AddressId string
	Text      string
}

func onlyIdTransactionDAO(userId models.UserID, addressId models.AddressID) *TransactionDAO {
	return &TransactionDAO{Id: string(userId) + "-" + string(addressId)}
}

func (d *TransactionDAO) toModel() *models.Transaction {
	return &models.Transaction{
		UserId:    models.UserID(d.UserId),
		AddressId: models.AddressID(d.AddressId),
		Text:      d.Text,
	}
}

func fromTransactionToDAO(from *models.Transaction) *TransactionDAO {
	return &TransactionDAO{
		Id:        string(from.UserId) + "-" + string(from.AddressId),
		UserId:    string(from.UserId),
		AddressId: string(from.AddressId),
		Text:      from.Text,
	}
}

type transactionRepository struct {
}

func NewTransactionRepository() repositories.TransactionRepository {
	return &transactionRepository{}
}

func (r *transactionRepository) Put(ctx context.Context, transaction *models.Transaction) error {
	g := goon.FromContext(ctx)

	if _, err := g.Put(fromTransactionToDAO(transaction)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *transactionRepository) Delete(ctx context.Context, userId models.UserID, addressId models.AddressID) error {
	g := goon.FromContext(ctx)

	transactionDAO := onlyIdTransactionDAO(userId, addressId)

	key := datastore.NewKey(g.Context, kindTransaction, transactionDAO.Id, 0, nil)

	if err := g.Delete(key); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

