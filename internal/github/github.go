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

// Scope represents a GitHub authorization scope.
// See https://docs.github.com/en/developers/apps/scopes-for-oauth-apps
type Scope string

const (
	// ScopeRepo grants full access to private and public repositories. It also grants ability to manage user projects.
	ScopeRepo Scope = "repo"
	// ScopeRepoStatus grants read/write access to public and private repository commit statuses.
	ScopeRepoStatus Scope = "repo:status"
	// ScopeRepoDeployment grants access to deployment statuses for public and private repositories.
	ScopeRepoDeployment Scope = "repo_deployment"
	// ScopePublicRepo grants access only to public repositories.
	ScopePublicRepo Scope = "public_repo"
	// ScopeRepoInvite grants accept/decline abilities for invitations to collaborate on a repository.
	ScopeRepoInvite Scope = "repo:invite"
	// ScopeSecurityEvents grants read and write access to security events in the code scanning API.
	ScopeSecurityEvents Scope = "security_events"

	// ScopeWritePackages grants access to upload or publish a package in GitHub Packages.
	ScopeWritePackages Scope = "write:packages"
	// ScopeReadPackages grants access to download or install packages from GitHub Packages.
	ScopeReadPackages Scope = "read:packages"
	// ScopeDeletePackages grants access to delete packages from GitHub Packages.
	ScopeDeletePackages Scope = "delete:packages"

	// ScopeAdminOrg grants access to fully manage the organization and its teams, projects, and memberships.
	ScopeAdminOrg Scope = "admin:org"
	// ScopeWriteOrg grants read and write access to organization membership, organization projects, and team membership.
	ScopeWriteOrg Scope = "write:org"
	// ScopeReadOrg grants read-only access to organization membership, organization projects, and team membership.
	ScopeReadOrg Scope = "read:org"

	// ScopeAdminPublicKey grants access to fully manage public keys.
	ScopeAdminPublicKey Scope = "admin:public_key"
	// ScopeWritePublicKey grants access to create, list, and view details for public keys.
	ScopeWritePublicKey Scope = "write:public_key"
	// ScopeReadPublicKey grants access to list and view details for public keys.
	ScopeReadPublicKey Scope = "read:public_key"

	// ScopeAdminRepoHook grants read, write, ping, and delete access to repository hooks in public and private repositories.
	ScopeAdminRepoHook Scope = "admin:repo_hook"
	// ScopeWriteRepoHook grants read, write, and ping access to hooks in public or private repositories.
	ScopeWriteRepoHook Scope = "write:repo_hook"
	// ScopeReadRepoHook grants read and ping access to hooks in public or private repositories.
	ScopeReadRepoHook Scope = "read:repo_hook"

	// ScopeAdminOrgHook grants read, write, ping, and delete access to organization hooks.
	ScopeAdminOrgHook Scope = "admin:org_hook"
	// ScopeGist grants write access to gists.
	ScopeGist Scope = "gist"
	// ScopeNotifications grants read access to a user's notifications and misc.
	ScopeNotifications Scope = "notifications"

	// ScopeUser grants read/write access to profile info only.
	ScopeUser Scope = "user"
	// ScopeReadUser grants access to read a user's profile data.
	ScopeReadUser Scope = "read:user"
	// ScopeUserEmail grants read access to a user's email addresses.
	ScopeUserEmail Scope = "user:email"
	// ScopeUserFollow grants access to follow or unfollow other users.
	ScopeUserFollow Scope = "user:follow"

	// ScopeDeleteRepo grants access to delete adminable repositories.
	ScopeDeleteRepo Scope = "delete_repo"

	// ScopeWriteDiscussion allows read and write access for team discussions.
	ScopeWriteDiscussion Scope = "write:discussion"
	// ScopeReadDiscussion allows read access for team discussions.
	ScopeReadDiscussion Scope = "read:discussion"

	// ScopeWorkflow grants the ability to add and update GitHub Actions workflow files.
	ScopeWorkflow Scope = "workflow"

	// ScopeAdminGPGKey grants access to fully manage GPG keys.
	ScopeAdminGPGKey Scope = "admin:gpg_key"
	// ScopeWriteGPGKey grants access to create, list, and view details for GPG keys.
	ScopeWriteGPGKey Scope = "write:gpg_key"
	// ScopeReadGPGKey grants access to list and view details for GPG keys.
	ScopeReadGPGKey Scope = "read:gpg_key"
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

// GitHub is a minimal client for GitHub REST API (also known as API v3).
// See https://docs.github.com/en/rest
type GitHub struct {
	client      *http.Client
	apiURL      string
	accessToken string
}

// New creates a new instance of GitHub.
// If the access token does not have any of the given scope, an error will be returned.
func New(accessToken string, scopes ...Scope) (*GitHub, error) {
	transport := &http.Transport{}
	client := &http.Client{
		Transport: transport,
	}

	g := &GitHub{
		client:      client,
		apiURL:      githubAPIURL,
		accessToken: accessToken,
	}

	if err := g.checkScopes(scopes...); err != nil {
		return nil, err
	}

	return g, nil
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

func (g *GitHub) checkScopes(scopes ...Scope) error {
	// Call an endpoint to get the OAuth scopes of the access token from the headers
	// See https://docs.github.com/en/developers/apps/scopes-for-oauth-apps

	if len(scopes) > 0 {
		req, err := g.createRequest(context.Background(), "HEAD", g.apiURL+"/user", "", nil)
		if err != nil {
			return err
		}

		resp, err := g.makeRequest(req, 200)
		if err != nil {
			return err
		}

		// Ensure the access token has all the required OAuth scopes
		oauthScopes := resp.Header.Get("X-OAuth-Scopes")
		for _, scope := range scopes {
			if !strings.Contains(oauthScopes, string(scope)) {
				return fmt.Errorf("access token does not have the scope: %s", scope)
			}
		}
	}

	return nil
}

// GetLatestRelease returns the latest GitHub release.
// The latest release is the most recent non-prerelease and non-draft release.
// See https://docs.github.com/en/rest/reference/repos#get-the-latest-release
func (g *GitHub) GetLatestRelease(ctx context.Context, owner, repo string) (*Release, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/releases/latest", g.apiURL, owner, repo)

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
func (g *GitHub) CreateRelease(ctx context.Context, owner, repo string, params ReleaseData) (*Release, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/releases", g.apiURL, owner, repo)

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
func (g *GitHub) UpdateRelease(ctx context.Context, owner, repo string, releaseID int, params ReleaseData) (*Release, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/releases/%d", g.apiURL, owner, repo, releaseID)

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
func (g *GitHub) EnableBranchProtection(ctx context.Context, owner, repo, branch string) error {
	url := fmt.Sprintf("%s/repos/%s/%s/branches/%s/protection/enforce_admins", g.apiURL, owner, repo, branch)

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
func (g *GitHub) DisableBranchProtection(ctx context.Context, owner, repo, branch string) error {
	url := fmt.Sprintf("%s/repos/%s/%s/branches/%s/protection/enforce_admins", g.apiURL, owner, repo, branch)

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
