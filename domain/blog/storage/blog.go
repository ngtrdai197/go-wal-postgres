package storage

type BlogStorage interface{}

type blogStorage struct{}

func NewBlogStorage() BlogStorage {
	return &blogStorage{}
}
