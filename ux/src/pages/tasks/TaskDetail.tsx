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
    Statistic,
    Tabs
} from 'antd';
import { Line } from '@ant-design/plots';
import { ArrowLeftOutlined, StopOutlined, ExclamationCircleOutlined } from '@ant-design/icons';
import taskService, { Task, ProgStatus } from '../../services/taskService';
import type { TableProps } from 'antd';

const { Title, Text } = Typography;
const { confirm } = Modal;

// 任务步骤枚举
enum TaskStep {
    Init = 0,
    Load = 1,
    Start = 2,
    Stats = 3,
    Metrics = 4,
    Stop = 5
}

// 任务步骤描述
const TaskStepDescriptions: Record<number, string> = {
    [TaskStep.Init]: '初始化',
    [TaskStep.Load]: '加载',
    [TaskStep.Start]: '启动',
    [TaskStep.Stats]: '统计',
    [TaskStep.Metrics]: '指标',
    [TaskStep.Stop]: '停止'
};

// 获取任务步骤描述
const getTaskStepDescription = (step: number): string => {
    return TaskStepDescriptions[step] || `步骤${step}`;
};

// 指标数据接口
interface MetricPoint {
    timestamp: string;
    value: number;
}

interface TaskMetrics {
    avg_run_time_ns: MetricPoint[];
    cpu_usage: MetricPoint[];
    events_per_second: MetricPoint[];
    period_ns: MetricPoint[];
    total_avg_run_time_ns: MetricPoint[];
}

const TaskDetail = () => {
    const { taskId } = useParams<{ taskId: string }>();
    const navigate = useNavigate();
    const [loading, setLoading] = useState(true);
    const [task, setTask] = useState<Task | null>(null);
    const [stopLoading, setStopLoading] = useState(false);
    const [stopModalVisible, setStopModalVisible] = useState(false);
    const [metrics, setMetrics] = useState<TaskMetrics | null>(null);
    const [metricsLoading, setMetricsLoading] = useState(false);

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

    // 获取指标数据
    const fetchMetrics = async () => {
        if (!taskId) return;

        try {
            setMetricsLoading(true);
            const response = await taskService.getTaskMetrics(parseInt(taskId));

            // 确保所有指标数据都存在
            const safeMetrics = {
                avg_run_time_ns: response?.avg_run_time_ns || [],
                cpu_usage: response?.cpu_usage || [],
                events_per_second: response?.events_per_second || [],
                period_ns: response?.period_ns || [],
                total_avg_run_time_ns: response?.total_avg_run_time_ns || []
            };

            setMetrics(safeMetrics);
        } catch (error) {
            console.error('获取任务指标失败:', error);
            message.error('获取任务指标失败');
        } finally {
            setMetricsLoading(false);
        }
    };

    // 初始化加载
    useEffect(() => {
        fetchTaskDetail();
        fetchMetrics();
    }, [taskId]);

    // 打开停止确认弹窗
    const handleStop = () => {
        setStopModalVisible(true);
    };

    // 确认停止任务
    const confirmStop = async () => {
        if (!task) return;

        try {
            setStopLoading(true);
            await taskService.stopTask(task.id);
            message.success(`任务 ${task.name} 已停止`);
            setStopModalVisible(false);
            // 重新获取任务详情
            fetchTaskDetail();
        } catch (error: any) {
            console.error('停止任务失败:', error);
            message.error(`停止失败: ${error.message || '未知错误'}`);
        } finally {
            setStopLoading(false);
        }
    };

    // 取消停止任务
    const cancelStop = () => {
        setStopModalVisible(false);
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

    // 渲染指标图表
    const renderMetricsChart = (data: MetricPoint[], title: string, yAxisLabel: string) => {
        // 过滤掉值为null的数据点并转换数据格式
        const validData = data
            .filter(point => point.value !== null && point.value !== undefined)
            .map(point => {
                // 确保时间戳是有效的日期对象
                const timestamp = new Date(point.timestamp);
                return {
                    date: timestamp.getTime(), // 使用时间戳数值而不是字符串
                    value: point.value,
                    category: title
                };
            })
            .filter(item => !isNaN(item.date)); // 过滤掉无效日期

        // 获取最新值和平均值
        const latestValue = validData.length > 0 ? validData[validData.length - 1].value : 0;
        const avgValue = validData.length > 0
            ? validData.reduce((sum, item) => sum + item.value, 0) / validData.length
            : 0;

        return (
            <Card
                title={title}
                style={{ marginBottom: 16 }}
                extra={
                    <Space>
                        <Statistic
                            title="当前值"
                            value={latestValue}
                            precision={yAxisLabel === '百分比' ? 4 : 0}
                            style={{ marginRight: 16 }}
                            valueStyle={{ fontSize: '14px' }}
                        />
                        <Statistic
                            title="平均值"
                            value={avgValue}
                            precision={yAxisLabel === '百分比' ? 4 : 0}
                            valueStyle={{ fontSize: '14px' }}
                        />
                    </Space>
                }
            >
                <div style={{ height: 200 }}>
                    {validData.length > 0 ? (
                        <Line
                            data={validData}
                            xField="date"
                            yField="value"
                            seriesField="category"
                            xAxis={{
                                type: 'time',
                                tickCount: 5,
                                label: {
                                    formatter: (v: string | number) => {
                                        const date = new Date(typeof v === 'string' ? parseInt(v) : v);
                                        return `${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`;
                                    }
                                }
                            }}
                            yAxis={{
                                label: {
                                    formatter: (v: string) => {
                                        if (yAxisLabel === '百分比') {
                                            return `${(parseFloat(v) * 100).toFixed(2)}%`;
                                        }
                                        return v;
                                    }
                                }
                            }}
                            tooltip={{
                                formatter: (datum: { date: number, value: number, category: string }) => {
                                    let displayValue: string;
                                    if (yAxisLabel === '百分比') {
                                        displayValue = `${(datum.value * 100).toFixed(2)}%`;
                                    } else if (yAxisLabel === '纳秒') {
                                        displayValue = `${datum.value.toLocaleString()} ns`;
                                    } else {
                                        displayValue = datum.value.toLocaleString();
                                    }

                                    const date = new Date(datum.date);
                                    const timeStr = `${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`;

                                    return {
                                        name: title,
                                        value: displayValue,
                                        title: timeStr
                                    };
                                }
                            }}
                            smooth={true}
                            animation={false}
                            lineStyle={{
                                lineWidth: 2,
                            }}
                            point={{
                                size: 3,
                                shape: 'circle',
                                style: {
                                    fill: '#5B8FF9',
                                    stroke: '#fff',
                                    lineWidth: 1,
                                },
                            }}
                            color="#5B8FF9"
                            areaStyle={{
                                fill: 'l(270) 0:#ffffff 0.5:#5B8FF9 1:#5B8FF9',
                                fillOpacity: 0.2,
                            }}
                        />
                    ) : (
                        <div style={{ height: '100%', display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
                            <Text type="secondary">暂无数据</Text>
                        </div>
                    )}
                </div>
            </Card>
        );
    };

    // 渲染指标区域
    const renderMetrics = () => {
        if (!metrics) return null;

        return (
            <>
                <Divider orientation="left">性能指标</Divider>
                <Spin spinning={metricsLoading}>
                    <Row gutter={[16, 16]}>
                        <Col span={12}>
                            {renderMetricsChart(
                                metrics.avg_run_time_ns,
                                '平均运行时间',
                                '纳秒'
                            )}
                        </Col>
                        <Col span={12}>
                            {renderMetricsChart(
                                metrics.cpu_usage,
                                'CPU 使用率',
                                '百分比'
                            )}
                        </Col>
                        <Col span={12}>
                            {renderMetricsChart(
                                metrics.events_per_second,
                                '每秒事件数',
                                '次数'
                            )}
                        </Col>
                        <Col span={12}>
                            {renderMetricsChart(
                                metrics.period_ns,
                                '周期',
                                '纳秒'
                            )}
                        </Col>
                        <Col span={24}>
                            {renderMetricsChart(
                                metrics.total_avg_run_time_ns,
                                '总平均运行时间',
                                '纳秒'
                            )}
                        </Col>
                    </Row>
                </Spin>
            </>
        );
    };

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
                                <Statistic
                                    title="步骤"
                                    value={task.step}
                                    formatter={(value) => getTaskStepDescription(value as number)}
                                />
                            </Col>
                        </Row>

                        <Tabs
                            defaultActiveKey="basic"
                            items={[
                                {
                                    key: 'basic',
                                    label: '基本信息',
                                    children: (
                                        <>
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
                                    ),
                                },
                                {
                                    key: 'metrics',
                                    label: '性能指标',
                                    children: renderMetrics(),
                                },
                            ]}
                        />
                    </>
                ) : (
                    <div style={{ textAlign: 'center', padding: '50px 0' }}>
                        <Text type="secondary">未找到任务信息</Text>
                    </div>
                )}
            </Card>

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
                <p>确定要停止任务 "{task?.name}" 吗？停止后任务将无法继续执行。</p>
            </Modal>
        </Spin>
    );
};

export default TaskDetail; 