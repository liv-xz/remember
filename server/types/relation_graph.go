package types

type RelationGraphResp struct {
	GraveID   uint        `json:"graveID"`
	GraveName string      `json:"graveName"`
	Nodes     []GraphNode `json:"nodes"`
	Edges     []GraphEdge `json:"edges"`
}

type GraphNode struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Name       string `json:"name"`
	Avatar     string `json:"avatar"`
	Gender     int8   `json:"gender"`
	Generation int    `json:"generation"`
	IsDead     bool   `json:"isDead"`
	BirthDate  string `json:"birthDate,omitempty"`
	DeathDate  string `json:"deathDate,omitempty"`
}

type GraphEdge struct {
	From               string `json:"from"`
	To                 string `json:"to"`
	RelationType       string `json:"relationType"`
	Confirmed          bool   `json:"confirmed"`
	Direction          int8   `json:"direction,omitempty"`
	ConfirmationStatus int8   `json:"confirmationStatus,omitempty"`
	ConfirmedBy        *uint  `json:"confirmedBy,omitempty"`
}
