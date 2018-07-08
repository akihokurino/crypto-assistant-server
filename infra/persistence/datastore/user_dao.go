package datastore

import (
	"github.com/akihokurino/crypto-assistant-server/domain/models"
	"net/url"
	"github.com/akihokurino/crypto-assistant-server/domain/repositories"
	"github.com/mjibson/goon"
	"google.golang.org/appengine/datastore"
	"context"
	"github.com/pkg/errors"
	"github.com/akihokurino/crypto-assistant-server/utils"
)

const kindUser = "User"

type UserDAO struct {
	_kind    string `goon:"kind,User"`
	Id       string `datastore:"-" goon:"id"`
	Username string
	IconURL  string
}

func onlyIdUserDAO(id models.UserID) *UserDAO {
	return &UserDAO{Id: string(id)}
}

func (d *UserDAO) toModel() *models.User {
	u, _ := url.Parse(d.IconURL)

	return &models.User{
		Id:       models.UserID(d.Id),
		Username: d.Username,
		IconURL:  u,
	}
}

func fromUserToDAO(from *models.User) *UserDAO {
	u := ""
	if from.IconURL != nil {
		u = from.IconURL.String()
	}

	return &UserDAO{
		Id:       string(from.Id),
		Username: from.Username,
		IconURL:  u,
	}
}

type userRepository struct {
	dateUtil utils.DateUtil
}

func NewUserRepository(dateUtil utils.DateUtil) repositories.UserRepository {
	return &userRepository{
		dateUtil: dateUtil,
	}
}

func (r *userRepository) GetAll(ctx context.Context) ([]*models.User, error) {
	g := goon.FromContext(ctx)
	q := datastore.NewQuery(kindUser)

	var userDAOList []*UserDAO
	if _, err := g.GetAll(q, &userDAOList); err != nil {
		return nil, errors.WithStack(err)
	}

	userList := make([]*models.User, len(userDAOList))
	for i, v := range userDAOList {
		userList[i] = v.toModel()
	}

	return userList, nil
}

func (r *userRepository) Get(ctx context.Context, id models.UserID) (*models.User, error) {
	g := goon.FromContext(ctx)

	userDAO := onlyIdUserDAO(id)

	if err := g.Get(userDAO); err != nil {
		return nil, errors.WithStack(err)
	}

	return userDAO.toModel(), nil
}

func (r *userRepository) Put(ctx context.Context, user *models.User) error {
	g := goon.FromContext(ctx)

	if _, err := g.Put(fromUserToDAO(user)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
