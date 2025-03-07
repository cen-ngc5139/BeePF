import { Breadcrumb } from 'antd'
import { Link, useLocation, useParams } from 'react-router-dom'
import { HomeOutlined } from '@ant-design/icons'
import { useEffect, useState } from 'react'
import clusterService from '../services/clusterService'
import componentService from '../services/componentService'

const breadcrumbNameMap: Record<string, string> = {
    '/clusters': '集群管理',
    '/clusters/list': '集群列表',
    '/clusters/create': '新建集群',
    '/clusters/edit': '编辑集群',
    '/components': '组件管理',
    '/components/list': '组件列表',
    '/components/upload': '上传组件',
    '/component': '组件详情',
    '/observability': '可观测',
    '/workflow': '工作流',
}

const PageBreadcrumb = () => {
    const location = useLocation()
    const params = useParams()
    const pathSnippets = location.pathname.split('/').filter((i) => i)
    const [clusterName, setClusterName] = useState<string>('')
    const [componentName, setComponentName] = useState<string>('')

    // 如果是编辑集群页面，获取集群名称
    useEffect(() => {
        const fetchClusterName = async () => {
            if (location.pathname.startsWith('/clusters/edit/') && params.id) {
                try {
                    const cluster = await clusterService.getCluster(parseInt(params.id))
                    setClusterName(cluster.cnname || cluster.name || '')
                } catch (error) {
                    console.error('获取集群名称失败:', error)
                }
            }
        }

        fetchClusterName()
    }, [location.pathname, params.id])

    // 如果是组件详情页面，获取组件名称
    useEffect(() => {
        const fetchComponentName = async () => {
            if (location.pathname.startsWith('/component/') && params.id) {
                try {
                    const component = await componentService.getComponent(parseInt(params.id))
                    setComponentName(component.name || '')
                } catch (error) {
                    console.error('获取组件名称失败:', error)
                }
            }
        }

        fetchComponentName()
    }, [location.pathname, params.id])

    // 处理编辑页面的面包屑
    if (location.pathname.startsWith('/clusters/edit/')) {
        return (
            <Breadcrumb
                items={[
                    {
                        title: (
                            <Link to="/">
                                <HomeOutlined /> 首页
                            </Link>
                        ),
                        key: 'home',
                    },
                    {
                        title: <Link to="/clusters/list">集群列表</Link>,
                        key: 'clusters-list',
                    },
                    {
                        title: `编辑集群${clusterName ? `: ${clusterName}` : ''}`,
                        key: 'edit-cluster',
                    },
                ]}
            />
        )
    }

    // 处理组件详情页面的面包屑
    if (location.pathname.startsWith('/component/')) {
        return (
            <Breadcrumb
                items={[
                    {
                        title: (
                            <Link to="/">
                                <HomeOutlined /> 首页
                            </Link>
                        ),
                        key: 'home',
                    },
                    {
                        title: <Link to="/components/list">组件列表</Link>,
                        key: 'components-list',
                    },
                    {
                        title: `组件详情${componentName ? `: ${componentName}` : ''}`,
                        key: 'component-detail',
                    },
                ]}
            />
        )
    }

    const extraBreadcrumbItems = pathSnippets.map((_, index) => {
        const url = `/${pathSnippets.slice(0, index + 1).join('/')}`
        return {
            key: url,
            title: <Link to={url}>{breadcrumbNameMap[url]}</Link>,
        }
    })

    const breadcrumbItems = [
        {
            title: (
                <Link to="/">
                    <HomeOutlined /> 首页
                </Link>
            ),
            key: 'home',
        },
    ].concat(extraBreadcrumbItems)

    return <Breadcrumb items={breadcrumbItems} />
}

export default PageBreadcrumb 