package services

import (
	"errors"
	"realworld-api/dto"
	"realworld-api/models"
	"realworld-api/repositories"
	"realworld-api/utils"
)

func formatArticleRes(article *models.Article) *dto.ArticleRes {
	var tagList []string
	for _, tag := range article.Tags {
		tagList = append(tagList, tag.Name)
	}

	res := &dto.ArticleRes{}
	res.Article.Slug = article.Slug
	res.Article.Title = article.Title
	res.Article.Description = article.Description
	res.Article.Body = article.Body
	res.Article.TagList = tagList
	res.Article.CreatedAt = article.CreatedAt
	res.Article.UpdatedAt = article.UpdatedAt
	res.Article.Author.Username = article.Author.Username
	res.Article.Author.Bio = article.Author.Bio
	res.Article.Author.Image = article.Author.Image
	return res
}

func CreateArticle(userID uint, req dto.CreateArticleReq) (*dto.ArticleRes, error) {
	tags := repositories.FindOrCreateTags(req.Article.TagList)

	article := models.Article{
		Title:       req.Article.Title,
		Slug:        utils.GenerateSlug(req.Article.Title),
		Description: req.Article.Description,
		Body:        req.Article.Body,
		AuthorID:    userID,
		Tags:        tags,
	}

	if err := repositories.CreateArticle(&article); err != nil {
		return nil, errors.New("could not create article")
	}

	newArticle, _ := repositories.FindArticleBySlug(article.Slug)
	return formatArticleRes(newArticle), nil
}

func UpdateArticle(userID uint, slug string, req dto.UpdateArticleReq) (*dto.ArticleRes, error) {
	article, err := repositories.FindArticleBySlug(slug)
	if err != nil {
		return nil, errors.New("article not found")
	}

	if article.AuthorID != userID {
		return nil, errors.New("you do not have permission to edit this article")
	}

	if req.Article.Title != "" {
		article.Title = req.Article.Title
		article.Slug = utils.GenerateSlug(req.Article.Title)
	}
	if req.Article.Description != "" {
		article.Description = req.Article.Description
	}
	if req.Article.Body != "" {
		article.Body = req.Article.Body
	}

	if err := repositories.SaveArticle(article); err != nil {
		return nil, errors.New("update failed")
	}
	return formatArticleRes(article), nil
}

func DeleteArticle(userID uint, slug string) error {
	article, err := repositories.FindArticleBySlug(slug)
	if err != nil {
		return errors.New("article not found")
	}
	if article.AuthorID != userID {
		return errors.New("you do not have permission to delete this article")
	}
	return repositories.DeleteArticle(article)
}

func formatArticleResFull(article *models.Article, currentUserID uint) *dto.ArticleRes {
	var tagList []string
	for _, tag := range article.Tags {
		tagList = append(tagList, tag.Name)
	}

	favCount := repositories.CountFavorites(article.ID)
	favorited := repositories.IsFavorited(currentUserID, article.ID)

	res := &dto.ArticleRes{}
	res.Article.Slug = article.Slug
	res.Article.Title = article.Title
	res.Article.Description = article.Description
	res.Article.Body = article.Body
	res.Article.TagList = tagList
	res.Article.CreatedAt = article.CreatedAt
	res.Article.UpdatedAt = article.UpdatedAt
	res.Article.Favorited = favorited
	res.Article.FavoritesCount = favCount
	res.Article.Author.Username = article.Author.Username
	res.Article.Author.Bio = article.Author.Bio
	res.Article.Author.Image = article.Author.Image
	return res
}

func FavoriteArticle(userID uint, slug string) (*dto.ArticleRes, error) {
	article, err := repositories.FindArticleBySlug(slug)
	if err != nil {
		return nil, errors.New("article not found")
	}
	repositories.AddFavorite(userID, article.ID)
	return formatArticleResFull(article, userID), nil
}

func UnfavoriteArticle(userID uint, slug string) (*dto.ArticleRes, error) {
	article, err := repositories.FindArticleBySlug(slug)
	if err != nil {
		return nil, errors.New("article not found")
	}
	repositories.RemoveFavorite(userID, article.ID)
	return formatArticleResFull(article, userID), nil
}

func GetArticles(currentUserID uint, tag, author, favorited string, limit, offset int) (*dto.MultipleArticlesRes, error) {
	articles, count, err := repositories.FindArticles(tag, author, favorited, limit, offset)
	if err != nil {
		return nil, errors.New("error retrieving article data")
	}

	var result []interface{}
	for _, article := range articles {
		res := formatArticleResFull(&article, currentUserID)
		result = append(result, res.Article)
	}

	if result == nil {
		result = []interface{}{}
	}

	return &dto.MultipleArticlesRes{Articles: result, ArticlesCount: count}, nil
}

func GetArticleBySlug(currentUserID uint, slug string) (*dto.ArticleRes, error) {
	article, err := repositories.FindArticleBySlug(slug)
	if err != nil {
		return nil, errors.New("article not found")
	}

	following := repositories.IsFollowing(currentUserID, article.AuthorID)
	res := formatArticleResFull(article, currentUserID)
	res.Article.Author.Following = following
	return res, nil
}

func GetFeed(currentUserID uint, limit, offset int) (*dto.MultipleArticlesRes, error) {
	articles, count, err := repositories.FindFeedArticles(currentUserID, limit, offset)
	if err != nil {
		return nil, errors.New("error retrieving feed")
	}

	var result []interface{}
	for _, article := range articles {
		res := formatArticleResFull(&article, currentUserID)
		result = append(result, res.Article)
	}

	if result == nil {
		result = []interface{}{}
	}

	return &dto.MultipleArticlesRes{Articles: result, ArticlesCount: count}, nil
}
