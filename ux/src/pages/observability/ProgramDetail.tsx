import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Card, Spin, Empty, Button, Descriptions, Table, Tag, Divider, Typography, Row, Col, Tooltip, Tabs } from 'antd';
import { ArrowLeftOutlined, ReloadOutlined } from '@ant-design/icons';
import { getProgramDetail, getProgramInstructions, ProgramDetail as ProgramDetailType, MapInfo } from '../../services/topo';
import './ProgramDetail.css';
import type { ColumnsType } from 'antd/es/table';

const { Title, Text } = Typography;
const { TabPane } = Tabs;

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

// eBPF Map 类型映射
const mapTypes: Record<number, string> = {
    1: 'HASH',
    2: 'ARRAY',
    3: 'PROG_ARRAY',
    4: 'PERF_EVENT_ARRAY',
    5: 'PERCPU_HASH',
    6: 'PERCPU_ARRAY',
    7: 'STACK_TRACE',
    8: 'CGROUP_ARRAY',
    9: 'LRU_HASH',
    10: 'LRU_PERCPU_HASH',
    11: 'LPM_TRIE',
    12: 'ARRAY_OF_MAPS',
    13: 'HASH_OF_MAPS',
    14: 'DEVMAP',
    15: 'SOCKMAP',
    16: 'CPUMAP',
    17: 'XSKMAP',
    18: 'SOCKHASH',
    19: 'CGROUP_STORAGE',
    20: 'REUSEPORT_SOCKARRAY',
    21: 'PERCPU_CGROUP_STORAGE',
    22: 'QUEUE',
    23: 'STACK',
    24: 'SK_STORAGE',
    25: 'DEVMAP_HASH',
    26: 'STRUCT_OPS',
    27: 'RINGBUF',
    28: 'INODE_STORAGE',
    29: 'TASK_STORAGE',
    30: 'BLOOM_FILTER',
};

const ProgramDetailPage: React.FC = () => {
    const { progId } = useParams<{ progId: string }>();
    const navigate = useNavigate();
    const [loading, setLoading] = useState<boolean>(true);
    const [programDetail, setProgramDetail] = useState<ProgramDetailType | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [instructions, setInstructions] = useState<string>('');
    const [instructionsLoading, setInstructionsLoading] = useState<boolean>(false);
    const [instructionsError, setInstructionsError] = useState<string | null>(null);

    const fetchProgramDetail = async () => {
        if (!progId) {
            setError('程序ID不能为空');
            setLoading(false);
            return;
        }

        setLoading(true);
        setError(null);
        try {
            const data = await getProgramDetail(parseInt(progId));
            setProgramDetail(data);
        } catch (err) {
            console.error('获取程序详情失败:', err);
            setError('获取程序详情失败，请稍后重试');
        } finally {
            setLoading(false);
        }
    };

    const fetchProgramInstructions = async () => {
        if (!progId) return;

        setInstructionsLoading(true);
        setInstructionsError(null);
        try {
            const data = await getProgramInstructions(parseInt(progId));
            setInstructions(data);
        } catch (err) {
            console.error('获取程序指令失败:', err);
            setInstructionsError('获取程序指令失败，请稍后重试');
        } finally {
            setInstructionsLoading(false);
        }
    };

    // 当在标签页切换到指令页时加载指令数据
    const handleTabChange = (key: string) => {
        if (key === 'instructions' && !instructions && !instructionsLoading) {
            fetchProgramInstructions();
        }
    };

    useEffect(() => {
        fetchProgramDetail();
    }, [progId]);

    const handleRefresh = () => {
        fetchProgramDetail();
    };

    const handleBack = () => {
        navigate('/observability/node-resources');
    };

    // 格式化日期时间
    const formatDateTime = (dateTimeString: string) => {
        if (!dateTimeString) return '-';

        try {
            const date = new Date(dateTimeString);
            return date.toLocaleString();
        } catch (e) {
            return dateTimeString;
        }
    };

    // 定义Maps表格列
    const mapColumns: ColumnsType<MapInfo> = [
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
                <Tag color="blue">{mapTypes[type] || `类型 ${type}`}</Tag>
            ),
        },
        {
            title: '键大小',
            dataIndex: 'KeySize',
            key: 'keySize',
            render: (size: number) => `${size} 字节`,
        },
        {
            title: '值大小',
            dataIndex: 'ValueSize',
            key: 'valueSize',
            render: (size: number) => `${size} 字节`,
        },
        {
            title: '最大条目数',
            dataIndex: 'MaxEntries',
            key: 'maxEntries',
        },
        {
            title: '状态',
            dataIndex: 'Frozen',
            key: 'frozen',
            render: (frozen: boolean) => (
                <Tag color={frozen ? 'red' : 'green'}>
                    {frozen ? '已冻结' : '活动中'}
                </Tag>
            ),
        },
    ];

    return (
        <div className="program-detail-page">
            <Card
                title={
                    <div className="detail-header">
                        <Button
                            icon={<ArrowLeftOutlined />}
                            onClick={handleBack}
                            style={{ marginRight: 16 }}
                        >
                            返回
                        </Button>
                        <span>程序详情</span>
                    </div>
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
                className="program-detail-card"
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
                ) : programDetail ? (
                    <div className="detail-content">
                        <Row gutter={[0, 24]}>
                            <Col span={24}>
                                <Title level={4}>
                                    {programDetail.Name || `程序 #${programDetail.ID}`}
                                    <Tag color="blue" style={{ marginLeft: 12 }}>
                                        {programTypes[programDetail.Type] || `类型 ${programDetail.Type}`}
                                    </Tag>
                                </Title>
                            </Col>

                            <Col span={24}>
                                <Tabs defaultActiveKey="info" onChange={handleTabChange}>
                                    <TabPane tab="基本信息" key="info">
                                        <Descriptions
                                            title="基本信息"
                                            bordered
                                            column={{ xxl: 4, xl: 3, lg: 3, md: 2, sm: 1, xs: 1 }}
                                        >
                                            <Descriptions.Item label="ID">
                                                {programDetail.ID}
                                            </Descriptions.Item>
                                            <Descriptions.Item label="名称">
                                                {programDetail.Name || '-'}
                                            </Descriptions.Item>
                                            <Descriptions.Item label="类型">
                                                {programTypes[programDetail.Type] || `类型 ${programDetail.Type}`}
                                            </Descriptions.Item>
                                            <Descriptions.Item label="标签">
                                                <Text copyable>{programDetail.Tag || '-'}</Text>
                                            </Descriptions.Item>
                                            <Descriptions.Item label="BTF ID">
                                                {programDetail.BTF || '-'}
                                            </Descriptions.Item>
                                            <Descriptions.Item label="创建者 UID">
                                                {programDetail.HaveCreatedByUID ? programDetail.CreatedByUID : '-'}
                                            </Descriptions.Item>
                                            <Descriptions.Item label="加载时间">
                                                {formatDateTime(programDetail.LoadTime) || '-'}
                                            </Descriptions.Item>
                                        </Descriptions>

                                        <Divider orientation="left">关联的 Map</Divider>
                                        <Table
                                            columns={mapColumns}
                                            dataSource={programDetail.MapsDetail.map((item) => ({ ...item, key: item.ID }))}
                                            rowKey="ID"
                                            scroll={{ x: 'max-content' }}
                                        />
                                    </TabPane>
                                    <TabPane tab="程序指令" key="instructions">
                                        {instructionsLoading ? (
                                            <div className="loading-container">
                                                <Spin size="default" tip="加载中..." />
                                            </div>
                                        ) : instructionsError ? (
                                            <div className="error-container">
                                                <Empty
                                                    description={instructionsError}
                                                    image={Empty.PRESENTED_IMAGE_SIMPLE}
                                                />
                                            </div>
                                        ) : (
                                            <pre className="instruction-code">
                                                {instructions || '暂无指令数据'}
                                            </pre>
                                        )}
                                    </TabPane>
                                </Tabs>
                            </Col>
                        </Row>
                    </div>
                ) : (
                    <Empty description="未找到程序数据" />
                )}
            </Card>
        </div>
    );
};

export default ProgramDetailPage; 