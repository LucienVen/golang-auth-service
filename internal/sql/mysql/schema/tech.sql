-- https://www.drawdb.app/editor?shareId=6828c697c5d567acf2360d9665f879dd

-- 班级表
CREATE TABLE `class` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `class_num` tinyint unsigned DEFAULT NULL COMMENT '班级序号',
  `grade_level` tinyint DEFAULT NULL COMMENT '年级',
  `main_teacher_id` int DEFAULT NULL COMMENT '班主任',
  `is_delete` tinyint DEFAULT '0' COMMENT '是否删除（逻辑删除标记）',
  `creator` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '创建者',
  `updater` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '更新者',
  `create_time` int DEFAULT NULL COMMENT '创建时间',
  `update_time` int DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;



-- 考试表
CREATE TABLE `exam` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '考试名称',
  `is_delete` tinyint DEFAULT '0' COMMENT '是否删除（逻辑删除标记）',
  `creator` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '创建者',
  `updater` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '更新者',
  `create_time` int DEFAULT NULL COMMENT '创建时间',
  `update_time` int DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;



-- 成绩表
CREATE TABLE `grade` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `student_id` int NOT NULL COMMENT '学生 id',
  `subject_id` int NOT NULL COMMENT '学科',
  `year` int NOT NULL COMMENT '考试年份',
  `score` decimal(5,2) NOT NULL COMMENT '成绩',
  `term` tinyint DEFAULT NULL COMMENT '上学期 1，下学期2',
  `exam_id` int NOT NULL COMMENT '试卷 ID',
  `is_delete` tinyint DEFAULT '0' COMMENT '是否删除（逻辑删除标记）',
  `creator` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '创建者',
  `updater` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '更新者',
  `create_time` int DEFAULT NULL COMMENT '创建时间',
  `update_time` int DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='学生成绩表';



-- 角色表
CREATE TABLE `role` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '角色名',
  `create_time` bigint NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` bigint NOT NULL DEFAULT '0' COMMENT '更新时间',
  `creator` varchar(100) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '创建人',
  `updater` varchar(100) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '更新人',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='角色表';



-- 学生表
CREATE TABLE `student` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '姓名',
  `gender` tinyint(1) DEFAULT NULL COMMENT '性别',
  `class_id` int NOT NULL,
  `phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT '' COMMENT '联系方式',
  `email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT '' COMMENT '邮箱',
  `passwd` varchar(255) COLLATE utf8mb4_bin DEFAULT '' COMMENT '登录密码',
  `is_delete` tinyint DEFAULT '0' COMMENT '是否删除（逻辑删除标记）',
  `creator` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '创建者',
  `updater` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '更新者',
  `create_time` int DEFAULT NULL COMMENT '创建时间',
  `update_time` int DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='学生信息表';


-- 学科表
CREATE TABLE `subject` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '学科',
  `description` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '介绍',
  `director_id` int DEFAULT NULL COMMENT '主任(关联教师 id)',
  `is_delete` tinyint NOT NULL DEFAULT '0' COMMENT '是否删除（逻辑删除标记）',
  `creator` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '创建者',
  `updater` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '更新者',
  `create_time` int DEFAULT NULL COMMENT '创建时间',
  `update_time` int DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='学科信息表';




-- 教师表
CREATE TABLE `teacher` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '教职人员ID',
  `name` varchar(100) NOT NULL COMMENT '姓名',
  `age` int NOT NULL DEFAULT '0' COMMENT '年龄',
  `gender` tinyint(1) NOT NULL DEFAULT '1' COMMENT '性别(1-男，2-女)',
  `subject_id` int DEFAULT NULL COMMENT '所教学科',
  `phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '联系方式',
  `email` varchar(100) DEFAULT NULL COMMENT '邮箱地址',
  `passwd` varchar(255) NOT NULL COMMENT '密码（已加密）',
  `level` int DEFAULT '0' COMMENT '教学年级',
  `is_delete` tinyint(1) DEFAULT '0' COMMENT '是否删除（逻辑删除标记）',
  `creator` varchar(100) DEFAULT NULL COMMENT '创建者',
  `updater` varchar(100) DEFAULT NULL COMMENT '更新者',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`),
  UNIQUE KEY `phone` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='教职人员表';


