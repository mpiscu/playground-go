package cli

import(
    "flag"
    "search/search-core/domain"
    "search/search-core/service"
)

const (
    delItemCmdName = "del"
)

type DelItemCommand struct {
    ids []string
    shortDescription string
    indexSvc *service.ItemIndexSvc
}

func NewDelItemCommand(indexSvc *service.ItemIndexSvc) Command {
    return &DelItemCommand{ indexSvc: indexSvc, shortDescription: delItemCmdName + " - Deletes an item which was already indexed", }
}

func (c *DelItemCommand) ShortDescription() string{
    return c.shortDescription
}

func (c *DelItemCommand) Parse(args []string) error {
    if len(args) <= 1 {
        return ErrCommandNoMatch 
    }
    if delItemCmdName != args[1] {
        return ErrCommandNoMatch
    }


    delItemCmd := flag.NewFlagSet(delItemCmdName, flag.ExitOnError)
    idPtr := delItemCmd.String("ids","","Comma separated ids of items to delete")

    delItemCmd.Parse(args[2:])

    c.ids =  TextSplitCSV(*idPtr)
    if err := ValidateNotEmpty(c.ids, "ids"); err != nil {
        return err
    }

    return nil

}

func (c *DelItemCommand) Execute() error {
    for _, id := range c.ids {
        err := c.indexSvc.Delete(id) 
        if err !=nil && err!=domain.ErrNotFound {
            return err
        }
        println("Item", id, "deleted")
    }
    return nil
}

