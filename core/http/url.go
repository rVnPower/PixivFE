package core

import "fmt"

func GetNewestArtworksURL(worktype, r18, lastID string) string {
	base := "https://www.pixiv.net/ajax/illust/new?limit=30&type=%s&r18=%s&lastId=%s"
	return fmt.Sprintf(base, worktype, r18, lastID)
}

func GetDiscoveryURL(mode string, limit int) string {
	base := "https://www.pixiv.net/ajax/discovery/artworks?mode=%s&limit=%d"
	return fmt.Sprintf(base, mode, limit)
}

func GetRankingURL(mode, content, date, page string) string {
	base := "https://www.pixiv.net/ranking.php?format=json&mode=%s&content=%s&date=%s&p=%s"
	baseNoDate := "https://www.pixiv.net/ranking.php?format=json&mode=%s&content=%s&p=%s"

	if date != "" {
		return fmt.Sprintf(base, mode, content, date, page)
	}

	return fmt.Sprintf(baseNoDate, mode, content, page)
}
