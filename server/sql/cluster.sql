-- 创建数据库   
create database if not exists beepf;

-- 集群表
create table if not exists beepf.cluster
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

-- Component 表
create table if not exists beepf.component
(
    id               bigint unsigned AUTO_INCREMENT NOT NULL PRIMARY KEY comment 'ID',
    name             varchar(100)                   NOT NULL comment '组件名称',
    cluster_id       bigint unsigned                NOT NULL comment '所属集群ID',
    binary_path      varchar(500)                   NOT NULL comment '组件二进制文件路径',
    deleted          tinyint                        NOT NULL DEFAULT '0' comment '是否删除',
    creator          varchar(50)                    NOT NULL DEFAULT '' comment '创建用户',
    created_time     datetime                       NOT NULL DEFAULT CURRENT_TIMESTAMP comment '创建时间',
    last_update_time datetime                       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP comment '更新时间',
    UNIQUE KEY `uk_name` (`name`)
) ENGINE = InnoDB
  DEFAULT CHARSET utf8mb4 COMMENT = 'eBPF组件表';

-- Program 表
create table if not exists beepf.program
(
    id               bigint unsigned AUTO_INCREMENT NOT NULL PRIMARY KEY comment 'ID',
    component_id     bigint unsigned                NOT NULL comment '所属组件ID',
    name             varchar(100)                   NOT NULL comment '程序名称',
    description      text                           NULL comment '程序描述',
    deleted          tinyint                        NOT NULL DEFAULT '0' comment '是否删除',
    creator          varchar(50)                    NOT NULL DEFAULT '' comment '创建用户',
    created_time     datetime                       NOT NULL DEFAULT CURRENT_TIMESTAMP comment '创建时间',
    last_update_time datetime                       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP comment '更新时间',
    KEY `idx_component_id` (`component_id`),
    FOREIGN KEY (`component_id`) REFERENCES beepf.component(`id`) ON DELETE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET utf8mb4 COMMENT = 'eBPF程序表';

-- ProgramSpec 表
create table if not exists beepf.program_spec
(
    id               bigint unsigned AUTO_INCREMENT NOT NULL PRIMARY KEY comment 'ID',
    program_id       bigint unsigned                NOT NULL comment '所属程序ID',
    name             varchar(100)                   NOT NULL comment '名称',
    type             int unsigned                   NOT NULL comment '程序类型',
    attach_type      int unsigned                   NOT NULL comment '附加类型',
    attach_to        varchar(255)                   NULL comment '附加目标',
    section_name     varchar(100)                   NULL comment 'ELF段名称',
    flags            int unsigned                   NOT NULL DEFAULT 0 comment '标志',
    license          varchar(50)                    NOT NULL comment '许可证',
    kernel_version   int unsigned                   NULL comment '内核版本',
    deleted          tinyint                        NOT NULL DEFAULT '0' comment '是否删除',
    created_time     datetime                       NOT NULL DEFAULT CURRENT_TIMESTAMP comment '创建时间',
    last_update_time datetime                       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP comment '更新时间',
    UNIQUE KEY `uk_program_id` (`program_id`),
    FOREIGN KEY (`program_id`) REFERENCES beepf.program(`id`) ON DELETE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET utf8mb4 COMMENT = 'eBPF程序规格表';

-- ProgramProperties 表
create table if not exists beepf.program_properties
(
    id               bigint unsigned AUTO_INCREMENT NOT NULL PRIMARY KEY comment 'ID',
    program_id       bigint unsigned                NOT NULL comment '所属程序ID',
    properties_json  json                           NOT NULL comment '程序属性JSON',
    deleted          tinyint                        NOT NULL DEFAULT '0' comment '是否删除',
    created_time     datetime                       NOT NULL DEFAULT CURRENT_TIMESTAMP comment '创建时间',
    last_update_time datetime                       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP comment '更新时间',
    UNIQUE KEY `uk_program_id` (`program_id`),
    FOREIGN KEY (`program_id`) REFERENCES beepf.program(`id`) ON DELETE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET utf8mb4 COMMENT = 'eBPF程序属性表';

-- Map 表
create table if not exists beepf.map
(
    id               bigint unsigned AUTO_INCREMENT NOT NULL PRIMARY KEY comment 'ID',
    component_id     bigint unsigned                NOT NULL comment '所属组件ID',
    name             varchar(100)                   NOT NULL comment 'Map名称',
    description      text                           NULL comment 'Map描述',
    deleted          tinyint                        NOT NULL DEFAULT '0' comment '是否删除',
    creator          varchar(50)                    NOT NULL DEFAULT '' comment '创建用户',
    created_time     datetime                       NOT NULL DEFAULT CURRENT_TIMESTAMP comment '创建时间',
    last_update_time datetime                       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP comment '更新时间',
    KEY `idx_component_id` (`component_id`),
    FOREIGN KEY (`component_id`) REFERENCES beepf.component(`id`) ON DELETE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET utf8mb4 COMMENT = 'eBPF Map表';

-- MapSpec 表
create table if not exists beepf.map_spec
(
    id               bigint unsigned AUTO_INCREMENT NOT NULL PRIMARY KEY comment 'ID',
    map_id           bigint unsigned                NOT NULL comment '所属Map ID',
    name             varchar(100)                   NOT NULL comment '名称',
    type             int unsigned                   NOT NULL comment 'Map类型',
    key_size         int unsigned                   NOT NULL comment '键大小',
    value_size       int unsigned                   NOT NULL comment '值大小',
    max_entries      int unsigned                   NOT NULL comment '最大条目数',
    flags            int unsigned                   NOT NULL DEFAULT 0 comment '标志',
    pinning          varchar(50)                    NULL comment '固定类型',
    deleted          tinyint                        NOT NULL DEFAULT '0' comment '是否删除',
    created_time     datetime                       NOT NULL DEFAULT CURRENT_TIMESTAMP comment '创建时间',
    last_update_time datetime                       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP comment '更新时间',
    UNIQUE KEY `uk_map_id` (`map_id`),
    FOREIGN KEY (`map_id`) REFERENCES beepf.map(`id`) ON DELETE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET utf8mb4 COMMENT = 'eBPF Map规格表';

-- MapProperties 表
create table if not exists beepf.map_properties
(
    id               bigint unsigned AUTO_INCREMENT NOT NULL PRIMARY KEY comment 'ID',
    map_id           bigint unsigned                NOT NULL comment '所属Map ID',
    properties_json  json                           NOT NULL comment 'Map属性JSON',
    deleted          tinyint                        NOT NULL DEFAULT '0' comment '是否删除',
    created_time     datetime                       NOT NULL DEFAULT CURRENT_TIMESTAMP comment '创建时间',
    last_update_time datetime                       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP comment '更新时间',
    UNIQUE KEY `uk_map_id` (`map_id`),
    FOREIGN KEY (`map_id`) REFERENCES beepf.map(`id`) ON DELETE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET utf8mb4 COMMENT = 'eBPF Map属性表';