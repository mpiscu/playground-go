package domain

type SearchCriteria struct {
    ID string
    IncludeAny []string
    IncludeAll []string
    ExcludeAny []string
    ExcludeAll []string
    Limit int
}

type SearchResult struct {
    URLItems []URLSearchItem
}
