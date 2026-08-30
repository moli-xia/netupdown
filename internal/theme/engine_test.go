package theme

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestBuiltinThemeLoads(t *testing.T) {
	e, err := New(t.TempDir(), false, "aurora")
	if err != nil {
		t.Fatal(err)
	}
	items, err := e.List()
	if err != nil || len(items) != 1 || !items[0].Builtin || !items[0].Active {
		t.Fatalf("unexpected themes: %#v, %v", items, err)
	}
	raw, contentType, err := e.Static("aurora", "css/main.css")
	if err != nil || len(raw) == 0 || contentType != "text/css; charset=utf-8" {
		t.Fatalf("static asset: %d %q %v", len(raw), contentType, err)
	}
}

func TestThemeInstallActivateAndRender(t *testing.T) {
	root := t.TempDir()
	archive := filepath.Join(root, "clean.zip")
	f, err := os.Create(archive)
	if err != nil {
		t.Fatal(err)
	}
	zw := zip.NewWriter(f)
	files := map[string]string{
		"theme.json":                  `{"id":"clean","name":"Clean","version":"1.0.0"}`,
		"templates/layouts/base.html": `{{define "base"}}<html>{{template "content" .}}</html>{{end}}`,
		"templates/index.html":        `{{define "content"}}clean theme{{end}}`,
		"static/css/main.css":         `body{color:green}`,
	}
	for name, body := range files {
		w, _ := zw.Create(name)
		_, _ = w.Write([]byte(body))
	}
	_ = zw.Close()
	_ = f.Close()
	e, err := New(filepath.Join(root, "data"), false, "aurora")
	if err != nil {
		t.Fatal(err)
	}
	meta, err := e.Install(archive)
	if err != nil || meta.ID != "clean" {
		t.Fatalf("install: %#v %v", meta, err)
	}
	if err := e.Activate("clean"); err != nil {
		t.Fatal(err)
	}
	rendered, err := e.Render("index", map[string]any{})
	if err != nil || !bytes.Contains(rendered, []byte("clean theme")) {
		t.Fatalf("rendered %q: %v", rendered, err)
	}
}

func TestThemeInstallRejectsZipSlip(t *testing.T) {
	root := t.TempDir()
	archive := filepath.Join(root, "unsafe.zip")
	f, err := os.Create(archive)
	if err != nil {
		t.Fatal(err)
	}
	zw := zip.NewWriter(f)
	meta, _ := zw.Create("theme.json")
	_, _ = meta.Write([]byte(`{"id":"unsafe","name":"Unsafe","version":"1.0.0"}`))
	evil, _ := zw.Create("../evil.txt")
	_, _ = evil.Write([]byte("no"))
	_ = zw.Close()
	_ = f.Close()
	e, err := New(filepath.Join(root, "data"), false, "aurora")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := e.Install(archive); err == nil {
		t.Fatal("zip-slip archive was accepted")
	}
	if _, err := os.Stat(filepath.Join(root, "data", "evil.txt")); !os.IsNotExist(err) {
		t.Fatal("zip-slip wrote outside theme directory")
	}
}
