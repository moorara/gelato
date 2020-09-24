package github

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

const (
	mockReleaseBody = `{
		"id": 1,
		"name": "v1.0.0",
		"tag_name": "v1.0.0",
		"target_commitish": "master",
		"draft": false,
		"prerelease": false,
		"body": "Description of the release",
		"author": {
			"id": 1,
			"login": "octocat",
			"type": "User"
		},
		"assets": [
			{
				"id": 1,
				"name": "example.zip",
				"label": "short description",
				"state": "uploaded",
				"content_type": "application/zip",
				"size": 1024,
				"uploader": {
					"id": 1,
					"login": "octocat",
					"type": "User"
				}
			}
		]
	}`

	mockReleaseAssetBody = `{
		"id": 1,
		"name": "example.zip",
		"label": "short description",
		"state": "uploaded",
		"content_type": "application/zip",
		"size": 1024,
		"uploader": {
			"id": 1,
			"login": "octocat",
			"type": "User"
		}
	}`
)

var (
	expectedRelease = &Release{
		ID:         1,
		Name:       "v1.0.0",
		TagName:    "v1.0.0",
		Target:     "master",
		Draft:      false,
		Prerelease: false,
		Body:       "Description of the release",
		Author: User{
			ID:    1,
			Login: "octocat",
			Type:  "User",
		},
		Assets: []ReleaseAsset{
			{
				ID:          1,
				Name:        "example.zip",
				Label:       "short description",
				State:       "uploaded",
				ContentType: "application/zip",
				Size:        1024,
				Uploader: User{
					ID:    1,
					Login: "octocat",
					Type:  "User",
				},
			},
		},
	}

	expectedReleaseAsset = &ReleaseAsset{
		ID:          1,
		Name:        "example.zip",
		Label:       "short description",
		State:       "uploaded",
		ContentType: "application/zip",
		Size:        1024,
		Uploader: User{
			ID:    1,
			Login: "octocat",
			Type:  "User",
		},
	}
)

type MockResponse struct {
	Method       string
	Path         string
	StatusCode   int
	ResponseBody string
}

func createMockHTTPServer(mocks ...MockResponse) *httptest.Server {
	r := mux.NewRouter()
	for _, m := range mocks {
		r.Methods(m.Method).Path(m.Path).HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(m.StatusCode)
			io.WriteString(w, m.ResponseBody)
		})
	}

	return httptest.NewServer(r)
}

func TestNew(t *testing.T) {
	tests := []struct {
		name          string
		accessToken   string
		scopes        []Scope
		expectedError string
	}{
		{
			name:          "NoScopeRequired",
			accessToken:   "github-token",
			scopes:        []Scope{},
			expectedError: "",
		},
		{
			name:          "ScopeRequired",
			accessToken:   "github-token",
			scopes:        []Scope{ScopeRepo},
			expectedError: "HEAD /user 401: ",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gh, err := New(tc.accessToken, tc.scopes...)

			if tc.expectedError != "" {
				assert.Nil(t, gh)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, gh)
				assert.Equal(t, githubAPIURL, gh.apiURL)
				assert.Equal(t, tc.accessToken, gh.accessToken)
			}
		})
	}
}

func TestGitHub_CreateRequest(t *testing.T) {
	tests := []struct {
		name          string
		accessToken   string
		ctx           context.Context
		method        string
		url           string
		contentType   string
		body          io.Reader
		expectedError string
	}{
		{
			name:          "RequestError",
			accessToken:   "github-token",
			ctx:           nil,
			method:        "",
			url:           "",
			contentType:   "",
			body:          nil,
			expectedError: "net/http: nil Context",
		},
		{
			name:          "Success",
			accessToken:   "github-token",
			ctx:           context.Background(),
			method:        "GET",
			url:           "https://api.github.com/users/octocat",
			contentType:   contentTypeJSON,
			body:          nil,
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gh := &GitHub{
				accessToken: tc.accessToken,
			}

			req, err := gh.createRequest(tc.ctx, tc.method, tc.url, tc.contentType, tc.body)

			if tc.expectedError != "" {
				assert.Nil(t, req)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, req)
				assert.Equal(t, "token "+tc.accessToken, req.Header.Get("Authorization"))
				assert.Equal(t, userAgentDefault, req.Header.Get("User-Agent"))
				assert.Equal(t, acceptDefault, req.Header.Get("Accept"))
				assert.Equal(t, tc.contentType, req.Header.Get("Content-Type"))
			}
		})
	}
}

func TestGitHub_MakeRequest(t *testing.T) {
	tests := []struct {
		name               string
		mockResponses      []MockResponse
		method             string
		url                string
		body               io.Reader
		expectedStatusCode int
		expectedError      string
	}{
		{
			name:               "ClientError",
			mockResponses:      []MockResponse{},
			method:             "GET",
			url:                "",
			body:               nil,
			expectedStatusCode: 200,
			expectedError:      `Get "": unsupported protocol scheme ""`,
		},
		{
			name: "UnexpectedStatusCode",
			mockResponses: []MockResponse{
				{"GET", "/users/{username}", 400, `bad request`},
			},
			method:             "GET",
			url:                "https://api.github.com/users/octocat",
			body:               nil,
			expectedStatusCode: 200,
			expectedError:      "GET /users/octocat 400: bad request",
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{
					"GET", "/users/{username}", 200, `{
						"id": 1,
						"type": "User",
						"login": "octocat",
						"email": "octocat@github.com",
						"name": "monalisa octocat"
					}`,
				},
			},
			method:             "GET",
			url:                "https://api.github.com/users/octocat",
			body:               nil,
			expectedStatusCode: 200,
			expectedError:      "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gh := &GitHub{
				client: new(http.Client),
			}

			if len(tc.mockResponses) > 0 {
				ts := createMockHTTPServer(tc.mockResponses...)
				defer ts.Close()
				tc.url = strings.Replace(tc.url, "https://api.github.com", ts.URL, 1)
			}

			req, err := http.NewRequest(tc.method, tc.url, tc.body)
			assert.NoError(t, err)

			resp, err := gh.makeRequest(req, tc.expectedStatusCode)

			if tc.expectedError != "" {
				assert.Nil(t, resp)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})
	}
}

func TestGitHub_CheckScopes(t *testing.T) {
	tests := []struct {
		name          string
		mockResponses []MockResponse
		accessToken   string
		scopes        []Scope
		expectedError string
	}{
		{
			name: "InvalidStatusCode",
			mockResponses: []MockResponse{
				{"HEAD", "/user", 401, `bad credentials`},
			},
			accessToken:   "github-token",
			scopes:        []Scope{ScopeRepo, ScopeUser},
			expectedError: "HEAD /user 401: ",
		},
		{
			name: "MissingScope",
			mockResponses: []MockResponse{
				{"HEAD", "/user", 200, ``},
			},
			accessToken:   "github-token",
			scopes:        []Scope{ScopeRepo, ScopeUser},
			expectedError: "access token does not have the scope: repo",
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"HEAD", "/user", 200, ``},
			},
			accessToken:   "github-token",
			scopes:        []Scope{},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gh := &GitHub{
				client:      new(http.Client),
				accessToken: tc.accessToken,
			}

			if len(tc.mockResponses) > 0 {
				ts := createMockHTTPServer(tc.mockResponses...)
				defer ts.Close()
				gh.apiURL = ts.URL
			}

			err := gh.checkScopes(tc.scopes...)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestGitHub_GetLatestRelease(t *testing.T) {
	tests := []struct {
		name            string
		mockResponses   []MockResponse
		accessToken     string
		ctx             context.Context
		owner, repo     string
		expectedRelease *Release
		expectedError   string
	}{
		{
			name:            "RequestError",
			mockResponses:   []MockResponse{},
			accessToken:     "github-token",
			ctx:             nil,
			owner:           "octocat",
			repo:            "Hello-World",
			expectedRelease: nil,
			expectedError:   "net/http: nil Context",
		},
		{
			name: "InvalidStatusCode",
			mockResponses: []MockResponse{
				{"GET", "/repos/{owner}/{repo}/releases/latest", 400, `bad request`},
			},
			accessToken:     "github-token",
			ctx:             context.Background(),
			owner:           "octocat",
			repo:            "Hello-World",
			expectedRelease: nil,
			expectedError:   "GET /repos/octocat/Hello-World/releases/latest 400: bad request",
		},
		{
			name: "InvalidResponse",
			mockResponses: []MockResponse{
				{"GET", "/repos/{owner}/{repo}/releases/latest", 200, `{`},
			},
			accessToken:     "github-token",
			ctx:             context.Background(),
			owner:           "octocat",
			repo:            "Hello-World",
			expectedRelease: nil,
			expectedError:   "unexpected EOF",
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"GET", "/repos/{owner}/{repo}/releases/latest", 200, mockReleaseBody},
			},
			accessToken:     "github-token",
			ctx:             context.Background(),
			owner:           "octocat",
			repo:            "Hello-World",
			expectedRelease: expectedRelease,
			expectedError:   "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gh := &GitHub{
				client:      new(http.Client),
				accessToken: tc.accessToken,
			}

			if len(tc.mockResponses) > 0 {
				ts := createMockHTTPServer(tc.mockResponses...)
				defer ts.Close()
				gh.apiURL = ts.URL
			}

			release, err := gh.GetLatestRelease(tc.ctx, tc.owner, tc.repo)

			if tc.expectedError != "" {
				assert.Nil(t, release)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedRelease, release)
			}
		})
	}
}

func TestGitHub_CreateRelease(t *testing.T) {
	releaseData := ReleaseData{
		Name:       "v1.0.0",
		TagName:    "v1.0.0",
		Target:     "master",
		Draft:      false,
		Prerelease: false,
		Body:       "Description of the release",
	}

	tests := []struct {
		name            string
		mockResponses   []MockResponse
		accessToken     string
		ctx             context.Context
		owner, repo     string
		params          ReleaseData
		expectedRelease *Release
		expectedError   string
	}{
		{
			name:            "RequestError",
			mockResponses:   []MockResponse{},
			accessToken:     "github-token",
			ctx:             nil,
			owner:           "octocat",
			repo:            "Hello-World",
			params:          releaseData,
			expectedRelease: nil,
			expectedError:   "net/http: nil Context",
		},
		{
			name: "InvalidStatusCode",
			mockResponses: []MockResponse{
				{"POST", "/repos/{owner}/{repo}/releases", 400, `bad request`},
			},
			accessToken:     "github-token",
			ctx:             context.Background(),
			owner:           "octocat",
			repo:            "Hello-World",
			params:          releaseData,
			expectedRelease: nil,
			expectedError:   "POST /repos/octocat/Hello-World/releases 400: bad request",
		},
		{
			name: "InvalidResponse",
			mockResponses: []MockResponse{
				{"POST", "/repos/{owner}/{repo}/releases", 201, `{`},
			},
			accessToken:     "github-token",
			ctx:             context.Background(),
			owner:           "octocat",
			repo:            "Hello-World",
			params:          releaseData,
			expectedRelease: nil,
			expectedError:   "unexpected EOF",
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"POST", "/repos/{owner}/{repo}/releases", 201, mockReleaseBody},
			},
			accessToken:     "github-token",
			ctx:             context.Background(),
			owner:           "octocat",
			repo:            "Hello-World",
			params:          releaseData,
			expectedRelease: expectedRelease,
			expectedError:   "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gh := &GitHub{
				client:      new(http.Client),
				accessToken: tc.accessToken,
			}

			if len(tc.mockResponses) > 0 {
				ts := createMockHTTPServer(tc.mockResponses...)
				defer ts.Close()
				gh.apiURL = ts.URL
			}

			release, err := gh.CreateRelease(tc.ctx, tc.owner, tc.repo, tc.params)

			if tc.expectedError != "" {
				assert.Nil(t, release)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedRelease, release)
			}
		})
	}
}

func TestGitHub_UpdateRelease(t *testing.T) {
	releaseData := ReleaseData{
		Name:       "v1.0.0",
		TagName:    "v1.0.0",
		Target:     "master",
		Draft:      false,
		Prerelease: false,
		Body:       "Description of the release",
	}

	tests := []struct {
		name            string
		mockResponses   []MockResponse
		accessToken     string
		ctx             context.Context
		owner, repo     string
		releaseID       int
		params          ReleaseData
		expectedRelease *Release
		expectedError   string
	}{
		{
			name:            "RequestError",
			mockResponses:   []MockResponse{},
			accessToken:     "github-token",
			ctx:             nil,
			owner:           "octocat",
			repo:            "Hello-World",
			releaseID:       1,
			params:          releaseData,
			expectedRelease: nil,
			expectedError:   "net/http: nil Context",
		},
		{
			name: "InvalidStatusCode",
			mockResponses: []MockResponse{
				{"PATCH", "/repos/{owner}/{repo}/releases/{release_id}", 400, `bad request`},
			},
			accessToken:     "github-token",
			ctx:             context.Background(),
			owner:           "octocat",
			repo:            "Hello-World",
			releaseID:       1,
			params:          releaseData,
			expectedRelease: nil,
			expectedError:   "PATCH /repos/octocat/Hello-World/releases/1 400: bad request",
		},
		{
			name: "InvalidResponse",
			mockResponses: []MockResponse{
				{"PATCH", "/repos/{owner}/{repo}/releases/{release_id}", 200, `{`},
			},
			accessToken:     "github-token",
			ctx:             context.Background(),
			owner:           "octocat",
			repo:            "Hello-World",
			releaseID:       1,
			params:          releaseData,
			expectedRelease: nil,
			expectedError:   "unexpected EOF",
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"PATCH", "/repos/{owner}/{repo}/releases/{release_id}", 200, mockReleaseBody},
			},
			accessToken:     "github-token",
			ctx:             context.Background(),
			owner:           "octocat",
			repo:            "Hello-World",
			releaseID:       1,
			params:          releaseData,
			expectedRelease: expectedRelease,
			expectedError:   "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gh := &GitHub{
				client:      new(http.Client),
				accessToken: tc.accessToken,
			}

			if len(tc.mockResponses) > 0 {
				ts := createMockHTTPServer(tc.mockResponses...)
				defer ts.Close()
				gh.apiURL = ts.URL
			}

			release, err := gh.UpdateRelease(tc.ctx, tc.owner, tc.repo, tc.releaseID, tc.params)

			if tc.expectedError != "" {
				assert.Nil(t, release)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedRelease, release)
			}
		})
	}
}

func TestGitHub_EnableBranchProtection(t *testing.T) {
	tests := []struct {
		name          string
		mockResponses []MockResponse
		accessToken   string
		ctx           context.Context
		owner, repo   string
		branch        string
		expectedError string
	}{
		{
			name:          "RequestError",
			mockResponses: []MockResponse{},
			accessToken:   "github-token",
			ctx:           nil,
			owner:         "octocat",
			repo:          "Hello-World",
			branch:        "master",
			expectedError: "net/http: nil Context",
		},
		{
			name: "InvalidStatusCode",
			mockResponses: []MockResponse{
				{"POST", "/repos/{owner}/{repo}/branches/{branch}/protection/enforce_admins", 400, `bad request`},
			},
			accessToken:   "github-token",
			ctx:           context.Background(),
			owner:         "octocat",
			repo:          "Hello-World",
			branch:        "master",
			expectedError: "POST /repos/octocat/Hello-World/branches/master/protection/enforce_admins 400: bad request",
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"POST", "/repos/{owner}/{repo}/branches/{branch}/protection/enforce_admins", 200, ``},
			},
			accessToken:   "github-token",
			ctx:           context.Background(),
			owner:         "octocat",
			repo:          "Hello-World",
			branch:        "master",
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gh := &GitHub{
				client:      new(http.Client),
				accessToken: tc.accessToken,
			}

			if len(tc.mockResponses) > 0 {
				ts := createMockHTTPServer(tc.mockResponses...)
				defer ts.Close()
				gh.apiURL = ts.URL
			}

			err := gh.EnableBranchProtection(tc.ctx, tc.owner, tc.repo, tc.branch)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestGitHub_DisableBranchProtection(t *testing.T) {
	tests := []struct {
		name          string
		mockResponses []MockResponse
		accessToken   string
		ctx           context.Context
		owner, repo   string
		branch        string
		expectedError string
	}{
		{
			name:          "RequestError",
			mockResponses: []MockResponse{},
			accessToken:   "github-token",
			ctx:           nil,
			owner:         "octocat",
			repo:          "Hello-World",
			branch:        "master",
			expectedError: "net/http: nil Context",
		},
		{
			name: "InvalidStatusCode",
			mockResponses: []MockResponse{
				{"DELETE", "/repos/{owner}/{repo}/branches/{branch}/protection/enforce_admins", 400, `bad request`},
			},
			accessToken:   "github-token",
			ctx:           context.Background(),
			owner:         "octocat",
			repo:          "Hello-World",
			branch:        "master",
			expectedError: "DELETE /repos/octocat/Hello-World/branches/master/protection/enforce_admins 400: bad request",
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"DELETE", "/repos/{owner}/{repo}/branches/{branch}/protection/enforce_admins", 204, ``},
			},
			accessToken:   "github-token",
			ctx:           context.Background(),
			owner:         "octocat",
			repo:          "Hello-World",
			branch:        "master",
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gh := &GitHub{
				client:      new(http.Client),
				accessToken: tc.accessToken,
			}

			if len(tc.mockResponses) > 0 {
				ts := createMockHTTPServer(tc.mockResponses...)
				defer ts.Close()
				gh.apiURL = ts.URL
			}

			err := gh.DisableBranchProtection(tc.ctx, tc.owner, tc.repo, tc.branch)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestGitHub_UploadReleaseAsset(t *testing.T) {
	tests := []struct {
		name                 string
		mockResponses        []MockResponse
		accessToken          string
		ctx                  context.Context
		uploadURL            string
		file                 string
		expectedReleaseAsset *ReleaseAsset
		expectedError        string
	}{
		{
			name:                 "NoFile",
			mockResponses:        []MockResponse{},
			accessToken:          "github-token",
			ctx:                  nil,
			uploadURL:            "https://uploads.github.com/repos/octocat/Hello-World/releases/1/assets",
			file:                 "unknown",
			expectedReleaseAsset: nil,
			expectedError:        "open unknown: no such file or directory",
		},
		{
			name:                 "BadFile",
			mockResponses:        []MockResponse{},
			accessToken:          "github-token",
			ctx:                  nil,
			uploadURL:            "https://uploads.github.com/repos/octocat/Hello-World/releases/1/assets",
			file:                 "/dev/null",
			expectedReleaseAsset: nil,
			expectedError:        "EOF",
		},
		{
			name:                 "RequestError",
			mockResponses:        []MockResponse{},
			accessToken:          "github-token",
			ctx:                  nil,
			uploadURL:            "https://uploads.github.com/repos/octocat/Hello-World/releases/1/assets",
			file:                 "github_test.go",
			expectedReleaseAsset: nil,
			expectedError:        "net/http: nil Context",
		},
		{
			name: "InvalidStatusCode",
			mockResponses: []MockResponse{
				{"POST", "/repos/{owner}/{repo}/releases/{release_id}/assets", 400, `bad request`},
			},
			accessToken:          "github-token",
			ctx:                  context.Background(),
			uploadURL:            "https://uploads.github.com/repos/octocat/Hello-World/releases/1/assets",
			file:                 "github_test.go",
			expectedReleaseAsset: nil,
			expectedError:        "POST /repos/octocat/Hello-World/releases/1/assets 400: bad request",
		},
		{
			name: "InvalidResponse",
			mockResponses: []MockResponse{
				{"POST", "/repos/{owner}/{repo}/releases/{release_id}/assets", 201, `{`},
			},
			accessToken:          "github-token",
			ctx:                  context.Background(),
			uploadURL:            "https://uploads.github.com/repos/octocat/Hello-World/releases/1/assets",
			file:                 "github_test.go",
			expectedReleaseAsset: nil,
			expectedError:        "unexpected EOF",
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"POST", "/repos/{owner}/{repo}/releases/{release_id}/assets", 201, mockReleaseAssetBody},
			},
			accessToken:          "github-token",
			ctx:                  context.Background(),
			uploadURL:            "https://uploads.github.com/repos/octocat/Hello-World/releases/1/assets",
			file:                 "github_test.go",
			expectedReleaseAsset: expectedReleaseAsset,
			expectedError:        "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gh := &GitHub{
				client:      new(http.Client),
				accessToken: tc.accessToken,
			}

			if len(tc.mockResponses) > 0 {
				ts := createMockHTTPServer(tc.mockResponses...)
				defer ts.Close()
				tc.uploadURL = strings.Replace(tc.uploadURL, "https://uploads.github.com", ts.URL, 1)
			}

			releaseAsset, err := gh.UploadReleaseAsset(tc.ctx, tc.uploadURL, tc.file)

			if tc.expectedError != "" {
				assert.Nil(t, releaseAsset)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedReleaseAsset, releaseAsset)
			}
		})
	}
}

func TestGitHub_DownloadReleaseAsset(t *testing.T) {
	tests := []struct {
		name          string
		mockResponses []MockResponse
		accessToken   string
		ctx           context.Context
		downloadURL   string
		expectedError string
	}{
		{
			name:          "RequestError",
			mockResponses: []MockResponse{},
			accessToken:   "github-token",
			ctx:           nil,
			downloadURL:   "https://github.com/octocat/Hello-World/releases/download/v1.0.0/example.zip",
			expectedError: "net/http: nil Context",
		},
		{
			name: "InvalidStatusCode",
			mockResponses: []MockResponse{
				{"GET", "/{owner}/{repo}/releases/download/{release_name}/{asset_name}", 400, `bad request`},
			},
			accessToken:   "github-token",
			ctx:           context.Background(),
			downloadURL:   "https://github.com/octocat/Hello-World/releases/download/v1.0.0/example.zip",
			expectedError: "GET /octocat/Hello-World/releases/download/v1.0.0/example.zip 400: bad request",
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"GET", "/{owner}/{repo}/releases/download/{release_name}/{asset_name}", 200, ``},
			},
			accessToken:   "github-token",
			ctx:           context.Background(),
			downloadURL:   "https://github.com/octocat/Hello-World/releases/download/v1.0.0/example.zip",
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gh := &GitHub{
				client:      new(http.Client),
				accessToken: tc.accessToken,
			}

			if len(tc.mockResponses) > 0 {
				ts := createMockHTTPServer(tc.mockResponses...)
				defer ts.Close()
				tc.downloadURL = strings.Replace(tc.downloadURL, "https://github.com", ts.URL, 1)
			}

			file, err := ioutil.TempFile("", "test")
			assert.NoError(t, err)
			file.Close()
			defer os.Remove(file.Name())

			err = gh.DownloadReleaseAsset(tc.ctx, tc.downloadURL, file.Name())

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
