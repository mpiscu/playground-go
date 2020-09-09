package domain

type UpdateRequest struct {
    ID string
    AddTags []string
    DelTags []string
    Note string
}
