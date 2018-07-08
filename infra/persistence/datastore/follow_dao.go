package datastore

import (
	"time"
	"github.com/mjibson/goon"
	"github.com/pkg/errors"
	"google.golang.org/appengine/datastore"
	"context"
	"github.com/akihokurino/crypto-assistant-server/domain/models"
)

const kindFollow = "Follow"

type FollowDAO struct {
	_kind      string    `goon:"kind,Follow"`
	Id         string    `datastore:"-" goon:"id"`
	FromUserId string
	ToUserId   string
	CreatedAt  time.Time
}

func onlyIdFollowDAO(fromUserId models.UserID, toUserId models.UserID) *FollowDAO {
	return &FollowDAO{
		Id: string(fromUserId) + "-" + string(toUserId),
	}
}

func newFollowDAO(fromUserId models.UserID, toUserId models.UserID, now time.Time) *FollowDAO {
	return &FollowDAO{
		Id: string(fromUserId) + "-" + string(toUserId),
		FromUserId: string(fromUserId),
		ToUserId: string(toUserId),
		CreatedAt: now,
	}
}

func (r *userRepository) AlreadyFollowed(ctx context.Context, fromUserId models.UserID, toUserId models.UserID) (bool, error) {
	g := goon.FromContext(ctx)

	followDAO := onlyIdFollowDAO(fromUserId, toUserId)

	err := g.Get(followDAO)

	if err == nil {
		return true, nil
	}

	if err == datastore.ErrNoSuchEntity {
		return false, nil
	} else {
		return false, errors.WithStack(err)
	}
}

func (r *userRepository) Follow(ctx context.Context, fromUserId models.UserID, toUserId models.UserID) error {
	g := goon.FromContext(ctx)

	followDAO := newFollowDAO(fromUserId, toUserId, r.dateUtil.CurrentTime())

	if _, err := g.Put(followDAO); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *userRepository) UnFollow(ctx context.Context, fromUserId models.UserID, toUserId models.UserID) error {
	g := goon.FromContext(ctx)

	followDAO := onlyIdFollowDAO(fromUserId, toUserId)

	key := datastore.NewKey(g.Context, kindFollow, followDAO.Id, 0, nil)

	if err := g.Delete(key); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *userRepository) GetFollows(ctx context.Context, userId models.UserID) ([]*models.User, error) {
	g := goon.FromContext(ctx)
	q := datastore.NewQuery(kindFollow).Filter("FromUserId =", userId)

	var followDAOList []*FollowDAO

	if _, err := g.GetAll(q, &followDAOList); err != nil {
		return nil, errors.WithStack(err)
	}

	userDAOList := make([]*UserDAO, len(followDAOList))

	for i, v := range followDAOList {
		userDAOList[i] = onlyIdUserDAO(models.UserID(v.ToUserId))
	}

	if err := g.GetMulti(&userDAOList); err != nil {
		return nil, errors.WithStack(err)
	}

	userList := make([]*models.User, len(userDAOList))

	for i, v := range userDAOList {
		userList[i] = v.toModel()
	}

	return userList, nil
}

func (r *userRepository) GetFollowers(ctx context.Context, userId models.UserID) ([]*models.User, error) {
	g := goon.FromContext(ctx)
	q := datastore.NewQuery(kindFollow).Filter("ToUserId =", userId)

	var followDAOList []*FollowDAO

	if _, err := g.GetAll(q, &followDAOList); err != nil {
		return nil, errors.WithStack(err)
	}

	userDAOList := make([]*UserDAO, len(followDAOList))

	for i, v := range followDAOList {
		userDAOList[i] = onlyIdUserDAO(models.UserID(v.FromUserId))
	}

	if err := g.GetMulti(&userDAOList); err != nil {
		return nil, errors.WithStack(err)
	}

	userList := make([]*models.User, len(userDAOList))

	for i, v := range userDAOList {
		userList[i] = v.toModel()
	}

	return userList, nil
}