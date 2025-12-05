package services

import (
	"fmt"
	"mime/multipart"

	"github.com/barzaevhalid/cloud_storage_backend/models"
	"github.com/barzaevhalid/cloud_storage_backend/repositories"
)

type FileService struct {
	FileRepository *repositories.FileRepository
}

func NewFileService(r *repositories.FileRepository) *FileService {
	return &FileService{
		FileRepository: r,
	}
}

func (s *FileService) SaveFileMetadata(userId int64, fileHeader *multipart.FileHeader, filename string) (int, error) {

	file := &models.File{
		UserID:       userId,
		Filename:     filename,
		OriginalName: fileHeader.Filename,
		MimeType:     fileHeader.Header.Get("Content-Type"),
		Size:         fileHeader.Size,
	}
	id, err := s.FileRepository.Save(file)

	if err != nil {
		return 0, fmt.Errorf("Canot save file : %w", err)
	}
	return id, nil
}

func (s *FileService) FindAllFiles(userId int64, fileType string) ([]*models.File, error) {

	files, err := s.FileRepository.FindAllFiles(userId, fileType)

	if err != nil {
		return nil, err
	}

	return files, nil

}

func (s *FileService) MarkDeleted(userId int64, ids []int64) error {
	return s.FileRepository.MarkDeleted(userId, ids)
}
