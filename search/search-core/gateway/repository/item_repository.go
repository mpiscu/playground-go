package repository 

import(
    "crypto/sha1"
    "encoding/hex"
    "fmt"
    "github.com/blevesearch/bleve"
    "github.com/blevesearch/bleve/analysis/analyzer/keyword"
    "os"
    "os/user"
    "strings"
    "search/search-core/domain"
)

const (
    itemIndexName = "items"

    itemTypeURL="url"
)

type searchItem struct {
    Tag []string `json:"tag"`
    Type string `json:"type"`
    URL string `json:"url"`
    Note string `json:"note"`
}

type ItemRepository struct {
    itemIndex bleve.Index
}

func NewItemRepository() *ItemRepository  {
    return &ItemRepository{}
}

func (r *ItemRepository) Open() error {

    user, err := user.Current()
    if err != nil {
        return err
    }

    baseDir := user.HomeDir + string(os.PathSeparator) + ".search" 
    indexDir := baseDir + string(os.PathSeparator) + itemIndexName
    
    if _, err := os.Stat(baseDir); os.IsNotExist(err) {
	os.MkdirAll(baseDir, os.ModePerm)
    }

    itemIndex, err := bleve.Open(indexDir)
 
    if err == bleve.ErrorIndexPathDoesNotExist {


        // mappings
        keywordFieldMapping := bleve.NewTextFieldMapping()
        keywordFieldMapping.Analyzer = keyword.Name

        textNotIndexFieldMapping := bleve.NewTextFieldMapping()
        textNotIndexFieldMapping.Name = "stringNotIndex"
        textNotIndexFieldMapping.Analyzer = "en"
        textNotIndexFieldMapping.Index=false
        textNotIndexFieldMapping.Store = true
        textNotIndexFieldMapping.IncludeTermVectors = true
        textNotIndexFieldMapping.IncludeInAll = true

        //document
        itemMapping := bleve.NewDocumentMapping()

        //fields
        itemMapping.AddFieldMappingsAt("type", keywordFieldMapping)
        itemMapping.AddFieldMappingsAt("tag", keywordFieldMapping)
        itemMapping.AddFieldMappingsAt("url", keywordFieldMapping)
        itemMapping.AddFieldMappingsAt("note", textNotIndexFieldMapping)

        //index mapping
        indexMapping := bleve.NewIndexMapping()
        indexMapping.AddDocumentMapping("item", itemMapping)

        //create index
        itemIndex, err := bleve.New(indexDir, indexMapping)
        if err != nil {
            return err
        }

        r.itemIndex = itemIndex

        return nil

    }

    r.itemIndex = itemIndex

    return err

}

func (r *ItemRepository) Close() error {
    return r.itemIndex.Close()
}

func (r *ItemRepository) IndexURLItem(item *domain.URLSearchItem) (string, error) {
    id := generateID(item.URL)
    item.ID = id
    return id, r.itemIndex.Index(id, searchItem {
        Type: itemTypeURL,
        Tag: item.Tags,
        Note: item.Note,
        URL: item.URL,
    })
}

func (r *ItemRepository) Delete(id string) error {
    return r.itemIndex.Delete(id)
}

func (r *ItemRepository) GetByID(id string) (*domain.SearchResult, error) {
    doc, err := r.itemIndex.Document(id)
    if err != nil {
        return nil, err
    }
    if doc == nil {
        return nil, domain.ErrNotFound
    }

    result := &domain.SearchResult{}
    var typeValue string
    var noteValue string
    var urlValue string
    var tagsValue []string
    for _, field := range doc.Fields {

        switch field.Name() {
            case "type": typeValue = string(field.Value())
            case "tag": tagsValue = append(tagsValue, string(field.Value()))
            case "url" : urlValue = string(field.Value())
            case "note" : noteValue = string(field.Value())
        }
    }

    searchItem := domain.SearchItem { ID: doc.ID, Tags: tagsValue, Note: noteValue} 


    switch typeValue {
        case itemTypeURL: result.URLItems = append(result.URLItems, domain.URLSearchItem{ SearchItem: searchItem, URL: urlValue})
        default: return nil, fmt.Errorf("Item type %s is not recognized", typeValue)
    }
 
    return result, nil
} 

func (r *ItemRepository) Search(criteria *domain.SearchCriteria) (*domain.SearchResult, error) {
    if criteria.ID != "" {
        return r.GetByID(criteria.ID)
    }

    searchQuery := bleve.NewBooleanQuery()

    includeQuery := bleve.NewDisjunctionQuery()
    searchQuery.AddMust(includeQuery)

    if len(criteria.IncludeAny)>0 {
        includeAnyQuery := bleve.NewDisjunctionQuery()
        for _, tag := range criteria.IncludeAny {
            includeTagQuery := bleve.NewMatchQuery(tag)
            includeTagQuery.SetField("tag")
            includeAnyQuery.AddQuery(includeTagQuery)
        }
        includeQuery.AddQuery(includeAnyQuery)
    }

    if len(criteria.IncludeAll)>0 {
        includeAllQuery := bleve.NewConjunctionQuery()
        for _, tag := range criteria.IncludeAll {
            includeTagQuery := bleve.NewMatchQuery(tag)
            includeTagQuery.SetField("tag")
            includeAllQuery.AddQuery(includeTagQuery)
        }
        includeQuery.AddQuery(includeAllQuery)
    }

    excludeQuery := bleve.NewDisjunctionQuery()
    searchQuery.AddMustNot(excludeQuery)

    if len(criteria.ExcludeAny)>0 {
        excludeAnyQuery := bleve.NewDisjunctionQuery()
        for _, tag := range criteria.ExcludeAny {
            excludeTagQuery := bleve.NewMatchQuery(tag)
            excludeTagQuery.SetField("tag")
            excludeAnyQuery.AddQuery(excludeTagQuery)
        }
        excludeQuery.AddQuery(excludeAnyQuery)
    }

    if len(criteria.ExcludeAll)>0 {
        excludeAllQuery := bleve.NewConjunctionQuery()
        for _, tag := range criteria.ExcludeAll {
            excludeTagQuery := bleve.NewMatchQuery(tag)
            excludeTagQuery.SetField("tag")
            excludeAllQuery.AddQuery(excludeTagQuery)
        }
        excludeQuery.AddQuery(excludeAllQuery)
    }

    searchRequest := bleve.NewSearchRequest(searchQuery)
    searchRequest.Fields = []string{"*"}
    searchRequest.Size = criteria.Limit
    searchResults, err := r.itemIndex.Search(searchRequest)

    if err != nil {
        return nil, err
    }

    result :=  &domain.SearchResult{}
    for _, hit := range searchResults.Hits {
        
        var tags []string
        if tagsValue, ok := hit.Fields["tag"]; ok {
            for _, tagValue := range tagsValue.([]interface{}) {
                tags = append(tags, tagValue.(string))
            }
        }

        var note string
        if noteValue, ok := hit.Fields["note"]; ok {
            note = noteValue.(string)
        }

        searchItem := domain.SearchItem { ID: hit.ID, Tags: tags, Note: note }

        switch itemType := hit.Fields["type"]; itemType.(string) {
            case itemTypeURL: result.URLItems = append(result.URLItems, domain.URLSearchItem{ SearchItem: searchItem, URL: hit.Fields["url"].(string) }) 
            default: return nil, fmt.Errorf("Item type %s is not reognized", itemType.(string))
        }
    }

    return result, nil
}

func generateID(parts ...string) string {
        h := sha1.New()
        h.Write([]byte(strings.Join(parts, "")))
        return hex.EncodeToString(h.Sum(nil))
}
