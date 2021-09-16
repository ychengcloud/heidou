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

func TestForeignKeyForHasMany(t *testing.T) {

	field := &Field{
		Name:     "many",
		JoinType: JoinTypeHasMany,
	}
	field.genName()

	assField := &Field{
		Name: "many_id",
	}
	assField.genName()

	assTable := &Table{
		Fields: []*Field{
			assField,
		},
	}
	handleAssociationForeignKey(assTable, "ManyId")
	if !assField.IsForeignKey {
		t.Fail()
	}
}
