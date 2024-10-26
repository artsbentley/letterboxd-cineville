import scrapy


class WatchlistSpider(scrapy.Spider):
    name = "watchlist"
    start_urls = ["https://letterboxd.com/deltore/watchlist/"]

    def parse(self, response):
        # Extract film names from the poster containers
        film_names = response.css("li.poster-container img::attr(alt)").getall()

        for name in film_names:
            print(name)  # Print each film name

        # Follow pagination links
        next_page = response.css("li.paginate-page a::attr(href)").re(
            r"/watchlist/page/\d+"
        )
        if next_page:
            yield response.follow(next_page[0], self.parse)
