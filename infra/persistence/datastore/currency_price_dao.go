package datastore

import (
	"github.com/akihokurino/crypto-assistant-server/domain/models"
	"github.com/akihokurino/crypto-assistant-server/domain/repositories"
	"context"
	"github.com/mjibson/goon"
	"google.golang.org/appengine/datastore"
	"time"
	"github.com/pkg/errors"
	"github.com/akihokurino/crypto-assistant-server/utils"
)

const kindCurrencyPrice = "CurrencyPrice"

type CurrencyPriceDAO struct {
	_kind        string `goon:"kind,CurrencyPrice"`
	Id           string `datastore:"-" goon:"id"`
	CurrencyCode string
	USD          float64
	JPY          float64
	Datetime     time.Time
}

func (d *CurrencyPriceDAO) toModel() *models.CurrencyPrice {
	return &models.CurrencyPrice{
		Id:           models.CurrencyPriceID(d.Id),
		CurrencyCode: models.CurrencyCode(d.CurrencyCode),
		USD:          d.USD,
		JPY:          d.JPY,
		Datetime:     d.Datetime,
	}
}

func fromCurrencyPriceToDAO(from *models.CurrencyPrice) *CurrencyPriceDAO {
	return &CurrencyPriceDAO{
		Id:           string(from.Id),
		CurrencyCode: string(from.CurrencyCode),
		USD:          from.USD,
		JPY:          from.JPY,
		Datetime:     from.Datetime,
	}
}

type currencyPriceRepository struct {
	dateUtil utils.DateUtil
}

func NewCurrencyPriceRepository(dateUtil utils.DateUtil) repositories.CurrencyPriceRepository {
	return &currencyPriceRepository{
		dateUtil: dateUtil,
	}
}

func (r *currencyPriceRepository) GetByCurrency(
	ctx context.Context,
	code models.CurrencyCode) ([]*models.CurrencyPrice, error) {
	g := goon.FromContext(ctx)
	q := datastore.NewQuery(kindCurrencyPrice).Filter("CurrencyCode =", code)

	var currencyPriceDAOList []*CurrencyPriceDAO
	if _, err := g.GetAll(q, &currencyPriceDAOList); err != nil {
		return nil, errors.WithStack(err)
	}

	currencyPriceList := make([]*models.CurrencyPrice, len(currencyPriceDAOList))
	for i, v := range currencyPriceDAOList {
		currencyPriceList[i] = v.toModel()
	}

	return currencyPriceList, nil
}

func (r *currencyPriceRepository) GetLastByCurrency(ctx context.Context, code models.CurrencyCode) (*models.CurrencyPrice, error) {
	g := goon.FromContext(ctx)
	q := datastore.NewQuery(kindCurrencyPrice).
		Filter("CurrencyCode =", code).Order("-Datetime").Limit(1)

	var currencyPriceDAOList []*CurrencyPriceDAO
	if _, err := g.GetAll(q, &currencyPriceDAOList); err != nil {
		return nil, errors.WithStack(err)
	}

	currencyPriceList := make([]*models.CurrencyPrice, len(currencyPriceDAOList))
	for i, v := range currencyPriceDAOList {
		currencyPriceList[i] = v.toModel()
	}

	return currencyPriceList[0], nil
}

func (r *currencyPriceRepository) GetLast24HourByCurrency(ctx context.Context, code models.CurrencyCode) ([]*models.CurrencyPrice, error) {
	g := goon.FromContext(ctx)

	target := r.dateUtil.CurrentTime().Add(-24 * time.Hour)

	q := datastore.NewQuery(kindCurrencyPrice).
		Filter("CurrencyCode =", code).
		Filter("Datetime >=", target).
		Order("-Datetime")

	var currencyPriceDAOList []*CurrencyPriceDAO
	if _, err := g.GetAll(q, &currencyPriceDAOList); err != nil {
		return nil, errors.WithStack(err)
	}

	currencyPriceList := make([]*models.CurrencyPrice, len(currencyPriceDAOList))
	for i, v := range currencyPriceDAOList {
		currencyPriceList[i] = v.toModel()
	}

	return currencyPriceList, nil
}

func (r *currencyPriceRepository) Put(ctx context.Context, currencyPrice *models.CurrencyPrice) error {
	g := goon.FromContext(ctx)

	if _, err := g.Put(fromCurrencyPriceToDAO(currencyPrice)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
