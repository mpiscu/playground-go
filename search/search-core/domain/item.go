package domain

type SearchItem struct {
    ID string
    Tags[] string
    Note string
}

type URLSearchItem struct {
    SearchItem
    URL string
}
