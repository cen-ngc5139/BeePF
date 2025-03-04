import { useState, useEffect } from 'react'
import { Card, Steps, Form, Input, Select, Button, Space, message } from 'antd'
import { CodeOutlined, SettingOutlined, DatabaseOutlined, InfoCircleOutlined } from '@ant-design/icons'
import MonacoEditor from 'react-monaco-editor'
import * as monaco from 'monaco-editor'

const { TextArea } = Input
const { Step } = Steps

// 添加 C 语言示例代码
const defaultCode = `#include <linux/bpf.h>
#include <bpf/libbpf.h>

SEC("cgroup/dev")
int bpf_prog(struct bpf_cgroup_dev_ctx *ctx) {
    return 1;
}

char LICENSE[] SEC("license") = "GPL";`

interface FormData {
    basicInfo: {
        name: string
        cluster: string
        description: string
    }
    codeInfo: {
        ebpfCode: string
    }
    progConfig: {
        cgroupPath: string
        pinPath: string
    }
    mapConfig: {
        mapPin: string
        exportType: string
    }
}

const CreateComponent = () => {
    const [currentStep, setCurrentStep] = useState(0)
    const [form] = Form.useForm()

    // 初始化表单数据
    useEffect(() => {
        form.setFieldsValue({
            codeInfo: {
                ebpfCode: defaultCode,
            },
        })
    }, [form])

    // 模拟数据
    const clusters = [
        { value: 'cluster1', label: '集群1' },
        { value: 'cluster2', label: '集群2' },
    ]

    const exportTypes = [
        { value: 'kafka', label: 'Kafka' },
        { value: 'prometheus', label: 'Prometheus' },
        { value: 'elasticsearch', label: 'Elasticsearch' },
    ]

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
    }

    const steps = [
        {
            title: '基础信息',
            icon: <InfoCircleOutlined />,
            content: (
                <Form.Item noStyle>
                    <Form.Item
                        name={['basicInfo', 'name']}
                        label="组件名称"
                        rules={[{ required: true, message: '请输入组件名称' }]}
                    >
                        <Input placeholder="请输入组件名称" />
                    </Form.Item>
                    <Form.Item
                        name={['basicInfo', 'cluster']}
                        label="所属集群"
                        rules={[{ required: true, message: '请选择所属集群' }]}
                    >
                        <Select options={clusters} placeholder="请选择所属集群" />
                    </Form.Item>
                    <Form.Item name={['basicInfo', 'description']} label="描述">
                        <TextArea rows={4} placeholder="请输入组件描述" />
                    </Form.Item>
                </Form.Item>
            ),
        },
        {
            title: '代码信息',
            icon: <CodeOutlined />,
            content: (
                <Form.Item
                    name={['codeInfo', 'ebpfCode']}
                    label="eBPF 代码"
                    rules={[{ required: true, message: '请输入 eBPF 代码' }]}
                >
                    <div style={{ border: '1px solid #d9d9d9', borderRadius: '2px' }}>
                        <MonacoEditor
                            height="500"
                            language="c"
                            theme="vs-dark"
                            value={defaultCode}
                            options={editorOptions}
                            onChange={(value) => {
                                form.setFieldsValue({
                                    codeInfo: {
                                        ebpfCode: value,
                                    },
                                })
                            }}
                        />
                    </div>
                </Form.Item>
            ),
        },
        {
            title: 'eBPF Prog 配置',
            icon: <SettingOutlined />,
            content: (
                <Form.Item noStyle>
                    <Form.Item
                        name={['progConfig', 'cgroupPath']}
                        label="CGroup Path"
                        rules={[{ required: true, message: '请输入 CGroup Path' }]}
                    >
                        <Input placeholder="请输入 CGroup Path" />
                    </Form.Item>
                    <Form.Item
                        name={['progConfig', 'pinPath']}
                        label="Pin Path"
                        rules={[{ required: true, message: '请输入 Pin Path' }]}
                    >
                        <Input placeholder="请输入 Pin Path" />
                    </Form.Item>
                </Form.Item>
            ),
        },
        {
            title: 'eBPF Map 配置',
            icon: <DatabaseOutlined />,
            content: (
                <Form.Item noStyle>
                    <Form.Item
                        name={['mapConfig', 'mapPin']}
                        label="Map Pin"
                        rules={[{ required: true, message: '请输入 Map Pin' }]}
                    >
                        <Input placeholder="请输入 Map Pin" />
                    </Form.Item>
                    <Form.Item
                        name={['mapConfig', 'exportType']}
                        label="导出类型"
                        rules={[{ required: true, message: '请选择导出类型' }]}
                    >
                        <Select options={exportTypes} placeholder="请选择导出类型" />
                    </Form.Item>
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
        console.log('Form values:', values)
        message.success('组件创建成功')
        // TODO: 处理表单提交
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
                        {currentStep < steps.length - 1 && (
                            <Button type="primary" onClick={next}>
                                下一步
                            </Button>
                        )}
                        {currentStep === steps.length - 1 && (
                            <Button type="primary" onClick={() => form.submit()}>
                                提交
                            </Button>
                        )}
                    </Space>
                </Form.Item>
            </Form>
        </Card>
    )
}

export default CreateComponent 