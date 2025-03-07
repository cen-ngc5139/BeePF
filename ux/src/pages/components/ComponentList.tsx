import { useState, useEffect } from 'react'
import {
    Card,
    Table,
    Select,
    Input,
    Button,
    Space,
    Dropdown,
    Modal,
    message,
    Spin,
    Tag
} from 'antd'
import { SearchOutlined, MoreOutlined, PlusOutlined, UploadOutlined } from '@ant-design/icons'
import type { TableProps } from 'antd'
import { useNavigate } from 'react-router-dom'
import componentService, { Component } from '../../services/componentService'
import clusterService, { Cluster } from '../../services/clusterService'

interface ComponentWithCluster extends Component {
    clusterName?: string;
}

const ComponentList = () => {
    const navigate = useNavigate();
    const [loading, setLoading] = useState(false);
    const [selectedCluster, setSelectedCluster] = useState<number>();
    const [searchKeyword, setSearchKeyword] = useState('');
    const [modalVisible, setModalVisible] = useState(false);
    const [selectedComponent, setSelectedComponent] = useState<ComponentWithCluster>();
    const [components, setComponents] = useState<ComponentWithCluster[]>([]);
    const [total, setTotal] = useState(0);
    const [clusters, setClusters] = useState<Cluster[]>([]);
    const [currentPage, setCurrentPage] = useState(1);
    const [pageSize, setPageSize] = useState(10);

    // 加载集群列表
    const loadClusters = async () => {
        try {
            const clusterList = await clusterService.getClustersByParams();
            setClusters(clusterList);
        } catch (error) {
            console.error('加载集群列表失败:', error);
            message.error('加载集群列表失败');
        }
    };

    // 加载组件列表
    const loadData = async (page = currentPage, size = pageSize) => {
        setLoading(true);
        try {
            const params: any = {
                pageNum: page,
                pageSize: size,
            };

            if (selectedCluster) {
                params.cluster_id = selectedCluster;
            }

            if (searchKeyword) {
                params.keyword = searchKeyword;
            }

            const result = await componentService.getComponentList(params);

            // 获取组件列表后，为每个组件添加集群名称
            const componentsWithCluster = await Promise.all(
                result.components.map(async (component) => {
                    let clusterName = '未知集群';

                    // 从已加载的集群列表中查找匹配的集群
                    const matchedCluster = clusters.find(cluster => cluster.id === component.cluster_id);

                    if (matchedCluster) {
                        clusterName = matchedCluster.name;
                    } else if (component.cluster_id) {
                        // 如果在已加载的集群列表中找不到，则单独请求该集群信息
                        try {
                            const clusterInfo = await clusterService.getCluster(component.cluster_id);
                            clusterName = clusterInfo.name;
                        } catch (error) {
                            console.error(`获取集群 ${component.cluster_id} 信息失败:`, error);
                        }
                    }

                    return {
                        ...component,
                        clusterName
                    };
                })
            );

            setComponents(componentsWithCluster);
            setTotal(result.total);
            setCurrentPage(page);
            setPageSize(size);
        } catch (error) {
            console.error('加载组件列表失败:', error);
            message.error('加载组件列表失败');
        } finally {
            setLoading(false);
        }
    };

    // 初始加载
    useEffect(() => {
        loadClusters();
    }, []);

    useEffect(() => {
        loadData();
    }, [selectedCluster, clusters]);

    const handleSearch = () => {
        loadData(1);
    };

    const handleDelete = (record: ComponentWithCluster) => {
        Modal.confirm({
            title: '确认删除',
            content: `确定要删除组件 ${record.name} 吗？`,
            onOk: () => {
                // 实现删除逻辑
                message.success('删除成功');
                loadData();
            },
        });
    };

    const columns: TableProps<ComponentWithCluster>['columns'] = [
        {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            width: 80,
        },
        {
            title: '组件名称',
            dataIndex: 'name',
            key: 'name',
            render: (text, record) => (
                <a onClick={() => navigate(`/component/${record.id}`)}>{text}</a>
            ),
        },
        {
            title: '所属集群',
            dataIndex: 'clusterName',
            key: 'clusterName',
            render: (text, record) => (
                <Tag color="blue">{text || '未知集群'}</Tag>
            ),
        },
        {
            title: '程序数量',
            key: 'programCount',
            render: (_, record) => (
                <span>{record.programs ? record.programs.length : 0}</span>
            ),
        },
        {
            title: 'Map数量',
            key: 'mapCount',
            render: (_, record) => (
                <span>{record.maps ? record.maps.length : 0}</span>
            ),
        },
        {
            title: '操作',
            key: 'action',
            width: 120,
            render: (_, record) => (
                <Dropdown
                    menu={{
                        items: [
                            {
                                key: 'view',
                                label: '查看详情',
                                onClick: () => navigate(`/component/${record.id}`),
                            },
                            {
                                key: 'delete',
                                label: '删除',
                                onClick: () => handleDelete(record),
                                danger: true,
                            },
                        ],
                    }}
                >
                    <Button type="text" icon={<MoreOutlined />} />
                </Dropdown>
            ),
        },
    ];

    return (
        <Card title="组件列表">
            <Space style={{ marginBottom: 16 }}>
                <Select
                    placeholder="选择集群"
                    style={{ width: 200 }}
                    options={clusters.map(cluster => ({ value: cluster.id, label: cluster.name }))}
                    value={selectedCluster}
                    onChange={setSelectedCluster}
                    allowClear
                />
                <Input
                    placeholder="请输入组件名称"
                    style={{ width: 200 }}
                    value={searchKeyword}
                    onChange={(e) => setSearchKeyword(e.target.value)}
                    onPressEnter={handleSearch}
                />
                <Button type="primary" icon={<SearchOutlined />} onClick={handleSearch}>
                    搜索
                </Button>
                <Button
                    type="primary"
                    icon={<PlusOutlined />}
                    onClick={() => navigate('/components/create')}
                >
                    新建组件
                </Button>
                <Button
                    type="primary"
                    icon={<UploadOutlined />}
                    onClick={() => navigate('/components/upload')}
                >
                    上传组件
                </Button>
            </Space>

            <Spin spinning={loading}>
                <Table
                    columns={columns}
                    dataSource={components}
                    rowKey="id"
                    pagination={{
                        current: currentPage,
                        pageSize: pageSize,
                        total: total,
                        onChange: (page, pageSize) => {
                            loadData(page, pageSize);
                        },
                        showSizeChanger: true,
                        showTotal: (total) => `共 ${total} 条记录`,
                    }}
                />
            </Spin>
        </Card>
    )
}

export default ComponentList 