import { Form, Input, Select, Button, Space, message } from 'antd'
import { useEffect } from 'react'
import { Cluster } from '../../services/clusterService'

interface ClusterFormProps {
    initialValues?: Cluster;
    onSubmit: (values: Cluster) => void;
    onCancel: () => void;
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

    const environments = [
        { value: 'prod', label: '生产环境' },
        { value: 'test', label: '测试环境' },
        { value: 'dev', label: '开发环境' },
    ]

    return (
        <Form
            form={form}
            layout="vertical"
            initialValues={{
                status: 0,
                environment: 'test',
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
                name="cnname"
                label="中文名称"
                rules={[{ required: true, message: '请输入集群中文名称' }]}
            >
                <Input placeholder="请输入集群中文名称" />
            </Form.Item>

            <Form.Item
                name="desc"
                label="描述"
                rules={[{ required: true, message: '请输入集群描述' }]}
            >
                <Input.TextArea rows={4} placeholder="请输入集群描述" />
            </Form.Item>

            <Form.Item
                name="master"
                label="主节点地址"
                rules={[{ required: true, message: '请输入主节点地址' }]}
            >
                <Input placeholder="请输入主节点地址，例如：https://192.168.1.100:6443" />
            </Form.Item>

            <Form.Item
                name="environment"
                label="环境"
                rules={[{ required: true, message: '请选择环境' }]}
            >
                <Select options={environments} placeholder="请选择环境" />
            </Form.Item>

            <Form.Item
                name="status"
                label="状态"
                rules={[{ required: true, message: '请选择状态' }]}
            >
                <Select
                    options={[
                        { value: 0, label: '正常' },
                        { value: 1, label: '停用' },
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