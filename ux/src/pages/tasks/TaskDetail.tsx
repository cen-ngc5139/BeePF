import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
    Card,
    Descriptions,
    Table,
    Button,
    Space,
    Tag,
    Spin,
    message,
    Modal,
    Divider,
    Typography,
    Row,
    Col,
    Statistic
} from 'antd';
import { ArrowLeftOutlined, StopOutlined, ExclamationCircleOutlined } from '@ant-design/icons';
import taskService, { Task, ProgStatus } from '../../services/taskService';
import type { TableProps } from 'antd';

const { Title, Text } = Typography;
const { confirm } = Modal;

const TaskDetail = () => {
    const { taskId } = useParams<{ taskId: string }>();
    const navigate = useNavigate();
    const [loading, setLoading] = useState(true);
    const [task, setTask] = useState<Task | null>(null);
    const [stopLoading, setStopLoading] = useState(false);

    // 获取任务详情
    const fetchTaskDetail = async () => {
        if (!taskId) return;

        try {
            setLoading(true);
            const taskData = await taskService.getTask(parseInt(taskId));
            setTask(taskData);
            console.log('任务详情:', taskData);
        } catch (error) {
            console.error('获取任务详情失败:', error);
            message.error('获取任务详情失败');
        } finally {
            setLoading(false);
        }
    };

    // 初始化加载
    useEffect(() => {
        fetchTaskDetail();
    }, [taskId]);

    // 停止任务
    const handleStop = () => {
        if (!task) return;

        confirm({
            title: '确认停止任务',
            icon: <ExclamationCircleOutlined />,
            content: `确定要停止任务 "${task.name}" 吗？停止后任务将无法继续执行。`,
            okText: '停止',
            okButtonProps: { danger: true },
            cancelText: '取消',
            onOk: async () => {
                try {
                    setStopLoading(true);
                    await taskService.stopTask(task.id);
                    message.success(`任务 ${task.name} 已停止`);
                    // 重新获取任务详情
                    fetchTaskDetail();
                } catch (error: any) {
                    console.error('停止任务失败:', error);
                    message.error(`停止失败: ${error.message || '未知错误'}`);
                } finally {
                    setStopLoading(false);
                }
            },
        });
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

    // 程序状态列表列定义
    const columns: TableProps<ProgStatus>['columns'] = [
        {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            width: 80,
        },
        {
            title: '程序名称',
            dataIndex: 'program_name',
            key: 'program_name',
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            render: (status) => getStatusTag(status),
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
            title: '更新时间',
            dataIndex: 'updated_at',
            key: 'updated_at',
            render: (text) => new Date(text).toLocaleString(),
        },
    ];

    return (
        <Spin spinning={loading}>
            <Card
                title={
                    <Space>
                        <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/tasks/list')}>
                            返回
                        </Button>
                        <span>任务详情</span>
                    </Space>
                }
                extra={
                    task && task.status === 1 ? (
                        <Button
                            type="primary"
                            danger
                            icon={<StopOutlined />}
                            onClick={handleStop}
                            loading={stopLoading}
                        >
                            停止任务
                        </Button>
                    ) : null
                }
            >
                {task ? (
                    <>
                        <Row gutter={16} style={{ marginBottom: 24 }}>
                            <Col span={6}>
                                <Statistic title="任务名称" value={task.name} />
                            </Col>
                            <Col span={6}>
                                <Statistic
                                    title="状态"
                                    value={task.status}
                                    formatter={(value) => getStatusTag(value as number)}
                                />
                            </Col>
                            <Col span={6}>
                                <Statistic title="组件" value={task.component_name} />
                            </Col>
                            <Col span={6}>
                                <Statistic title="步骤" value={task.step} />
                            </Col>
                        </Row>

                        <Descriptions bordered column={2} style={{ marginBottom: 24 }}>
                            <Descriptions.Item label="任务ID">{task.id}</Descriptions.Item>
                            <Descriptions.Item label="描述">{task.description || '-'}</Descriptions.Item>
                            <Descriptions.Item label="组件ID">{task.component_id}</Descriptions.Item>
                            <Descriptions.Item label="错误信息">{task.error || '-'}</Descriptions.Item>
                            <Descriptions.Item label="创建时间">{new Date(task.created_at).toLocaleString()}</Descriptions.Item>
                            <Descriptions.Item label="更新时间">{new Date(task.updated_at).toLocaleString()}</Descriptions.Item>
                        </Descriptions>

                        <Divider orientation="left">程序状态列表</Divider>
                        <Table
                            columns={columns}
                            dataSource={task.prog_status || []}
                            rowKey="id"
                            pagination={false}
                        />
                    </>
                ) : (
                    <div style={{ textAlign: 'center', padding: '50px 0' }}>
                        <Text type="secondary">未找到任务信息</Text>
                    </div>
                )}
            </Card>
        </Spin>
    );
};

export default TaskDetail; 