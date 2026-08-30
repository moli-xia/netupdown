package service

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/moli-xia/netupdown/internal/model"
	"github.com/moli-xia/netupdown/internal/repo"
	"gorm.io/gorm"
)

func TestVersionLess(t *testing.T) {
	cases := []struct {
		cur, target string
		want        bool
	}{{"1.0.0", "1.1.0", true}, {"v2.0.0", "1.9.0", false}, {"1.0.0-beta.1", "1.0.0", true}, {"same", "same", false}, {"old", "new", true}}
	for _, tc := range cases {
		if got := versionLess(tc.cur, nil, tc.target, nil); got != tc.want {
			t.Errorf("%s < %s = %v", tc.cur, tc.target, got)
		}
	}
}

func TestPublicCatalogOnlyReturnsSelfDevelopedApps(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:catalog-self-only?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.Category{}, &model.Tag{}, &model.App{}, &model.Release{}, &model.Asset{}, &model.DownloadSource{}); err != nil {
		t.Fatal(err)
	}
	rows := []model.App{
		{Name: "自研产品", Slug: "self-product", Type: model.AppTypeSelf, Status: model.StatusPublished},
		{Name: "第三方产品", Slug: "third-product", Type: model.AppTypeThird, Status: model.StatusPublished},
	}
	if err := db.Create(&rows).Error; err != nil {
		t.Fatal(err)
	}
	catalog := NewCatalog(repo.New(db), "http://localhost")
	result, err := catalog.Apps(context.Background(), AppQuery{Page: 1, PageSize: 20})
	if err != nil {
		t.Fatal(err)
	}
	apps := result.List.([]model.App)
	if result.Total != 1 || len(apps) != 1 || apps[0].Slug != "self-product" {
		t.Fatalf("public catalog = total %d, apps %#v", result.Total, apps)
	}
	if _, err := catalog.AppBySlug(context.Background(), "third-product", true); err == nil {
		t.Fatal("published third-party app must not be publicly addressable")
	}
	thirdResult, err := catalog.Apps(context.Background(), AppQuery{Page: 1, PageSize: 20, Type: "third"})
	if err != nil {
		t.Fatal(err)
	}
	if thirdResult.Total != 0 {
		t.Fatalf("legacy third-party filter returned %d public rows", thirdResult.Total)
	}
}
