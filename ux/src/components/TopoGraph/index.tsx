import React, { useEffect, useRef, useState } from 'react';
import G6 from '@antv/g6';
import { Topology } from '../../services/topo';
import { Radio, Space, Tooltip, Badge } from 'antd';
import { PartitionOutlined, RadarChartOutlined, NodeIndexOutlined, ColumnWidthOutlined, AppstoreOutlined } from '@ant-design/icons';
import './index.css';

// 导入 dagre 布局
import 'dagre';

interface TopoGraphProps {
    data: Topology;
    loading?: boolean;
}

type LayoutType = 'dagre' | 'force' | 'radial' | 'grid' | 'map-centric';

const TopoGraph: React.FC<TopoGraphProps> = ({ data, loading = false }) => {
    const containerRef = useRef<HTMLDivElement>(null);
    const graphRef = useRef<any>(null);
    const [layout, setLayout] = useState<LayoutType>('map-centric');

    // 将后端数据转换为 G6 可用的格式
    const transformData = (topology: Topology) => {
        const nodes: any[] = [];
        const edges: any[] = [];

        // 计算每个 Map 被引用的次数
        const mapRefCount: Record<string, number> = {};
        topology.Edges.forEach(edge => {
            const mapId = `map-${edge.MapID}`;
            mapRefCount[mapId] = (mapRefCount[mapId] || 0) + 1;
        });

        // 添加程序节点
        topology.ProgNodes.forEach(prog => {
            nodes.push({
                id: `prog-${prog.ID}`,
                label: prog.Name || `Program ${prog.ID}`,
                type: 'program-node',
                style: {
                    fill: '#91d5ff',
                    stroke: '#40a9ff'
                },
                // 添加类型标记，用于筛选
                nodeType: 'program'
            });
        });

        // 添加 Map 节点
        topology.MapNodes.forEach(map => {
            const mapId = `map-${map.ID}`;
            const refCount = mapRefCount[mapId] || 0;
            nodes.push({
                id: mapId,
                label: map.Name || `Map ${map.ID}`,
                type: 'map-node',
                style: {
                    fill: '#d3f261',
                    stroke: '#7cb305'
                },
                // 添加类型标记，用于筛选
                nodeType: 'map',
                // 添加引用计数，用于中心度计算
                refCount: refCount,
                // 设置节点大小，根据引用次数调整
                size: Math.max(40, Math.min(80, 40 + refCount * 5))
            });
        });

        // 添加边
        topology.Edges.forEach(edge => {
            edges.push({
                id: `edge-${edge.ProgID}-${edge.MapID}`,
                source: `prog-${edge.ProgID}`,
                target: `map-${edge.MapID}`,
                style: {
                    endArrow: true
                }
            });
        });

        return {
            nodes,
            edges
        };
    };

    // 获取布局配置
    const getLayoutConfig = (type: LayoutType) => {
        switch (type) {
            case 'dagre':
                return {
                    type: 'dagre',
                    rankdir: 'LR', // 水平布局
                    nodesep: 50,
                    ranksep: 70,
                    preventOverlap: true,
                };
            case 'force':
                return {
                    type: 'force',
                    preventOverlap: true,
                    linkDistance: 100,
                    nodeStrength: -50,
                    edgeStrength: 0.1,
                    nodeSpacing: 50,
                };
            case 'map-centric':
                return {
                    type: 'force',
                    preventOverlap: true,
                    linkDistance: (edge: any) => 100,
                    nodeStrength: (node: any) => {
                        // Map 节点的引力与引用次数成正比
                        if (node.nodeType === 'map') {
                            return -30 + node.refCount * -10;
                        }
                        // 程序节点的引力较小
                        return -10;
                    },
                    // 中心引力，引用次数越多的 Map 节点越靠近中心
                    center: [0, 0],
                    gravity: 0.1,
                    // 根据节点类型和引用次数调整质量
                    nodeSize: (node: any) => {
                        if (node.nodeType === 'map') {
                            return node.refCount > 0 ? node.refCount * 5 + 20 : 20;
                        }
                        return 20;
                    },
                    edgeStrength: 0.5,
                    nodeSpacing: 50,
                    // 迭代次数增加，使布局更稳定
                    iterations: 500,
                    // 初始位置，Map 节点更靠近中心
                    getCenter: (node: any) => {
                        if (node.nodeType === 'map' && node.refCount > 0) {
                            return [0, 0];
                        }
                        return [Math.random() * 200 - 100, Math.random() * 200 - 100];
                    },
                };
            case 'radial':
                return {
                    type: 'radial',
                    preventOverlap: true,
                    unitRadius: 100,
                };
            case 'grid':
                return {
                    type: 'grid',
                    preventOverlap: true,
                    nodeSize: 40,
                    sortBy: 'nodeType',
                };
            default:
                return {
                    type: 'dagre',
                    rankdir: 'LR',
                    nodesep: 50,
                    ranksep: 70,
                    preventOverlap: true,
                };
        }
    };

    // 更新布局
    const updateLayout = (type: LayoutType) => {
        if (!graphRef.current) return;

        setLayout(type);

        const layoutConfig = getLayoutConfig(type);
        graphRef.current.updateLayout(layoutConfig);
        graphRef.current.fitView();
    };

    useEffect(() => {
        if (!containerRef.current || loading) return;

        // 注册自定义节点
        G6.registerNode(
            'program-node',
            {
                draw(cfg: any, group: any) {
                    const width = 120;
                    const height = 40;
                    const keyShape = group.addShape('rect', {
                        attrs: {
                            x: -width / 2,
                            y: -height / 2,
                            width,
                            height,
                            radius: 4,
                            fill: '#91d5ff',
                            stroke: '#40a9ff',
                            lineWidth: 2,
                            cursor: 'pointer',
                            ...cfg.style
                        },
                        name: 'program-node-keyshape',
                    });

                    // 添加标签
                    group.addShape('text', {
                        attrs: {
                            text: cfg.label || '',
                            x: 0,
                            y: 0,
                            fontSize: 12,
                            textAlign: 'center',
                            textBaseline: 'middle',
                            fill: '#000',
                            cursor: 'pointer',
                        },
                        name: 'program-node-label',
                    });

                    // 添加图标
                    group.addShape('text', {
                        attrs: {
                            text: 'P',
                            x: -width / 2 + 12,
                            y: -height / 2 + 12,
                            fontSize: 10,
                            fontWeight: 'bold',
                            fill: '#1890ff',
                            cursor: 'pointer',
                        },
                        name: 'program-node-icon',
                    });

                    return keyShape;
                },
            },
            'single-node',
        );

        G6.registerNode(
            'map-node',
            {
                draw(cfg: any, group: any) {
                    const width = 120;
                    const height = 40;
                    const keyShape = group.addShape('rect', {
                        attrs: {
                            x: -width / 2,
                            y: -height / 2,
                            width,
                            height,
                            radius: 4,
                            fill: '#d3f261',
                            stroke: '#7cb305',
                            lineWidth: 2,
                            cursor: 'pointer',
                            ...cfg.style
                        },
                        name: 'map-node-keyshape',
                    });

                    // 添加标签
                    group.addShape('text', {
                        attrs: {
                            text: cfg.label || '',
                            x: 0,
                            y: 0,
                            fontSize: 12,
                            textAlign: 'center',
                            textBaseline: 'middle',
                            fill: '#000',
                            cursor: 'pointer',
                        },
                        name: 'map-node-label',
                    });

                    // 添加图标
                    group.addShape('text', {
                        attrs: {
                            text: 'M',
                            x: -width / 2 + 12,
                            y: -height / 2 + 12,
                            fontSize: 10,
                            fontWeight: 'bold',
                            fill: '#7cb305',
                            cursor: 'pointer',
                        },
                        name: 'map-node-icon',
                    });

                    // 如果有引用计数，显示在右上角
                    if (cfg.refCount && cfg.refCount > 0) {
                        group.addShape('circle', {
                            attrs: {
                                x: width / 2 - 10,
                                y: -height / 2 + 10,
                                r: 8,
                                fill: '#fa8c16',
                                stroke: '#fff',
                                lineWidth: 1,
                            },
                            name: 'map-node-badge-bg',
                        });

                        group.addShape('text', {
                            attrs: {
                                text: cfg.refCount,
                                x: width / 2 - 10,
                                y: -height / 2 + 10,
                                fontSize: 10,
                                fontWeight: 'bold',
                                fill: '#fff',
                                textAlign: 'center',
                                textBaseline: 'middle',
                            },
                            name: 'map-node-badge-text',
                        });
                    }

                    return keyShape;
                },
            },
            'single-node',
        );

        // 如果图已经存在，先销毁
        if (graphRef.current) {
            graphRef.current.destroy();
            graphRef.current = null;
        }

        // 创建图实例
        const width = containerRef.current.scrollWidth;
        const height = containerRef.current.scrollHeight || 500;

        const layoutConfig = getLayoutConfig(layout);

        const graph = new G6.Graph({
            container: containerRef.current,
            width,
            height,
            layout: layoutConfig,
            defaultNode: {
                size: [120, 40],
            },
            defaultEdge: {
                style: {
                    stroke: '#91d5ff',
                    lineWidth: 2,
                    endArrow: {
                        path: G6.Arrow.triangle(8, 10, 0),
                        fill: '#91d5ff',
                    },
                },
            },
            modes: {
                default: ['drag-canvas', 'zoom-canvas', 'drag-node', 'click-select'],
            },
            fitView: true,
            animate: true,
            // 添加节点交互效果
            nodeStateStyles: {
                hover: {
                    lineWidth: 3,
                    shadowColor: '#ccc',
                    shadowBlur: 10
                },
                selected: {
                    lineWidth: 3,
                    shadowColor: '#1890ff',
                    shadowBlur: 10
                }
            },
            // 添加边交互效果
            edgeStateStyles: {
                hover: {
                    lineWidth: 3
                },
                selected: {
                    lineWidth: 3,
                    stroke: '#1890ff'
                }
            }
        });

        // 添加节点交互
        graph.on('node:mouseenter', (evt: any) => {
            const { item } = evt;
            graph.setItemState(item, 'hover', true);
            // 高亮相关边
            const edges = item.getEdges();
            edges.forEach((edge: any) => {
                graph.setItemState(edge, 'hover', true);
            });
        });

        graph.on('node:mouseleave', (evt: any) => {
            const { item } = evt;
            graph.setItemState(item, 'hover', false);
            // 取消高亮相关边
            const edges = item.getEdges();
            edges.forEach((edge: any) => {
                graph.setItemState(edge, 'hover', false);
            });
        });

        // 添加节点点击事件，显示详细信息
        graph.on('node:click', (evt: any) => {
            const { item } = evt;
            const model = item.getModel();

            // 显示节点详细信息
            console.log('节点详细信息:', model);

            // 高亮选中节点
            graph.getNodes().forEach((node: any) => {
                graph.clearItemStates(node);
            });
            graph.setItemState(item, 'selected', true);

            // 高亮相关边和节点
            const edges = item.getEdges();
            edges.forEach((edge: any) => {
                graph.setItemState(edge, 'selected', true);
                const otherNode = edge.getSource() === item ? edge.getTarget() : edge.getSource();
                graph.setItemState(otherNode, 'selected', true);
            });
        });

        // 添加画布点击事件，清除选中状态
        graph.on('canvas:click', () => {
            graph.getNodes().forEach((node: any) => {
                graph.clearItemStates(node);
            });
            graph.getEdges().forEach((edge: any) => {
                graph.clearItemStates(edge);
            });
        });

        graphRef.current = graph;

        // 渲染数据
        if (data && data.ProgNodes && data.MapNodes) {
            const graphData = transformData(data);
            graph.data(graphData);
            graph.render();
        }

        // 监听窗口大小变化
        const handleResize = () => {
            if (graphRef.current) {
                graphRef.current.changeSize(
                    containerRef.current?.scrollWidth || 800,
                    containerRef.current?.scrollHeight || 500
                );
                graphRef.current.fitView();
            }
        };

        window.addEventListener('resize', handleResize);

        return () => {
            window.removeEventListener('resize', handleResize);
        };
    }, [data, loading]);

    return (
        <div className="topo-graph-container">
            {loading ? (
                <div className="loading-container">加载中...</div>
            ) : (
                <>
                    <div className="layout-controls">
                        <Radio.Group value={layout} onChange={(e) => updateLayout(e.target.value)} buttonStyle="solid">
                            <Tooltip title="层次布局">
                                <Radio.Button value="dagre"><ColumnWidthOutlined /></Radio.Button>
                            </Tooltip>
                            <Tooltip title="力导向布局">
                                <Radio.Button value="force"><NodeIndexOutlined /></Radio.Button>
                            </Tooltip>
                            <Tooltip title="Map 中心布局">
                                <Radio.Button value="map-centric"><AppstoreOutlined /></Radio.Button>
                            </Tooltip>
                            <Tooltip title="辐射布局">
                                <Radio.Button value="radial"><RadarChartOutlined /></Radio.Button>
                            </Tooltip>
                            <Tooltip title="网格布局">
                                <Radio.Button value="grid"><PartitionOutlined /></Radio.Button>
                            </Tooltip>
                        </Radio.Group>
                    </div>
                    <div className="layout-legend">
                        <div className="legend-item">
                            <div className="legend-color program-color"></div>
                            <span>程序节点</span>
                        </div>
                        <div className="legend-item">
                            <div className="legend-color map-color"></div>
                            <span>Map 节点</span>
                        </div>
                        <div className="legend-item">
                            <div className="legend-badge"></div>
                            <span>引用次数</span>
                        </div>
                    </div>
                    <div ref={containerRef} className="graph-container" />
                </>
            )}
        </div>
    );
};

export default TopoGraph; 