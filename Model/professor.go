package professor

type Professor struct {
	ID         string        `json:"id,omitempty"`
	Name       string        `json:"nome,omitempty"`
	Discipline *[]Discipline `json:"disciplina,omitempty"`
}
