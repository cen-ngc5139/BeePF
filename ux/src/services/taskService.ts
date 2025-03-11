import api from './api';

export interface Task {
    id: number;
    name: string;
    description?: string;
    component_id: number;
    component_name: string;
    step: number;
    status: number;
    error?: string;
    creator?: string;
    created_at: string;
    updated_at: string;
    prog_status?: ProgStatus[];
}

export interface ProgStatus {
    id: number;
    task_id: number;
    component_id: number;
    component_name: string;
    program_id: number;
    program_name: string;
    status: number;
    error: string;
    created_at: string;
    updated_at: string;
}

class TaskService {
    /**
     * 创建并运行组件任务
     * @param componentId 组件ID
     * @returns 创建的任务信息
     */
    async runComponent(componentId: number): Promise<Task> {
        try {
            const response = await api.post(`/task/component/${componentId}`);
            const responseData = response as any;
            if (responseData.success && responseData.data) {
                return responseData.data;
            }
            throw new Error(responseData.errorMsg || '运行组件任务失败');
        } catch (error) {
            console.error('运行组件任务失败:', error);
            throw error;
        }
    }

    /**
     * 获取任务列表
     * @param params 查询参数
     * @returns 任务列表
     */
    async getTasks(params?: any): Promise<{ list: Task[], total: number }> {
        try {
            const response = await api.get('/task', { params });
            // 确保返回的数据格式正确
            const responseData = response as any;
            if (responseData.success && responseData.data) {
                return {
                    list: Array.isArray(responseData.data.list) ? responseData.data.list : [],
                    total: typeof responseData.data.total === 'number' ? responseData.data.total : 0
                };
            }
            return { list: [], total: 0 };
        } catch (error) {
            console.error('获取任务列表失败:', error);
            throw error;
        }
    }

    /**
     * 获取任务详情
     * @param taskId 任务ID
     * @returns 任务详情
     */
    async getTask(taskId: number): Promise<Task> {
        try {
            const response = await api.get(`/task/${taskId}`);
            // 确保返回的数据格式正确
            const responseData = response as any;
            if (responseData.success && responseData.data) {
                return responseData.data;
            }
            return {} as Task;
        } catch (error) {
            console.error('获取任务详情失败:', error);
            throw error;
        }
    }

    /**
     * 获取正在运行的任务
     * @returns 运行中的任务列表
     */
    async getRunningTasks(): Promise<{ list: Task[], total: number }> {
        try {
            const response = await api.get('/task/running');
            // 确保返回的数据格式正确
            const responseData = response as any;
            if (responseData.success && responseData.data) {
                return {
                    list: Array.isArray(responseData.data.list) ? responseData.data.list : [],
                    total: typeof responseData.data.total === 'number' ? responseData.data.total : 0
                };
            }
            return { list: [], total: 0 };
        } catch (error) {
            console.error('获取运行中任务失败:', error);
            throw error;
        }
    }

    /**
     * 停止任务
     * @param taskId 任务ID
     * @returns 操作结果
     */
    async stopTask(taskId: number): Promise<any> {
        try {
            const response = await api.post(`/task/${taskId}/stop`);
            const responseData = response as any;
            if (responseData.success) {
                return responseData.data;
            }
            throw new Error(responseData.errorMsg || '停止任务失败');
        } catch (error) {
            console.error('停止任务失败:', error);
            throw error;
        }
    }
}

export default new TaskService(); 