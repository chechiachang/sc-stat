package fetcher

import (
	"fmt"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

const YilanGymDir = "data/yilan/gym"
const YilanSwimmingPoolDir = "data/yilan/swimmingpool"

func Yilan() {

	c := colly.NewCollector(
		colly.AllowedDomains("yilansports.com.tw"),
	)

	// Find and visit all links
	//c.OnHTML("a[href]", func(e *colly.HTMLElement) {
	//	e.Request.Visit(e.Attr("href"))
	//})

	r, _ := regexp.Compile("[0-9]+")

	c.OnHTML("div[class=\"pcount\"]", func(e *colly.HTMLElement) {
		switch {
		case strings.Contains(e.Text, "健身房") && strings.Contains(e.Text, "容留"):
			nums := r.FindAllString(e.Text, -1)
			now := time.Now()

			if len(nums) == 2 {
				log.Info(fmt.Sprintf("健身房: %s / %s\n", nums[0], nums[1]))
				datas := []Data{}
				datas = append(datas, Data{
					Time:     now,
					Name:     "gym",
					Number:   nums[0],
					Capacity: nums[1],
				})
				path := path.Join(
					YilanGymDir,
					getPathFromTime(now),
				)
				if err := writeCSV(path, DataToCsv(datas)); err != nil {
					log.Error(fmt.Sprintf("writeCSV error: %v", err))
				}
			} else {
				log.Warning(fmt.Sprintf("Invalid data: %v", nums))
			}

		case strings.Contains(e.Text, "游泳池") && strings.Contains(e.Text, "容留"):
			nums := r.FindAllString(e.Text, -1)
			now := time.Now()

			if len(nums) == 2 {
				log.Info(fmt.Sprintf("游泳池: %s / %s\n", nums[0], nums[1]))
				datas := []Data{}
				datas = append(datas, Data{
					Time:     now,
					Name:     "swimmingpool",
					Number:   nums[0],
					Capacity: nums[1],
				})
				path := path.Join(
					YilanSwimmingPoolDir,
					getPathFromTime(now),
				)
				if err := writeCSV(path, DataToCsv(datas)); err != nil {
					log.Error(fmt.Sprintf("writeCSV error: %v", err))
				}
			} else {
				log.Warning(fmt.Sprintf("Invalid data: %v", nums))
			}
		default:
			// do nothing
		}
	})

	c.OnRequest(func(r *colly.Request) {
		log.Infof("Visiting %v", r.URL)
	})

	c.Visit("https://yilansports.com.tw/")
}
