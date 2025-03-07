import { useState, useEffect } from 'react'
import { Card, Form, Input, Select, Button, Upload, message, Row, Col, Divider, Space, Modal } from 'antd'
import { UploadOutlined, InfoCircleOutlined, CheckCircleOutlined, CloseCircleOutlined } from '@ant-design/icons'
import type { UploadFile, UploadProps } from 'antd/es/upload/interface'
import axios from 'axios'
import { useNavigate } from 'react-router-dom'
import clusterService, { Cluster } from '../../services/clusterService'

const { TextArea } = Input
const { Option } = Select

interface ComponentFormData {
    name: string
    cluster_id: number
    description?: string
}

const UploadComponent = () => {
    const [form] = Form.useForm()
    const [fileList, setFileList] = useState<UploadFile[]>([])
    const [uploading, setUploading] = useState(false)
    const [clusters, setClusters] = useState<Cluster[]>([])
    const [loading, setLoading] = useState(false)
    const [successModalVisible, setSuccessModalVisible] = useState(false)
    const [errorModalVisible, setErrorModalVisible] = useState(false)
    const [errorMessage, setErrorMessage] = useState('')
    const navigate = useNavigate()

    // 从后端获取集群列表
    useEffect(() => {
        const fetchClusters = async () => {
            setLoading(true)
            try {
                const clusterList = await clusterService.getClustersByParams()
                setClusters(clusterList)

                // 如果有集群数据，设置默认选中第一个
                if (clusterList.length > 0) {
                    form.setFieldsValue({ cluster_id: clusterList[0].id })
                }
            } catch (error) {
                console.error('获取集群列表失败:', error)
                message.error('获取集群列表失败，请刷新页面重试')
            } finally {
                setLoading(false)
            }
        }

        fetchClusters()
    }, [form])

    const handleUpload = async () => {
        try {
            // 验证表单
            let values;
            try {
                values = await form.validateFields()

                if (fileList.length === 0) {
                    // 显示错误弹窗
                    setErrorMessage('请选择要上传的二进制文件')
                    setErrorModalVisible(true)
                    return
                }
            } catch (validationError) {
                // 表单验证失败，显示错误弹窗
                setErrorMessage('表单验证失败，请检查必填字段')
                setErrorModalVisible(true)
                return
            }

            const formData = new FormData()

            // 添加二进制文件
            formData.append('binary', fileList[0] as any)

            // 将组件元信息作为JSON字符串添加到data字段
            const componentData = {
                name: values.name,
                cluster_id: values.cluster_id,
                description: values.description || ''
            }
            formData.append('data', JSON.stringify(componentData))

            setUploading(true)

            // 发送上传请求
            const response = await axios.post('/api/v1/component/upload', formData, {
                headers: {
                    'Content-Type': 'multipart/form-data',
                },
            })

            setUploading(false)
            setFileList([])
            form.resetFields()

            // 调试日志
            console.log('上传响应数据:', response.data)

            // 根据后端接口返回格式判断成功或失败
            if (response.data.success === true) {
                // 显示成功弹窗
                setSuccessModalVisible(true)
                // 不再立即跳转，而是在用户确认后跳转
            } else {
                // 显示错误弹窗
                setErrorMessage(response.data.errorMsg || '未知错误')
                setErrorModalVisible(true)
            }
        } catch (error: any) {
            setUploading(false)
            console.error('上传失败:', error)

            // 提取更详细的错误信息
            let errorMsg = '上传失败，请检查表单和文件'
            if (error.response) {
                // 服务器返回了错误状态码
                const { status, data } = error.response
                // 尝试从后端返回的错误格式中提取错误信息
                if (data && typeof data === 'object') {
                    if (data.errorMsg) {
                        // 使用后端返回的错误信息
                        errorMsg = `上传失败: ${data.errorMsg}`
                    } else if (data.message || data.error) {
                        // 兼容其他可能的错误格式
                        errorMsg = `上传失败 (${status}): ${data.message || data.error || '服务器错误'}`
                    } else {
                        errorMsg = `上传失败 (${status}): 服务器错误`
                    }
                } else {
                    errorMsg = `上传失败 (${status}): 服务器错误`
                }
            } else if (error.request) {
                // 请求已发送但没有收到响应
                errorMsg = '上传失败: 服务器无响应，请检查网络连接'
            } else if (error.message) {
                // 请求设置触发的错误
                errorMsg = `上传失败: ${error.message}`
            }

            // 显示错误弹窗
            setErrorMessage(errorMsg)
            setErrorModalVisible(true)
        }
    }

    const uploadProps: UploadProps = {
        onRemove: (file) => {
            setFileList([])
        },
        beforeUpload: (file) => {
            // 检查文件类型，这里假设只接受.o文件
            const isValidType = file.name.endsWith('.o')
            if (!isValidType) {
                // 显示错误弹窗
                setErrorMessage('文件类型错误：只能上传.o格式的二进制文件！')
                setErrorModalVisible(true)
                return Upload.LIST_IGNORE
            }

            // 限制文件大小为10MB
            const isLt10M = file.size / 1024 / 1024 < 10
            if (!isLt10M) {
                // 显示错误弹窗
                setErrorMessage('文件过大：文件大小不能超过10MB！')
                setErrorModalVisible(true)
                return Upload.LIST_IGNORE
            }

            // 文件验证通过，不需要弹窗提示，直接添加到文件列表
            setFileList([file])
            return false
        },
        fileList,
    }

    return (
        <>
            <Card title="上传组件" bordered={false}>
                <Form
                    form={form}
                    layout="vertical"
                    initialValues={{ cluster_id: clusters.length > 0 ? clusters[0].id : undefined }}
                >
                    <Row gutter={24}>
                        <Col span={12}>
                            <Form.Item
                                name="name"
                                label="组件名称"
                                rules={[{ required: true, message: '请输入组件名称' }]}
                            >
                                <Input placeholder="请输入组件名称" />
                            </Form.Item>
                        </Col>
                        <Col span={12}>
                            <Form.Item
                                name="cluster_id"
                                label="所属集群"
                                rules={[{ required: true, message: '请选择所属集群' }]}
                            >
                                <Select placeholder="请选择所属集群" loading={loading}>
                                    {clusters.map(cluster => (
                                        <Option key={cluster.id} value={cluster.id}>{cluster.name}</Option>
                                    ))}
                                </Select>
                            </Form.Item>
                        </Col>
                    </Row>

                    <Form.Item
                        name="description"
                        label="组件描述"
                    >
                        <TextArea rows={4} placeholder="请输入组件描述" />
                    </Form.Item>

                    <Divider />

                    <Form.Item
                        label="上传二进制文件"
                        required
                        tooltip={{ title: '支持.o格式的eBPF二进制文件，大小不超过10MB', icon: <InfoCircleOutlined /> }}
                    >
                        <Upload {...uploadProps} maxCount={1}>
                            <Button icon={<UploadOutlined />}>选择文件</Button>
                        </Upload>
                    </Form.Item>

                    <Form.Item>
                        <Space>
                            <Button type="primary" onClick={handleUpload} loading={uploading}>
                                上传组件
                            </Button>
                            <Button onClick={() => navigate('/components/list')}>
                                取消
                            </Button>
                        </Space>
                    </Form.Item>
                </Form>
            </Card>

            {/* 成功弹窗 */}
            <Modal
                title={<div style={{ display: 'flex', alignItems: 'center' }}><CheckCircleOutlined style={{ color: '#52c41a', marginRight: 8 }} />上传成功</div>}
                open={successModalVisible}
                onOk={() => {
                    setSuccessModalVisible(false);
                    navigate('/components/list');
                }}
                onCancel={() => {
                    setSuccessModalVisible(false);
                    navigate('/components/list');
                }}
                okText="确定"
                cancelText="取消"
                centered
            >
                <p style={{ fontSize: '16px', margin: '20px 0' }}>组件上传成功！</p>
                <p>点击确定返回组件列表页面。</p>
            </Modal>

            {/* 错误弹窗 */}
            <Modal
                title={<div style={{ display: 'flex', alignItems: 'center' }}><CloseCircleOutlined style={{ color: '#ff4d4f', marginRight: 8 }} />上传失败</div>}
                open={errorModalVisible}
                onOk={() => setErrorModalVisible(false)}
                onCancel={() => setErrorModalVisible(false)}
                okText="确定"
                cancelText="取消"
                centered
            >
                <p style={{ fontSize: '16px', margin: '20px 0' }}>组件上传失败！</p>
                <p style={{ color: '#ff4d4f' }}>{errorMessage}</p>
                <p>请检查表单和文件后重试。</p>
            </Modal>
        </>
    )
}

export default UploadComponent 