package model

const (
	RelationSpouse int8 = 1
	RelationParent int8 = 2
)

const (
	DirSymmetric    int8 = 0
	DirAIsParentOfB int8 = 1
	DirBIsParentOfA int8 = 2
)

func IsSymmetricRelation(rt int8) bool {
	return rt == RelationSpouse
}

func IsValidDirection(rt int8, dir int8) bool {
	if IsSymmetricRelation(rt) {
		return dir == DirSymmetric
	}
	return dir == DirAIsParentOfB || dir == DirBIsParentOfA
}
