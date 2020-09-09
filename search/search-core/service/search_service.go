package service

import (
    "search/search-core/domain"
    "search/search-core/gateway/repository"
)

type ItemSearchSvc struct {
    itemRepository *repository.ItemRepository
}

func NewItemSearchSvc(itemRepository *repository.ItemRepository) *ItemSearchSvc {
    return &ItemSearchSvc{ itemRepository: itemRepository}
}

func (svc *ItemSearchSvc) Search(criteria *domain.SearchCriteria) (*domain.SearchResult, error) {
    if criteria.ID != "" {
        return svc.itemRepository.GetByID(criteria.ID)
    }

    return svc.itemRepository.Search(criteria)
}
