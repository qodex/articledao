package articledao

import (
	"strings"
	"testing"
)

func TestDynamoGet(t *testing.T) {
	articleID := "20190922-ae918a60-1c95-11ea-88b6-8f5b5674661e"
	articlesDao := new(ArticleDAODynamoDB)

	article, e := articlesDao.GetArticle(articleID)

	if e != nil {
		t.Errorf("GetArticle() = %q", e.Error())
	}

	if article.ID != articleID {
		t.Errorf("SaveArticle() did not return expected article object")
	}
}

func TestDynamoSave(t *testing.T) {
	articlesDao := new(ArticleDAODynamoDB)
	articleID, _ := articlesDao.SaveArticle(Article{
		Title: "testing 1",
		Date:  "2019-09-23",
		Body:  "testing....",
	})
	if len(articleID) < 1 || !strings.Contains(articleID, "20190923") {
		t.Errorf("SaveArticle() has not generated expeced article ID")
	}
}

func TestDynamoFind(t *testing.T) {
	articlesDao := new(ArticleDAODynamoDB)
	response, _ := articlesDao.FindByTagAndDate("birds", "20190922")
	if response.Count < 1 {
		t.Error("FindByTagAndDate() has not fount anything..")
	}
}
