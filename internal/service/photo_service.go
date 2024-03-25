package service

import (
	"context"

	"github.com/geedotrar/mygram/internal/model"
	"github.com/geedotrar/mygram/internal/repository"
)

type PhotoService interface {
	GetPhotos(ctx context.Context) ([]model.Photo, error)
	GetPhotoByID(ctx context.Context, id uint64) (model.UpdatePhoto, error)
	GetPhotoByUserID(ctx context.Context, photoID uint64) ([]model.GetPhoto, error)
	CreatePhoto(ctx context.Context, photo model.CreatePhoto, userID uint64) (model.CreatePhoto, error)
	UpdatePhoto(ctx context.Context, id uint64, photo model.UpdatePhoto) (model.UpdatePhoto, error)
	DeletePhotoByID(ctx context.Context, id uint64) (model.UpdatePhoto, error)
}

type photoServiceImpl struct {
	repoPhoto repository.PhotoQuery
	repoUser  repository.UserQuery
}

func NewPhotoService(repoPhoto repository.PhotoQuery, repoUser repository.UserQuery) PhotoService {
	return &photoServiceImpl{
		repoPhoto: repoPhoto,
		repoUser:  repoUser,
	}
}

func (p *photoServiceImpl) GetPhotos(ctx context.Context) ([]model.Photo, error) {
	photos, err := p.repoPhoto.GetPhotos(ctx)
	if err != nil {
		return []model.Photo{}, err
	}

	for i, photo := range photos {
		user, err := p.repoUser.GetUsersByID(ctx, photo.UserID)

		if err != nil {
			return nil, err
		}

		photos[i].User.ID = user.ID
		photos[i].User.Email = user.Email
		photos[i].User.Username = user.Username
	}

	return photos, nil
}

func (p *photoServiceImpl) GetPhotoByID(ctx context.Context, id uint64) (model.UpdatePhoto, error) {
	photo, err := p.repoPhoto.GetPhotoByID(ctx, id)
	if err != nil {
		return model.UpdatePhoto{}, err
	}
	return photo, err
}

func (p *photoServiceImpl) DeletePhotoByID(ctx context.Context, id uint64) (model.UpdatePhoto, error) {
	photo, err := p.repoPhoto.GetPhotoByID(ctx, id)
	if err != nil {
		return model.UpdatePhoto{}, err
	}

	if photo.ID == 0 {
		return model.UpdatePhoto{}, nil
	}

	err = p.repoPhoto.DeletePhotoByID(ctx, id)
	if err != nil {
		return model.UpdatePhoto{}, err
	}

	return photo, err
}

func (p *photoServiceImpl) CreatePhoto(ctx context.Context, CreatePhoto model.CreatePhoto, userID uint64) (model.CreatePhoto, error) {
	photo := model.CreatePhoto{
		Title:    CreatePhoto.Title,
		Caption:  CreatePhoto.Caption,
		PhotoURL: CreatePhoto.PhotoURL,
		UserID:   userID,
	}

	createdPhoto, err := p.repoPhoto.CreatePhoto(ctx, photo)
	if err != nil {
		return model.CreatePhoto{}, err
	}
	return createdPhoto, nil
}

func (p *photoServiceImpl) UpdatePhoto(ctx context.Context, id uint64, photo model.UpdatePhoto) (model.UpdatePhoto, error) {
	updatedPhoto, err := p.repoPhoto.UpdatePhoto(ctx, id, photo)
	if err != nil {
		return model.UpdatePhoto{}, err
	}
	return updatedPhoto, nil
}

func (p *photoServiceImpl) GetPhotoByUserID(ctx context.Context, photoID uint64) ([]model.GetPhoto, error) {
	photo, err := p.repoPhoto.GetPhotoByUserID(ctx, photoID)
	if err != nil {
		return []model.GetPhoto{}, err
	}
	return photo, nil
}
