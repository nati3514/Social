package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound  = errors.New("recored not found")
	ErrConflict  = errors.New("record already exists")
	QueryTimeout = time.Second * 5
)

type Storage struct {
	Posts interface {
		GetByID(context.Context, int64) (*Post, error)
		Create(context.Context, *Post) error
		Update(context.Context, *Post) error
		Delete(context.Context, int64) error
		GetUserFeed(context.Context, int64) ([]PostWithMetadata, error)
	}

	Users interface {
		GetByID(context.Context, int64) (*User, error)
		Create(context.Context, *User) error
		CreateAndInvite(ctx context.Context, user *User, token string) error
	}
	Comments interface {
		Create(context.Context, *Comment) error
		GetByPostsID(context.Context, int64) ([]Comment, error)
	}
	Followers interface {
		Follow(ctx context.Context, followerID, userID int64) error
		Unfollow(ctx context.Context, follwerID, userID int64) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostStore{db},
		Users:     &UserStore{db},
		Comments:  &CommentStore{db},
		Followers: &FollowerStore{db},
	}
}
