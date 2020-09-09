package service

import (
    "search/search-core/domain"
    "search/search-core/gateway/repository"
)

type ItemIndexSvc struct {
    itemRepository *repository.ItemRepository
}

func NewItemIndexSvc(itemRepository *repository.ItemRepository) *ItemIndexSvc {
    return &ItemIndexSvc{ itemRepository: itemRepository}
}

func (svc *ItemIndexSvc) IndexURLItem(item *domain.URLSearchItem) (string, error) {
    return svc.itemRepository.IndexURLItem(item)
}

func (svc *ItemIndexSvc) Update(r domain.UpdateRequest) error {
    if r.ID == "" {
        return domain.ErrNotFound
    }

    searchResult, err := svc.itemRepository.GetByID(r.ID)
    if err != nil {
        return err
    }


    //update url items
    for _, item := range searchResult.URLItems {
        changed := false
        for _, tag := range r.DelTags {
            p := textSlicePos(item.Tags, tag) 
            if p != -1 {
                item.Tags = textSliceRemovePos(item.Tags, p)
                changed = true
            }
        }
        for _, tag := range r.AddTags {
            p := textSlicePos(item.Tags, tag) 
            if p == -1 {
                item.Tags = append(item.Tags, tag)
                changed = true
            }
        }
        if r.Note != "" {
            item.Note = r.Note
            changed = true
        }
        if (changed) {
            _,err := svc.itemRepository.IndexURLItem(&item)
            if err != nil {
                return err
            }
        }
    }

    return nil

}

func (svc *ItemIndexSvc) Delete(id string) error {
    return svc.itemRepository.Delete(id)
}

func textSliceRemovePos(slice []string, pos int) []string {
    if len(slice) == 0 || pos<0 || pos>=len(slice) {
        return slice
    }
    if pos == len(slice)-1 {
        return slice[:pos]
    }
    if pos == 0 {
        return slice[1:]
    }
    return append(slice[:pos], slice[pos+1:]...)
}

func textSlicePos(slice []string, item string) int {
    for i, el := range slice {
        if el == item {
            return i
        }
    }
    return -1
}
