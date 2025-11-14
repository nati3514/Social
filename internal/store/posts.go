package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int32     `json:"version"`
	Comments  []Comment `json:"comments"`
	User      User      `json:"user"`
}

type PostWithMetadata struct {
	Post
	CommentCount int `json:"comment_count"`
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]PostWithMetadata, error) {
	query := `
    SELECT p.id, p.content, p.title, p.user_id, p.tags, p.created_at, p.updated_at, p.version, 
           COUNT(c.id) AS comment_count, u.username
    FROM posts p
    LEFT JOIN comments c ON p.id = c.post_id
    LEFT JOIN users u ON p.user_id = u.id
    WHERE p.user_id = $1 
       OR p.user_id IN (
           SELECT user_id
           FROM followers
           WHERE follower_id = $1
       )
    GROUP BY p.id, u.username, p.content, p.title, p.user_id, p.tags, p.created_at, p.updated_at, p.version
    ORDER BY p.created_at ` + fq.Sort + `
    LIMIT $2 OFFSET $3
    `
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userID, fq.Limit, fq.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feed []PostWithMetadata
	for rows.Next() {
		var p PostWithMetadata
		var username sql.NullString // Handle potential NULL usernames
		// Initialize the embedded Post and User structs
		p.Post = Post{}
		p.Post.User = User{}

		if err := rows.Scan(
			&p.ID,
			&p.Content,
			&p.Title,
			&p.UserID,
			pq.Array(&p.Tags),
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.Version,
			&p.CommentCount,
			&username, // Scan the username
		); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		// Set the username in the embedded User struct
		if username.Valid {
			p.Post.User.Username = username.String
		}

		feed = append(feed, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}
	return feed, nil
}

var (
	ErrEditConflict = errors.New("edit conflict: post has been modified by another user")
)

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
	INSERT INTO posts (content, title, user_id, tags) 
	VALUES  ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil

}

func (s *PostStore) GetByID(ctx context.Context, id int64) (*Post, error) {
	query := `
	SELECT id, content, title, user_id, tags, created_at, updated_at, version
	FROM posts
	WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	row := s.db.QueryRowContext(ctx, query, id)
	var post Post
	if err := row.Scan(
		&post.ID,
		&post.Content,
		&post.Title,
		&post.UserID,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Version,
	); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &post, nil
}

func (s *PostStore) Delete(ctx context.Context, postID int64) error {
	query := `
    DELETE FROM posts
    WHERE id = $1
    `
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, postID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *PostStore) Update(ctx context.Context, post *Post) error {
	// Validate the post
	if err := validatePost(post); err != nil {
		return err
	}

	query := `
        UPDATE posts 
        SET title = $1, content = $2, tags = $3, version = version + 1, updated_at = NOW() 
        WHERE id = $4 AND version = $5
        RETURNING version, updated_at
    `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	originalVersion := post.Version
	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		pq.Array(post.Tags),
		post.ID,
		originalVersion,
	).Scan(&post.Version, &post.UpdatedAt)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return fmt.Errorf("error updating post: %w", err)
		}
	}
	return nil
}

// Add this validation helper function
func validatePost(post *Post) error {
	if post.Title == "" {
		return errors.New("title is required")
	}
	if len(post.Title) > 100 {
		return errors.New("title must be less than 100 characters")
	}
	if post.Content == "" {
		return errors.New("content is required")
	}
	if len(post.Content) > 1000 {
		return errors.New("content must be less than 1000 characters")
	}
	return nil
}
