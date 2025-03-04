import { useState } from 'react'
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
} from 'antd'
import { SearchOutlined, PlusOutlined, MoreOutlined } from '@ant-design/icons'
import type { TableProps } from 'antd'
import { useNavigate } from 'react-router-dom'

interface ClusterItem {
    id: string
    name: string
    description: string
    region: string
    status: 'active' | 'inactive'
}

const ClusterList = () => {
    const [selectedRegion, setSelectedRegion] = useState<string>()
    const [searchKeyword, setSearchKeyword] = useState('')
    const navigate = useNavigate()

    // 模拟数据
    const regions = [
        { value: 'cn-north-1', label: '华北1（北京）' },
        { value: 'cn-south-1', label: '华南1（广州）' },
        { value: 'cn-east-1', label: '华东1（上海）' },
    ]

    const handleEdit = (record: ClusterItem) => {
        navigate(`/clusters/edit/${record.id}`)
    }

    const handleDelete = (record: ClusterItem) => {
        Modal.confirm({
            title: '确认删除',
            content: `确定要删除集群 ${record.name} 吗？`,
            onOk: () => {
                message.success('删除成功')
            },
        })
    }

    const handleAdd = () => {
        navigate('/clusters/create')
    }

    const columns: TableProps<ClusterItem>['columns'] = [
        {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            width: 100,
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
        {
            title: '区域',
            dataIndex: 'region',
            key: 'region',
            render: (region: string) => {
                const regionMap = {
                    'cn-north-1': '华北1（北京）',
                    'cn-south-1': '华南1（广州）',
                    'cn-east-1': '华东1（上海）',
                }
                return <Tag color="blue">{regionMap[region as keyof typeof regionMap]}</Tag>
            },
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            render: (status: string) => {
                const statusMap = {
                    active: { text: '运行中', color: 'success' },
                    inactive: { text: '已停止', color: 'default' },
                }
                const { text, color } = statusMap[status as keyof typeof statusMap]
                return <Tag color={color}>{text}</Tag>
            },
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

    const data: ClusterItem[] = [
        {
            id: '1',
            name: '生产集群',
            description: '用于生产环境的集群',
            region: 'cn-north-1',
            status: 'active',
        },
        {
            id: '2',
            name: '测试集群',
            description: '用于测试环境的集群',
            region: 'cn-south-1',
            status: 'inactive',
        },
    ]

    return (
        <Card>
            <Space style={{ marginBottom: 16 }}>
                <Select
                    placeholder="选择区域"
                    style={{ width: 200 }}
                    options={regions}
                    value={selectedRegion}
                    onChange={setSelectedRegion}
                    allowClear
                />
                <Input.Search
                    placeholder="请输入搜索关键词"
                    style={{ width: 200 }}
                    value={searchKeyword}
                    onChange={(e) => setSearchKeyword(e.target.value)}
                    onSearch={() => {
                        // 实现搜索逻辑
                    }}
                />
                <Button type="primary" icon={<SearchOutlined />}>
                    搜索
                </Button>
                <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
                    新增集群
                </Button>
            </Space>

            <Table columns={columns} dataSource={data} rowKey="id" />
        </Card>
    )
}

export default ClusterList 