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

// 定义节点资源数据接口
export interface ProgramInfo {
    ID: number;
    Type: number;
    Tag: string;
    Name: string;
    CreatedByUID: number;
    HaveCreatedByUID: boolean;
    BTF: number;
    LoadTime: string;
    Maps: number[];
}

// 定义Map信息接口
export interface MapInfo {
    Type: number;
    KeySize: number;
    ValueSize: number;
    MaxEntries: number;
    Flags: number;
    Name: string;
    ID: number;
    BTF: number;
    MapExtra: number;
    Memlock: number;
    Frozen: boolean;
}

// 定义程序详情接口
export interface ProgramDetail extends ProgramInfo {
    MapsDetail: MapInfo[] | null;
}

// 获取拓扑数据
export const getTopology = async (): Promise<Topology> => {
    const response = await axios.get('/api/v1/observability/topo');
    return response.data;
};

// 获取节点资源数据
export const getPrograms = async (): Promise<ProgramInfo[]> => {
    const response = await axios.get('/api/v1/observability/topo/prog');
    return response.data;
};

// 获取程序详情数据
export const getProgramDetail = async (progId: number): Promise<ProgramDetail> => {
    const response = await axios.get(`/api/v1/observability/topo/prog/${progId}`);
    return response.data;
};

// 获取程序指令信息
export const getProgramInstructions = async (progId: number, type: 'xlated' | 'jited' = 'xlated'): Promise<string> => {
    const response = await axios.get(`/api/v1/observability/topo/prog/${progId}/dump`, {
        params: { type },
        responseType: 'text'
    });
    return response.data;
};

// 获取程序源代码信息
export const getProgramSourceCode = async (progId: number): Promise<string> => {
    return getProgramInstructions(progId, 'jited');
}; 