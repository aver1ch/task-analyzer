package structs

type DBStatusChanges struct {
	IssueId    int
	AuthorId   int
	ChangeTime string //mayby time
	FromStatus string
	ToStatus   string
}

type DBAuthor struct {
	Id   int
	Name string
}

type DBProject struct {
	Id    int
	Title string
}

type DBIssue struct {
	Id          int
	ProjectId   int
	AuthorId    int
	AssigneeId  int
	Key         string
	Summary     string
	Description string
	Type        string
	Priority    string
	Status      string
	CreatedTime string //maybe time
	ClosedTime  string //maybe time
	UpdatedTime string //maybe time
	TimeSpent   int
}
