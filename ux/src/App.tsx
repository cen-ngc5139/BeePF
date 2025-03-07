import { useState } from 'react'
import { Layout, Menu, theme, Typography } from 'antd'
import {
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  ClusterOutlined,
  AppstoreOutlined,
  LineChartOutlined,
  NodeIndexOutlined,
  UnorderedListOutlined,
  PlusOutlined,
  UploadOutlined,
} from '@ant-design/icons'
import { BrowserRouter as Router, Route, Routes, Link, Navigate } from 'react-router-dom'
import ComponentList from './pages/components/ComponentList'
import ComponentDetail from './pages/components/ComponentDetail'
import CreateComponent from './pages/components/CreateComponent'
import UploadComponent from './pages/components/UploadComponent'
import ClusterList from './pages/clusters/ClusterList'
import CreateCluster from './pages/clusters/CreateCluster'
import PageBreadcrumb from './components/PageBreadcrumb'

const { Header, Sider, Content } = Layout
const { Title } = Typography

function App() {
  const [collapsed, setCollapsed] = useState(false)
  const {
    token: { colorBgContainer, borderRadiusLG },
  } = theme.useToken()

  return (
    <Router>
      <Layout style={{ minHeight: '100vh' }}>
        <Sider trigger={null} collapsible collapsed={collapsed}>
          <div className="logo">
            <Title level={4} style={{ color: '#fff', margin: '16px', textAlign: 'center' }}>
              {collapsed ? 'BPF' : 'BeePF'}
            </Title>
          </div>
          <Menu
            theme="dark"
            mode="inline"
            defaultSelectedKeys={['1']}
            defaultOpenKeys={['components', 'clusters']}
            items={[
              {
                key: 'clusters',
                icon: <ClusterOutlined />,
                label: '集群管理',
                children: [
                  {
                    key: 'clusters-list',
                    icon: <UnorderedListOutlined />,
                    label: <Link to="/clusters/list">集群列表</Link>,
                  },
                  {
                    key: 'clusters-create',
                    icon: <PlusOutlined />,
                    label: <Link to="/clusters/create">新建集群</Link>,
                  },
                ],
              },
              {
                key: 'components',
                icon: <AppstoreOutlined />,
                label: '组件管理',
                children: [
                  {
                    key: '2-1',
                    icon: <UnorderedListOutlined />,
                    label: <Link to="/components/list">组件列表</Link>,
                  },
                  {
                    key: '2-2',
                    icon: <PlusOutlined />,
                    label: <Link to="/components/create">新建组件</Link>,
                  },
                  {
                    key: '2-3',
                    icon: <UploadOutlined />,
                    label: <Link to="/components/upload">上传组件</Link>,
                  },
                ],
              },
              {
                key: '3',
                icon: <LineChartOutlined />,
                label: <Link to="/observability">可观测</Link>,
              },
              {
                key: '4',
                icon: <NodeIndexOutlined />,
                label: <Link to="/workflow">工作流</Link>,
              },
            ]}
          />
        </Sider>
        <Layout>
          <Header style={{ padding: 0, background: colorBgContainer }}>
            <div
              style={{
                padding: '0 16px',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between',
                height: '100%',
              }}
            >
              <div style={{ display: 'flex', alignItems: 'center' }}>
                {collapsed ? (
                  <MenuUnfoldOutlined onClick={() => setCollapsed(!collapsed)} />
                ) : (
                  <MenuFoldOutlined onClick={() => setCollapsed(!collapsed)} />
                )}
                <Title level={4} style={{ margin: '0 0 0 16px' }}>
                  BeePF 管理平台
                </Title>
              </div>
            </div>
          </Header>
          <div style={{ padding: '16px 24px 0', background: colorBgContainer }}>
            <PageBreadcrumb />
          </div>
          <Content
            style={{
              margin: '0 16px 16px',
              padding: 24,
              background: colorBgContainer,
              borderRadius: borderRadiusLG,
            }}
          >
            <Routes>
              <Route path="/" element={<Navigate to="/components/list" replace />} />
              <Route path="/clusters/list" element={<ClusterList />} />
              <Route path="/clusters/create" element={<CreateCluster />} />
              <Route path="/clusters/edit/:id" element={<CreateCluster />} />
              <Route path="/components/list" element={<ComponentList />} />
              <Route path="/components/create" element={<CreateComponent />} />
              <Route path="/components/upload" element={<UploadComponent />} />
              <Route path="/component/:id" element={<ComponentDetail />} />
              <Route path="/observability" element={<div>可观测</div>} />
              <Route path="/workflow" element={<div>工作流</div>} />
            </Routes>
          </Content>
        </Layout>
      </Layout>
    </Router>
  )
}

export default App
