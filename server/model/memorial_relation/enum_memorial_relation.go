package model

const (
	MemRelationSon         int8 = 1
	MemRelationDaughter    int8 = 2
	MemRelationSpouse      int8 = 3
	MemRelationBrother     int8 = 4
	MemRelationSister      int8 = 5
	MemRelationGrandchild  int8 = 6
	MemRelationFather      int8 = 7
	MemRelationMother      int8 = 8
	MemRelationGrandfather int8 = 9
	MemRelationGrandmother int8 = 10
	MemRelationElder       int8 = 11
	MemRelationFollower    int8 = 99
)

const (
	MemConfirmationPending   int8 = 0
	MemConfirmationConfirmed int8 = 1
	MemConfirmationRejected  int8 = 2
)

var autoConfirmedTypes = map[int8]bool{
	MemRelationFollower: true,
}

func IsAutoConfirmed(relationType int8) bool {
	return autoConfirmedTypes[relationType]
}
