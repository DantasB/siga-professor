package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Discipline struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"nome" bson:"nome,omitempty"`
	Code        string             `json:"codigo" bson:"codigo,omitempty"`
	Class       string             `json:"turma" bson:"turma,omitempty"`
	Datetime    Datetime           `json:"dias"  bson:"dias,omitempty"`
	ProfessorID primitive.ObjectID `json:"professor_id" bson:"professor_id,omitempty"`
}
