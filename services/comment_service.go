package services

import (
	"errors"
	"realworld-api/dto"
	"realworld-api/models"
	"realworld-api/repositories"
)

func formatCommentRes(comment *models.Comment) *dto.CommentRes {
	res := &dto.CommentRes{}
	res.Comment.ID = comment.ID
	res.Comment.CreatedAt = comment.CreatedAt
	res.Comment.UpdatedAt = comment.UpdatedAt
	res.Comment.Body = comment.Body
	res.Comment.Author.Username = comment.Author.Username
	res.Comment.Author.Bio = comment.Author.Bio
	res.Comment.Author.Image = comment.Author.Image
	return res
}

func CreateComment(userID uint, slug string, req dto.CommentReq) (*dto.CommentRes, error) {
	article, err := repositories.FindArticleBySlug(slug)
	if err != nil {
		return nil, errors.New("article not found")
	}

	comment := models.Comment{
		Body:      req.Comment.Body,
		ArticleID: article.ID,
		AuthorID:  userID,
	}

	if err := repositories.CreateComment(&comment); err != nil {
		return nil, errors.New("could not create comment")
	}

	newComment, _ := repositories.FindCommentByID(comment.ID)
	user, _ := repositories.FindUserByID(newComment.AuthorID)
	newComment.Author = *user

	return formatCommentRes(newComment), nil
}

func GetComments(slug string) (*dto.MultipleCommentsRes, error) {
	article, err := repositories.FindArticleBySlug(slug)
	if err != nil {
		return nil, errors.New("article not found")
	}

	comments, _ := repositories.FindCommentsByArticleID(article.ID)
	var result []interface{}
	for _, c := range comments {
		result = append(result, formatCommentRes(&c).Comment)
	}

	return &dto.MultipleCommentsRes{Comments: result}, nil
}

func DeleteComment(userID uint, slug string, commentID uint) error {
	comment, err := repositories.FindCommentByID(commentID)
	if err != nil {
		return errors.New("comment not found")
	}
	if comment.AuthorID != userID {
		return errors.New("you do not have permission to delete this comment")
	}
	return repositories.DeleteComment(comment)
}
