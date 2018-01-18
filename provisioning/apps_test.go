package provisioning

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/nextcloud/nextcloudgo"
)

func TestGetList(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/ocs/v2.php/cloud/apps?filter=enabled" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"ocs":{"meta":{"status":"ok","statuscode":200,"message":"OK"},"data":{"apps":["files","dav"]}}}`)
			return
		} else if r.RequestURI == "/ocs/v2.php/cloud/apps?filter=disabled" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"ocs":{"meta":{"status":"ok","statuscode":200,"message":"OK"},"data":{"apps":["comments","twofactor_backupcodes"]}}}`)
			return
		} else if r.RequestURI == "/ocs/v2.php/cloud/apps" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"ocs":{"meta":{"status":"ok","statuscode":200,"message":"OK"},"data":{"apps":["files","dav","comments","twofactor_backupcodes","activity"]}}}`)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	nc := nextcloudgo.NextcloudGo{ServerURL: ts.URL, User: "admin", Password: "admin"}
	api := New(nc)

	l, err := api.GetApps("all")
	if err != nil {
		t.Error(err.Error())
	}
	if !reflect.DeepEqual(l, []string{"files", "dav", "comments", "twofactor_backupcodes", "activity"}) {
		t.Error("List of all apps didn't match")
	}

	l, err = api.GetApps("enabled")
	if err != nil {
		t.Error(err.Error())
	}
	if !reflect.DeepEqual(l, []string{"files", "dav"}) {
		t.Error("List of enabled apps didn't match")
	}

	l, err = api.GetApps("disabled")
	if err != nil {
		t.Error(err.Error())
	}
	if !reflect.DeepEqual(l, []string{"comments", "twofactor_backupcodes"}) {
		t.Error("List of disabled apps didn't match")
	}
}

func TestGetListInvalidFilter(t *testing.T) {
	nc := nextcloudgo.NextcloudGo{ServerURL: "", User: "admin", Password: "admin"}
	api := New(nc)
	l, err := api.GetApps("foo")
	if err == nil {
		t.Error("Missing expected error on invalid filter")
	} else if err.Error() != "Invalid filter given" {
		t.Error(err.Error())
	}
	if !reflect.DeepEqual(l, []string{}) {
		t.Error("List of invalid filter should be empty")
	}
}

func TestGetListError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	nc := nextcloudgo.NextcloudGo{ServerURL: ts.URL, User: "admin", Password: "admin"}
	api := New(nc)

	l, err := api.GetApps("all")
	if err == nil {
		t.Error("Missing expected error on server failure")
	} else if err.Error() != "An error occured while searching for apps" {
		t.Error(err.Error())
	}
	if !reflect.DeepEqual(l, []string{}) {
		t.Error("List should be empty on server error")
	}
}

func TestGetListWrongRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, `{"ocs":{"meta":{"status":"error","statuscode":400,"message":"failure"},"data":[]}}`)
	}))
	defer ts.Close()

	nc := nextcloudgo.NextcloudGo{ServerURL: ts.URL, User: "admin", Password: "admin"}
	api := New(nc)

	l, err := api.GetApps("all")
	if err == nil {
		t.Error("Missing expected error on server failure")
	} else if err.Error() != "An error occured while searching for apps" {
		t.Error(err.Error())
	}
	if !reflect.DeepEqual(l, []string{}) {
		t.Error("List should be empty on server error")
	}
}

func TestIsAppEnabled(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.RequestURI == "/ocs/v2.php/cloud/apps?filter=enabled" {
			fmt.Fprintln(w, `{"ocs":{"meta":{"status":"ok","statuscode":200,"message":"OK"},"data":{"apps":["files","dav"]}}}`)
		}
	}))
	defer ts.Close()

	nc := nextcloudgo.NextcloudGo{ServerURL: ts.URL, User: "admin", Password: "admin"}
	api := New(nc)

	if ok, err := api.IsAppEnabled("files"); err != nil || !ok {
		t.Error("Files app should be enabled")
	}

	if ok, err := api.IsAppEnabled("comments"); err != nil || ok {
		t.Error("Comments app should not be enabled")
	}
}

func TestIsAppDisabled(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.RequestURI == "/ocs/v2.php/cloud/apps?filter=disabled" {
			fmt.Fprintln(w, `{"ocs":{"meta":{"status":"ok","statuscode":200,"message":"OK"},"data":{"apps":["comments","twofactor_backupcodes"]}}}`)
		}
	}))
	defer ts.Close()

	nc := nextcloudgo.NextcloudGo{ServerURL: ts.URL, User: "admin", Password: "admin"}
	api := New(nc)

	if ok, err := api.IsAppDisabled("comments"); err != nil || !ok {
		t.Error("Comments app should be disabled")
	}

	if ok, err := api.IsAppDisabled("files"); err != nil || ok {
		t.Error("Files app should not be disabled")
	}

}

func TestIsAppAvailable(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.RequestURI == "/ocs/v2.php/cloud/apps" {
			fmt.Fprintln(w, `{"ocs":{"meta":{"status":"ok","statuscode":200,"message":"OK"},"data":{"apps":["files","dav","comments","twofactor_backupcodes","activity"]}}}`)
		}
	}))
	defer ts.Close()

	nc := nextcloudgo.NextcloudGo{ServerURL: ts.URL, User: "admin", Password: "admin"}
	api := New(nc)

	if ok, err := api.IsAppAvailable("files"); err != nil || !ok {
		t.Error("Files app should be available")
	}

	if ok, err := api.IsAppAvailable("comments"); err != nil || !ok {
		t.Error("Comments app should be available")
	}

	if ok, err := api.IsAppAvailable("activity"); err != nil || !ok {
		t.Error("Activity app should be available")
	}

	if ok, err := api.IsAppAvailable("not-available"); err != nil || ok {
		t.Error("NotAvailable app should not be available")
	}
}

func TestIsAppAvailableError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	nc := nextcloudgo.NextcloudGo{ServerURL: ts.URL, User: "admin", Password: "admin"}
	api := New(nc)

	ok, err := api.IsAppAvailable("files")
	if ok || err == nil {
		t.Error("Should receive an error on server error")
	} else if err.Error() != "An error occured while searching for apps" {
		t.Error(err.Error())
	}
}
