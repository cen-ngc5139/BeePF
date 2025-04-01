import React, { useEffect, useState } from 'react';
import { Card, Spin, Empty, Button, Table, Typography, Badge, Tooltip, Tag } from 'antd';
import { ReloadOutlined } from '@ant-design/icons';
import { getPrograms, ProgramInfo } from '../../services/topo';
import './NodeResources.css';
import type { ColumnsType } from 'antd/es/table';

const { Title } = Typography;

interface ProgramData extends ProgramInfo {
    key: string;
}

// eBPF 程序类型映射
const programTypes: Record<number, string> = {
    0: '未指定',
    1: 'SOCKET_FILTER',
    2: 'KPROBE',
    3: 'SCHED_CLS',
    4: 'SCHED_ACT',
    5: 'TRACEPOINT',
    6: 'XDP',
    7: 'PERF_EVENT',
    8: 'CGROUP_SKB',
    9: 'CGROUP_SOCK',
    10: 'LWT_IN',
    11: 'LWT_OUT',
    12: 'LWT_XMIT',
    13: 'SOCK_OPS',
    14: 'SK_SKB',
    15: 'CGROUP_DEVICE',
    16: 'SK_MSG',
    17: 'RAW_TRACEPOINT',
    18: 'CGROUP_SOCK_ADDR',
    19: 'LWT_SEG6LOCAL',
    20: 'LIRC_MODE2',
    21: 'SK_REUSEPORT',
    22: 'FLOW_DISSECTOR',
    23: 'CGROUP_SYSCTL',
    24: 'RAW_TRACEPOINT_WRITABLE',
    25: 'CGROUP_SOCKOPT',
    26: 'TRACING',
    27: 'STRUCT_OPS',
    28: 'EXT',
    29: 'LSM',
    30: 'SK_LOOKUP',
    31: 'SYSCALL',
};

const NodeResourcesPage: React.FC = () => {
    const [loading, setLoading] = useState<boolean>(true);
    const [programs, setPrograms] = useState<ProgramInfo[]>([]);
    const [error, setError] = useState<string | null>(null);
    const [pagination, setPagination] = useState({
        current: 1,
        pageSize: 10,
        showSizeChanger: true,
        showQuickJumper: true,
    });

    const fetchPrograms = async () => {
        setLoading(true);
        setError(null);
        try {
            const data = await getPrograms();
            setPrograms(data);
            setPagination(prev => ({
                ...prev,
                current: 1
            }));
        } catch (err) {
            console.error('获取节点资源数据失败:', err);
            setError('获取节点资源数据失败，请稍后重试');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchPrograms();
    }, []);

    const handleRefresh = () => {
        fetchPrograms();
    };

    // 将数组转换为表格数据
    const programsArray: ProgramData[] = programs.map((info, index) => ({
        key: `${info.ID}-${index}`,
        ...info
    }));

    // 格式化纳秒为可读时间
    const formatTime = (nanoseconds: number) => {
        if (nanoseconds === 0) return '0';

        if (nanoseconds < 1000) {
            return `${nanoseconds}ns`;
        } else if (nanoseconds < 1000000) {
            return `${(nanoseconds / 1000).toFixed(2)}µs`;
        } else if (nanoseconds < 1000000000) {
            return `${(nanoseconds / 1000000).toFixed(2)}ms`;
        } else {
            return `${(nanoseconds / 1000000000).toFixed(2)}s`;
        }
    };

    const columns: ColumnsType<ProgramData> = [
        {
            title: 'ID',
            dataIndex: 'ID',
            key: 'id',
            sorter: (a, b) => a.ID - b.ID,
        },
        {
            title: '名称',
            dataIndex: 'Name',
            key: 'name',
            sorter: (a, b) => (a.Name || '').localeCompare(b.Name || ''),
            render: (name: string) => (
                <Tooltip title={name}>
                    <span className="ellipsis-text">{name || '未命名'}</span>
                </Tooltip>
            ),
        },
        {
            title: '类型',
            dataIndex: 'Type',
            key: 'type',
            render: (type: number) => (
                <Badge
                    status="processing"
                    text={programTypes[type] || `类型 ${type}`}
                />
            ),
            filters: Array.from(new Set(programsArray.map(p => p.Type)))
                .filter(Boolean)
                .map(type => ({
                    text: programTypes[type] || `类型 ${type}`,
                    value: type
                })),
            onFilter: (value, record) => record.Type === value,
        },
        {
            title: '标签',
            dataIndex: 'Tag',
            key: 'tag',
            render: (tag: string) => (
                <Tag color="blue">{tag}</Tag>
            ),
        },
        {
            title: '加载时间',
            dataIndex: 'LoadTime',
            key: 'loadTime',
            sorter: (a, b) => (a.LoadTime || 0) - (b.LoadTime || 0),
            render: (time: number) => (
                <Tooltip title={`${time} 纳秒`}>
                    {formatTime(time || 0)}
                </Tooltip>
            ),
        },
        {
            title: '关联的 Maps',
            dataIndex: 'Maps',
            key: 'maps',
            render: (maps: number[]) => (
                <Tooltip title={maps?.join(', ') || '无'}>
                    {maps?.length || 0} 个
                </Tooltip>
            ),
        },
        {
            title: 'BTF ID',
            dataIndex: 'BTF',
            key: 'btf',
            width: 100,
        },
    ];

    const handleTableChange = (newPagination: any) => {
        setPagination({
            ...pagination,
            current: newPagination.current,
            pageSize: newPagination.pageSize,
        });
    };

    return (
        <div className="node-resources-page">
            <Card
                title="eBPF 节点资源"
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
                className="node-resources-card"
            >
                {loading ? (
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
                ) : programsArray.length > 0 ? (
                    <Table
                        dataSource={programsArray}
                        columns={columns}
                        onChange={handleTableChange}
                        pagination={{
                            ...pagination,
                            total: programsArray.length,
                            showTotal: (total) => `共 ${total} 个 eBPF 程序`,
                            position: ['bottomRight']
                        }}
                        rowKey="key"
                        scroll={{ x: 'max-content' }}
                    />
                ) : (
                    <Empty description="暂无节点资源数据" />
                )}
            </Card>
        </div>
    );
};

export default NodeResourcesPage; 