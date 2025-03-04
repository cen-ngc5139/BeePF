-- 创建数据库   
create database beepf;

-- 集群表
create table beepf.cluster
(
    id               bigint unsigned AUTO_INCREMENT NOT NULL PRIMARY KEY comment 'ID',
    cluster_name     varchar(100)                   NOT NULL DEFAULT '' comment '集群名称',
    cn_name          varchar(64)                    NOT NULL DEFAULT '' comment '别名',
    cluster_master   varchar(100)                   NOT NULL DEFAULT '' comment '主节点地址',
    kube_config      text                           NOT NULL comment 'kubeconf 配置文件',
    cluster_status   int(11)                        NOT NULL DEFAULT '1' comment '集群状态',
    deleted          tinyint                        NOT NULL DEFAULT '0' comment '是否删除',
    cluster_desc     varchar(100)                   NOT NULL DEFAULT '' comment '集群描述',
    environment      varchar(64)                    NOT NULL DEFAULT '' comment '集群所属环境',
    creator          varchar(50)                    NOT NULL DEFAULT '' comment '创建用户',
    created_time     datetime                       NOT NULL DEFAULT CURRENT_TIMESTAMP comment '创建时间',
    last_update_time datetime                       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP comment '更新时间',
    location_id      int                            NOT NULL DEFAULT '0' comment '集群所属数据中心信息',
    KEY `idx_last_update_time` (`last_update_time`),
    UNIQUE KEY `uk_name` (`cluster_name`)
) ENGINE = InnoDB
  DEFAULT CHARSET utf8mb4 COMMENT = '集群表';