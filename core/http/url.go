package core

import "fmt"

func GetNewestArtworksURL(worktype, r18, lastID string) string {
	base := "https://www.pixiv.net/ajax/illust/new?limit=30&type=%s&r18=%s&lastId=%s"
	return fmt.Sprintf(base, worktype, r18, lastID)
}
