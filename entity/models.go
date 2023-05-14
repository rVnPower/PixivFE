package entity

type Illust struct {
	ID      int               `json:"id"`
	Title   string            `json:"title"`
	Caption string            `json:"caption"`
	Images  map[string]string `json:"image_urls"`
	Artist  UserBrief         `json:"user"`
	// Tags Tag[];
	Date      string `json:"create_date"`
	Pages     int    `json:"page_count"`
	Views     int    `json:"total_view"`
	Bookmarks int    `json:"total_bookmarks"`
}

type Spotlight struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Thumbnail string `json:"thumbnail"`
	URL       string `json:"article_url"`
	Date      string `json:"publish_date"`
}

type UserBrief struct {
	ID      int               `json:"id"`
	Name    string            `json:"name"`
	Account string            `json:"account"`
	Avatar  map[string]string `json:"profile_image_urls"`
}
