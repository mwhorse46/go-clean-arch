package repository

import (
	"database/sql"
	"fmt"

	author "github.com/bxcodec/go-clean-arch/models"

	"github.com/sirupsen/logrus"

	article "github.com/bxcodec/go-clean-arch/article"
	models "github.com/bxcodec/go-clean-arch/models"
)

type mysqlArticleRepository struct {
	Conn *sql.DB
}

func NewMysqlArticleRepository(Conn *sql.DB) article.ArticleRepository {

	return &mysqlArticleRepository{Conn}
}

func (m *mysqlArticleRepository) fetch(query string, args ...interface{}) ([]*models.Article, error) {

	rows, err := m.Conn.Query(query, args...)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer rows.Close()
	result := make([]*models.Article, 0)
	for rows.Next() {
		t := new(models.Article)
		authorID := int64(0)
		err = rows.Scan(
			&t.ID,
			&t.Title,
			&t.Content,
			&authorID,
			&t.UpdatedAt,
			&t.CreatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		t.Author = author.Author{
			ID: authorID,
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlArticleRepository) Fetch(cursor string, num int64) ([]*models.Article, error) {

	query := `SELECT id,title,content, author_id, updated_at, created_at
  						FROM article WHERE ID > ? LIMIT ?`

	return m.fetch(query, cursor, num)

}
func (m *mysqlArticleRepository) GetByID(id int64) (*models.Article, error) {
	query := `SELECT id,title,content, author_id, updated_at, created_at
  						FROM article WHERE ID = ?`

	list, err := m.fetch(query, id)
	if err != nil {
		return nil, err
	}

	a := &models.Article{}
	if len(list) > 0 {
		a = list[0]
	} else {
		return nil, models.NOT_FOUND_ERROR
	}

	return a, nil
}

func (m *mysqlArticleRepository) GetByTitle(title string) (*models.Article, error) {
	query := `SELECT id,title,content, author_id, updated_at, created_at
  						FROM article WHERE title = ?`

	list, err := m.fetch(query, title)
	if err != nil {
		return nil, err
	}

	a := &models.Article{}
	if len(list) > 0 {
		a = list[0]
	} else {
		return nil, models.NOT_FOUND_ERROR
	}
	return a, nil
}

func (m *mysqlArticleRepository) Store(a *models.Article) (int64, error) {

	query := `INSERT  article SET title=? , content=? , author_id=?, updated_at=? , created_at=?`
	stmt, err := m.Conn.Prepare(query)
	if err != nil {

		return 0, err
	}

	logrus.Debug("Created At: ", a.CreatedAt)
	res, err := stmt.Exec(a.Title, a.Content, a.Author.ID, a.UpdatedAt, a.CreatedAt)
	if err != nil {

		return 0, err
	}
	return res.LastInsertId()
}

func (m *mysqlArticleRepository) Delete(id int64) (bool, error) {
	query := "DELETE FROM article WHERE id = ?"

	stmt, err := m.Conn.Prepare(query)
	if err != nil {
		return false, err
	}
	res, err := stmt.Exec(id)
	if err != nil {

		return false, err
	}
	rowsAfected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	if rowsAfected != 1 {
		err = fmt.Errorf("Weird  Behaviour. Total Affected: %d", rowsAfected)
		logrus.Error(err)
		return false, err
	}

	return true, nil
}
func (m *mysqlArticleRepository) Update(ar *models.Article) (*models.Article, error) {
	query := `UPDATE article set title=?, content=?, author_id=?, updated_at=? WHERE ID = ?`

	stmt, err := m.Conn.Prepare(query)
	if err != nil {
		return nil, nil
	}

	res, err := stmt.Exec(ar.Title, ar.Content, ar.Author.ID, ar.UpdatedAt, ar.ID)
	if err != nil {
		return nil, err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affect != 1 {
		err = fmt.Errorf("Weird  Behaviour. Total Affected: %d", affect)
		logrus.Error(err)
		return nil, err
	}

	return ar, nil
}
