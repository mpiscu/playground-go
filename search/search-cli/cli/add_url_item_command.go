package cli

import(
    "flag"
    "strings"
    "search/search-core/domain"
    "search/search-core/service"
)

const (
    addUrlItemCmdName = "add-url"
)

type AddUrlItemCommand struct {
    url string
    tags []string
    note string
    shortDescription string
    indexSvc *service.ItemIndexSvc
}

func NewAddUrlItemCommand(indexSvc *service.ItemIndexSvc) Command {
    return &AddUrlItemCommand{ indexSvc: indexSvc, shortDescription: addUrlItemCmdName + " - Indexes a new url item", }
}

func (c *AddUrlItemCommand) ShortDescription() string{
    return c.shortDescription
}

func (c *AddUrlItemCommand) Parse(args []string) error {
    if len(args) <= 1 {
        return ErrCommandNoMatch 
    }
    if addUrlItemCmdName != args[1] {
        return ErrCommandNoMatch
    }


    addItemCmd := flag.NewFlagSet(addUrlItemCmdName, flag.ExitOnError)
    urlPtr := addItemCmd.String("url","","URL which will be indexed")
    tagsPtr := addItemCmd.String("tags", "", "comma separated seach tags")
    notePtr := addItemCmd.String("note", "", "item notes")

    addItemCmd.Parse(args[2:])

    c.url =  *urlPtr
    if err := ValidateMandatory(c.url, "url", ""); err != nil {
        return err
    }
    if err := ValidateURL(c.url, "url"); err!=nil {
        return err
    }

    c.tags = strings.Split(*tagsPtr, ",")

    c.note = *notePtr

    return nil

}

func (c *AddUrlItemCommand) Execute() error {
    id, err := c.indexSvc.IndexURLItem(&domain.URLSearchItem{
        SearchItem: domain.SearchItem {
            Tags: c.tags,
            Note: c.note,
        },
        URL: c.url,
    })

    if err != nil {
        return err
    }

    println("ID:", id)

    return nil
}

