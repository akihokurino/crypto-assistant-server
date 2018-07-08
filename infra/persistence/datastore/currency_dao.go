package datastore

import (
	"github.com/akihokurino/crypto-assistant-server/domain/repositories"
	"context"
	"github.com/akihokurino/crypto-assistant-server/domain/models"
	"github.com/mjibson/goon"
	"google.golang.org/appengine/datastore"
	"github.com/pkg/errors"
)

const kindCurrency = "Currency"

type CurrencyDAO struct {
	_kind string `goon:"kind,Currency"`
	Code  string `datastore:"-" goon:"id"`
	Name  string
}

func (d *CurrencyDAO) toModel() *models.Currency {
	return &models.Currency{
		Code: models.CurrencyCode(d.Code),
		Name: d.Name,
	}
}

func fromCurrencyToDAO(from *models.Currency) *CurrencyDAO {
	return &CurrencyDAO{
		Code: string(from.Code),
		Name: from.Name,
	}
}

type currencyRepository struct {
}

func NewCurrencyRepository() repositories.CurrencyRepository {
	return &currencyRepository{}
}

func (r *currencyRepository) GetAll(ctx context.Context) ([]*models.Currency, error) {
	g := goon.FromContext(ctx)
	q := datastore.NewQuery(kindCurrency)

	var currencyDAOList []*CurrencyDAO
	if _, err := g.GetAll(q, &currencyDAOList); err != nil {
		return nil, errors.WithStack(err)
	}

	currencyList := make([]*models.Currency, len(currencyDAOList))
	for i, v := range currencyDAOList {
		currencyList[i] = v.toModel()
	}

	return currencyList, nil
}

func (r *currencyRepository) Put(ctx context.Context, currency *models.Currency) error {
	g := goon.FromContext(ctx)

	if _, err := g.Put(fromCurrencyToDAO(currency)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
