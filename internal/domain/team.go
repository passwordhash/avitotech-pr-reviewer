package domain

type Team struct {
	ID string
	Name string
	Members []User
}

type TeamNormalized struct {
	ID string
	Name string
	MembersID []string
}
