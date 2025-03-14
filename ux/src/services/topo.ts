import axios from 'axios';

// 定义拓扑数据的接口
export interface ProgNode {
    GUID: string;
    ID: number;
    Name: string;
}

export interface MapNode {
    GUID: string;
    ID: number;
    Name: string;
}

export interface TopologyEdge {
    ProgGUID: string;
    MapGUID: string;
    ProgID: number;
    MapID: number;
}

export interface Topology {
    ProgNodes: ProgNode[];
    MapNodes: MapNode[];
    Edges: TopologyEdge[];
}

// 获取拓扑数据
export const getTopology = async (): Promise<Topology> => {
    const response = await axios.get('/api/v1/observability/topo');
    return response.data;
}; 