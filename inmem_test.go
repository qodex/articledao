package articledao

import (
	"testing"
)

func TestSave(t *testing.T) {
	articlesDao := new(ArticleDAOInMem)
	want := "1"
	got, _ := articlesDao.SaveArticle(Article{
		ID:    "1",
		Title: "testing 1",
		Date:  "20190923",
		Body:  "testing....",
	})
	if got != want {
		t.Errorf("SaveArticle() = %q, want %q", got, want)
	}
}

func TestGet(t *testing.T) {
	articlesDao := new(ArticleDAOInMem)

	articlesDao.SaveArticle(Article{
		ID:    "1",
		Title: "testing 1",
		Date:  "20190923",
		Body:  "testing....",
	})

	_, e := articlesDao.GetArticle("1")

	if e != nil {
		t.Errorf("GetArticle() = %q", e.Error())
	}
}
