# Simple URL Shortener in Go using Chi

This project is a simple URL shortener written in Go and uses the Chi library for handling HTTP requests. It allows you to create short aliases for long URLs and perform redirection from the short alias to the original URL.

## Installation Instructions

1. Install Go if you haven't already. You can download it from the [official website](https://golang.org/).

2. Clone the repository:

    ```shell
    git clone https://github.com/RomanVassilchenko/simple-url-shortener.git
    ```

3. Change to the project directory:

    ```shell
    cd simple-url-shortener
    ```

4. Install dependencies:

    ```shell
    go mod download
    ```
> :warning: **Don't forget to add CONFIG_PATH to your env like this**: ```export CONFIG_PATH=./config/local.yaml```

5. Build and Run the project:

    ```shell
    go build -o simple-url-shortener ./cmd/simple-url-shortener
    ```

The application should be accessible at `http://localhost:8082`.

## Usage

### Creating a Short Alias for a URL

To create a short alias for a URL, send a POST request to `/` with a JSON body in the following format:

```json
{
  "alias": "myalias",
  "url": "https://example.com/long-url"
}
```

### Getting Redirection
To get a redirection to the original URL, send a GET request to /alias, where alias is the short alias you created earlier. The application will perform a redirect to the corresponding URL.

### Deleting an Alias
To delete an alias from the database, send a DELETE request to `/myalias`, where myalias is the alias you want to delete.

## Database
The application uses an SQLite database to store aliases and their corresponding URLs. The database is created automatically on the first run of the application.

## Dependencies
This project uses the following libraries:

[Chi](https://github.com/go-chi/chi) - HTTP request routing.

## License
This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.


Author: Roman Vassilchenko
