package meta

type ProgNode struct {
	GUID string
	ID   uint32
	Name string
}

type ProgNodes []*ProgNode

type MapNode struct {
	GUID string
	ID   uint32
	Name string
}

type MapNodes []*MapNode

type TopologyEdge struct {
	ProgGUID string
	MapGUID  string
	ProgID   uint32
	MapID    uint32
}

type TopologyMap []*TopologyEdge

type Topology struct {
	ProgNodes ProgNodes
	MapNodes  MapNodes
	Edges     TopologyMap
}

func NewTopology() Topology {
	return Topology{
		ProgNodes: make(ProgNodes, 0),
		MapNodes:  make(MapNodes, 0),
		Edges:     make(TopologyMap, 0),
	}
}

func NewTopologyMap() TopologyMap {
	return make(TopologyMap, 0)
}

func (t *TopologyMap) AddEdge(edge TopologyEdge) {
	*t = append(*t, &edge)
}

func (t *ProgNodes) AddProgNode(node ProgNode) {
	*t = append(*t, &node)
}

func (t *MapNodes) AddMapNode(node MapNode) {
	*t = append(*t, &node)
}
