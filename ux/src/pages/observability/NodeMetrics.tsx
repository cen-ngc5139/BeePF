import React, { useEffect, useState } from 'react';
import { Card, Table, Typography, Spin, Empty, Tag, Button, Tooltip, Row, Col } from 'antd';
import { ReloadOutlined } from '@ant-design/icons';
import { getNodeMetrics, ProgramMetric } from '../../services/nodeMetrics';
import './NodeMetrics.css';
import type { ColumnsType } from 'antd/es/table';
import { useNavigate } from 'react-router-dom';

const { Title } = Typography;

interface TableProgramMetric extends ProgramMetric {
    key: string;
}

// eBPF 程序类型映射
const programTypes: Record<string, string> = {
    'SchedCLS': '流量控制',
    'CGroupSockAddr': '套接字地址控制',
    'Tracing': '系统追踪',
    'XDP': '快速数据路径',
    'KPROBE': '内核探针',
    'UPROBE': '用户空间探针',
    'PERF_EVENT': '性能事件',
    'SchedACT': '流量动作',
    'SocketFilter': '套接字过滤',
    'CGROUP_DEVICE': '设备控制',
    'CGROUP_SKB': '套接字缓冲控制',
    'CGROUP_SOCK': '套接字控制',
    'LWT_IN': '轻量级隧道入口',
    'LWT_OUT': '轻量级隧道出口',
    'LWT_XMIT': '轻量级隧道发送',
    'SOCK_OPS': '套接字操作',
    'SK_SKB': 'SK 缓冲区',
    'SK_MSG': 'SK 消息',
    'LIRC_MODE2': '红外遥控',
    'SK_REUSEPORT': '端口复用',
    'FLOW_DISSECTOR': '流解析器',
    'CGROUP_SYSCTL': 'Sysctl 控制',
    'RAW_TRACEPOINT': '原始跟踪点',
    'CGROUP_SOCKOPT': '套接字选项控制',
    'TRACING': '跟踪',
    'STRUCT_OPS': '结构操作',
    'EXT': '扩展',
    'LSM': 'Linux 安全模块',
    'SK_LOOKUP': 'SK 查询',
    'SYSCALL': '系统调用',
};

const NodeMetricsPage: React.FC = () => {
    const navigate = useNavigate();
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);
    const [metricsData, setMetricsData] = useState<TableProgramMetric[]>([]);

    const fetchMetrics = async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await getNodeMetrics();

            if (response.success) {
                const metrics = response.data.metrics;
                const metricsArray: TableProgramMetric[] = Object.keys(metrics).map(id => ({
                    ...metrics[id],
                    key: id
                }));

                // 按CPU使用率降序排序
                metricsArray.sort((a, b) => b.stats.cpu_time_percent - a.stats.cpu_time_percent);

                setMetricsData(metricsArray);
            } else {
                setError(response.errorMsg || '获取节点指标失败');
            }
        } catch (err) {
            console.error('获取节点指标失败:', err);
            setError('获取节点指标失败，请稍后重试');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchMetrics();

        // 设置定时刷新（每10秒）
        const intervalId = setInterval(fetchMetrics, 10000);

        // 组件卸载时清除定时器
        return () => clearInterval(intervalId);
    }, []);

    const handleRefresh = () => {
        fetchMetrics();
    };

    const handleProgramClick = (id: number) => {
        navigate(`/observability/program-detail/${id}`);
    };

    // 格式化百分比
    const formatPercent = (value: number) => {
        return value.toFixed(1) + '%';
    };

    // 格式化纳秒
    const formatNanoseconds = (ns: number) => {
        if (ns < 1000) {
            return `${ns} ns`;
        } else if (ns < 1000000) {
            return `${(ns / 1000).toFixed(2)} μs`;
        } else {
            return `${(ns / 1000000).toFixed(2)} ms`;
        }
    };

    const columns: ColumnsType<TableProgramMetric> = [
        {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            width: 100,
            render: (id: number) => (
                <Button type="link" onClick={() => handleProgramClick(id)}>
                    {id}
                </Button>
            ),
        },
        {
            title: '类型',
            dataIndex: 'type',
            key: 'type',
            width: 150,
            render: (type: string) => (
                <Tag color="blue">{type}</Tag>
            ),
        },
        {
            title: '名称',
            dataIndex: 'name',
            key: 'name',
            ellipsis: true,
            render: (name: string) => (
                <Tooltip title={name}>
                    <span>{name}</span>
                </Tooltip>
            ),
        },
        {
            title: '每周期运行时间 (ns)',
            dataIndex: ['stats', 'avg_run_time_ns'],
            key: 'avg_run_time_ns',
            width: 180,
            render: (value: number) => formatNanoseconds(value),
        },
        {
            title: '平均运行时间 (ns)',
            dataIndex: ['stats', 'total_avg_run_time_ns'],
            key: 'total_avg_run_time_ns',
            width: 180,
            render: (value: number) => formatNanoseconds(value),
        },
        {
            title: '每秒事件数',
            dataIndex: ['stats', 'events_per_second'],
            key: 'events_per_second',
            width: 130,
            sorter: (a, b) => a.stats.events_per_second - b.stats.events_per_second,
        },
        {
            title: 'CPU 使用率',
            dataIndex: ['stats', 'cpu_time_percent'],
            key: 'cpu_time_percent',
            width: 130,
            defaultSortOrder: 'descend',
            sorter: (a, b) => a.stats.cpu_time_percent - b.stats.cpu_time_percent,
            render: (value: number) => {
                let color = 'green';
                if (value > 0.01) {
                    color = 'red';
                } else if (value > 0.005) {
                    color = 'orange';
                }
                return <Tag color={color}>{value}</Tag>;
            },
        },
    ];

    return (
        <div className="node-metrics-page">
            <Card
                title={
                    <Row>
                        <Col span={12}>
                            <Title level={4}>eBPF 程序性能指标</Title>
                        </Col>
                    </Row>
                }
                extra={
                    <Button
                        type="primary"
                        icon={<ReloadOutlined />}
                        onClick={handleRefresh}
                        loading={loading}
                    >
                        刷新
                    </Button>
                }
                className="metrics-card"
            >
                {loading && metricsData.length === 0 ? (
                    <div className="loading-container">
                        <Spin size="large" tip="加载中..." />
                    </div>
                ) : error ? (
                    <div className="error-container">
                        <Empty
                            description={error}
                            image={Empty.PRESENTED_IMAGE_SIMPLE}
                        />
                    </div>
                ) : metricsData.length > 0 ? (
                    <Table
                        columns={columns}
                        dataSource={metricsData}
                        rowKey="key"
                        pagination={false}
                        scroll={{ x: 'max-content', y: 'calc(100vh - 300px)' }}
                    />
                ) : (
                    <Empty description="暂无性能指标数据" />
                )}
            </Card>
        </div>
    );
};

export default NodeMetricsPage; 