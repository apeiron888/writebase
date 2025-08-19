package domain

import "testing"

func TestValidatePagination_Defaults(t *testing.T) {
	p := &Pagination{}
	p.ValidatePagination()
	if p.Page != 1 || p.PageSize != 10 || p.SortOrder != "" {
		t.Fatalf("unexpected %+v", *p)
	}
}

func TestValidatePagination_MaxAndOrder(t *testing.T) {
	p := &Pagination{Page: -5, PageSize: 1000, SortOrder: "wrong"}
	p.ValidatePagination()
	if p.Page != 1 {
		t.Fatal("page default")
	}
	if p.PageSize != 100 {
		t.Fatalf("page size capped: %d", p.PageSize)
	}
	if p.SortOrder != "desc" {
		t.Fatalf("sort order fixed: %s", p.SortOrder)
	}
}
