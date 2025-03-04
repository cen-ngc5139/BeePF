import api from './api';

export interface Cluster {
    id?: number;
    name: string;
    cnname?: string;      // 中文名称
    description?: string; // 兼容旧字段
    desc?: string;        // 新的描述字段
    region?: string;      // 兼容旧字段
    master?: string;      // 主节点地址
    status: number | 'active' | 'inactive'; // 状态改为数字类型
    environment?: string; // 环境
    creator?: string;     // 创建者
    kubeconfig?: string;  // kubeconfig内容
    deleted?: boolean;    // 是否删除
    createdat?: string;   // 创建时间
    updateat?: string;    // 更新时间
    locationID?: number;  // 位置ID
}

// 响应数据结构
export interface ApiResponse<T> {
    success: boolean;
    errorCode: number;
    errorMsg: string;
    data: T;
}

export interface ClusterData {
    cluster?: Cluster;  // 可能不存在
    detail?: Cluster;   // 可能使用detail字段
}

export interface ClustersData {
    list: Cluster[];  // 修改为list字段
    total: number;
}

export interface ClusterListParams {
    pageSize?: number;
    pageNum?: number;
    environment?: string;
    keyword?: string;
}

export interface ClusterListResponse {
    total: number;
    clusters: Cluster[];
}

/**
 * 获取集群列表
 */
export const getClusterList = async (params: ClusterListParams): Promise<ClusterListResponse> => {
    const { pageSize, pageNum, environment, keyword } = params;
    const queryParams: Record<string, any> = {};

    if (pageSize) queryParams.pageSize = pageSize;
    if (pageNum) queryParams.pageNum = pageNum;
    if (environment) queryParams.environment = environment;
    if (keyword) queryParams.keyword = keyword;

    const response = await api.get('/cluster', { params: queryParams }) as ApiResponse<ClustersData>;
    // 调试日志
    console.log('集群列表接口返回数据:', response);

    // 处理响应数据
    if (response.success && response.data) {
        return {
            total: response.data.total || 0,
            clusters: response.data.list || []  // 修改为从list字段获取数据
        };
    }
    return { total: 0, clusters: [] };
};

/**
 * 获取单个集群详情
 */
export const getCluster = async (clusterId: number): Promise<Cluster> => {
    const response = await api.get(`/cluster/${clusterId}`) as ApiResponse<ClusterData>;
    // 调试日志
    console.log('集群详情接口返回数据:', response);

    // 处理响应数据
    if (response.success && response.data) {
        // 尝试从不同可能的字段获取数据
        const clusterData = response.data.cluster || response.data.detail;
        if (clusterData) {
            return clusterData;
        }
    }
    throw new Error('获取集群详情失败');
};

/**
 * 创建集群
 */
export const createCluster = async (clusterData: Cluster): Promise<any> => {
    try {
        // 检查必要字段
        if (!clusterData.name || !clusterData.cnname || !clusterData.master) {
            console.error('创建集群缺少必要字段:', clusterData);
            throw new Error('缺少必要字段：名称、中文名称或主节点地址');
        }

        // 调试日志
        console.log('创建集群请求数据:', JSON.stringify(clusterData, null, 2));

        const response = await api.post('/cluster', clusterData) as ApiResponse<any>;
        console.log('创建集群响应数据:', response);

        if (response.success) {
            return response.data;
        } else {
            throw new Error(response.errorMsg || '创建集群失败');
        }
    } catch (error) {
        console.error('创建集群请求错误:', error);
        throw error;
    }
};

/**
 * 更新集群
 */
export const updateCluster = async (clusterId: number, clusterData: Cluster): Promise<any> => {
    try {
        // 调试日志
        console.log('更新集群请求数据:', { clusterId, clusterData });

        const response = await api.put(`/cluster/${clusterId}`, clusterData) as ApiResponse<any>;
        console.log('更新集群响应数据:', response);

        if (response.success) {
            return response.data;
        } else {
            throw new Error(response.errorMsg || '更新集群失败');
        }
    } catch (error) {
        console.error('更新集群请求错误:', error);
        throw error;
    }
};

/**
 * 删除集群
 */
export const deleteCluster = async (clusterId: number): Promise<any> => {
    return api.delete(`/cluster/${clusterId}`);
};

/**
 * 根据参数获取集群列表（不分页）
 */
export const getClustersByParams = async (params: Record<string, any> = {}): Promise<Cluster[]> => {
    const response = await api.get('/clusterList', { params }) as ApiResponse<ClustersData>;
    // 处理响应数据
    if (response.success && response.data) {
        return response.data.list || [];  // 修改为从list字段获取数据
    }
    return [];
};

export default {
    getClusterList,
    getCluster,
    createCluster,
    updateCluster,
    deleteCluster,
    getClustersByParams,
}; 