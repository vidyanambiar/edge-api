package imagebuilder

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/redhatinsights/edge-api/config"
	"github.com/redhatinsights/edge-api/pkg/models"
)

func setUp() {
	config.Init()
	config.Get().Debug = true
}

func tearDown() {

}

func TestMain(m *testing.M) {
	setUp()
	retCode := m.Run()
	tearDown()
	os.Exit(retCode)
}
func TestInitClient(t *testing.T) {
	InitClient()
	if ClientInstance == nil {
		t.Errorf("Client shouldnt be nil")
	}
}

func TestComposeImage(t *testing.T) {
	config.Init()

	InitClient()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, `{"id": "compose-job-id-returned-from-image-builder"}`)
	}))
	defer ts.Close()
	config.Get().ImageBuilderConfig.URL = ts.URL

	pkgs := []models.Package{
		{
			Name: "vim",
		},
		{
			Name: "ansible",
		},
	}
	img := &models.Image{Distribution: "rhel-8", Commit: &models.Commit{
		Arch:     "x86_64",
		Packages: pkgs,
	}}
	headers := make(map[string]string)
	img, err := ClientInstance.ComposeCommit(img, headers)
	if err != nil {
		t.Errorf("Shouldnt throw error")
	}
	if img == nil {
		t.Errorf("Image shouldnt be nil")
	}
	if img != nil && img.Commit.ComposeJobID != "compose-job-id-returned-from-image-builder" {
		t.Error("Compose job is not correct")
	}
}
