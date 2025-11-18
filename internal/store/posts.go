package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
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
	// Build the query with dynamic conditions
	query := `
    SELECT 
        p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags, 
        u.username,
        COUNT(DISTINCT c.id) AS comment_count
    FROM posts p
    LEFT JOIN comments c ON c.post_id = p.id
    LEFT JOIN users u ON p.user_id = u.id
    LEFT JOIN followers f ON f.user_id = p.user_id
    WHERE 
        (f.follower_id = $1 OR p.user_id = $1) 
    `

	// Track parameters for the query
	var params []interface{}
	params = append(params, userID)

	// Add search condition if needed
	searchTerm := ""
	if fq.Search != "" {
		query += ` AND (p.title ILIKE $` + strconv.Itoa(len(params)+1) + ` OR p.content ILIKE $` + strconv.Itoa(len(params)+1) + `)`
		searchTerm = "%" + fq.Search + "%"
		params = append(params, searchTerm)
	}

	// Add tags condition if needed
	if len(fq.Tags) > 0 {
		// Use EXISTS with LIKE for partial tag matching
		for _, tag := range fq.Tags {
			query += ` AND EXISTS (
				SELECT 1 FROM unnest(p.tags) AS t 
				WHERE t ILIKE $` + strconv.Itoa(len(params)+1) + `
			)`
			params = append(params, tag)
		}
	}

	// Add GROUP BY, ORDER BY, and LIMIT/OFFSET
	query += `
    GROUP BY p.id, u.username, p.user_id, p.title, p.content, p.created_at, p.version, p.tags
    ORDER BY p.created_at ` + fq.Sort + `
    LIMIT ` + strconv.Itoa(fq.Limit) + ` OFFSET ` + strconv.Itoa(fq.Offset) + `
    `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	// Log the query and parameters for debugging
	log.Printf("Executing query:\n%s\nWith params: %+v\n", query, params)

	// Execute the query with parameters
	rows, err := s.db.QueryContext(ctx, query, params...)
	if err != nil {
		log.Printf("Query execution error: %v\nQuery: %s\nParams: %+v", err, query, params)
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var posts []PostWithMetadata
	for rows.Next() {
		var post PostWithMetadata
		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.Version,
			pq.Array(&post.Tags),
			&post.User.Username,
			&post.CommentCount,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return posts, nil
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
