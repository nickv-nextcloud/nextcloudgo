package nextcloudgo

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConnect(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.RequestURI == "/status.php" {
			fmt.Fprintln(w, `{"installed":true,"maintenance":false,"needsDbUpgrade":false,"version":"13.0.0.6","versionstring":"13.0.0 Beta 1","edition":"","productname":"Nextcloud"}`)
			return
		}
	}))
	defer ts.Close()

	nc := NextcloudGo{ServerURL: ts.URL, User: "admin", Password: "admin"}
	if !nc.isConnected() {
		t.Error("Server URL not set")
	}
	if !nc.isLoggedIn() {
		t.Error("User name or password not set")
	}
}

func TestStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.RequestURI == "/status.php" {
			fmt.Fprintln(w, `{"installed":true,"maintenance":false,"needsDbUpgrade":false,"version":"13.0.0.6","versionstring":"13.0.0 Beta 1","edition":"","productname":"Nextcloud"}`)
			return
		}
	}))
	defer ts.Close()

	nc := NextcloudGo{ServerURL: ts.URL, User: "admin", Password: "admin"}
	s, err := nc.Status()
	if err != nil {
		t.Error("Could not get status from server")
	}
	if !s.Installed {
		t.Error("Server should be installed")
	}
	if s.Maintenance {
		t.Error("Server should be not be in maintenance")
	}
	if s.Version != "13.0.0.6" {
		t.Error("Version was not extracted correctly")
	}
}

func TestStatusNotConnected(t *testing.T) {
	nc := NextcloudGo{ServerURL: "", User: "admin", Password: "admin"}
	s, err := nc.Status()
	if err != ErrNotConnected {
		t.Error("Should receive ErrNotConnected error when calling status without ServerURL")
	}
	if s.Installed {
		t.Error("Server is not installed")
	}
	if !s.Maintenance {
		t.Error("Server is in maintenance")
	}
	if s.Version != "" {
		t.Error("Version should be empty")
	}
}

func TestStatusFailed(t *testing.T) {
	nc := NextcloudGo{ServerURL: "https://ihopenooneeverregistersthissubdomain.nextcloud.com", User: "admin", Password: "admin"}
	s, err := nc.Status()
	if err == nil {
		t.Error("Should receive an error when calling status with an invalid ServerURL")
	}
	if s.Installed {
		t.Error("Server is not installed")
	}
	if !s.Maintenance {
		t.Error("Server is in maintenance")
	}
	if s.Version != "" {
		t.Error("Version should be empty")
	}
}
