package articledao

//ArticleDAO interface defines persistence for Articles
type ArticleDAO interface {
	GetArticle(articleID string) (Article, error)

	SaveArticle(Article) (string, error)

	FindByTagAndDate(tag string, date string) (FindResponse, error)
}
