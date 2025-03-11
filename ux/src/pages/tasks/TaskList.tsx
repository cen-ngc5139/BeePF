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
import { SearchOutlined, MoreOutlined, StopOutlined, ExclamationCircleOutlined } from '@ant-design/icons'
import type { TableProps } from 'antd'
import { useNavigate } from 'react-router-dom'
import taskService, { Task } from '../../services/taskService'
import componentService from '../../services/componentService'

const { confirm } = Modal;

const TaskList = () => {
    const navigate = useNavigate();
    const [loading, setLoading] = useState(false);
    const [selectedComponent, setSelectedComponent] = useState<number>();
    const [searchKeyword, setSearchKeyword] = useState('');
    const [tasks, setTasks] = useState<Task[]>([]);
    const [total, setTotal] = useState(0);
    const [components, setComponents] = useState<any[]>([]);
    const [currentPage, setCurrentPage] = useState(1);
    const [pageSize, setPageSize] = useState(10);
    const [componentsLoaded, setComponentsLoaded] = useState(false);

    // 停止任务相关状态
    const [stopModalVisible, setStopModalVisible] = useState(false);
    const [taskToStop, setTaskToStop] = useState<Task | null>(null);
    const [stopLoading, setStopLoading] = useState(false);
    const [stopForm] = Form.useForm();

    // 加载组件列表
    const loadComponents = useCallback(async () => {
        try {
            setLoading(true);
            const response = await componentService.getComponentList({});
            setComponents(response.components || []);
            setComponentsLoaded(true);
        } catch (error) {
            console.error('加载组件列表失败:', error);
            message.error('加载组件列表失败');
        } finally {
            setLoading(false);
        }
    }, []);

    // 加载任务列表
    const loadData = useCallback(async (page = currentPage, size = pageSize) => {
        if (!componentsLoaded) {
            return; // 等待组件数据加载完成
        }

        setLoading(true);
        try {
            const params: any = {
                pageNum: page,
                pageSize: size,
            };

            if (selectedComponent) {
                params.component_id = selectedComponent;
            }

            if (searchKeyword) {
                params.keyword = searchKeyword;
            }

            const response = await taskService.getTasks(params);
            setTasks(response.list || []);
            setTotal(response.total || 0);
        } catch (error) {
            console.error('加载任务列表失败:', error);
            message.error('加载任务列表失败');
        } finally {
            setLoading(false);
        }
    }, [currentPage, pageSize, selectedComponent, searchKeyword, componentsLoaded]);

    // 初始化加载
    useEffect(() => {
        loadComponents();
    }, [loadComponents]);

    // 当组件列表加载完成后，加载任务列表
    useEffect(() => {
        if (componentsLoaded) {
            loadData();
        }
    }, [loadData, componentsLoaded]);

    // 搜索
    const handleSearch = () => {
        setCurrentPage(1);
        loadData(1);
    };

    // 处理停止任务
    const handleStop = (record: Task) => {
        console.log('点击停止按钮，任务ID:', record.id, '任务名称:', record.name);

        // 设置要停止的任务并显示确认对话框
        setTaskToStop(record);
        setStopModalVisible(true);
        stopForm.resetFields();
    };

    // 确认停止任务
    const confirmStop = async () => {
        try {
            await stopForm.validateFields();

            if (!taskToStop) {
                message.error('未选择要停止的任务');
                return;
            }

            setStopLoading(true);
            try {
                await taskService.stopTask(taskToStop.id);
                message.success(`任务 ${taskToStop.name} 已停止`);
                setStopModalVisible(false);

                // 重新加载数据，确保刷新任务列表
                setTimeout(() => {
                    loadData(currentPage, pageSize);
                    console.log('刷新任务列表，当前页码:', currentPage, '每页条数:', pageSize);
                }, 500);
            } catch (error: any) {
                console.error('停止任务失败:', error);
                message.error(`停止失败: ${error.message || '未知错误'}`);
            } finally {
                setStopLoading(false);
                setTaskToStop(null);
            }
        } catch (error) {
            // 表单验证失败
            console.log('表单验证失败:', error);
        }
    };

    // 取消停止任务
    const cancelStop = () => {
        setStopModalVisible(false);
        setTaskToStop(null);
    };

    // 获取任务状态标签
    const getStatusTag = (status: number) => {
        if (status === undefined || status === null) {
            return <Tag>未知</Tag>;
        }

        switch (status) {
            case 0:
                return <Tag color="blue">初始化</Tag>;
            case 1:
                return <Tag color="green">运行中</Tag>;
            case 2:
                return <Tag color="blue">已完成</Tag>;
            case 3:
                return <Tag color="red">失败</Tag>;
            case 4:
                return <Tag color="orange">已停止</Tag>;
            default:
                return <Tag>{`状态${status}`}</Tag>;
        }
    };

    const columns: TableProps<Task>['columns'] = [
        {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            width: 80,
        },
        {
            title: '任务名称',
            dataIndex: 'name',
            key: 'name',
            render: (text, record) => (
                <a onClick={() => navigate(`/task/${record.id}`)}>{text}</a>
            ),
        },
        {
            title: '组件名称',
            dataIndex: 'component_name',
            key: 'component_name',
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            render: (status) => getStatusTag(status || 0),
        },
        {
            title: '错误信息',
            dataIndex: 'error',
            key: 'error',
            ellipsis: true,
            render: (text) => text || '-',
        },
        {
            title: '创建时间',
            dataIndex: 'created_at',
            key: 'created_at',
            render: (text) => new Date(text).toLocaleString(),
        },
        {
            title: '最后更新时间',
            dataIndex: 'updated_at',
            key: 'updated_at',
            render: (text) => new Date(text).toLocaleString(),
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
                                onClick: () => navigate(`/task/${record.id}`),
                            },
                            {
                                key: 'stop',
                                label: '停止任务',
                                icon: <StopOutlined />,
                                onClick: () => handleStop(record),
                                disabled: record.status !== 1,
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
        <Card title="任务列表">
            <Space style={{ marginBottom: 16 }}>
                <Select
                    placeholder="选择组件"
                    style={{ width: 200 }}
                    options={components.map(component => ({ value: component.id, label: component.name }))}
                    value={selectedComponent}
                    onChange={setSelectedComponent}
                    allowClear
                />
                <Input
                    placeholder="请输入任务名称"
                    style={{ width: 200 }}
                    value={searchKeyword}
                    onChange={(e) => setSearchKeyword(e.target.value)}
                    onPressEnter={handleSearch}
                />
                <Button type="primary" icon={<SearchOutlined />} onClick={handleSearch}>
                    搜索
                </Button>
            </Space>

            <Spin spinning={loading}>
                <Table
                    columns={columns}
                    dataSource={tasks}
                    rowKey="id"
                    pagination={{
                        current: currentPage,
                        pageSize: pageSize,
                        total: total,
                        showSizeChanger: true,
                        showQuickJumper: true,
                        showTotal: (total) => `共 ${total} 条记录`,
                        onChange: (page, pageSize) => {
                            setCurrentPage(page);
                            setPageSize(pageSize);
                            loadData(page, pageSize);
                        },
                    }}
                />
            </Spin>

            {/* 停止任务确认弹窗 */}
            <Modal
                title={
                    <div>
                        <ExclamationCircleOutlined style={{ color: '#ff4d4f', marginRight: 8 }} />
                        确认停止任务
                    </div>
                }
                open={stopModalVisible}
                onOk={confirmStop}
                onCancel={cancelStop}
                confirmLoading={stopLoading}
                okText="停止"
                cancelText="取消"
                okButtonProps={{ danger: true }}
            >
                <p>停止任务后，任务将无法继续执行，请谨慎操作！</p>
            </Modal>
        </Card>
    );
};

export default TaskList; 