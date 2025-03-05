import api from './api';
import { ApiResponse } from './clusterService';

export interface Component {
    id: number;
    name: string;
    cluster_id: number;
    programs: Program[] | null;
    maps: Map[] | null;
}

export interface Program {
    id: number;
    name: string;
    description: string;
    spec: any;
    properties: any;
}

export interface Map {
    id: number;
    name: string;
    description: string;
    spec: any;
    properties: any;
}

export interface ComponentsData {
    list: Component[];
    total: number;
}

export interface ComponentData {
    component: Component;
}

export interface ComponentListParams {
    pageSize?: number;
    pageNum?: number;
    keyword?: string;
    cluster_id?: number;
}

export interface ComponentListResponse {
    total: number;
    components: Component[];
}

/**
 * 获取组件列表
 */
export const getComponentList = async (params: ComponentListParams): Promise<ComponentListResponse> => {
    const { pageSize, pageNum, keyword, cluster_id } = params;
    const queryParams: Record<string, any> = {};

    if (pageSize) queryParams.pageSize = pageSize;
    if (pageNum) queryParams.pageNum = pageNum;
    if (keyword) queryParams.keyword = keyword;
    if (cluster_id) queryParams.cluster_id = cluster_id;

    const response = await api.get('/component', { params: queryParams }) as ApiResponse<ComponentsData>;
    // 调试日志
    console.log('组件列表接口返回数据:', response);

    // 处理响应数据
    if (response.success && response.data) {
        return {
            total: response.data.total || 0,
            components: response.data.list || []
        };
    }
    return { total: 0, components: [] };
};

/**
 * 获取单个组件详情
 */
export const getComponent = async (componentId: number): Promise<Component> => {
    const response = await api.get(`/component/${componentId}`) as ApiResponse<ComponentData>;
    // 调试日志
    console.log('组件详情接口返回数据:', response);

    // 处理响应数据
    if (response.success && response.data && response.data.component) {
        return response.data.component;
    }
    throw new Error('获取组件详情失败');
};

/**
 * 创建组件
 */
export const createComponent = async (componentData: Component): Promise<any> => {
    try {
        // 检查必要字段
        if (!componentData.name || !componentData.cluster_id) {
            console.error('创建组件缺少必要字段:', componentData);
            throw new Error('缺少必要字段：名称或集群ID');
        }

        // 调试日志
        console.log('创建组件请求数据:', JSON.stringify(componentData, null, 2));

        const response = await api.post('/component', componentData) as ApiResponse<any>;
        console.log('创建组件响应数据:', response);

        if (response.success) {
            return response.data;
        } else {
            throw new Error(response.errorMsg || '创建组件失败');
        }
    } catch (error) {
        console.error('创建组件请求错误:', error);
        throw error;
    }
};

export default {
    getComponentList,
    getComponent,
    createComponent
}; 