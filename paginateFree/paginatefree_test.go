package paginateFree

import "testing"

func TestPaginator_Paginate(t *testing.T) {
	 p := NewPaginator()
	 var (
	 	page = 0
	 	countLinks = 1
	 	countPage = 2
	 	a = []int{1,2,3,4,5,6,7,8,9,10,11,12}
	 )
	 if pr, err := p.Paginate(page, countPage, countLinks, a); err != nil {
	 	t.Error(err)
	 } else {
		 t.Logf("Result paginate: %v\nListPage: %v\n", pr, pr.ListPage)
	 }
	if pr, err := p.Paginate(page + 1 , countPage, countLinks, a); err != nil {
		t.Error(err)
	} else {
		t.Logf("Result paginate: %v\nListPage: %v\n", pr, pr.ListPage)
	}
	if pr, err := p.Paginate(page + 2 , countPage, countLinks, a); err != nil {
		t.Error(err)
	} else {
		t.Logf("Result paginate: %v\nListPage: %v\n", pr, pr.ListPage)
	}
}
