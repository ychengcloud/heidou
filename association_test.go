package heidou

import (
	"testing"
)

func TestForeignKey(t *testing.T) {
	field := &Field{
		Name: "parent_id",
	}
	field.genName()

	table := &Table{
		Fields: []*Field{
			field,
		},
	}

	handleAssociationForeignKey(table, "ParentId")
	if !field.IsForeignKey {
		t.Fail()
	}

}
