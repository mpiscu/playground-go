# Description

A command line utility which can index various piece of data. The data will be searchable by tags. Each information can have specific properites. E.g. an url with a description.

# Requirements

1. AddURLItem : As an user I want to index an url with specified tags. The url must be searchable by those tags. The url can have user notes.
2. AddTag : As an user I want to be able to add a tag to an existing item.
3. DelTag: As an user I want to be able to delete a tag from an existing item.
4. DelURLItem: As an user I want to be able to delete an item. 
5. UpdateItemNote: As an user I want to be able to update an existing item note.
6. SearchItem: As an user I want to be able to search an item by its tags. Exact and partial matching will be supported.

# Minimal design

- cli layer - contains all commands and cli application
- service layer - contains logic for commands
- repository layer - access to bleve index

# Technical stack
- [golang](https://golang.org/) - programming language
- [bleve](https://blevesearch.com/)   - indexing engine

# Terms

*Item* - An item is a piece of information which must be searchable. The item will be searchable by using tags. The item can have some notes which describe the item. E.g. UrlItem has an url.

*Tag* - Is an identifier which belongs to an item. Tags are used as criteria in search operations.

*Note* - A description of the item which should provide more context to the user

# Dev

Prerequisites : golang installed, modules enabled.

Setup : `source setup-dev.sh` prepares environment for development. Utility scripts available in `scripts` folder are introduced in path.


# Install/Uninstall

Install   : `sudo cp .go/bin/search-cli /usr/bin/`  
Uninstall : `sudo rm .go/bin/search-cli && sudo rm -rf ~/.search`

# Usage

See help

# Improvements

- Add search filter by item type
- List tags
- Match tags partially
- Add support for logging
- Make index location configurable
- Add support for pagination
- If scale is needed add partitioning by using multiple bleve indexes and aggregate results on client side
- Split repository for each type of item
- Make it easier to plug in other item types
- Create a package search-crawler which have different crawlers e.g. music, document, wiki which automatically indexes their targets
