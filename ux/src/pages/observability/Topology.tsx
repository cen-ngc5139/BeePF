import React, { useEffect, useState } from 'react';
import { Card, Spin, Empty, Button } from 'antd';
import { ReloadOutlined } from '@ant-design/icons';
import TopoGraph from '../../components/TopoGraph';
import { getTopology, Topology } from '../../services/topo';
import './Topology.css';

const TopologyPage: React.FC = () => {
    const [loading, setLoading] = useState<boolean>(true);
    const [topology, setTopology] = useState<Topology | null>(null);
    const [error, setError] = useState<string | null>(null);

    const fetchTopology = async () => {
        setLoading(true);
        setError(null);
        try {
            const data = await getTopology();
            setTopology(data);
        } catch (err) {
            console.error('获取拓扑数据失败:', err);
            setError('获取拓扑数据失败，请稍后重试');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchTopology();
    }, []);

    const handleRefresh = () => {
        fetchTopology();
    };

    return (
        <div className="topology-page">
            <Card
                title="eBPF 程序和 Map 拓扑关系"
                extra={
                    <Button
                        type="primary"
                        icon={<ReloadOutlined />}
                        onClick={handleRefresh}
                        loading={loading}
                    >
                        刷新
                    </Button>
                }
                className="topology-card"
            >
                {loading ? (
                    <div className="loading-container">
                        <Spin size="large" tip="加载中..." />
                    </div>
                ) : error ? (
                    <div className="error-container">
                        <Empty
                            description={error}
                            image={Empty.PRESENTED_IMAGE_SIMPLE}
                        />
                    </div>
                ) : topology && (topology.ProgNodes.length > 0 || topology.MapNodes.length > 0) ? (
                    <TopoGraph data={topology} />
                ) : (
                    <Empty description="暂无拓扑数据" />
                )}
            </Card>
        </div>
    );
};

export default TopologyPage; 