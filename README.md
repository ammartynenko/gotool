# gotool - paginate 
pagination pages based on ORM gorm

Example use:
___
```go
//inset to your instanse http server  instanse 
	pagparams := &paginate.Params{
		Limit:      s.Config.PaginateCountOnPage,
		DBS:        s.DB,
		DebugQuery: s.Config.PaginateDebug,
		SortTypes:  s.Config.PaginateSortType,
	}
	s.Paginator = paginate.NewPaginate(pagparams)
```
In handler  use for get paginate result 
```go 
func (s *Server) HandlerAdminUsers(w http.ResponseWriter, r *http.Request) {
	data := s.GetDataFromContext(r)
	command := chi.URLParam(r, "command")
	id := chi.URLParam(r, "id")


	switch command {
	case "", "page":
		var users []User
		s.Paginator.MakePaginate(s.Convert.DirectStringtoInt(id), &users)
		data.Stock.Put("stock","paginate", s.Paginator)
```

In template use simple logic + css style for show result 
```html 
<!-- pagination -->
<div class="pagination">
        <a href="/admin/users/page/{{$paginate.Help.Predpage}}">&laquo;</a>
            {{range $i, $x := $paginate.Help.List}}
                {{if or (and (eq $id "") (eq $x "1")) (eq $id $x)}}
                    <a href="/admin/users/page/{{$x}}" class="active">{{$x}}</a>
                {{else}}
                    <a href="/admin/users/page/{{$x}}">{{$x}}</a>
                {{end}}
            {{end}}
        <a href="/admin/users/page/{{$paginate.Help.Nextpage}}">&raquo;</a>
</div>
```
and simple css style for pagination  for example 

```css 
/******************************************************************
 PAGINATE
 *******************************************************************/
 .pagination {
     display: flex;
     justify-content: center;
 }

 .pagination a {
     color: black;
     float: left;
     padding: 8px 16px;
     text-decoration: none;
     border: 1px solid #ddd;
 }

 .pagination a.active {
     background-color: #4CAF50;
     color: white;
     border: 1px solid #4CAF50;
 }

 .pagination a:hover:not(.active) {background-color: #ddd;}

 .pagination a:first-child {
     border-top-left-radius: 5px;
     border-bottom-left-radius: 5px;
 }

 .pagination a:last-child {
     border-top-right-radius: 5px;
     border-bottom-right-radius: 5px;
 }

```
Enjoy!
wbr//spouk
___