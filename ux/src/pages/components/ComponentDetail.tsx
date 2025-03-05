import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
    Card,
    Descriptions,
    Tabs,
    Table,
    Tag,
    Button,
    Spin,
    message,
    Typography,
    Collapse,
    Space,
    Divider,
    Row,
    Col
} from 'antd';
import { ArrowLeftOutlined, CodeOutlined, DatabaseOutlined } from '@ant-design/icons';
import componentService, { Component, Program, Map } from '../../services/componentService';
import clusterService from '../../services/clusterService';

const { Title, Text } = Typography;
const { TabPane } = Tabs;
const { Panel } = Collapse;

const ComponentDetail = () => {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const [loading, setLoading] = useState(true);
    const [component, setComponent] = useState<Component | null>(null);
    const [clusterName, setClusterName] = useState<string>('未知集群');

    // 加载组件详情
    const loadComponentDetail = async () => {
        if (!id) return;

        setLoading(true);
        try {
            const componentId = parseInt(id);
            const componentData = await componentService.getComponent(componentId);
            setComponent(componentData);

            // 获取集群名称
            if (componentData.cluster_id) {
                try {
                    const clusterInfo = await clusterService.getCluster(componentData.cluster_id);
                    setClusterName(clusterInfo.name);
                } catch (error) {
                    console.error(`获取集群 ${componentData.cluster_id} 信息失败:`, error);
                }
            }
        } catch (error) {
            console.error('加载组件详情失败:', error);
            message.error('加载组件详情失败');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        loadComponentDetail();
    }, [id]);

    // 渲染程序表格
    const renderProgramsTable = () => {
        const columns = [
            {
                title: 'ID',
                dataIndex: 'id',
                key: 'id',
                width: 80,
            },
            {
                title: '名称',
                dataIndex: 'name',
                key: 'name',
            },
            {
                title: '描述',
                dataIndex: 'description',
                key: 'description',
            },
        ];

        const expandedRowRender = (record: Program) => {
            return (
                <div style={{ padding: '0 20px' }}>
                    <Row gutter={[16, 16]}>
                        <Col span={12}>
                            <Card title="程序规格 (Spec)" size="small">
                                <Descriptions column={1} size="small" bordered>
                                    <Descriptions.Item label="名称">{record.spec.Name}</Descriptions.Item>
                                    <Descriptions.Item label="类型">{getProgramTypeName(record.spec.Type)}</Descriptions.Item>
                                    <Descriptions.Item label="附加类型">{getAttachTypeName(record.spec.AttachType)}</Descriptions.Item>
                                    <Descriptions.Item label="附加目标">{record.spec.AttachTo || '-'}</Descriptions.Item>
                                    <Descriptions.Item label="段名称">{record.spec.SectionName || '-'}</Descriptions.Item>
                                    <Descriptions.Item label="标志">{record.spec.Flags}</Descriptions.Item>
                                    <Descriptions.Item label="许可证">{record.spec.License}</Descriptions.Item>
                                    <Descriptions.Item label="内核版本">{record.spec.KernelVersion || '自动检测'}</Descriptions.Item>
                                </Descriptions>
                            </Card>
                        </Col>
                        <Col span={12}>
                            <Card title="程序属性 (Properties)" size="small">
                                <Descriptions column={1} size="small" bordered>
                                    <Descriptions.Item label="CGroup路径">{record.properties.CGroupPath || '-'}</Descriptions.Item>
                                    <Descriptions.Item label="固定路径">{record.properties.PinPath || '-'}</Descriptions.Item>
                                    <Descriptions.Item label="链接固定路径">{record.properties.LinkPinPath || '-'}</Descriptions.Item>
                                    <Descriptions.Item label="TC配置">{record.properties.Tc ? JSON.stringify(record.properties.Tc) : '-'}</Descriptions.Item>
                                </Descriptions>
                            </Card>
                        </Col>
                    </Row>
                </div>
            );
        };

        return (
            <Table
                columns={columns}
                dataSource={component?.programs || []}
                rowKey="id"
                expandable={{
                    expandedRowRender,
                    expandRowByClick: true,
                }}
                pagination={false}
            />
        );
    };

    // 渲染Map表格
    const renderMapsTable = () => {
        const columns = [
            {
                title: 'ID',
                dataIndex: 'id',
                key: 'id',
                width: 80,
            },
            {
                title: '名称',
                dataIndex: 'name',
                key: 'name',
            },
            {
                title: '描述',
                dataIndex: 'description',
                key: 'description',
            },
        ];

        const expandedRowRender = (record: Map) => {
            return (
                <div style={{ padding: '0 20px' }}>
                    <Row gutter={[16, 16]}>
                        <Col span={12}>
                            <Card title="Map规格 (Spec)" size="small">
                                <Descriptions column={1} size="small" bordered>
                                    <Descriptions.Item label="名称">{record.spec.Name}</Descriptions.Item>
                                    <Descriptions.Item label="类型">{getMapTypeName(record.spec.Type)}</Descriptions.Item>
                                    <Descriptions.Item label="键大小">{record.spec.KeySize} 字节</Descriptions.Item>
                                    <Descriptions.Item label="值大小">{record.spec.ValueSize} 字节</Descriptions.Item>
                                    <Descriptions.Item label="最大条目数">{record.spec.MaxEntries}</Descriptions.Item>
                                    <Descriptions.Item label="标志">{record.spec.Flags}</Descriptions.Item>
                                    <Descriptions.Item label="固定类型">{getPinTypeName(record.spec.Pinning)}</Descriptions.Item>
                                </Descriptions>
                            </Card>
                        </Col>
                        <Col span={12}>
                            <Card title="Map属性 (Properties)" size="small">
                                <Descriptions column={1} size="small" bordered>
                                    <Descriptions.Item label="固定路径">{record.properties.PinPath || '-'}</Descriptions.Item>
                                </Descriptions>
                            </Card>
                        </Col>
                    </Row>
                </div>
            );
        };

        return (
            <Table
                columns={columns}
                dataSource={component?.maps || []}
                rowKey="id"
                expandable={{
                    expandedRowRender,
                    expandRowByClick: true,
                }}
                pagination={false}
            />
        );
    };

    // 获取程序类型名称
    const getProgramTypeName = (type: number): string => {
        const typeMap: Record<number, string> = {
            0: '未指定',
            1: 'Kprobe',
            2: 'Tracepoint',
            3: 'SocketFilter',
            4: 'XDP',
            5: 'PerfEvent',
            6: 'CGroupSKB',
            7: 'CGroupSock',
            8: 'LWTIn',
            9: 'LWTOut',
            10: 'LWTXmit',
            11: 'SockOps',
            12: 'SK_SKB',
            13: 'CGroupDevice',
            14: 'SK_MSG',
            15: 'RawTracepoint',
            16: 'CGroupSockAddr',
            17: 'LWTSeg6Local',
            18: 'LircMode2',
            19: 'SkReuseport',
            20: 'FlowDissector',
            21: 'CGroupSysctl',
            22: 'RawTracepointWritable',
            23: 'CGroupSockopt',
            24: 'Tracing',
            25: 'StructOps',
            26: 'Extension',
            27: 'LSM',
            28: 'SkLookup',
            29: 'Syscall',
        };
        return typeMap[type] || `未知类型(${type})`;
    };

    // 获取附加类型名称
    const getAttachTypeName = (type: number): string => {
        const typeMap: Record<number, string> = {
            0: '未指定',
            1: 'CGroupInetIngress',
            2: 'CGroupInetEgress',
            3: 'CGroupInetSockCreate',
            4: 'CGroupSockOps',
            5: 'SkSKBStreamParser',
            6: 'SkSKBStreamVerdict',
            7: 'CGroupDevice',
            8: 'SkMsgVerdict',
            9: 'CGroupInet4Bind',
            10: 'CGroupInet6Bind',
            11: 'CGroupInet4Connect',
            12: 'CGroupInet6Connect',
            13: 'CGroupInet4PostBind',
            14: 'CGroupInet6PostBind',
            15: 'CGroupUDP4Sendmsg',
            16: 'CGroupUDP6Sendmsg',
            17: 'LircMode2',
            18: 'FlowDissector',
            19: 'CGroupSysctl',
            20: 'CGroupUDP4Recvmsg',
            21: 'CGroupUDP6Recvmsg',
            22: 'CGroupGetsockopt',
            23: 'CGroupSetsockopt',
            24: 'TraceRawTp',
            25: 'TraceFentry',
            26: 'TraceFexit',
            27: 'ModifyReturn',
            28: 'LSMMac',
            29: 'TraceIter',
            30: 'CgroupInet4Getpeername',
            31: 'CgroupInet6Getpeername',
            32: 'CgroupInet4Getsockname',
            33: 'CgroupInet6Getsockname',
            34: 'XdpDevmap',
            35: 'CgroupInetSockRelease',
            36: 'XdpCpumap',
            37: 'SkLookup',
            38: 'XDP',
            39: 'SkSKBVerdict',
        };
        return typeMap[type] || `未知类型(${type})`;
    };

    // 获取Map类型名称
    const getMapTypeName = (type: number): string => {
        const typeMap: Record<number, string> = {
            0: '未指定',
            1: 'Hash',
            2: 'Array',
            3: 'ProgramArray',
            4: 'PerfEventArray',
            5: 'PerCPUHash',
            6: 'PerCPUArray',
            7: 'StackTrace',
            8: 'CGroupArray',
            9: 'LRUHash',
            10: 'LRUPerCPUHash',
            11: 'LPMTrie',
            12: 'ArrayOfMaps',
            13: 'HashOfMaps',
            14: 'DevMap',
            15: 'SockMap',
            16: 'CPUMap',
            17: 'XSKMap',
            18: 'SockHash',
            19: 'CGroupStorage',
            20: 'ReusePortSockArray',
            21: 'PerCPUCGroupStorage',
            22: 'Queue',
            23: 'Stack',
            24: 'SkStorage',
            25: 'DevMapHash',
            26: 'StructOps',
            27: 'RingBuf',
            28: 'InodeStorage',
            29: 'TaskStorage',
        };
        return typeMap[type] || `未知类型(${type})`;
    };

    // 获取固定类型名称
    const getPinTypeName = (type: number): string => {
        const typeMap: Record<number, string> = {
            0: '不固定',
            1: '按路径固定',
            2: '按名称固定',
        };
        return typeMap[type] || `未知类型(${type})`;
    };

    if (loading) {
        return (
            <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100%' }}>
                <Spin size="large" tip="加载组件详情..." />
            </div>
        );
    }

    if (!component) {
        return (
            <Card>
                <div style={{ textAlign: 'center' }}>
                    <Title level={4}>未找到组件信息</Title>
                    <Button type="primary" onClick={() => navigate('/components/list')}>
                        返回组件列表
                    </Button>
                </div>
            </Card>
        );
    }

    return (
        <div>
            <div style={{ marginBottom: 16 }}>
                <Button
                    type="link"
                    icon={<ArrowLeftOutlined />}
                    onClick={() => navigate('/components/list')}
                    style={{ paddingLeft: 0 }}
                >
                    返回组件列表
                </Button>
            </div>

            <Card title={<Title level={4}>{component.name} 详情</Title>}>
                <Descriptions title="基本信息" bordered>
                    <Descriptions.Item label="ID">{component.id}</Descriptions.Item>
                    <Descriptions.Item label="名称">{component.name}</Descriptions.Item>
                    <Descriptions.Item label="所属集群">
                        <Tag color="blue">{clusterName}</Tag>
                    </Descriptions.Item>
                    <Descriptions.Item label="程序数量">{component.programs?.length || 0}</Descriptions.Item>
                    <Descriptions.Item label="Map数量">{component.maps?.length || 0}</Descriptions.Item>
                </Descriptions>

                <Divider />

                <Tabs defaultActiveKey="programs">
                    <TabPane
                        tab={
                            <span>
                                <CodeOutlined />
                                程序信息 ({component.programs?.length || 0})
                            </span>
                        }
                        key="programs"
                    >
                        {renderProgramsTable()}
                    </TabPane>
                    <TabPane
                        tab={
                            <span>
                                <DatabaseOutlined />
                                Map信息 ({component.maps?.length || 0})
                            </span>
                        }
                        key="maps"
                    >
                        {renderMapsTable()}
                    </TabPane>
                </Tabs>
            </Card>
        </div>
    );
};

export default ComponentDetail; 