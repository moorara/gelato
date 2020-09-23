package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	netURL "net/url"

	"github.com/moorara/gelato/pkg/xhttp"
)

const (
	githubURL        = "https://github.com"
	githubAPIURL     = "https://api.github.com"
	userAgentDefault = "gelato"
	acceptDefault    = "application/vnd.github.v3+json"
	contentTypeJSON  = "application/json"
)

type (
	// User is a GitHub user.
	User struct {
		ID         int       `json:"id"`
		Login      string    `json:"login"`
		Type       string    `json:"type"`
		Email      string    `json:"email"`
		Name       string    `json:"name"`
		Company    string    `json:"company"`
		Location   string    `json:"location"`
		URL        string    `json:"url"`
		HTMLURL    string    `json:"html_url"`
		AvatarURL  string    `json:"avatar_url"`
		GravatarID string    `json:"gravatar_id"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
	}

	// ReleaseData is used for creating or updating a GitHub release.
	ReleaseData struct {
		Name       string `json:"name"`
		TagName    string `json:"tag_name"`
		Target     string `json:"target_commitish"`
		Draft      bool   `json:"draft"`
		Prerelease bool   `json:"prerelease"`
		Body       string `json:"body"`
	}

	// Release is a GitHub release.
	Release struct {
		ID          int            `json:"id"`
		Name        string         `json:"name"`
		TagName     string         `json:"tag_name"`
		Target      string         `json:"target_commitish"`
		Draft       bool           `json:"draft"`
		Prerelease  bool           `json:"prerelease"`
		Body        string         `json:"body"`
		URL         string         `json:"url"`
		HTMLURL     string         `json:"html_url"`
		AssetsURL   string         `json:"assets_url"`
		UploadURL   string         `json:"upload_url"`
		TarballURL  string         `json:"tarball_url"`
		ZipballURL  string         `json:"zipball_url"`
		CreatedAt   time.Time      `json:"created_at"`
		PublishedAt time.Time      `json:"published_at"`
		Author      User           `json:"author"`
		Assets      []ReleaseAsset `json:"assets"`
	}

	// ReleaseAsset is a Github release asset.
	ReleaseAsset struct {
		ID            int       `json:"id"`
		Name          string    `json:"name"`
		Label         string    `json:"label"`
		State         string    `json:"state"`
		ContentType   string    `json:"content_type"`
		Size          int       `json:"size"`
		DownloadCount int       `json:"download_count"`
		URL           string    `json:"url"`
		DownloadURL   string    `json:"browser_download_url"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
		Uploader      User      `json:"uploader"`
	}
)

// GitHub is a client for GitHub REST API (also known as API v3).
// See https://docs.github.com/en/rest
type GitHub struct {
	client      *http.Client
	apiURL      string
	accessToken string
	owner, repo string
}

// New creates a new instance of GitHub.
func New(accessToken, owner, repo string) *GitHub {
	transport := &http.Transport{}
	client := &http.Client{
		Transport: transport,
	}

	return &GitHub{
		client:      client,
		apiURL:      githubAPIURL,
		accessToken: accessToken,
		owner:       owner,
		repo:        repo,
	}
}

func (g *GitHub) createRequest(ctx context.Context, method, url, contentType string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+g.accessToken)
	req.Header.Set("User-Agent", userAgentDefault) // See https://docs.github.com/en/rest/overview/resources-in-the-rest-api#user-agent-required
	req.Header.Set("Accept", acceptDefault)        // See https://docs.github.com/en/rest/overview/media-types

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	return req, nil
}

func (g *GitHub) makeRequest(req *http.Request, expectedStatusCode int) (*http.Response, error) {
	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != expectedStatusCode {
		return nil, xhttp.NewClientError(resp)
	}

	return resp, nil
}

// GetLatestRelease returns the latest GitHub release.
// The latest release is the most recent non-prerelease and non-draft release.
// See https://docs.github.com/en/rest/reference/repos#get-the-latest-release
func (g *GitHub) GetLatestRelease(ctx context.Context) (*Release, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/releases/latest", g.apiURL, g.owner, g.repo)

	req, err := g.createRequest(ctx, "GET", url, "", nil)
	if err != nil {
		return nil, err
	}

	resp, err := g.makeRequest(req, 200)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	release := &Release{}
	if err = json.NewDecoder(resp.Body).Decode(release); err != nil {
		return nil, err
	}

	return release, nil
}

// CreateRelease creates a new GitHub release.
// See https://docs.github.com/en/rest/reference/repos#create-a-release
func (g *GitHub) CreateRelease(ctx context.Context, params ReleaseData) (*Release, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/releases", g.apiURL, g.owner, g.repo)

	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(params)

	req, err := g.createRequest(ctx, "POST", url, contentTypeJSON, body)
	if err != nil {
		return nil, err
	}

	resp, err := g.makeRequest(req, 201)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	release := &Release{}
	if err = json.NewDecoder(resp.Body).Decode(release); err != nil {
		return nil, err
	}

	return release, nil
}

// UpdateRelease updates an existing GitHub release.
// See https://docs.github.com/en/rest/reference/repos#update-a-release
func (g *GitHub) UpdateRelease(ctx context.Context, releaseID int, params ReleaseData) (*Release, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/releases/%d", g.apiURL, g.owner, g.repo, releaseID)

	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(params)

	req, err := g.createRequest(ctx, "PATCH", url, contentTypeJSON, body)
	if err != nil {
		return nil, err
	}

	resp, err := g.makeRequest(req, 200)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	release := &Release{}
	if err = json.NewDecoder(resp.Body).Decode(release); err != nil {
		return nil, err
	}

	return release, nil
}

// EnableBranchProtection enables a branch protection for administrator users.
// See https://docs.github.com/en/rest/reference/repos#set-admin-branch-protection
func (g *GitHub) EnableBranchProtection(ctx context.Context, branch string) error {
	url := fmt.Sprintf("%s/repos/%s/%s/branches/%s/protection/enforce_admins", g.apiURL, g.owner, g.repo, branch)

	req, err := g.createRequest(ctx, "POST", url, "", nil)
	if err != nil {
		return err
	}

	resp, err := g.makeRequest(req, 200)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// DisableBranchProtection disables a branch protection for administrator users.
// See https://docs.github.com/en/rest/reference/repos#delete-admin-branch-protection
func (g *GitHub) DisableBranchProtection(ctx context.Context, branch string) error {
	url := fmt.Sprintf("%s/repos/%s/%s/branches/%s/protection/enforce_admins", g.apiURL, g.owner, g.repo, branch)

	req, err := g.createRequest(ctx, "DELETE", url, "", nil)
	if err != nil {
		return err
	}

	resp, err := g.makeRequest(req, 204)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// UploadReleaseAsset uploads a file to a GitHub release.
// See https://docs.github.com/en/rest/reference/repos#upload-a-release-asset
func (g *GitHub) UploadReleaseAsset(ctx context.Context, uploadURL, file string) (*ReleaseAsset, error) {
	assetPath := filepath.Clean(file)
	assetName := filepath.Base(assetPath)

	assetFile, err := os.Open(assetPath)
	if err != nil {
		return nil, err
	}
	defer assetFile.Close()

	stat, err := assetFile.Stat()
	if err != nil {
		return nil, err
	}
	length := stat.Size()

	// Read the first 512 bytes of file to determine the mime type of the asset file
	buff := make([]byte, 512)
	if _, err := assetFile.Read(buff); err != nil {
		return nil, err
	}
	mimeType := http.DetectContentType(buff) // http.DetectContentType will return "application/octet-stream" if it cannot determine a more specific one

	// Reset the offset back to the beginning of the file
	_, err = assetFile.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	url := strings.Replace(uploadURL, "{?name,label}", "", 1)
	url = fmt.Sprintf("%s?name=%s", url, netURL.QueryEscape(assetName))

	req, err := g.createRequest(ctx, "POST", url, mimeType, assetFile)
	if err != nil {
		return nil, err
	}
	req.ContentLength = length

	resp, err := g.makeRequest(req, 201)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	releaseAsset := &ReleaseAsset{}
	if err = json.NewDecoder(resp.Body).Decode(releaseAsset); err != nil {
		return nil, err
	}

	return releaseAsset, nil
}

// DownloadReleaseAsset downloads a file from a GitHub release.
func (g *GitHub) DownloadReleaseAsset(ctx context.Context, downloadURL, file string) error {
	req, err := g.createRequest(ctx, "GET", downloadURL, "", nil)
	if err != nil {
		return err
	}
	req.Header.Del("Accept")

	resp, err := g.makeRequest(req, 200)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.OpenFile(file, os.O_WRONLY, 0755)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
