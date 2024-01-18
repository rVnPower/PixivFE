package pages

import (
	"fmt"
	"strconv"
	"time"

	core "codeberg.org/vnpower/pixivfe/v2/core/webapi"
	"github.com/gofiber/fiber/v2"
)

func RankingCalendarPage(c *fiber.Ctx) error {
	mode := c.Query("mode", "daily")
	date := c.Query("date", "")

	var year int
	var month int
	var monthLit string

	// If the user supplied a date
	if len(date) == 6 {
		var err error
		year, err = strconv.Atoi(date[:4])
		if err != nil {
			return err
		}
		month, err = strconv.Atoi(date[4:])
		if err != nil {
			return err
		}
	} else {
		now := time.Now()
		year = now.Year()
		month = int(now.Month())
	}

	realDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	monthLit = realDate.Month().String()

	monthBefore := realDate.AddDate(0, -1, 0)
	monthAfter := realDate.AddDate(0, 1, 0)
	thisMonthLink := fmt.Sprintf("%d%02d", realDate.Year(), realDate.Month())
	monthBeforeLink := fmt.Sprintf("%d%02d", monthBefore.Year(), monthBefore.Month())
	monthAfterLink := fmt.Sprintf("%d%02d", monthAfter.Year(), monthAfter.Month())

	render, err := core.GetRankingCalendar(c, mode, year, month)
	if err != nil {
		return err
	}

	return c.Render("pages/rankingCalendar", fiber.Map{"Title": "Ranking calendar", "Render": render, "Mode": mode, "Month": monthLit, "Year": year, "MonthBefore": monthBeforeLink, "MonthAfter": monthAfterLink, "ThisMonth": thisMonthLink})
}
