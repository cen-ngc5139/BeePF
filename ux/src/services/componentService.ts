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
    spec: ProgramSpec;
    properties: ProgramProperties;
}

export interface ProgramSpec {
    Name: string;
    Type: number;
    AttachType: number;
    AttachTo: string;
    SectionName: string;
    Flags: number;
    License: string;
    KernelVersion: number;
}

export interface ProgramProperties {
    CGroupPath: string;
    PinPath: string;
    LinkPinPath: string;
    Tc: any;
}

export interface Map {
    id: number;
    name: string;
    description: string;
    spec: MapSpec;
    properties: MapProperties;
}

export interface MapSpec {
    Name: string;
    Type: number;
    KeySize: number;
    ValueSize: number;
    MaxEntries: number;
    Flags: number;
    Pinning: number;
}

export interface MapProperties {
    PinPath: string;
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

/**
 * 删除组件
 */
export const deleteComponent = async (componentId: number): Promise<boolean> => {
    try {
        // 调试日志
        console.log(`删除组件请求，ID: ${componentId}`);
        console.log(`完整API路径: /component/${componentId}`);

        // 确保componentId是数字
        if (typeof componentId !== 'number' || isNaN(componentId)) {
            throw new Error(`无效的组件ID: ${componentId}`);
        }

        // 确保使用正确的API路径
        const apiPath = `/component/${componentId}`;
        console.log(`发送DELETE请求到: ${apiPath}`);

        // 使用更详细的错误处理
        try {
            const response = await api.delete(apiPath) as ApiResponse<any>;
            console.log('删除组件响应数据:', response);

            if (response.success) {
                return true;
            } else {
                throw new Error(response.errorMsg || '删除组件失败');
            }
        } catch (axiosError: any) {
            console.error('Axios错误:', axiosError);
            if (axiosError.response) {
                console.error('响应状态:', axiosError.response.status);
                console.error('响应数据:', axiosError.response.data);
            }
            throw axiosError;
        }
    } catch (error) {
        console.error('删除组件请求错误:', error);
        throw error;
    }
};

export default {
    getComponentList,
    getComponent,
    createComponent,
    deleteComponent
}; 