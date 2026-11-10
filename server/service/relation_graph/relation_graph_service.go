package relation_graph

import (
	"server/applog"
	"server/errors"
	deceasedmodel "server/model/deceased"
	gravemodel "server/model/grave"
	kdr "server/model/kinship_deceased_relation"
	krmodel "server/model/kinship_relation"
	deceasedrepo "server/repository/deceased"
	graverepo "server/repository/grave"
	memberrepo "server/repository/grave_member"
	kdrrepo "server/repository/kinship_deceased_relation"
	krrepo "server/repository/kinship_relation"
	"server/types"
	"server/utils"
	"strconv"
	"time"
)

// RelationGraphService 关系图服务
// 负责组装墓园的关系图数据，包括节点和边
// 节点包括：逝者节点、用户节点
// 边包括：逝者间关系、用户-逝者关系
type RelationGraphService struct {
	graveRepo            *graverepo.GraveRepository
	deceasedRepo         *deceasedrepo.DeceasedRepository
	memberRepo           *memberrepo.GraveMemberRepository
	deceasedRelationRepo *kdrrepo.DeceasedRelationRepository
	krRepo               *krrepo.KinshipRelationRepository
}

func NewRelationGraphService(
	graveRepo *graverepo.GraveRepository,
	deceasedRepo *deceasedrepo.DeceasedRepository,
	memberRepo *memberrepo.GraveMemberRepository,
	deceasedRelationRepo *kdrrepo.DeceasedRelationRepository,
	krRepo *krrepo.KinshipRelationRepository,
) *RelationGraphService {
	return &RelationGraphService{
		graveRepo:            graveRepo,
		deceasedRepo:         deceasedRepo,
		memberRepo:           memberRepo,
		deceasedRelationRepo: deceasedRelationRepo,
		krRepo:               krRepo,
	}
}

// GetGraph 获取关系图数据
// 关系图用于展示家族谱系，包括逝者之间的关系和用户与逝者的关系
// 参数:
//   - graveID: 墓地ID
//   - userID: 用户ID（用于权限验证）
//
// 返回:
//   - RelationGraphResp: 关系图数据（包含节点和边）
//   - error: 错误信息
func (s *RelationGraphService) GetGraph(graveID, userID uint) (*types.RelationGraphResp, error) {
	// 1. 验证墓地是否存在
	grave, err := s.graveRepo.GetByID(graveID)
	if err != nil {
		return nil, err
	}

	// 2. 纪念模式不展示关系图
	if grave.GraveMode == gravemodel.GraveModeMemorial {
		return &types.RelationGraphResp{
			GraveID:   graveID,
			GraveName: grave.Name,
			Nodes:     []types.GraphNode{},
			Edges:     []types.GraphEdge{},
		}, nil
	}

	if !s.memberRepo.IsMember(graveID, userID) {
		return nil, errors.ErrGraveAccessDenied
	}

	// 3. 获取墓园所有逝者
	deceasedList, _, err := s.deceasedRepo.ListByGrave(graveID, 1, 0)
	if err != nil {
		return nil, err
	}

	// 4. 获取逝者间关系（构建家族谱系的骨架）
	deceasedRelations, _, err := s.deceasedRelationRepo.ListByGrave(graveID, 1, 0)
	if err != nil {
		return nil, err
	}

	// 5. 获取用户-逝者关系（展示当前家族成员与逝者的关系）
	userDeceasedRelations, _, err := s.krRepo.ListByGrave(graveID, 1, 0)
	if err != nil {
		return nil, err
	}

	// 6. 组装节点和边
	nodes := s.buildNodes(deceasedList, deceasedRelations, userDeceasedRelations)
	edges := s.buildEdges(deceasedRelations, userDeceasedRelations)

	// 7. 返回完整的关系图数据
	return &types.RelationGraphResp{
		GraveID:   graveID,
		GraveName: grave.Name,
		Nodes:     nodes,
		Edges:     edges,
	}, nil
}

// buildNodes 构建关系图的节点
// 节点包括两类：逝者节点和用户节点
// 参数:
//   - deceasedList: 逝者列表
//   - userRelations: 用户-逝者关系列表
//
// 返回:
//   - []GraphNode: 节点列表
func (s *RelationGraphService) buildNodes(deceasedList []deceasedmodel.Deceased, deceasedRelations []kdr.DeceasedRelation, userRelations []krmodel.KinshipRelation) []types.GraphNode {
	nodes := make([]types.GraphNode, 0)

	parentMap := s.buildParentMap(deceasedRelations)
	s.detectAndBreakCycles(parentMap, deceasedList)
	generationMap := s.computeGenerations(len(deceasedList), parentMap)

	for _, d := range deceasedList {
		gen, ok := generationMap[d.ID]
		if !ok {
			if spouseGen := s.getSpouseGeneration(d.ID, deceasedRelations, generationMap); spouseGen >= 0 {
				gen = spouseGen
			} else {
				gen = s.calculateGeneration(d.BirthDate)
			}
		}
		nodes = append(nodes, types.GraphNode{
			ID:         s.formatNodeID("deceased", d.ID),
			Type:       "deceased",
			Name:       d.Name,
			Avatar:     d.PortraitURL,
			Gender:     d.Gender,
			Generation: gen,
			IsDead:     true,
			BirthDate:  d.BirthDate,
			DeathDate:  d.DeathDate,
		})
	}

	userGenMap := s.computeUserGenerations(userRelations, generationMap)
	userSet := make(map[uint]bool)
	for _, r := range userRelations {
		if r.ConfirmationStatus == krmodel.ConfirmationConfirmed && !userSet[r.UserID] {
			userSet[r.UserID] = true
			displayName := r.User.Nickname
			if displayName == "" {
				displayName = r.User.Username
			}
			nodes = append(nodes, types.GraphNode{
				ID:         s.formatNodeID("user", r.UserID),
				Type:       "user",
				Name:       displayName,
				Avatar:     r.User.Avatar,
				Gender:     0,
				Generation: userGenMap[r.UserID],
				IsDead:     false,
			})
		}
	}

	return nodes
}

func (s *RelationGraphService) detectAndBreakCycles(parentMap map[uint]uint, deceasedList []deceasedmodel.Deceased) {
	deceasedIDs := make(map[uint]bool)
	for _, d := range deceasedList {
		deceasedIDs[d.ID] = true
	}

	whitelist := map[[2]uint]bool{}
	degree := make(map[uint]int)
	children := make(map[uint][]uint)
	for child, parent := range parentMap {
		if deceasedIDs[child] && deceasedIDs[parent] {
			degree[parent]++
			children[child] = append(children[child], parent)
		}
	}

	queue := make([]uint, 0)
	for _, d := range deceasedList {
		if degree[d.ID] == 0 {
			queue = append(queue, d.ID)
		}
	}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		for _, parent := range children[node] {
			degree[parent]--
			whitelist[[2]uint{parent, node}] = true
			if degree[parent] == 0 {
				queue = append(queue, parent)
			}
		}
	}

	removed := 0
	for child, parent := range parentMap {
		if !deceasedIDs[child] || !deceasedIDs[parent] {
			continue
		}
		if !whitelist[[2]uint{parent, child}] {
			applog.Warnf("关系图检测到环形引用，已移除边: parent=%d child=%d", parent, child)
			delete(parentMap, child)
			removed++
		}
	}
	if removed > 0 {
		applog.Warnf("关系图共移除 %d 条环形引用边", removed)
	}
}

func (s *RelationGraphService) buildParentMap(relations []kdr.DeceasedRelation) map[uint]uint {
	parentMap := make(map[uint]uint)
	for _, r := range relations {
		if r.RelationType != kdr.RelationParent {
			continue
		}
		switch r.Direction {
		case kdr.DirAIsParentOfB:
			parentMap[r.DeceasedID2] = r.DeceasedID1
		case kdr.DirBIsParentOfA:
			parentMap[r.DeceasedID1] = r.DeceasedID2
		}
	}
	return parentMap
}

func (s *RelationGraphService) computeGenerations(nodeCount int, parentMap map[uint]uint) map[uint]int {
	generationMap := make(map[uint]int)
	if len(parentMap) == 0 {
		return generationMap
	}

	childrenMap := make(map[uint][]uint)
	for child, parent := range parentMap {
		childrenMap[parent] = append(childrenMap[parent], child)
	}

	roots := make([]uint, 0)
	for child := range parentMap {
		if _, hasParent := parentMap[parentMap[child]]; !hasParent {
			found := false
			for _, r := range roots {
				if r == parentMap[child] {
					found = true
					break
				}
			}
			if !found {
				roots = append(roots, parentMap[child])
			}
		}
	}
	if len(roots) == 0 {
		for parent := range childrenMap {
			if _, isChild := parentMap[parent]; !isChild {
				roots = append(roots, parent)
			}
		}
	}

	queue := make([]uint, 0, nodeCount)
	visited := make(map[uint]bool)
	for _, root := range roots {
		generationMap[root] = 0
		queue = append(queue, root)
		visited[root] = true
	}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		for _, child := range childrenMap[node] {
			if visited[child] {
				continue
			}
			generationMap[child] = generationMap[node] + 1
			queue = append(queue, child)
			visited[child] = true
		}
	}

	return generationMap
}

func (s *RelationGraphService) getSpouseGeneration(deceasedID uint, relations []kdr.DeceasedRelation, generationMap map[uint]int) int {
	for _, r := range relations {
		if r.RelationType != kdr.RelationSpouse {
			continue
		}
		if r.DeceasedID1 == deceasedID {
			if gen, ok := generationMap[r.DeceasedID2]; ok {
				return gen
			}
		}
		if r.DeceasedID2 == deceasedID {
			if gen, ok := generationMap[r.DeceasedID1]; ok {
				return gen
			}
		}
	}
	return -1
}

func (s *RelationGraphService) computeUserGenerations(userRelations []krmodel.KinshipRelation, generationMap map[uint]int) map[uint]int {
	result := make(map[uint]int)
	grouped := make(map[uint][]krmodel.KinshipRelation)
	for _, r := range userRelations {
		if r.ConfirmationStatus == krmodel.ConfirmationConfirmed {
			grouped[r.UserID] = append(grouped[r.UserID], r)
		}
	}
	for uid, rels := range grouped {
		bestGen := -1
		for _, r := range rels {
			decGen, ok := generationMap[r.DeceasedID]
			if !ok {
				continue
			}
			var offset int
			switch r.RelationType {
			case krmodel.KinshipRelationBrother, krmodel.KinshipRelationSister, krmodel.KinshipRelationSpouse:
				offset = 0
			case krmodel.KinshipRelationSon, krmodel.KinshipRelationDaughter:
				offset = 1
			case krmodel.KinshipRelationGrandchild:
				offset = 2
			case krmodel.KinshipRelationFather, krmodel.KinshipRelationMother, krmodel.KinshipRelationElder:
				offset = -1
			case krmodel.KinshipRelationGrandfather, krmodel.KinshipRelationGrandmother:
				offset = -2
			default:
				continue
			}
			userGen := decGen + offset
			if userGen > bestGen {
				bestGen = userGen
			}
		}
		if bestGen < 0 {
			bestGen = 0
		}
		result[uid] = bestGen
	}
	return result
}

// buildEdges 构建关系图的边
// 边包括两类：逝者间关系、用户-逝者关系
// 参数:
//   - deceasedRelations: 逝者间关系列表
//   - userRelations: 用户-逝者关系列表
//
// 返回:
//   - []GraphEdge: 边列表
func (s *RelationGraphService) buildEdges(deceasedRelations []kdr.DeceasedRelation, userRelations []krmodel.KinshipRelation) []types.GraphEdge {
	edges := make([]types.GraphEdge, 0)

	// 1. 添加逝者间关系的边
	// 这些关系构建了家族谱系的骨架
	for _, r := range deceasedRelations {
		relationType := s.mapRelationType(r.RelationType)
		fromID := s.formatNodeID("deceased", r.DeceasedID1)
		toID := s.formatNodeID("deceased", r.DeceasedID2)
		if r.Direction == kdr.DirBIsParentOfA {
			fromID, toID = toID, fromID
		}
		edges = append(edges, types.GraphEdge{
			From:         fromID,
			To:           toID,
			RelationType: relationType,
			Confirmed:    true,
			Direction:    r.Direction,
		})
	}

	// 2. 添加用户-逝者关系的边
	// 只添加已确认的关系
	for _, r := range userRelations {
		if r.ConfirmationStatus == krmodel.ConfirmationConfirmed {
			relationType := s.mapUserRelationType(r.RelationType)
			edges = append(edges, types.GraphEdge{
				From:               s.formatNodeID("user", r.UserID),
				To:                 s.formatNodeID("deceased", r.DeceasedID),
				RelationType:       relationType,
				Confirmed:          true,
				ConfirmationStatus: r.ConfirmationStatus,
				ConfirmedBy:        r.ConfirmedBy,
			})
		}
	}

	return edges
}

// formatNodeID 格式化节点ID
// 使用类型前缀避免ID冲突（如：deceased_1 和 user_1）
// 参数:
//   - nodeType: 节点类型（"deceased" 或 "user"）
//   - id: 实体ID
//
// 返回:
//   - string: 格式化后的节点ID
func (s *RelationGraphService) formatNodeID(nodeType string, id uint) string {
	return nodeType + "_" + strconv.Itoa(int(id))
}

// calculateGeneration 根据出生日期估算辈分（后备方案）
// 当逝者间没有 parent 关系时，用出生年份估算
func (s *RelationGraphService) calculateGeneration(birthDate string) int {
	// 解析出生日期
	birthTime, err := time.Parse(utils.DateFormat, birthDate)
	if err != nil {
		return 0 // 解析失败返回默认值
	}

	// 以2024年为基准，每30年为一代
	baseYear := 2024
	birthYear := birthTime.Year()
	generation := (birthYear - baseYear) / 30

	return generation
}

// mapRelationType 映射逝者关系类型
// 将数据库中的关系类型码转换为前端可读的类型字符串
// 参数:
//   - relationType: 关系类型码
//
// 返回:
//   - string: 关系类型字符串
func (s *RelationGraphService) mapRelationType(relationType int8) string {
	// 关系类型定义参考：server/model/deceased_relation/enum_deceased_relation.go
	types := map[int8]string{
		1: "spouse", // 夫妻
		2: "parent", // 亲子
	}
	return types[relationType]
}

// mapUserRelationType 映射用户-逝者关系类型
// 将数据库中的关系类型码转换为前端可读的类型字符串
// 参数:
//   - relationType: 关系类型码
//
// 返回:
//   - string: 关系类型字符串
func (s *RelationGraphService) mapUserRelationType(relationType int8) string {
	// 关系类型定义参考：server/model/kinship_relation/enum_kinship_relation.go
	types := map[int8]string{
		1:  "son",         // 儿子
		2:  "daughter",    // 女儿
		3:  "spouse",      // 配偶
		4:  "brother",     // 兄弟
		5:  "sister",      // 姐妹
		6:  "grandchild",  // 孙辈
		7:  "father",      // 父亲
		8:  "mother",      // 母亲
		9:  "grandfather", // 祖父
		10: "grandmother", // 祖母
		11: "elder",       // 长辈
		99: "follower",    // 追随者/其他
	}
	return types[relationType]
}
