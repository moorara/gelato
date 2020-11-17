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
	mockPermissionBody = `{
		"permission": "admin",
		"user": {
			"id": 1,
			"login": "octocat",
			"type": "User"
		}
	}`

	mockUserBody = `{
		"id": 1,
		"login": "octocat",
		"type": "User",
		"email": "octocat@github.com",
		"name": "monalisa octocat",
		"company": "GitHub",
		"location": "San Francisco",
		"url": "https://api.github.com/users/octocat",
		"html_url": "https://github.com/octocat",
		"avatar_url": "https://github.com/images/error/octocat_happy.gif",
		"gravatar_id": "",
		"repos_url": "https://api.github.com/users/octocat/repos",
		"organizations_url": "https://api.github.com/users/octocat/orgs"
	}`

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
	expectedUser = &User{
		ID:         1,
		Login:      "octocat",
		Type:       "User",
		Email:      "octocat@github.com",
		Name:       "monalisa octocat",
		Company:    "GitHub",
		Location:   "San Francisco",
		URL:        "https://api.github.com/users/octocat",
		HTMLURL:    "https://github.com/octocat",
		AvatarURL:  "https://github.com/images/error/octocat_happy.gif",
		GravatarID: "",
		ReposURL:   "https://api.github.com/users/octocat/repos",
		OrgsURL:    "https://api.github.com/users/octocat/orgs",
	}

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
	Method         string
	Path           string
	StatusCode     int
	ResponseHeader http.Header
	ResponseBody   string
}

func createMockHTTPServer(mocks ...MockResponse) *httptest.Server {
	r := mux.NewRouter()
	for _, m := range mocks {
		m := m
		r.Methods(m.Method).Path(m.Path).HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			for k, vals := range m.ResponseHeader {
				for _, v := range vals {
					w.Header().Add(k, v)
				}
			}

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
		expectedError string
	}{
		{
			name:          "OK",
			accessToken:   "github-token",
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gh, err := New(tc.accessToken)

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
				{"GET", "/users/{username}", 400, http.Header{}, `bad request`},
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
					"GET", "/users/{username}", 200, http.Header{}, `{
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

func TestGitHub_GetScopes(t *testing.T) {
	tests := []struct {
		name           string
		mockResponses  []MockResponse
		accessToken    string
		ctx            context.Context
		expectedScopes []Scope
		expectedError  string
	}{
		{
			name:           "RequestError",
			mockResponses:  []MockResponse{},
			accessToken:    "github-token",
			ctx:            nil,
			expectedScopes: nil,
			expectedError:  "net/http: nil Context",
		},
		{
			name: "InvalidStatusCode",
			mockResponses: []MockResponse{
				{"HEAD", "/user", 401, http.Header{}, `bad credentials`},
			},
			accessToken:    "github-token",
			ctx:            context.Background(),
			expectedScopes: nil,
			expectedError:  "HEAD /user 401: ",
		},
		{
			name: "SuccessWithoutScope",
			mockResponses: []MockResponse{
				{"HEAD", "/user", 200, http.Header{}, ``},
			},
			accessToken:    "github-token",
			ctx:            context.Background(),
			expectedScopes: []Scope{},
			expectedError:  "",
		},
		{
			name: "SuccessWithScopes",
			mockResponses: []MockResponse{
				{"HEAD", "/user", 200, http.Header{"X-OAuth-Scopes": []string{"repo, user"}}, ``},
			},
			accessToken:    "github-token",
			ctx:            context.Background(),
			expectedScopes: []Scope{ScopeRepo, ScopeUser},
			expectedError:  "",
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

			scopes, err := gh.GetScopes(tc.ctx)

			if tc.expectedError != "" {
				assert.Nil(t, scopes)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedScopes, scopes)
			}
		})
	}
}

func TestGitHub_GetRepoPermission(t *testing.T) {
	tests := []struct {
		name               string
		mockResponses      []MockResponse
		accessToken        string
		ctx                context.Context
		owner, repo        string
		username           string
		expectedPermission Permission
		expectedError      string
	}{
		{
			name:               "RequestError",
			mockResponses:      []MockResponse{},
			accessToken:        "github-token",
			ctx:                nil,
			owner:              "octocat",
			repo:               "Hello-World",
			username:           "octocat",
			expectedPermission: "",
			expectedError:      "net/http: nil Context",
		},
		{
			name: "InvalidStatusCode",
			mockResponses: []MockResponse{
				{"GET", "/repos/{owner}/{repo}/collaborators/{username}/permission", 401, http.Header{}, `bad credentials`},
			},
			accessToken:        "github-token",
			ctx:                context.Background(),
			owner:              "octocat",
			repo:               "Hello-World",
			username:           "octocat",
			expectedPermission: "",
			expectedError:      "GET /repos/octocat/Hello-World/collaborators/octocat/permission 401: bad credentials",
		},
		{
			name: "InvalidResponse",
			mockResponses: []MockResponse{
				{"GET", "/repos/{owner}/{repo}/collaborators/{username}/permission", 200, http.Header{}, `{`},
			},
			accessToken:        "github-token",
			ctx:                context.Background(),
			owner:              "octocat",
			repo:               "Hello-World",
			username:           "octocat",
			expectedPermission: "",
			expectedError:      "unexpected EOF",
		},
		{
			name: "SuccessWithoutScope",
			mockResponses: []MockResponse{
				{"GET", "/repos/{owner}/{repo}/collaborators/{username}/permission", 200, http.Header{}, mockPermissionBody},
			},
			accessToken:        "github-token",
			ctx:                context.Background(),
			owner:              "octocat",
			repo:               "Hello-World",
			username:           "octocat",
			expectedPermission: PermissionAdmin,
			expectedError:      "",
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

			permission, err := gh.GetRepoPermission(tc.ctx, tc.owner, tc.repo, tc.username)

			if tc.expectedError != "" {
				assert.Empty(t, permission)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedPermission, permission)
			}
		})
	}
}

func TestGitHub_GetUser(t *testing.T) {
	tests := []struct {
		name          string
		mockResponses []MockResponse
		accessToken   string
		ctx           context.Context
		expectedUser  *User
		expectedError string
	}{
		{
			name:          "RequestError",
			mockResponses: []MockResponse{},
			accessToken:   "github-token",
			ctx:           nil,
			expectedUser:  nil,
			expectedError: "net/http: nil Context",
		},
		{
			name: "InvalidStatusCode",
			mockResponses: []MockResponse{
				{"GET", "/user", 401, http.Header{}, `bad credentials`},
			},
			accessToken:   "github-token",
			ctx:           context.Background(),
			expectedUser:  nil,
			expectedError: "GET /user 401: bad credentials",
		},
		{
			name: "InvalidResponse",
			mockResponses: []MockResponse{
				{"GET", "/user", 200, http.Header{}, `{`},
			},
			accessToken:   "github-token",
			ctx:           context.Background(),
			expectedUser:  nil,
			expectedError: "unexpected EOF",
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"GET", "/user", 200, http.Header{}, mockUserBody},
			},
			accessToken:   "github-token",
			ctx:           context.Background(),
			expectedUser:  expectedUser,
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

			user, err := gh.GetUser(tc.ctx)

			if tc.expectedError != "" {
				assert.Nil(t, user)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedUser, user)
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
				{"GET", "/repos/{owner}/{repo}/releases/latest", 400, http.Header{}, `bad request`},
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
				{"GET", "/repos/{owner}/{repo}/releases/latest", 200, http.Header{}, `{`},
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
				{"GET", "/repos/{owner}/{repo}/releases/latest", 200, http.Header{}, mockReleaseBody},
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
				{"POST", "/repos/{owner}/{repo}/releases", 400, http.Header{}, `bad request`},
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
				{"POST", "/repos/{owner}/{repo}/releases", 201, http.Header{}, `{`},
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
				{"POST", "/repos/{owner}/{repo}/releases", 201, http.Header{}, mockReleaseBody},
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
				{"PATCH", "/repos/{owner}/{repo}/releases/{release_id}", 400, http.Header{}, `bad request`},
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
				{"PATCH", "/repos/{owner}/{repo}/releases/{release_id}", 200, http.Header{}, `{`},
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
				{"PATCH", "/repos/{owner}/{repo}/releases/{release_id}", 200, http.Header{}, mockReleaseBody},
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
				{"POST", "/repos/{owner}/{repo}/branches/{branch}/protection/enforce_admins", 400, http.Header{}, `bad request`},
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
				{"POST", "/repos/{owner}/{repo}/branches/{branch}/protection/enforce_admins", 200, http.Header{}, ``},
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
				{"DELETE", "/repos/{owner}/{repo}/branches/{branch}/protection/enforce_admins", 400, http.Header{}, `bad request`},
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
				{"DELETE", "/repos/{owner}/{repo}/branches/{branch}/protection/enforce_admins", 204, http.Header{}, ``},
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
				{"POST", "/repos/{owner}/{repo}/releases/{release_id}/assets", 400, http.Header{}, `bad request`},
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
				{"POST", "/repos/{owner}/{repo}/releases/{release_id}/assets", 201, http.Header{}, `{`},
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
				{"POST", "/repos/{owner}/{repo}/releases/{release_id}/assets", 201, http.Header{}, mockReleaseAssetBody},
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
				{"GET", "/{owner}/{repo}/releases/download/{release_name}/{asset_name}", 400, http.Header{}, `bad request`},
			},
			accessToken:   "github-token",
			ctx:           context.Background(),
			downloadURL:   "https://github.com/octocat/Hello-World/releases/download/v1.0.0/example.zip",
			expectedError: "GET /octocat/Hello-World/releases/download/v1.0.0/example.zip 400: bad request",
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"GET", "/{owner}/{repo}/releases/download/{release_name}/{asset_name}", 200, http.Header{}, ``},
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
