package handler

const (
	ArtworkInformationURL   = "https://www.pixiv.net/ajax/illust/%s"
	ArtworkImagesURL        = "https://www.pixiv.net/ajax/illust/%s/pages"
	ArtworkRelatedURL       = "https://www.pixiv.net/ajax/illust/%s/recommend/init?limit=%d"
	ArtworkCommentsURL      = "https://www.pixiv.net/ajax/illusts/comments/roots?illust_id=%s&limit=100"
	ArtworkNewestURL        = "https://www.pixiv.net/ajax/illust/new?limit=30&type=%s&r18=%s&lastId=%s"
	ArtworkRankingURL       = "https://www.pixiv.net/ranking.php?format=json&mode=%s&content=%s&p=%s"
	ArtworkDiscoveryURL     = "https://www.pixiv.net/ajax/discovery/artworks?mode=%s&limit=%d"
	SearchTagURL            = "https://www.pixiv.net/ajax/search/tags/%s"
	SearchArtworksURL       = "https://www.pixiv.net/ajax/search/%s/%s?order=%s&mode=%s&p=%s"
	SearchTopURL            = "https://www.pixiv.net/ajax/search/top/%s"
	UserInformationURL      = "https://www.pixiv.net/ajax/user/%s?full=1"
	UserBasicInformationURL = "https://www.pixiv.net/ajax/user/%s"
	UserArtworksURL         = "https://www.pixiv.net/ajax/user/%s/profile/all"
	UserArtworksFullURL     = "https://www.pixiv.net/ajax/user/%s/profile/illusts?work_category=illustManga&is_first_page=0&lang=en%s"
	FrequentTagsURL         = "https://www.pixiv.net/ajax/tags/frequent/illust?%s"
)
