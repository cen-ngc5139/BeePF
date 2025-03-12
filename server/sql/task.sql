-- 创建任务表
CREATE TABLE IF NOT EXISTS `beepf`.`task` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '任务ID',
  `name` VARCHAR(255) NOT NULL COMMENT '任务名称',
  `description` TEXT NULL COMMENT '任务描述',
  `component_id` BIGINT UNSIGNED NOT NULL COMMENT '组件ID',
  `component_name` VARCHAR(255) NOT NULL COMMENT '组件名称',
  `step` INT NOT NULL COMMENT '任务步骤: 0-初始化, 1-加载, 2-启动, 3-统计, 4-指标, 5-停止',
  `status` INT NOT NULL COMMENT '任务状态: 0-等待中, 1-运行中, 2-成功, 3-失败',
  `error` TEXT NULL COMMENT '错误信息',
  `deleted` TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否删除: 0-否, 1-是',
  `creator` VARCHAR(255) NULL COMMENT '创建者',
  `created_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `last_update_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
  PRIMARY KEY (`id`),
  INDEX `idx_component_id` (`component_id`),
  INDEX `idx_status` (`status`),
  INDEX `idx_created_time` (`created_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='任务表';

-- 创建任务程序状态表
CREATE TABLE IF NOT EXISTS `beepf`.`task_program_status` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '状态ID',
  `task_id` BIGINT UNSIGNED NOT NULL COMMENT '任务ID',
  `component_id` BIGINT UNSIGNED NOT NULL COMMENT '组件ID',
  `component_name` VARCHAR(255) NOT NULL COMMENT '组件名称',
  `program_id` BIGINT UNSIGNED NOT NULL COMMENT '程序ID',
  `program_name` VARCHAR(255) NOT NULL COMMENT '程序名称',
  `attach_id` BIGINT UNSIGNED NOT NULL COMMENT '挂载ID',
  `status` INT NOT NULL COMMENT '状态: 0-等待中, 1-运行中, 2-成功, 3-失败',
  `error` TEXT NULL COMMENT '错误信息',
  `deleted` TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否删除: 0-否, 1-是',
  `created_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `last_update_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
  PRIMARY KEY (`id`),
  INDEX `idx_task_id` (`task_id`),
  INDEX `idx_component_id` (`component_id`),
  INDEX `idx_program_id` (`program_id`),
  INDEX `idx_status` (`status`),
  INDEX `idx_created_time` (`created_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='任务程序状态表';

-- 添加外键约束（如果需要）
-- ALTER TABLE `beepf`.`task_program_status` 
--   ADD CONSTRAINT `fk_task_program_status_task`
--   FOREIGN KEY (`task_id`) REFERENCES `beepf`.`task` (`id`)
--   ON DELETE CASCADE
--   ON UPDATE CASCADE; 