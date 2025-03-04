import { Form, Input, Select, Button, Space, message } from 'antd'
import { useEffect } from 'react'

interface ClusterFormProps {
    initialValues?: {
        id?: string
        name: string
        description: string
        region: string
        status: 'active' | 'inactive'
    }
    onSubmit: (values: any) => void
    onCancel: () => void
}

const ClusterForm = ({ initialValues, onSubmit, onCancel }: ClusterFormProps) => {
    const [form] = Form.useForm()

    useEffect(() => {
        if (initialValues) {
            form.setFieldsValue(initialValues)
        }
    }, [initialValues, form])

    const handleSubmit = async () => {
        try {
            const values = await form.validateFields()
            onSubmit(values)
        } catch (error) {
            console.error('表单验证失败:', error)
        }
    }

    const regions = [
        { value: 'cn-north-1', label: '华北1（北京）' },
        { value: 'cn-south-1', label: '华南1（广州）' },
        { value: 'cn-east-1', label: '华东1（上海）' },
    ]

    return (
        <Form
            form={form}
            layout="vertical"
            initialValues={{
                status: 'active',
            }}
        >
            <Form.Item
                name="name"
                label="集群名称"
                rules={[{ required: true, message: '请输入集群名称' }]}
            >
                <Input placeholder="请输入集群名称" />
            </Form.Item>

            <Form.Item
                name="description"
                label="描述"
                rules={[{ required: true, message: '请输入集群描述' }]}
            >
                <Input.TextArea rows={4} placeholder="请输入集群描述" />
            </Form.Item>

            <Form.Item
                name="region"
                label="区域"
                rules={[{ required: true, message: '请选择区域' }]}
            >
                <Select options={regions} placeholder="请选择区域" />
            </Form.Item>

            <Form.Item
                name="status"
                label="状态"
                rules={[{ required: true, message: '请选择状态' }]}
            >
                <Select
                    options={[
                        { value: 'active', label: '运行中' },
                        { value: 'inactive', label: '已停止' },
                    ]}
                />
            </Form.Item>

            <Form.Item>
                <Space>
                    <Button type="primary" onClick={handleSubmit}>
                        保存
                    </Button>
                    <Button onClick={onCancel}>取消</Button>
                </Space>
            </Form.Item>
        </Form>
    )
}

export default ClusterForm 