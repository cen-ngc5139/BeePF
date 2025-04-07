import api from './api';

// 定义节点指标数据接口
export interface ProgramStats {
    cpu_time_percent: number;
    events_per_second: number;
    avg_run_time_ns: number;
    total_avg_run_time_ns: number;
    period_ns: number;
    last_update: string;
}

export interface ProgramMetric {
    stats: ProgramStats;
    id: number;
    type: string;
    name: string;
}

export interface NodeMetricsResponse {
    success: boolean;
    errorCode: number;
    errorMsg: string;
    data: {
        metrics: Record<string, ProgramMetric>;
    };
}

// 获取节点指标数据
export const getNodeMetrics = async (): Promise<NodeMetricsResponse> => {
    return api.get('/observability/node/metrics');
}; 