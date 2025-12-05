package repositories

import (
	"context"

	"github.com/barzaevhalid/cloud_storage_backend/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FileRepository struct {
	db *pgxpool.Pool
}

func NewFileRepository(db *pgxpool.Pool) *FileRepository {
	return &FileRepository{db: db}
}

func (r *FileRepository) Save(file *models.File) (int, error) {
	var id int
	err := r.db.QueryRow(
		context.Background(),
		`INSERT INTO files (user_id, filename, originalname, mimetype, size)
     VALUES($1, $2, $3, $4, $5)
     RETURNING id`,
		file.UserID, file.Filename, file.OriginalName, file.MimeType, file.Size,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil

}

func (s *FileRepository) FindAllFiles(userId int64, fileType string) ([]*models.File, error) {

	baseQuery := `SELECT id, filename, originalname, mimetype, size, user_id, deletedat
		FROM files
		WHERE user_id =$1`

	var query string

	switch fileType {
	case "image":
		query = baseQuery + ` AND mimetype ILIKE '%image%' AND deletedat IS NULL`
	case "trash":
		query = baseQuery + ` AND deletedat IS NOT NULL`
	default:
		query = baseQuery + ` AND deletedat IS NULL`
	}

	rows, err := s.db.Query(context.Background(), query, userId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := make([]*models.File, 0)

	for rows.Next() {
		var f models.File

		if err := rows.Scan(&f.ID, &f.Filename, &f.OriginalName, &f.MimeType, &f.Size, &f.UserID, &f.DeletedAt); err != nil {
			return nil, err
		}
		files = append(files, &f)
	}
	return files, nil

}

func (r *FileRepository) MarkDeleted(userId int64, ids []int64) error {
	query := `
		UPDATE files
		SET deletedat = NOW()
		WHERE user_id = $1 AND id = ANY($2)
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		userId,
		ids,
	)

	return err

}
