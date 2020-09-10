package cli

import(
    "flag"
    "errors"
    "search/search-core/domain"
    "search/search-core/service"
)

const (
    updateItemCmdName = "update"
)

type UpdateItemCommand struct {
    id string
    addTags []string
    delTags []string
    note string
    shortDescription string
    indexSvc *service.ItemIndexSvc
    searchSvc *service.ItemSearchSvc
}

func NewUpdateItemCommand(indexSvc *service.ItemIndexSvc, searchSvc *service.ItemSearchSvc) Command {
    return &UpdateItemCommand{ indexSvc: indexSvc, shortDescription: updateItemCmdName + " - Updates an item which was already indexed", }
}

func (c *UpdateItemCommand) ShortDescription() string{
    return c.shortDescription
}

func (c *UpdateItemCommand) Parse(args []string) error {
    if len(args) <= 1 {
        return ErrCommandNoMatch 
    }
    if updateItemCmdName != args[1] {
        return ErrCommandNoMatch
    }


    updateItemCmd := flag.NewFlagSet(updateItemCmdName, flag.ExitOnError)
    idPtr := updateItemCmd.String("id","","id of the item which will be updated")
    addTagsPtr := updateItemCmd.String("add-tags", "", "tags to be added")
    delTagsPtr := updateItemCmd.String("del-tags", "", "tags to be removed")
    notePtr := updateItemCmd.String("note", "", "item notes")

    updateItemCmd.Parse(args[2:])

    c.id =  *idPtr
    if err := ValidateMandatory(c.id, "id", ""); err != nil {
        return err
    }

    c.addTags = TextSplitCSV(*addTagsPtr)
    c.delTags = TextSplitCSV(*delTagsPtr)

    c.note = *notePtr

    if len(c.addTags) == 0 && len(c.delTags) == 0 && c.note == "" {
        return errors.New("no change specified")
    }

    return nil

}

func (c *UpdateItemCommand) Execute() error {

    request := domain.UpdateRequest {
        ID: c.id,
        AddTags: c.addTags,
        DelTags: c.delTags,
        Note: c.note,
    }

    err := c.indexSvc.Update(request)
    if err != nil {
        return err
    }
    
    println("Updated", c.id)

    return nil
}

