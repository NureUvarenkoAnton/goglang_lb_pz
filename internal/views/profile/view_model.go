package profile

type ProfildPageData struct {
	UserData           UserData
	Walks              []Walk
	PendingWalks       []Walk
	Pets               []Pet
	TheMostWalkeblePet Pet
}

type UserData struct {
	Email              string
	Name               string
	TheMostWalkeblePet IDContainer
}

type Walk struct {
	ID         string
	WalkerName string
	OwnerName  string
	PetName    string
}

type Pet struct {
	ID             string
	Name           string
	Age            string
	AdditionalInfo string
}

type IDContainer struct {
	ID   string
	Name string
}

type WalkFormData struct {
	Walkers []IDContainer
	Pets    []IDContainer
}
