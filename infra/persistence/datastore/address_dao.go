package datastore

import (
	"github.com/akihokurino/crypto-assistant-server/domain/models"
	"github.com/akihokurino/crypto-assistant-server/domain/repositories"
	"context"
	"github.com/mjibson/goon"
	"google.golang.org/appengine/datastore"
	"github.com/pkg/errors"
)

const kindAddress = "Address"

type AddressDAO struct {
	_kind        string `goon:"kind,Address"`
	Id           string `datastore:"-" goon:"id"`
	UserId       string
	CurrencyCode string
	Value        string
}

func onlyIdAddressDAO(id models.AddressID) *AddressDAO {
	return &AddressDAO{Id: string(id)}
}

func (d *AddressDAO) toModel() *models.Address {
	return &models.Address{
		Id:           models.AddressID(d.Id),
		UserId:       models.UserID(d.UserId),
		CurrencyCode: models.CurrencyCode(d.CurrencyCode),
		Value:        d.Value,
	}
}

func fromAddressToDAO(from *models.Address) *AddressDAO {
	return &AddressDAO{
		Id:           string(from.Id),
		UserId:       string(from.UserId),
		CurrencyCode: string(from.CurrencyCode),
		Value:        from.Value,
	}
}

type addressRepository struct {
}

func NewAddressRepository() repositories.AddressRepository {
	return &addressRepository{}
}

func (r *addressRepository) GetByUser(ctx context.Context, userId models.UserID) ([]*models.Address, error) {
	g := goon.FromContext(ctx)
	q := datastore.NewQuery(kindAddress).Filter("UserId =", userId)

	var addressDAOList []*AddressDAO
	if _, err := g.GetAll(q, &addressDAOList); err != nil {
		return nil, errors.WithStack(err)
	}

	addressList := make([]*models.Address, len(addressDAOList))
	for i, v := range addressDAOList {
		addressList[i] = v.toModel()
	}

	return addressList, nil
}

func (r *addressRepository) Get(ctx context.Context, addressId models.AddressID) (*models.Address, error) {
	g := goon.FromContext(ctx)

	addressDAO := onlyIdAddressDAO(addressId)

	if err := g.Get(addressDAO); err != nil {
		return nil, errors.WithStack(err)
	}

	return addressDAO.toModel(), nil
}

func (r *addressRepository) ExistAddress(ctx context.Context, addressText string) (bool, error) {
	g := goon.FromContext(ctx)
	q := datastore.NewQuery(kindAddress).Filter("Value =", addressText)

	var addressDAOList []*AddressDAO

	if _, err := g.GetAll(q, &addressDAOList); err != nil {
		return false, errors.WithStack(err)
	}

	return len(addressDAOList) > 0, nil
}

func (r *addressRepository) Put(ctx context.Context, address *models.Address) error {
	g := goon.FromContext(ctx)

	if _, err := g.Put(fromAddressToDAO(address)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *addressRepository) Delete(ctx context.Context, addressId models.AddressID) error {
	g := goon.FromContext(ctx)

	addressDAO := onlyIdAddressDAO(addressId)

	key := datastore.NewKey(g.Context, kindAddress, addressDAO.Id, 0, nil)

	if err := g.Delete(key); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
