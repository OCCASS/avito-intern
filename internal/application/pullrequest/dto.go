package pullrequest

type CreatePullRequestDto struct {
	Id       string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorId string `json:"author_id"`
}
