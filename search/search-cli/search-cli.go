package main

import (
    "fmt"
    "os"
    "search/search-core/gateway/repository"
    "search/search-core/service"
    "search/search-cli/cli"
)

const(
    exitCodeInternal=3
    exitCodeInput=2
    exitCodeProcessing=1
)

func main() {
    /*
    add-url url -tags "a,b,c" -note "Ana are mere"
    del -id ID
    search
    update-tag -id ID -add "a,b,c" -del "a,b,c"
    update-note -id ID -note "note"
    */
    //fmt.Println("Search Cli")
    //fmt.Println()
    indexRepo := repository.NewItemRepository()
    if err :=indexRepo.Open(); err!=nil {
        fmt.Println("Storage error: ", err.Error())
        os.Exit(exitCodeInternal)
    }
    defer indexRepo.Close()

    indexSvc := service.NewItemIndexSvc(indexRepo)
    searchSvc := service.NewItemSearchSvc(indexRepo)

    cmds := []cli.Command {
        cli.NewAddUrlItemCommand(indexSvc),
        cli.NewUpdateItemCommand(indexSvc,searchSvc),
        cli.NewDelItemCommand(indexSvc),
        cli.NewSearchItemCommand(searchSvc),
    }

    for i:= range cmds {
        
        err:=cmds[i].Parse(os.Args); 

        if err == cli.ErrCommandNoMatch {
            continue
        }

        if err != nil {
            fmt.Println("Input error:", err.Error())
            os.Exit(exitCodeInput)
        } 

        err = cmds[i].Execute(); 
        if err != nil {
            fmt.Println("Processing error: ", err.Error())
            os.Exit(exitCodeProcessing)
        } else {
            return
        }
        
    }

    fmt.Println("No command specified, available commands are: ")
    for i:= range cmds {
        fmt.Println("\t-", cmds[i].ShortDescription())
    }

    fmt.Println();
}
