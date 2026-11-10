package model

const (
	KinshipRelationSon         int8 = 1
	KinshipRelationDaughter    int8 = 2
	KinshipRelationSpouse      int8 = 3
	KinshipRelationBrother     int8 = 4
	KinshipRelationSister      int8 = 5
	KinshipRelationGrandchild  int8 = 6
	KinshipRelationFather      int8 = 7
	KinshipRelationMother      int8 = 8
	KinshipRelationGrandfather int8 = 9
	KinshipRelationGrandmother int8 = 10
	KinshipRelationElder       int8 = 11
	KinshipRelationFollower    int8 = 99
)

const (
	ConfirmationPending   int8 = 0
	ConfirmationConfirmed int8 = 1
	ConfirmationRejected  int8 = 2
)
