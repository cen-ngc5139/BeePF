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
} from 'antd'
import { SearchOutlined, MoreOutlined } from '@ant-design/icons'
import type { TableProps } from 'antd'

interface ComponentItem {
    id: string
    name: string
    cluster: string
    status: 'running' | 'stopped' | 'error'
}

const ComponentList = () => {
    const [selectedCluster, setSelectedCluster] = useState<string>()
    const [selectedRule, setSelectedRule] = useState<string>()
    const [searchKeyword, setSearchKeyword] = useState('')
    const [modalVisible, setModalVisible] = useState(false)
    const [selectedComponent, setSelectedComponent] = useState<ComponentItem>()

    // 模拟数据
    const clusters = [
        { value: 'cluster1', label: '集群1' },
        { value: 'cluster2', label: '集群2' },
    ]

    const rules = [
        { value: 'rule1', label: '规则1' },
        { value: 'rule2', label: '规则2' },
    ]

    const handleStop = (record: ComponentItem) => {
        setSelectedComponent(record)
        setModalVisible(true)
    }

    const handleDelete = (record: ComponentItem) => {
        Modal.confirm({
            title: '确认删除',
            content: `确定要删除组件 ${record.name} 吗？`,
            onOk: () => {
                message.success('删除成功')
            },
        })
    }

    const columns: TableProps<ComponentItem>['columns'] = [
        {
            title: '组件名称',
            dataIndex: 'name',
            key: 'name',
        },
        {
            title: '所属集群',
            dataIndex: 'cluster',
            key: 'cluster',
        },
        {
            title: '部署状态',
            dataIndex: 'status',
            key: 'status',
            render: (status: string) => {
                const statusMap = {
                    running: '运行中',
                    stopped: '已停止',
                    error: '异常',
                }
                return statusMap[status as keyof typeof statusMap]
            },
        },
        {
            title: '操作',
            key: 'action',
            render: (_, record) => (
                <Dropdown
                    menu={{
                        items: [
                            {
                                key: 'stop',
                                label: '停止',
                                onClick: () => handleStop(record),
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

    const data: ComponentItem[] = [
        {
            id: '1',
            name: '组件1',
            cluster: '集群1',
            status: 'running',
        },
        {
            id: '2',
            name: '组件2',
            cluster: '集群2',
            status: 'stopped',
        },
    ]

    return (
        <Card>
            <Space style={{ marginBottom: 16 }}>
                <Select
                    placeholder="选择集群"
                    style={{ width: 200 }}
                    options={clusters}
                    value={selectedCluster}
                    onChange={setSelectedCluster}
                    allowClear
                />
                <Select
                    placeholder="选择规则"
                    style={{ width: 200 }}
                    options={rules}
                    value={selectedRule}
                    onChange={setSelectedRule}
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
            </Space>

            <Table columns={columns} dataSource={data} rowKey="id" />

            <Modal
                title="确认停止"
                open={modalVisible}
                onOk={() => {
                    message.success('停止成功')
                    setModalVisible(false)
                }}
                onCancel={() => setModalVisible(false)}
            >
                <p>确定要停止组件 {selectedComponent?.name} 吗？</p>
            </Modal>
        </Card>
    )
}

export default ComponentList 