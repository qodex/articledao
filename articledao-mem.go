package articledao

import (
	"errors"
	"strings"
)

//ArticleDAOInMem is In-memory implementation of ArticleDAO
type ArticleDAOInMem struct {
	articles []Article
}

//GetArticle finds one article by id
func (dao *ArticleDAOInMem) GetArticle(articleID string) (Article, error) {
	var article Article
	for _, a := range dao.articles {
		if a.ID == articleID {
			article = a
		}
	}
	var e error
	if article.ID == "" {
		e = errors.New("Article with ID '" + articleID + "' is not found")
	} else {
		e = nil
	}
	return article, e
}

//SaveArticle adds the article to an in-memory storage, returns article id
func (dao *ArticleDAOInMem) SaveArticle(article Article) (string, error) {
	dao.articles = append(dao.articles, article)
	return article.ID, nil
}

//FindByTagAndDate finds articles by tag and date
func (dao *ArticleDAOInMem) FindByTagAndDate(tag string, date string) (FindResponse, error) {

	if !IsDateParamValid(date) {
		return FindResponse{}, errors.New("date path variable is invalid")
	}

	response := FindResponse{}
	response.Tag = tag
	response.Count = 0
	response.Articles = make([]string, 0, 100)
	response.RelatedTags = make([]string, 0, 100)

	for _, a := range dao.articles {
		dateStr := strings.ReplaceAll(a.Date, "-", "")
		if dateStr == date {
			response.RelatedTags = MergeUnique(response.RelatedTags, a.Tags, tag)
			if Contains(a.Tags, tag) {
				response.Count = response.Count + 1
				response.Articles = append(response.Articles, a.ID)
			}
		}
	}

	return response, nil
}
