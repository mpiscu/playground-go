package cli

import(
    "flag"
    "errors"
    "strings"
    "search/search-core/domain"
    "search/search-core/service"
)

const (
    searchItemCmdName = "search"
)

type SearchItemCommand struct {
    searchService *service.ItemSearchSvc
    shortDescription string
    searchCriteria domain.SearchCriteria
}

func NewSearchItemCommand(searchService *service.ItemSearchSvc) Command {
    return &SearchItemCommand{ searchService: searchService, shortDescription: searchItemCmdName + " - Searches for matching items", }
}

func (c *SearchItemCommand) ShortDescription() string{
    return c.shortDescription
}

func (c *SearchItemCommand) Parse(args []string) error {
    if len(args) <= 1 {
        return ErrCommandNoMatch 
    }
    if searchItemCmdName != args[1] {
        return ErrCommandNoMatch
    }


    searchItemCmd := flag.NewFlagSet(searchItemCmdName, flag.ExitOnError)
    searchIDPtr := searchItemCmd.String("id", "", "search item by a specific id")
    searchTypePtr := searchItemCmd.String("type","", "comma separated type, one of: url")
    searchIncludePtr := searchItemCmd.String("include", "", "include item if has any of the comma separated tags")
    searchIncludeAllPtr := searchItemCmd.String("include-all", "", "include item if has all comma separated tags")
    searchExcludePtr := searchItemCmd.String("exclude", "", "exclude item if has any of the comma separated tags")
    searchExcludeAllPtr := searchItemCmd.String("exclude-all", "", "exclude item if has all comma separated tags")
    sizePtr := searchItemCmd.Int("size", 10, "results limit for each item type")
    searchItemCmd.Parse(args[2:])

    if *searchIDPtr=="" && *searchTypePtr=="" && *searchIncludePtr=="" && *searchIncludeAllPtr=="" && *searchExcludePtr=="" && *searchExcludeAllPtr =="" {
        return errors.New("at least one filtering criteria must be specified")
    }

    c.searchCriteria = domain.SearchCriteria{
        ID: *searchIDPtr,
        Limit: *sizePtr,
        Types: TextSplitCSV(*searchTypePtr),
        IncludeAny: TextSplitCSV(*searchIncludePtr),
        IncludeAll: TextSplitCSV(*searchIncludeAllPtr),
        ExcludeAny: TextSplitCSV(*searchExcludePtr),
        ExcludeAll: TextSplitCSV(*searchExcludeAllPtr),
    }

    return nil

}

func (c *SearchItemCommand) Execute() error {

    result, err := c.searchService.Search(&c.searchCriteria)

    if err != nil {
        return err
    }

    println( "URL item match : ")
    println()
    if len(result.URLItems) == 0 {
        println("\tNone")
        println()
    } else {
        for i := range result.URLItems {
            println("\tID  :", result.URLItems[i].ID) 
            println("\tTags:", strings.Join(result.URLItems[i].Tags,","))
            println("\tURL :", result.URLItems[i].URL)
            println("\tNote:", result.URLItems[i].Note)
            println()
        }
    }


    return nil
}

