# Web Crawler

This project is a web crawler that allows users to enter a URL and request the crawling of a webpage. It checks if the URL has been crawled in the last 60 minutes. If the page is available, it returns the stored crawled page. Otherwise, it crawls the page in real-time and returns it to the user.

## Getting Started

To get started with this project, you'll need to set up both the client and server components.

### Prerequisites

- Go (for the server)
- HTML/JavaScript (for the client)

### Client

The client-side code is responsible for providing a user interface to enter URLs and request crawling. It's a simple HTML page with a search bar and a crawl button.

1. Clone the repository to your local machine.

2. Open the `client.html` file in a web browser.

3. Enter the URL you want to crawl and click the "Crawl" button.

4. The client sends a request to the server, which processes the request and returns the crawled page.

### Server

The server-side code is responsible for handling incoming requests, checking if the page has been crawled, and performing real-time crawling when necessary.

1. Clone the repository to your server or local machine.

2. Ensure you have Go installed.

3. Run the server code by executing the following command:


4. The server listens on port 8080 by default.

5. Send HTTP requests to the server's endpoints to request crawling. Paying customers can be identified using query parameters.

## Functionality

- The server can handle multiple crawling requests concurrently.
- Paying customers receive priority in crawling over non-paying customers.
- The server retries in case a webpage is not available.
- Crawled pages are stored on disk for faster retrieval if requested within 60 minutes.


