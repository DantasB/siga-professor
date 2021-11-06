package discipline

type Discipline struct {
	Name     string      `json:"nome,omitempty"`
	Code     string      `json:"codigo,omitempty"`
	Class    string      `json:"turma,omitempty"`
	Datetime *[]Datetime `json:"dias,omitempty"`
}
