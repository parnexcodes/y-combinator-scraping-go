package main

import (
	"fmt"

	"github.com/gocolly/colly"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Response struct {
	Message string
	Code    int
}

type Jobs struct {
	CompanyName string
	JobTitle    string
	JobType     string
	JobLocation string
	JobRole     string
	JobLink     string
}

func main() {

	app := fiber.New()
	app.Use(cors.New())

	jobs := []Jobs{}

	scrape := colly.NewCollector(
		colly.AllowedDomains("www.ycombinator.com", "ycombinator.com"),
	)

	app.Get("/", func(c *fiber.Ctx) error {
		data := Response{
			Message: "Welcome to api.",
			Code:    200,
		}
		return c.JSON(data)
	})

	app.Get("/jobs", func(c *fiber.Ctx) error {

		scrape.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL.String())
		})

		scrape.OnHTML("ul.space-y-2.overflow-hidden", func(h *colly.HTMLElement) {
			h.ForEach("li", func(i int, h *colly.HTMLElement) {
				j := Jobs{
					CompanyName: h.ChildText("a.justify-start.leading-loose span:nth-child(1)"),
					JobTitle:    h.ChildText("a.font-semibold.text-linkColor"),
					JobType:     h.ChildText("div.flex.flex-row.flex-wrap.justify-center div:nth-child(1)"),
					JobLocation: h.ChildText("div.flex.flex-row.flex-wrap.justify-center div:nth-child(2)"),
					JobRole:     h.ChildText("div.flex.flex-row.flex-wrap.justify-center div:nth-child(3)"),
					JobLink:     h.ChildAttr("div.mt-3.shrink-0.grow-0 a", "href"),
				}
				jobs = append(jobs, j)
			})
		})

		scrape.Visit("https://www.ycombinator.com/jobs/role/software-engineer")
		return c.JSON(jobs)
	})

	app.Listen(":3000")
}
