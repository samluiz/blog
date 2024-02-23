package comment

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/samluiz/blog/pkg/article"
	"github.com/samluiz/blog/pkg/types"
	"github.com/samluiz/blog/pkg/user"
)

type Repository interface {
	FindCommentsByArticleId(articleId int) ([]*types.Comment, error)
	FindCommentById(id int) (*types.Comment, error)
	FindCommentsByUserId(userId int) ([]*types.Comment, error)
	CreateComment(input *types.CreateCommentInput) (*types.Comment, error)
	UpdateComment(id int, input *types.UpdateCommentInput) (*types.Comment, error)
	DeleteComment(id int) error
	DeleteCommentsByArticleId(articleId int) error
	CommentExists(id int) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db}
}

func (r *repository) FindCommentsByArticleId(articleId int) ([]*types.Comment, error) {
	var comments []*types.Comment
	err := r.db.Select(&comments, "SELECT * FROM comments WHERE article_id = ?", articleId)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *repository) FindCommentById(id int) (*types.Comment, error) {
	var comment types.Comment
	err := r.db.Get(&comment, "SELECT * FROM comments WHERE id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrCommentNotFound
		}
		return nil, err
	}
	return &comment, nil
}

func (r *repository) FindCommentsByUserId(userId int) ([]*types.Comment, error) {

	userRepo := user.NewRepository(r.db)

	if err := userRepo.UserExistsById(userId); err != nil {
		return nil, err
	}

	var comments []*types.Comment
	err := r.db.Select(&comments, "SELECT * FROM comments WHERE author_id = ?", userId)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *repository) CreateComment(input *types.CreateCommentInput) (*types.Comment, error) {

	userRepo := user.NewRepository(r.db)

	if err := userRepo.UserExistsById(input.AuthorID); err != nil {
		return nil, err
	}

	var comment types.Comment
	res := r.db.MustExec("INSERT INTO comments (author_id, article_id, content) VALUES (?, ?, ?) RETURNING *", input.AuthorID, input.ArticleID, input.Content)

	idCreated, err := res.LastInsertId()

	if err != nil {
		return nil, err
	}

	err = r.db.Get(&comment, "SELECT * FROM comments WHERE id = ?", idCreated)

	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (r *repository) UpdateComment(id int, input *types.UpdateCommentInput) (*types.Comment, error) {
	var comment types.Comment
	res := r.db.MustExec("UPDATE comments SET content = ? WHERE id = ? RETURNING *", input.Content, id)

	_, err := res.RowsAffected()

	if err != nil {
		return nil, err
	}

	err = r.db.Get(&comment, "SELECT * FROM comments WHERE id = ?", id)

	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (r *repository) DeleteComment(id int) error {

	if err := r.CommentExists(id); err != nil {
		return err
	}

	_, err := r.db.Exec("DELETE FROM comments WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) DeleteCommentsByArticleId(articleId int) error {

	articleRepo := article.NewRepository(r.db)

	if err := articleRepo.ArticleExists(articleId); err != nil {
		if err == sql.ErrNoRows {
			return types.ErrArticleNotFound
		}
		return err
	}

	_, err := r.db.Exec("DELETE FROM comments WHERE article_id = ?", articleId)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) CommentExists(id int) error {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM comments WHERE id = ?", id)

	if err != nil {
		if err == sql.ErrNoRows {
			return types.ErrCommentNotFound
		}
		return err
	}

	return nil
}
