package service

import (
	"context"

	"github.com/geedotrar/mygram/internal/model"
	"github.com/geedotrar/mygram/internal/repository"
)

type CommentService interface {
	GetCommentByID(ctx context.Context, id uint64) (model.GetCommentByID, error)
	DeleteCommentByID(ctx context.Context, id uint64) (model.UpdateComment, error)
	CreateComment(ctx context.Context, comment model.CreateComment, user uint64) (model.CreateComment, error)
	UpdateComment(ctx context.Context, id uint64, comment model.UpdateComment) (model.UpdateComment, error)
	GetCommentsByPhotoID(ctx context.Context, photoID uint64) ([]model.Comment, error)
	GetCommentByID1(ctx context.Context, id uint64) (model.UpdateComment, error)
	GetComments(ctx context.Context) ([]model.GetCommentByID, error)
}

type commentServiceImpl struct {
	repoComment repository.CommentQuery
	repoUser    repository.UserQuery
	repoPhoto   repository.PhotoQuery
}

func NewCommentService(repoComment repository.CommentQuery, repoUser repository.UserQuery, repoPhoto repository.PhotoQuery) CommentService {
	return &commentServiceImpl{
		repoComment: repoComment,
		repoUser:    repoUser,
		repoPhoto:   repoPhoto,
	}
}

func (c *commentServiceImpl) GetCommentByID(ctx context.Context, id uint64) (model.GetCommentByID, error) {
	comment, err := c.repoComment.GetCommentByID(ctx, id)
	if err != nil {
		return model.GetCommentByID{}, err
	}
	user, err := c.repoUser.GetUsersByID(ctx, comment.UserID)
	if err != nil {
		return model.GetCommentByID{}, err
	}

	comment.User.ID = user.ID
	comment.User.Email = user.Email
	comment.User.Username = user.Username

	photo, err := c.repoPhoto.GetPhotoByID(ctx, comment.PhotoID)
	if err != nil {
		return model.GetCommentByID{}, err
	}
	comment.Photo.ID = photo.ID
	comment.Photo.Title = photo.Title
	comment.Photo.Caption = photo.Caption
	comment.Photo.PhotoURL = photo.PhotoURL
	comment.Photo.UserID = photo.UserID

	return comment, err
}
func (c *commentServiceImpl) GetCommentByID1(ctx context.Context, id uint64) (model.UpdateComment, error) {
	comment, err := c.repoComment.GetCommentByID1(ctx, id)
	if err != nil {
		return model.UpdateComment{}, err
	}
	return comment, err
}
func (c *commentServiceImpl) GetComments(ctx context.Context) ([]model.GetCommentByID, error) {
	comments, err := c.repoComment.GetComments(ctx)
	if err != nil {
		return nil, err
	}

	for i, comment := range comments {
		user, err := c.repoUser.GetUsersByID(ctx, comment.UserID)

		if err != nil {
			return nil, err
		}

		comments[i].User.Email = user.Email
		comments[i].User.Username = user.Username
		comments[i].User.ID = user.ID

		photo, err := c.repoPhoto.GetPhotoByID(ctx, comment.PhotoID)
		if err != nil {
			return nil, err
		}
		comments[i].Photo.ID = photo.ID
		comments[i].Photo.Title = photo.Title
		comments[i].Photo.Caption = photo.Caption
		comments[i].Photo.PhotoURL = photo.PhotoURL
		comments[i].Photo.UserID = photo.UserID
	}

	return comments, nil
}
func (c *commentServiceImpl) DeleteCommentByID(ctx context.Context, id uint64) (model.UpdateComment, error) {
	comment, err := c.repoComment.GetCommentByID1(ctx, id)
	if err != nil {
		return model.UpdateComment{}, err
	}

	if comment.ID == 0 {
		return model.UpdateComment{}, nil
	}

	err = c.repoComment.DeleteCommentByID(ctx, id)
	if err != nil {
		return model.UpdateComment{}, err
	}

	return comment, err
}

func (c *commentServiceImpl) CreateComment(ctx context.Context, CreateComment model.CreateComment, userID uint64) (model.CreateComment, error) {
	comment := model.CreateComment{
		Message: CreateComment.Message,
		PhotoID: CreateComment.PhotoID,
		UserID:  userID,
	}
	createdComment, err := c.repoComment.CreateComment(ctx, comment)
	if err != nil {
		return model.CreateComment{}, err
	}
	return createdComment, nil
}

func (c *commentServiceImpl) UpdateComment(ctx context.Context, id uint64, comment model.UpdateComment) (model.UpdateComment, error) {
	updatedComment, err := c.repoComment.UpdateComment(ctx, id, comment)
	if err != nil {
		return model.UpdateComment{}, err
	}
	return updatedComment, nil
}
func (c *commentServiceImpl) GetCommentsByPhotoID(ctx context.Context, photoID uint64) ([]model.Comment, error) {
	comments, err := c.repoComment.GetCommentsByPhotoID(ctx, photoID)
	if err != nil {
		return nil, err
	}
	for i, comment := range comments {
		user, err := c.repoUser.GetUsersByID(ctx, comment.UserID)

		if err != nil {
			return nil, err
		}

		comments[i].User.Email = user.Email
		comments[i].User.Username = user.Username
		comments[i].User.ID = user.ID

		photo, err := c.repoPhoto.GetPhotoByID(ctx, comment.PhotoID)
		if err != nil {
			return nil, err
		}
		comments[i].Photo.ID = photo.ID
		comments[i].Photo.Title = photo.Title
		comments[i].Photo.Caption = photo.Caption
		comments[i].Photo.PhotoURL = photo.PhotoURL
		comments[i].Photo.UserID = photo.UserID
	}
	return comments, nil
}
