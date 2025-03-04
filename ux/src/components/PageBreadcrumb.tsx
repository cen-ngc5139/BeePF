import { Breadcrumb } from 'antd'
import { Link, useLocation, useParams } from 'react-router-dom'
import { HomeOutlined } from '@ant-design/icons'

const breadcrumbNameMap: Record<string, string> = {
    '/clusters': '集群管理',
    '/clusters/list': '集群列表',
    '/clusters/create': '新建集群',
    '/clusters/edit': '编辑集群',
    '/components': '组件管理',
    '/components/list': '组件列表',
    '/components/create': '新建组件',
    '/observability': '可观测',
    '/workflow': '工作流',
}

const PageBreadcrumb = () => {
    const location = useLocation()
    const params = useParams()
    const pathSnippets = location.pathname.split('/').filter((i) => i)

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
                        title: `编辑集群 ${params.id}`,
                        key: 'edit-cluster',
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