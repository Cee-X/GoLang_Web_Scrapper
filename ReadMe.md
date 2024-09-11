 # RSS Aggregator

## Overview
RSS Aggregator is a Go-based application that fetches and processes RSS feeds from various sources. It stores the feed data in a PostgreSQL database and provides an API to access the aggregated data.

## Features
- Fetches RSS feeds from multiple sources
- Stores feed data in a PostgreSQL database
- Provides an API to access the aggregated data
- Handles duplicate entries and errors gracefully

## Prerequisites
- Go 1.18 or later
- PostgreSQL 12 or later

## Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/Cee-X/rss-aggregator.git
   cd rss-aggregator

2. install dependencies:
    ```sh
    go mod tidy

3. Create a PostgreSQL database:


4. Run the application:
    ```sh
    go build && ./rssagg

Usage 
GET/users  - Get all users
POST/users - Create a user
GET/feeds  - Get all feeds
GET/feed_follows - Get all feed follows
POST/feed_follows - Create a feed follow
DELETE/feed_follows - Delete a feed follow
GET/posts - Get all posts   


Acknowledgements
This project uses the following libraries:
.go-chi for the http router
.sqlc for the SQL query generation
.goose for the database migration
