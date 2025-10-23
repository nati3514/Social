package store 

import (
	"database/sql"
	"context"
)

type Storage struct {
	
	Posts interface {
	  Create(context.Context, *Post) error 
     }
	
	Users interface {
		Create(context.Context, *User) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage {
		Posts: &PostsStore{db},
		Users: &UsersStore{db},
		
	}
}

