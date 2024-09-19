


# Novel Chapter Scraper

**novel-chapter-scraper** is a basic web crawler and scraper that gathers the latest chapters of novels from various websites and stores them in a database and notifies you with a telegram bot. This tool allows you to easily monitor and keep track of new releases without manually checking the websites.

## Features

- Scrapes new chapters from supported novel websites.
- Stores novel information and chapters in a database.
- Simple and easy-to-extend codebase for adding new scraping sources.
- Notifies the registered users via Telegram bot

## Table of Contents
- [Supported-Sites](#supported-sites)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Contributing](#contributing)
- [License](#license)

## Supported-sites

- **novelbin.me**
- **allnovelbin.net**
- **novel-next.net**


## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yuri-ncs/novel-chapter-scraper.git
   cd novel-chapter-scraper
   ```

2. Install the required dependencies:
   ```golang
   go mod tidy
   ```

3. Set up your database:
   
   The code can run with an empty database and you can register new novels and sites through telegram, if you need to thest right away, there is a method to populate the database.
   
   ```bash
   cd novel-chapter-scraper/docker-compose
   docker-compose up
   ```
 5. Run the code:
    ```bash
     go run main
    ```
 

## Usage

To run the scraper, execute the following command:
```bash
go run main
```

This will start the process of scraping the configured novel websites and saving the latest chapters to the database.

### Adding New Sources

You can add new novel sources by extending the scraping logic inside the `src/scrapers` directory. Each website should have its own scraper that adheres to the existing scraping format.

## Configuration

The configuration for the scraper, including database settings and scraping intervals, can be adjusted in the `.env` file or a configuration file within the project. Make sure to update the following fields:

- **Database Connection**: Update your database credentials.
- **Scraping Intervals**: Define how often the scraper should check for new chapters.

Example `.env` file:
```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_username
DB_PASS=your_password
DB_NAME=novel_scraper
SCRAPE_INTERVAL=60 # in minutes
```

## Contributing

Contributions are welcome! If you'd like to add a new feature, fix a bug, or improve the documentation, feel free to open a pull request. Please make sure your changes are well-tested and documented.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.
