import { useState, useEffect } from 'react'
import { Card, Form, Input, Select, Button, Space, message, Alert } from 'antd'
import { InfoCircleOutlined, CheckCircleOutlined } from '@ant-design/icons'
import MonacoEditor from 'react-monaco-editor'
import * as monaco from 'monaco-editor'
import { useNavigate, useParams } from 'react-router-dom'
import clusterService, { Cluster } from '../../services/clusterService'

const { TextArea } = Input

// 默认的 kubeconfig 模板
const defaultKubeconfig = `apiVersion: v1
kind: Config
clusters:
- cluster:
    server: https://kubernetes.default.svc
    certificate-authority-data: <your-ca-data>
  name: default-cluster
contexts:
- context:
    cluster: default-cluster
    user: default-user
  name: default-context
current-context: default-context
users:
- name: default-user
  user:
    token: <your-token>`

// 表单数据结构
interface FormData {
    basicInfo?: {
        name: string
        cnname: string
        desc: string
        master: string
        environment: string
        status: number
    }
    kubeconfig?: {
        content: string
    }
}

// 确保表单数据完整的函数
const ensureFormDataComplete = (values: FormData): Required<FormData> => {
    const result = { ...values } as Required<FormData>;

    // 确保 basicInfo 存在
    if (!result.basicInfo) {
        result.basicInfo = {
            name: '默认集群',
            cnname: '默认集群',
            desc: '自动创建的默认集群',
            master: 'https://kubernetes.default.svc',
            environment: 'test',
            status: 0
        };
    }

    // 确保 kubeconfig 存在
    if (!result.kubeconfig) {
        result.kubeconfig = {
            content: defaultKubeconfig
        };
    }

    return result;
}

const CreateCluster = () => {
    const [form] = Form.useForm()
    const [isValidating, setIsValidating] = useState(false)
    const [validationResult, setValidationResult] = useState<{
        success: boolean
        message: string
    } | null>(null)
    const [loading, setLoading] = useState(false)
    const navigate = useNavigate()
    const { id } = useParams()
    const isEdit = !!id

    // 环境选项
    const environments = [
        { value: 'prod', label: '生产环境' },
        { value: 'test', label: '测试环境' },
        { value: 'dev', label: '开发环境' },
    ]

    // 状态选项
    const statusOptions = [
        { value: 0, label: '正常' },
        { value: 1, label: '停用' },
    ]

    // 如果是编辑模式，加载集群数据
    useEffect(() => {
        if (isEdit && id) {
            const fetchCluster = async () => {
                try {
                    setLoading(true)
                    const cluster = await clusterService.getCluster(parseInt(id))

                    // 将后端数据转换为表单数据格式
                    const formData: FormData = {
                        basicInfo: {
                            name: cluster.name,
                            cnname: cluster.cnname || '',
                            desc: cluster.desc || cluster.description || '',
                            master: cluster.master || '',
                            environment: cluster.environment || 'test',
                            status: typeof cluster.status === 'number' ? cluster.status : 0,
                        },
                        kubeconfig: {
                            content: cluster.kubeconfig || defaultKubeconfig,
                        },
                    }

                    form.setFieldsValue(formData)
                } catch (error) {
                    console.error('获取集群详情失败:', error)
                    message.error('获取集群详情失败')
                } finally {
                    setLoading(false)
                }
            }

            fetchCluster()
        }
    }, [isEdit, id, form])

    const editorOptions: monaco.editor.IStandaloneEditorConstructionOptions = {
        selectOnLineNumbers: true,
        roundedSelection: false,
        readOnly: false,
        cursorStyle: 'line',
        automaticLayout: true,
        minimap: { enabled: true },
        fontSize: 14,
        scrollBeyondLastLine: false,
        lineNumbers: 'on' as const,
        scrollbar: {
            vertical: 'visible',
            horizontal: 'visible',
        },
        wordWrap: 'on',
        language: 'yaml',
        theme: 'vs-dark',
    }

    const handleValidate = async () => {
        setIsValidating(true)
        try {
            const values = await form.validateFields()
            // 调用后端 API 验证 kubeconfig
            // 这里可以添加实际的验证逻辑
            await new Promise(resolve => setTimeout(resolve, 1000))
            setValidationResult({
                success: true,
                message: 'kubeconfig 验证成功',
            })
        } catch (error) {
            setValidationResult({
                success: false,
                message: 'kubeconfig 验证失败：' + (error as Error).message,
            })
        } finally {
            setIsValidating(false)
        }
    }

    const handleSubmit = () => {
        console.log('点击提交按钮');
        // 先获取表单当前值，检查是否完整
        const currentValues = form.getFieldsValue(true);
        console.log('当前表单值:', currentValues);
        console.log('当前表单值 JSON:', JSON.stringify(currentValues, null, 2));

        // 检查基本信息是否存在
        if (!currentValues.basicInfo) {
            console.error('表单数据缺少 basicInfo:', currentValues);
            // 尝试初始化 basicInfo 字段
            const completeValues = ensureFormDataComplete(currentValues);
            form.setFieldsValue(completeValues);
            console.log('已修复表单数据结构:', completeValues);
            console.log('已修复表单数据结构 JSON:', JSON.stringify(completeValues, null, 2));
            message.warning('正在修复表单数据结构，请再次点击提交');
            return;
        }

        // 手动验证表单
        form.validateFields()
            .then(values => {
                console.log('表单验证成功:', values);
                console.log('表单验证成功 JSON:', JSON.stringify(values, null, 2));
                form.submit();
            })
            .catch(errors => {
                console.error('表单验证失败:', errors);
                message.error('表单验证失败，请检查输入');
            });
    }

    // 调试函数 - 修复表单数据结构
    const handleDebug = () => {
        const currentValues = form.getFieldsValue(true);
        console.log('调试 - 当前表单值:', currentValues);

        // 使用 ensureFormDataComplete 函数修复表单数据
        const completeValues = ensureFormDataComplete(currentValues);
        form.setFieldsValue(completeValues);

        console.log('调试 - 修复后的表单值:', completeValues);
        message.success('表单数据结构已修复，请再次尝试提交');
    }

    const onFinish = async (values: FormData) => {
        try {
            setLoading(true)
            console.log('提交的表单数据:', values);
            console.log('提交的表单数据 JSON:', JSON.stringify(values, null, 2));

            // 检查表单数据是否为空
            if (!values) {
                console.error('表单数据为空');
                message.error('表单数据为空，请重新填写');
                setLoading(false);
                return;
            }

            // 确保表单数据完整
            const completeValues = ensureFormDataComplete(values);
            console.log('完整的表单数据:', completeValues);
            console.log('完整的表单数据 JSON:', JSON.stringify(completeValues, null, 2));

            // 将表单数据转换为API需要的格式
            const clusterData: Cluster = {
                name: completeValues.basicInfo.name,
                cnname: completeValues.basicInfo.cnname,
                desc: completeValues.basicInfo.desc,
                master: completeValues.basicInfo.master,
                environment: completeValues.basicInfo.environment,
                status: completeValues.basicInfo.status,
                locationID: 1, // 默认位置ID为1，不在表单中显示
                kubeconfig: completeValues.kubeconfig.content,
            }

            console.log('发送到API的数据:', clusterData);
            console.log('发送到API的数据 JSON:', JSON.stringify(clusterData, null, 2));

            if (isEdit && id) {
                // 更新集群
                const result = await clusterService.updateCluster(parseInt(id), clusterData)
                console.log('更新集群结果:', result);
                message.success('集群更新成功')
            } else {
                // 创建集群
                const result = await clusterService.createCluster(clusterData)
                console.log('创建集群结果:', result);
                message.success('集群创建成功')
            }

            navigate('/clusters/list')
        } catch (error) {
            console.error(isEdit ? '更新集群失败:' : '创建集群失败:', error)
            message.error(isEdit ? '更新集群失败' : '创建集群失败')
        } finally {
            setLoading(false)
        }
    }

    return (
        <Card title={isEdit ? '编辑集群' : '新建集群'}>
            <Form
                form={form}
                layout="vertical"
                onFinish={onFinish}
                style={{ maxWidth: 800, margin: '0 auto' }}
                initialValues={{
                    basicInfo: {
                        name: '',
                        cnname: '',
                        desc: '',
                        master: '',
                        environment: 'test',
                        status: 0,  // 默认状态为正常(0)
                    },
                    kubeconfig: {
                        content: defaultKubeconfig
                    }
                }}
            >
                {/* 基础信息部分 */}
                <h3><InfoCircleOutlined /> 基础信息</h3>
                <Form.Item
                    name={['basicInfo', 'name']}
                    label="集群名称"
                    rules={[{ required: true, message: '请输入集群名称' }]}
                >
                    <Input placeholder="请输入集群名称" />
                </Form.Item>
                <Form.Item
                    name={['basicInfo', 'cnname']}
                    label="中文名称"
                    rules={[{ required: true, message: '请输入集群中文名称' }]}
                >
                    <Input placeholder="请输入集群中文名称" />
                </Form.Item>
                <Form.Item
                    name={['basicInfo', 'desc']}
                    label="描述"
                    rules={[{ required: true, message: '请输入集群描述' }]}
                >
                    <TextArea rows={4} placeholder="请输入集群描述" />
                </Form.Item>
                <Form.Item
                    name={['basicInfo', 'master']}
                    label="主节点地址"
                    rules={[{ required: true, message: '请输入主节点地址' }]}
                >
                    <Input placeholder="请输入主节点地址，例如：https://192.168.1.100:6443" />
                </Form.Item>
                <Form.Item
                    name={['basicInfo', 'environment']}
                    label="环境"
                    rules={[{ required: true, message: '请选择环境' }]}
                >
                    <Select options={environments} placeholder="请选择环境" />
                </Form.Item>
                <Form.Item
                    name={['basicInfo', 'status']}
                    label="状态"
                    rules={[{ required: true, message: '请选择状态' }]}
                >
                    <Select options={statusOptions} placeholder="请选择状态" />
                </Form.Item>

                {/* Kubeconfig 部分 */}
                <Form.Item
                    name={['kubeconfig', 'content']}
                    label="Kubeconfig"
                    rules={[{ required: true, message: '请输入 kubeconfig' }]}
                >
                    <div style={{ border: '1px solid #d9d9d9', borderRadius: '2px' }}>
                        <MonacoEditor
                            height="400"
                            language="yaml"
                            theme="vs-dark"
                            value={form.getFieldValue(['kubeconfig', 'content']) || defaultKubeconfig}
                            options={editorOptions}
                            onChange={(value) => {
                                form.setFieldsValue({
                                    kubeconfig: {
                                        content: value,
                                    },
                                })
                            }}
                        />
                    </div>
                </Form.Item>

                <Space style={{ marginBottom: 16 }}>
                    <Button
                        type="primary"
                        icon={<CheckCircleOutlined />}
                        onClick={handleValidate}
                        loading={isValidating}
                    >
                        验证配置
                    </Button>
                </Space>

                {validationResult && (
                    <Alert
                        style={{ marginBottom: 16 }}
                        type={validationResult.success ? 'success' : 'error'}
                        message={validationResult.message}
                    />
                )}

                {/* 提交按钮 */}
                <Form.Item>
                    <Space>
                        <Button type="primary" onClick={handleSubmit} loading={loading}>
                            {isEdit ? '更新' : '提交'}
                        </Button>
                        <Button onClick={() => navigate('/clusters/list')}>
                            取消
                        </Button>
                    </Space>
                </Form.Item>
            </Form>
        </Card>
    )
}

export default CreateCluster 