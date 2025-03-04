import { useState, useEffect } from 'react'
import { Card, Steps, Form, Input, Select, Button, Space, message, Alert } from 'antd'
import { InfoCircleOutlined, SettingOutlined, CheckCircleOutlined } from '@ant-design/icons'
import MonacoEditor from 'react-monaco-editor'
import * as monaco from 'monaco-editor'
import { useNavigate, useParams } from 'react-router-dom'

const { TextArea } = Input
const { Step } = Steps

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

interface FormData {
    basicInfo: {
        name: string
        description: string
        region: string
    }
    kubeconfig: {
        content: string
    }
}

const CreateCluster = () => {
    const [currentStep, setCurrentStep] = useState(0)
    const [form] = Form.useForm()
    const [isValidating, setIsValidating] = useState(false)
    const [validationResult, setValidationResult] = useState<{
        success: boolean
        message: string
    } | null>(null)
    const navigate = useNavigate()
    const { id } = useParams()
    const isEdit = !!id

    // 模拟数据
    const regions = [
        { value: 'cn-north-1', label: '华北1（北京）' },
        { value: 'cn-south-1', label: '华南1（广州）' },
        { value: 'cn-east-1', label: '华东1（上海）' },
    ]

    // 如果是编辑模式，加载集群数据
    useEffect(() => {
        if (isEdit) {
            // TODO: 调用后端 API 获取集群数据
            const mockData = {
                basicInfo: {
                    name: '生产集群',
                    description: '用于生产环境的集群',
                    region: 'cn-north-1',
                },
                kubeconfig: {
                    content: defaultKubeconfig,
                },
            }
            form.setFieldsValue(mockData)
        }
    }, [isEdit, form])

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
            // TODO: 调用后端 API 验证 kubeconfig
            // 模拟验证过程
            await new Promise(resolve => setTimeout(resolve, 2000))
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

    const steps = [
        {
            title: '基础信息',
            icon: <InfoCircleOutlined />,
            content: (
                <Form.Item noStyle>
                    <Form.Item
                        name={['basicInfo', 'name']}
                        label="集群名称"
                        rules={[{ required: true, message: '请输入集群名称' }]}
                    >
                        <Input placeholder="请输入集群名称" />
                    </Form.Item>
                    <Form.Item
                        name={['basicInfo', 'description']}
                        label="描述"
                        rules={[{ required: true, message: '请输入集群描述' }]}
                    >
                        <TextArea rows={4} placeholder="请输入集群描述" />
                    </Form.Item>
                    <Form.Item
                        name={['basicInfo', 'region']}
                        label="区域"
                        rules={[{ required: true, message: '请选择区域' }]}
                    >
                        <Select options={regions} placeholder="请选择区域" />
                    </Form.Item>
                </Form.Item>
            ),
        },
        {
            title: 'Kubeconfig 配置',
            icon: <SettingOutlined />,
            content: (
                <Form.Item noStyle>
                    <Form.Item
                        name={['kubeconfig', 'content']}
                        label="Kubeconfig 内容"
                        rules={[{ required: true, message: '请输入 kubeconfig 内容' }]}
                    >
                        <div style={{ border: '1px solid #d9d9d9', borderRadius: '2px' }}>
                            <MonacoEditor
                                height="500"
                                language="yaml"
                                theme="vs-dark"
                                value={defaultKubeconfig}
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
                    {validationResult && (
                        <Alert
                            style={{ marginTop: 16 }}
                            type={validationResult.success ? 'success' : 'error'}
                            message={validationResult.message}
                        />
                    )}
                </Form.Item>
            ),
        },
    ]

    const next = async () => {
        try {
            await form.validateFields()
            setCurrentStep(currentStep + 1)
        } catch (error) {
            console.error('Validation failed:', error)
        }
    }

    const prev = () => {
        setCurrentStep(currentStep - 1)
    }

    const onFinish = async (values: FormData) => {
        try {
            // TODO: 调用后端 API 创建/更新集群
            console.log('Form values:', values)
            message.success(isEdit ? '集群更新成功' : '集群创建成功')
            navigate('/clusters/list')
        } catch (error) {
            message.error(isEdit ? '集群更新失败' : '集群创建失败')
        }
    }

    return (
        <Card>
            <Steps current={currentStep} items={steps} style={{ marginBottom: 24 }} />
            <Form
                form={form}
                layout="vertical"
                onFinish={onFinish}
                style={{ maxWidth: 800, margin: '0 auto' }}
            >
                {steps[currentStep].content}
                <Form.Item>
                    <Space>
                        {currentStep > 0 && <Button onClick={prev}>上一步</Button>}
                        {currentStep === 1 && (
                            <Button
                                type="primary"
                                icon={<CheckCircleOutlined />}
                                onClick={handleValidate}
                                loading={isValidating}
                            >
                                验证配置
                            </Button>
                        )}
                        {currentStep < steps.length - 1 && (
                            <Button type="primary" onClick={next}>
                                下一步
                            </Button>
                        )}
                        {currentStep === steps.length - 1 && (
                            <Button type="primary" onClick={() => form.submit()}>
                                {isEdit ? '更新' : '提交'}
                            </Button>
                        )}
                    </Space>
                </Form.Item>
            </Form>
        </Card>
    )
}

export default CreateCluster 