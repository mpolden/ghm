package main

import (
    "encoding/json"
    "fmt"
    "github.com/docopt/docopt-go"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "os"
    "path/filepath"
    "regexp"
    "strconv"
)

type Client struct {
    Username string
    Path     string
}

type Repo struct {
    Name string
}

func findNextPage(link string) int {
    re := regexp.MustCompile("page=(\\d+)")
    matches := re.FindStringSubmatch(link)
    if len(matches) != 2 {
        return -1
    }
    page, err := strconv.Atoi(matches[1])
    if err != nil {
        return -1
    }
    return page
}

func (c *Client) GetRepos(page int) ([]Repo, int) {
    u := fmt.Sprintf("https://api.github.com/users/%s/repos?", c.Username)
    params := url.Values{"page": {strconv.Itoa(page)}}
    resp, err := http.Get(u + params.Encode())
    if err != nil {
        log.Fatalf("GitHub API call failed: %s", err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatalf("Failed to read response body: %s", err)
    }

    var repos []Repo
    if err := json.Unmarshal(body, &repos); err != nil {
        log.Fatalf("Failed to unmarshal JSON response: %s", err)
    }
    next := findNextPage(resp.Header.Get("Link"))
    return repos, next
}

func (c *Client) GetAllRepos() []Repo {
    current := 1
    repos, nextPage := c.GetRepos(current)
    for nextPage > current {
        page, next := c.GetRepos(nextPage)
        repos = append(repos, page...)
        current = nextPage
        nextPage = next
    }
    return repos
}

func (c *Client) Run() {
    for _, repo := range c.GetAllRepos() {
        if !repo.IsExist(c.Path) {
            fmt.Printf("Not found: %s\n", repo.Name)
        }
    }
}

func (r *Repo) IsExist(path string) bool {
    repoPath := filepath.Join(path, r.Name)
    if _, err := os.Stat(repoPath); err == nil {
        return true
    }
    repoPath += ".git"
    if _, err := os.Stat(repoPath); err == nil {
        return true
    }
    return false
}

func main() {
    usage := `Check if local copies of GitHub repos exists in given path.

    Usage:
      check-repo <username> <path>
      check-repo -h | --help

    Options:
      -h --help             Show help.`

    arguments, _ := docopt.Parse(usage, nil, true, "", false)
    username := arguments["<username>"].(string)
    path := arguments["<path>"].(string)

    c := Client{Username: username, Path: path}
    c.Run()
}
