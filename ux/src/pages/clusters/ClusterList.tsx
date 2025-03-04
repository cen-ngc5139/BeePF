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
    Tag,
    Form,
} from 'antd'
import { SearchOutlined, PlusOutlined, MoreOutlined, ExclamationCircleOutlined } from '@ant-design/icons'
import type { TableProps } from 'antd'
import { useNavigate } from 'react-router-dom'
import clusterService, { Cluster, ClusterListParams } from '../../services/clusterService'

const ClusterList = () => {
    const [selectedEnvironment, setSelectedEnvironment] = useState<string>()
    const [searchKeyword, setSearchKeyword] = useState('')
    const [loading, setLoading] = useState(false)
    const [data, setData] = useState<Cluster[]>([])
    const [total, setTotal] = useState(0)
    const [currentPage, setCurrentPage] = useState(1)
    const [pageSize, setPageSize] = useState(10)
    const [deleteForm] = Form.useForm()
    const [deleteModalVisible, setDeleteModalVisible] = useState(false)
    const [clusterToDelete, setClusterToDelete] = useState<Cluster | null>(null)
    const [deleteLoading, setDeleteLoading] = useState(false)
    const navigate = useNavigate()

    // 环境选项
    const environments = [
        { value: 'prod', label: '生产环境' },
        { value: 'test', label: '测试环境' },
        { value: 'dev', label: '开发环境' },
    ]

    // 加载集群列表数据
    const loadData = async () => {
        setLoading(true)
        try {
            const params: ClusterListParams = {
                pageSize,
                pageNum: currentPage,
            }

            if (selectedEnvironment) {
                params.environment = selectedEnvironment
            }

            if (searchKeyword) {
                params.keyword = searchKeyword
            }

            console.log('发送请求参数:', params)

            const response = await clusterService.getClusterList(params)
            console.log('处理后的响应数据:', response)

            if (response.clusters && response.clusters.length > 0) {
                setData(response.clusters)
                setTotal(response.total)
                console.log('设置表格数据:', response.clusters)
            } else {
                console.warn('接口返回的集群列表为空')
                setData([])
                setTotal(0)
            }
        } catch (error) {
            console.error('获取集群列表失败:', error)
            message.error('获取集群列表失败')
        } finally {
            setLoading(false)
        }
    }

    // 首次加载和筛选条件变化时重新加载数据
    useEffect(() => {
        loadData()
    }, [currentPage, pageSize])

    const handleSearch = () => {
        setCurrentPage(1) // 重置到第一页
        loadData()
    }

    const handleEdit = (record: Cluster) => {
        navigate(`/clusters/edit/${record.id}`)
    }

    const handleDelete = (record: Cluster) => {
        setClusterToDelete(record)
        setDeleteModalVisible(true)
        deleteForm.resetFields()
    }

    const confirmDelete = async () => {
        try {
            await deleteForm.validateFields()
            const values = deleteForm.getFieldsValue()

            if (!clusterToDelete) {
                message.error('未选择要删除的集群')
                return
            }

            if (values.confirmName !== clusterToDelete.name) {
                message.error('输入的集群名称不匹配')
                return
            }

            setDeleteLoading(true)

            try {
                await clusterService.deleteCluster(clusterToDelete.id as number)
                message.success('删除成功')
                setDeleteModalVisible(false)
                loadData() // 重新加载数据
            } catch (error) {
                console.error('删除集群失败:', error)
                message.error('删除集群失败')
            } finally {
                setDeleteLoading(false)
            }
        } catch (error) {
            // 表单验证失败
        }
    }

    const cancelDelete = () => {
        setDeleteModalVisible(false)
        setClusterToDelete(null)
    }

    const handleAdd = () => {
        navigate('/clusters/create')
    }

    const handleTableChange = (pagination: any) => {
        setCurrentPage(pagination.current)
        setPageSize(pagination.pageSize)
    }

    // 状态映射
    const getStatusTag = (status: number | string) => {
        if (typeof status === 'number') {
            switch (status) {
                case 0:
                    return <Tag color="success">正常</Tag>
                case 1:
                    return <Tag color="default">停用</Tag>
                default:
                    return <Tag color="default">未知</Tag>
            }
        } else {
            const statusMap = {
                active: { text: '正常', color: 'success' },
                inactive: { text: '停用', color: 'default' },
            }
            const { text, color } = statusMap[status as keyof typeof statusMap] || { text: status, color: 'default' }
            return <Tag color={color}>{text}</Tag>
        }
    }

    const columns: TableProps<Cluster>['columns'] = [
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
            width: 120,
        },
        {
            title: '中文名称',
            dataIndex: 'cnname',
            key: 'cnname',
            width: 120,
        },
        {
            title: '主节点地址',
            dataIndex: 'master',
            key: 'master',
            width: 180,
        },
        {
            title: '环境',
            dataIndex: 'environment',
            key: 'environment',
            width: 100,
            render: (env: string) => {
                const envMap = {
                    prod: { text: '生产环境', color: 'red' },
                    test: { text: '测试环境', color: 'green' },
                    dev: { text: '开发环境', color: 'blue' },
                }
                const { text, color } = envMap[env as keyof typeof envMap] || { text: env, color: 'default' }
                return <Tag color={color}>{text}</Tag>
            },
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            width: 100,
            render: (status) => getStatusTag(status),
        },
        {
            title: '创建者',
            dataIndex: 'creator',
            key: 'creator',
            width: 100,
        },
        {
            title: '创建时间',
            dataIndex: 'createdat',
            key: 'createdat',
            width: 180,
            render: (time: string) => time ? new Date(time).toLocaleString() : '-',
        },
        {
            title: '更新时间',
            dataIndex: 'updateat',
            key: 'updateat',
            width: 180,
            render: (time: string) => time ? new Date(time).toLocaleString() : '-',
        },
        {
            title: '操作',
            key: 'action',
            width: 120,
            fixed: 'right',
            render: (_, record) => (
                <Dropdown
                    menu={{
                        items: [
                            {
                                key: 'edit',
                                label: '编辑',
                                onClick: () => handleEdit(record),
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
    ]

    return (
        <Card>
            <Space style={{ marginBottom: 16 }}>
                <Select
                    placeholder="选择环境"
                    style={{ width: 200 }}
                    options={environments}
                    value={selectedEnvironment}
                    onChange={setSelectedEnvironment}
                    allowClear
                />
                <Input.Search
                    placeholder="请输入搜索关键词"
                    style={{ width: 200 }}
                    value={searchKeyword}
                    onChange={(e) => setSearchKeyword(e.target.value)}
                    onSearch={handleSearch}
                />
                <Button type="primary" icon={<SearchOutlined />} onClick={handleSearch}>
                    搜索
                </Button>
                <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
                    新增集群
                </Button>
            </Space>

            <Table
                columns={columns}
                dataSource={data}
                rowKey="id"
                loading={loading}
                pagination={{
                    current: currentPage,
                    pageSize: pageSize,
                    total: total,
                    showSizeChanger: true,
                    showQuickJumper: true,
                    showTotal: (total) => `共 ${total} 条记录`,
                }}
                onChange={handleTableChange}
                scroll={{ x: 1300 }}
            />

            {/* 删除确认弹窗 */}
            <Modal
                title={
                    <div>
                        <ExclamationCircleOutlined style={{ color: '#ff4d4f', marginRight: 8 }} />
                        确认删除集群
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
                <p>请输入集群名称 <strong>{clusterToDelete?.name}</strong> 以确认删除：</p>
                <Form form={deleteForm}>
                    <Form.Item
                        name="confirmName"
                        rules={[
                            { required: true, message: '请输入集群名称' },
                            {
                                validator: (_, value) => {
                                    if (value && clusterToDelete && value !== clusterToDelete.name) {
                                        return Promise.reject(new Error('集群名称不匹配'));
                                    }
                                    return Promise.resolve();
                                }
                            }
                        ]}
                    >
                        <Input placeholder="请输入集群名称" />
                    </Form.Item>
                </Form>
            </Modal>
        </Card>
    )
}

export default ClusterList 