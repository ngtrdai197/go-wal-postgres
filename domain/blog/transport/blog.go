package transport

type BlogTransport interface{}

type blogTransport struct{}

func NewBlogTransport() BlogTransport {
	return blogTransport{}
}
