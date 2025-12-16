package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nikoskechris/gofetch"
	"github.com/nikoskechris/gofetch/domain/errors"
)

// User represents a user from the API.
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

// Post represents a blog post from the API.
type Post struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func main() {
	// Example 1: Basic GET request
	basicGetExample()

	// Example 2: GET with path parameters
	pathParametersExample()

	// Example 3: POST request
	postExample()

	// Example 4: Using interceptors
	interceptorExample()

	// Example 5: Error handling
	errorHandlingExample()

	// Example 6: Creating derived clients
	derivedClientExample()

	// Example 7: Progress tracking
	progressTrackingExample()
}

func basicGetExample() {
	fmt.Println("\n=== Example 1: Basic GET Request ===")

	client := gofetch.NewClient().
		SetBaseURL("https://jsonplaceholder.typicode.com").
		SetTimeout(10 * time.Second)

	var users []User
	resp, err := client.Get(context.Background(), "/users", nil, &users)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Retrieved %d users\n", len(users))
	if len(users) > 0 {
		fmt.Printf("First user: %s (%s)\n", users[0].Name, users[0].Email)
	}
}

func pathParametersExample() {
	fmt.Println("\n=== Example 2: GET with Path Parameters ===")

	client := gofetch.NewClient().
		SetBaseURL("https://jsonplaceholder.typicode.com")

	var user User
	params := map[string]interface{}{
		"id": 1,
	}

	resp, err := client.Get(context.Background(), "/users/:id", params, &user)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("User: %s (%s)\n", user.Name, user.Email)
}

func postExample() {
	fmt.Println("\n=== Example 3: POST Request ===")

	client := gofetch.NewClient().
		SetBaseURL("https://jsonplaceholder.typicode.com").
		SetHeader("Content-Type", "application/json")

	newPost := Post{
		UserID: 1,
		Title:  "GoFetch Example",
		Body:   "This is a test post created with GoFetch",
	}

	var createdPost Post
	resp, err := client.Post(context.Background(), "/posts", nil, newPost, &createdPost)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Created post with ID: %d\n", createdPost.ID)
	fmt.Printf("Title: %s\n", createdPost.Title)
}

func interceptorExample() {
	fmt.Println("\n=== Example 4: Using Interceptors ===")

	client := gofetch.NewClient().
		SetBaseURL("https://jsonplaceholder.typicode.com").
		AddRequestInterceptor(func(req *http.Request) (*http.Request, error) {
			fmt.Printf("Request interceptor: %s %s\n", req.Method, req.URL.String())
			req.Header.Set("X-Custom-Header", "GoFetch-Example")
			return req, nil
		}).
		AddResponseInterceptor(func(resp *http.Response) (*http.Response, error) {
			fmt.Printf("Response interceptor: Status %d\n", resp.StatusCode)
			return resp, nil
		})

	var user User
	_, err := client.Get(context.Background(), "/users/1", nil, &user)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("User retrieved: %s\n", user.Name)
}

func errorHandlingExample() {
	fmt.Println("\n=== Example 5: Error Handling ===")

	client := gofetch.NewClient().
		SetBaseURL("https://jsonplaceholder.typicode.com")

	var user User
	_, err := client.Get(context.Background(), "/users/99999", nil, &user)
	if err != nil {
		if httpErr, ok := err.(*errors.HTTPError); ok {
			fmt.Printf("HTTP Error: Status %d\n", httpErr.StatusCode)
			fmt.Printf("Response body: %s\n", string(httpErr.Body))
		} else {
			fmt.Printf("Error: %v\n", err)
		}
		return
	}
}

func derivedClientExample() {
	fmt.Println("\n=== Example 6: Creating Derived Clients ===")

	baseClient := gofetch.NewClient().
		SetBaseURL("https://jsonplaceholder.typicode.com").
		SetHeader("X-App-Version", "1.0.0")

	// Create a derived client with additional auth header
	authClient := baseClient.NewInstance().
		SetHeader("Authorization", "Bearer token123")

	var posts []Post
	_, err := authClient.Get(context.Background(), "/posts", map[string]interface{}{
		"userId": 1,
	}, &posts)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Retrieved %d posts for user 1\n", len(posts))
}

func progressTrackingExample() {
	fmt.Println("\n=== Example 7: Progress Tracking ===")

	client := gofetch.NewClient().
		SetBaseURL("https://jsonplaceholder.typicode.com").
		SetDownloadProgress(func(bytesTransferred, totalBytes int64) {
			if totalBytes > 0 {
				percentage := float64(bytesTransferred) / float64(totalBytes) * 100
				fmt.Printf("\rDownload progress: %.2f%%", percentage)
			}
		})

	var users []User
	_, err := client.Get(context.Background(), "/users", nil, &users)
	if err != nil {
		log.Printf("\nError: %v", err)
		return
	}

	fmt.Printf("\nDownload complete! Retrieved %d users\n", len(users))
}
