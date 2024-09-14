package business

type BlogBusiness interface{}

type blogBusiness struct{}

func NewBlogBusiness() BlogBusiness {
	return &blogBusiness{}
}
