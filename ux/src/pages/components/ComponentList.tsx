import { useState, useEffect, useCallback } from 'react'
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
    Tag,
    Form
} from 'antd'
import { SearchOutlined, MoreOutlined, PlusOutlined, UploadOutlined, ExclamationCircleOutlined } from '@ant-design/icons'
import type { TableProps } from 'antd'
import { useNavigate } from 'react-router-dom'
import componentService, { Component } from '../../services/componentService'
import clusterService, { Cluster } from '../../services/clusterService'

interface ComponentWithCluster extends Component {
    clusterName?: string;
}

// 创建集群缓存映射
interface ClusterCache {
    [key: number]: string;
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
    const [clusterCache, setClusterCache] = useState<ClusterCache>({});
    const [currentPage, setCurrentPage] = useState(1);
    const [pageSize, setPageSize] = useState(10);
    const [clustersLoaded, setClustersLoaded] = useState(false);

    // 删除相关状态
    const [deleteForm] = Form.useForm();
    const [deleteModalVisible, setDeleteModalVisible] = useState(false);
    const [componentToDelete, setComponentToDelete] = useState<ComponentWithCluster | null>(null);
    const [deleteLoading, setDeleteLoading] = useState(false);

    // 加载集群列表并创建缓存
    const loadClusters = useCallback(async () => {
        try {
            setLoading(true);
            const clusterList = await clusterService.getClustersByParams();
            setClusters(clusterList);

            // 创建集群ID到名称的映射缓存
            const cache: ClusterCache = {};
            clusterList.forEach(cluster => {
                if (cluster.id) {
                    cache[cluster.id] = cluster.name;
                }
            });
            setClusterCache(cache);
            setClustersLoaded(true);
        } catch (error) {
            console.error('加载集群列表失败:', error);
            message.error('加载集群列表失败');
        } finally {
            setLoading(false);
        }
    }, []);

    // 从缓存中获取集群名称
    const getClusterName = useCallback((clusterId: number): string => {
        return clusterCache[clusterId] || '未知集群';
    }, [clusterCache]);

    // 加载组件列表
    const loadData = useCallback(async (page = currentPage, size = pageSize) => {
        if (!clustersLoaded) {
            return; // 等待集群数据加载完成
        }

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

            // 使用缓存的集群数据为组件添加集群名称
            const componentsWithCluster = result.components.map(component => ({
                ...component,
                clusterName: component.cluster_id ? getClusterName(component.cluster_id) : '未知集群'
            }));

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
    }, [selectedCluster, searchKeyword, clustersLoaded, getClusterName]);

    // 初始加载集群数据
    useEffect(() => {
        loadClusters();
    }, [loadClusters]);

    // 当集群数据加载完成后，加载组件列表
    useEffect(() => {
        if (clustersLoaded) {
            loadData(1);
        }
    }, [clustersLoaded, loadData]);

    // 当选择的集群或搜索关键词变化时，重新加载组件列表
    useEffect(() => {
        if (clustersLoaded && (selectedCluster !== undefined || searchKeyword)) {
            loadData(1);
        }
    }, [selectedCluster, searchKeyword, clustersLoaded, loadData]);

    const handleSearch = () => {
        loadData(1);
    };

    const handleDelete = (record: ComponentWithCluster) => {
        // 调试日志
        console.log('点击删除按钮，组件ID:', record.id, '组件名称:', record.name);

        // 设置要删除的组件并显示确认对话框
        setComponentToDelete(record);
        setDeleteModalVisible(true);
        deleteForm.resetFields();
    };

    // 确认删除
    const confirmDelete = async () => {
        try {
            await deleteForm.validateFields();
            const values = deleteForm.getFieldsValue();

            if (!componentToDelete) {
                message.error('未选择要删除的组件');
                return;
            }

            if (values.confirmName !== componentToDelete.name) {
                message.error('输入的组件名称不匹配');
                return;
            }

            setDeleteLoading(true);

            try {
                await componentService.deleteComponent(componentToDelete.id);
                message.success(`删除组件 ${componentToDelete.name} 成功`);
                setDeleteModalVisible(false);
                // 重新加载数据
                loadData(currentPage);
            } catch (error: any) {
                console.error('删除组件失败:', error);
                message.error(`删除失败: ${error.message || '未知错误'}`);
            } finally {
                setDeleteLoading(false);
            }
        } catch (error) {
            // 表单验证失败
            console.log('表单验证失败:', error);
        }
    };

    // 取消删除
    const cancelDelete = () => {
        setDeleteModalVisible(false);
        setComponentToDelete(null);
    };

    // 直接删除方法，用于测试
    const handleDirectDelete = (record: ComponentWithCluster) => {
        console.log('直接删除，组件ID:', record.id, '组件名称:', record.name);
        setLoading(true);

        componentService.deleteComponent(record.id)
            .then(() => {
                console.log('删除API调用成功');
                message.success(`删除组件 ${record.name} 成功`);
                loadData(currentPage);
            })
            .catch((error) => {
                console.error('删除组件失败:', error);
                if (error.response) {
                    console.error('错误响应:', error.response);
                }
                message.error(`删除失败: ${error.message || '未知错误'}`);
            })
            .finally(() => {
                setLoading(false);
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
                            loadData(page, pageSize || 10);
                        },
                        showSizeChanger: true,
                        showTotal: (total) => `共 ${total} 条记录`,
                    }}
                />
            </Spin>

            {/* 删除确认弹窗 */}
            <Modal
                title={
                    <div>
                        <ExclamationCircleOutlined style={{ color: '#ff4d4f', marginRight: 8 }} />
                        确认删除组件
                    </div>
                }
                open={deleteModalVisible}
                onOk={confirmDelete}
                onCancel={cancelDelete}
                confirmLoading={deleteLoading}
                okText="删除"
                cancelText="取消"
                okButtonProps={{ danger: true }}
            >
                <p>删除操作不可恢复，请谨慎操作！</p>
                <p>请输入组件名称 <strong>{componentToDelete?.name}</strong> 以确认删除：</p>
                <Form form={deleteForm}>
                    <Form.Item
                        name="confirmName"
                        rules={[
                            { required: true, message: '请输入组件名称' },
                            {
                                validator: (_, value) => {
                                    if (value && componentToDelete && value !== componentToDelete.name) {
                                        return Promise.reject(new Error('组件名称不匹配'));
                                    }
                                    return Promise.resolve();
                                }
                            }
                        ]}
                    >
                        <Input placeholder="请输入组件名称" />
                    </Form.Item>
                </Form>
            </Modal>
        </Card>
    )
}

export default ComponentList 