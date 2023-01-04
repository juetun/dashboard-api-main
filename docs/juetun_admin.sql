/*
 Navicat Premium Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 50716
 Source Host           : localhost:3306
 Source Schema         : juetun_admin

 Target Server Type    : MySQL
 Target Server Version : 50716
 File Encoding         : 65001

 Date: 03/01/2023 21:55:46
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for admin_app
-- ----------------------------
DROP TABLE IF EXISTS `admin_app`;
CREATE TABLE `admin_app`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `unique_key` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '系统唯一KEY',
  `port` int(10) UNSIGNED NOT NULL DEFAULT 80 COMMENT '端口号',
  `hosts` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'host配置json字符串',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '系统名称',
  `desc` varchar(300) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '系统描述',
  `is_stop` tinyint(3) UNSIGNED NOT NULL DEFAULT 1 COMMENT '是否停用 1:使用 2:停用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 7 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_app
-- ----------------------------
INSERT INTO `admin_app` VALUES (1, 'api-user', 8093, '{\"dev\":\"dev-api.juetun.com\",\"local\":\"localhost-api.juetun.com\",\"pre\":\"pre-api.juetun.com\",\"release\":\"api.juetun.com\",\"test\":\"test-api.juetun.com\"}', '用户', '用户用户', 1, '2021-05-23 18:53:52', '2021-10-30 00:00:00', NULL);
INSERT INTO `admin_app` VALUES (2, 'api-upload', 8090, '{\"dev\":\"dev-api.juetun.com\",\"pre\":\"pre-api.juetun.com\",\"release\":\"api.juetun.com\",\"test\":\"test-api.juetun.com\"}', '上传', '', 2, '2021-05-23 22:07:10', '2021-06-01 22:20:42', NULL);
INSERT INTO `admin_app` VALUES (3, 'admin-car', 8091, '{\"dev\":\"dev-api.juetun.com\",\"pre\":\"pre-api.juetun.com\",\"release\":\"api.juetun.com\",\"test\":\"test-api.juetun.com\"}', '汽车', '', 1, '2021-05-27 00:20:27', '2021-05-28 23:45:32', NULL);
INSERT INTO `admin_app` VALUES (4, 'admin-scheduler', 8193, '{\"dev\":\"dev-api.juetun.com\",\"pre\":\"pre-api.juetun.com\",\"release\":\"api.juetun.com\",\"test\":\"test-api.juetun.com\"}', '定时任务', '定时任务调度器', 1, '2021-05-27 00:27:20', '2021-05-28 23:45:38', NULL);
INSERT INTO `admin_app` VALUES (5, 'web', 8092, '{\"dev\":\"dev-api.juetun.com\",\"pre\":\"pre-api.juetun.com\",\"release\":\"api.juetun.com\",\"test\":\"test-api.juetun.com\"}', '网站界面', '', 1, '2021-05-27 00:27:58', '2021-05-28 23:45:19', NULL);
INSERT INTO `admin_app` VALUES (6, 'admin-main', 8089, '{\"dev\":\"dev-api.juetun.com\",\"pre\":\"pre-api.juetun.com\",\"release\":\"api.juetun.com\",\"test\":\"test-api.juetun.com\"}', '后台管理', '后台管理代码', 1, '2021-05-27 23:59:05', '2021-05-28 23:45:23', NULL);

-- ----------------------------
-- Table structure for admin_base_sys
-- ----------------------------
DROP TABLE IF EXISTS `admin_base_sys`;
CREATE TABLE `admin_base_sys`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` datetime NULL DEFAULT NULL,
  `theme` int(10) NULL DEFAULT 0 COMMENT '主题',
  `title` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '标题',
  `keywords` varchar(300) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '网站关键词',
  `description` varchar(5000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '网站描述',
  `record_number` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT '' COMMENT '网站备案号',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '系统设置' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_base_sys
-- ----------------------------

-- ----------------------------
-- Table structure for admin_group
-- ----------------------------
DROP TABLE IF EXISTS `admin_group`;
CREATE TABLE `admin_group`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '组名',
  `parent_id` int(11) NOT NULL DEFAULT 0 COMMENT '上级用户组',
  `group_code` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '机构号',
  `is_super_admin` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否是超级管理员 0:否    1：是',
  `is_admin_group` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否为管理员组',
  `last_child_group_code` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 11 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_group
-- ----------------------------
INSERT INTO `admin_group` VALUES (1, '系统基础组', 0, '2MLKJIHGF', 1, 1, 'EDC', '2021-05-15 00:18:49', '2021-12-18 22:15:05', NULL);
INSERT INTO `admin_group` VALUES (2, '财务', 1, '2MLKJIHGFEDC', 0, 1, '2MLKJIHGFEDCEDD', '2021-06-01 22:54:31', '2021-12-23 21:54:23', NULL);
INSERT INTO `admin_group` VALUES (5, '财务组1', 2, '2MLKJIHGFEDCEDD', 0, 1, '', '2021-12-12 21:11:01', '2021-12-18 22:34:06', NULL);
INSERT INTO `admin_group` VALUES (6, '财务组2', 2, '2MLKJIHGFEDCEDC', 0, 1, '', '2021-12-12 21:11:07', '2021-12-18 22:47:25', NULL);

-- ----------------------------
-- Table structure for admin_help_document
-- ----------------------------
DROP TABLE IF EXISTS `admin_help_document`;
CREATE TABLE `admin_help_document`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `p_key` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '唯一的key',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '文档内容',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_pk`(`p_key`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '帮助文档' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_help_document
-- ----------------------------
INSERT INTO `admin_help_document` VALUES (1, 'user_agreement', '', '2022-07-11 20:46:19', '2022-07-11 20:46:19', NULL);

-- ----------------------------
-- Table structure for admin_help_document_relate
-- ----------------------------
DROP TABLE IF EXISTS `admin_help_document_relate`;
CREATE TABLE `admin_help_document_relate`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `biz_code` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '业务场景码',
  `display` tinyint(4) NOT NULL DEFAULT 1 COMMENT '是否在列表页展示 1-展示 0-不展示',
  `parent_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '上级文档ID',
  `label` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '名称',
  `is_leaf_node` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否叶子节点',
  `doc_key` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '文档唯一的key',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_pk`(`biz_code`, `parent_id`, `label`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '帮助文档关系描述' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_help_document_relate
-- ----------------------------
INSERT INTO `admin_help_document_relate` VALUES (1, 'help', 1, 0, '帮助', 0, '', '2022-07-11 14:48:32', '2022-07-11 14:48:32', NULL);
INSERT INTO `admin_help_document_relate` VALUES (2, 'help', 1, 1, '用户', 0, '', '2022-07-11 14:48:39', '2022-07-11 14:48:39', NULL);
INSERT INTO `admin_help_document_relate` VALUES (3, 'help', 1, 1, '电商', 0, '', '2022-07-11 14:48:57', '2022-07-11 14:48:57', NULL);
INSERT INTO `admin_help_document_relate` VALUES (4, 'help', 1, 2, '用户协议', 1, 'user_agreement', '2022-07-11 14:49:09', '2022-07-11 14:49:09', NULL);

-- ----------------------------
-- Table structure for admin_import
-- ----------------------------
DROP TABLE IF EXISTS `admin_import`;
CREATE TABLE `admin_import`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `permit_key` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `app_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '应用英文名',
  `app_version` varchar(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '1.0' COMMENT '引用版本号',
  `url_path` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'URL路径后缀',
  `request_method` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '请求方法 (GET)',
  `sort_value` int(11) NOT NULL DEFAULT 0 COMMENT '排序值 值越大越靠前',
  `default_open` tinyint(3) UNSIGNED NOT NULL DEFAULT 1 COMMENT '1:开启权限 2：不开启权限',
  `need_login` tinyint(3) UNSIGNED NOT NULL DEFAULT 1 COMMENT '是否需要登录 2:不需要;1:需要',
  `need_sign` tinyint(3) UNSIGNED NOT NULL DEFAULT 1 COMMENT '是否需要签名验证 1:需要 ;2:不需要',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `app_name`(`app_name`, `url_path`, `request_method`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 21 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '接口表数据' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_import
-- ----------------------------
INSERT INTO `admin_import` VALUES (1, 'M0Xqxvnx1J9B', 'admin-main', '1.0', 'admin/permit/menu', 'GET', 1, 1, 2, 2, '2021-09-21 22:12:26', '2021-06-01 22:45:05', NULL);
INSERT INTO `admin_import` VALUES (3, '7OWZR3qpL9a0', 'admin-main', '1.0', 'admin/export/list', 'GET', 2, 1, 1, 2, '2021-09-21 21:49:16', '2021-06-03 22:55:53', NULL);
INSERT INTO `admin_import` VALUES (5, 'Y3VdpOqR20Nn', 'admin-main', '1.0', 'admin/permit/admin_user', 'POST,GET', 2, 1, 1, 1, '2021-06-15 23:47:10', '2021-06-03 22:58:27', NULL);
INSERT INTO `admin_import` VALUES (6, 'qyE757npKLwb', 'admin-main', '1.0', 'admin/permit/admin_user_group_release', 'POST', 2, 1, 1, 1, '2021-06-10 22:16:07', '2021-06-03 23:00:16', NULL);
INSERT INTO `admin_import` VALUES (7, 'dDqz5z4xjPGo', 'admin-main', '1.0', 'admin/permit/admin_user_group_add', 'POST', 0, 1, 1, 1, '2021-06-15 23:47:07', '2021-06-03 23:01:03', NULL);
INSERT INTO `admin_import` VALUES (8, 'XE1959m5jg2l', 'admin-main', '1.0', 'admin/permit/admin_user_delete', 'POST', 2, 1, 1, 1, '2021-12-07 00:03:05', '2021-06-03 23:01:35', NULL);
INSERT INTO `admin_import` VALUES (9, 'AYq45PZxj6aG', 'admin-main', '1.0', 'admin/permit/admin_group', 'GET', 7, 1, 1, 1, '2021-06-10 22:16:09', '2021-06-03 23:02:16', NULL);
INSERT INTO `admin_import` VALUES (10, 'BejyRD1xG0L3', 'admin-main', '1.0', 'admin/console/home', 'GET', 0, 1, 1, 1, '2021-06-09 23:51:40', '2021-06-09 23:51:26', NULL);
INSERT INTO `admin_import` VALUES (11, 'Pq3d5rV5k2Nj', 'admin-main', '1.0', 'admin/console/cate', 'GET,POST', 0, 1, 1, 1, '2021-06-15 23:47:03', '2021-06-09 23:52:42', NULL);
INSERT INTO `admin_import` VALUES (12, 'K4Nd5mBR9E8n', 'admin-main', '1.0', 'admin/console/cate/edit/:id', 'GET', 0, 1, 1, 1, '2021-06-15 23:47:02', '2021-06-09 23:53:07', NULL);
INSERT INTO `admin_import` VALUES (13, '87LVpwopYvEy', 'admin-main', '1.0', 'admin/console/cate/:id', 'PUT', 0, 1, 2, 2, '2021-12-18 23:16:41', '2021-06-09 23:53:25', NULL);
INSERT INTO `admin_import` VALUES (14, '87LVpwopYvEZ', 'admin-main', '1.0', 'admin/permit/get_app_config', 'GET', 0, 1, 1, 1, '2021-09-21 21:10:45', '2021-09-05 11:52:34', NULL);
INSERT INTO `admin_import` VALUES (16, 'jQXZxKZRNVr8', 'admin-main', '1.0', 'admin/export/init', 'POST', 0, 1, 1, 1, '2021-12-22 23:05:36', '2021-09-21 21:49:46', NULL);
INSERT INTO `admin_import` VALUES (17, 'YGNWxbnxELln', 'admin-main', '1.0', 'admin/export/progress', 'GET', 0, 1, 1, 1, '2021-09-21 23:12:50', '2021-09-21 21:50:11', NULL);
INSERT INTO `admin_import` VALUES (18, '64K05g9xabWw', 'admin-main', '1.0', 'admin/export/cancel', 'GET', 0, 1, 1, 1, '2021-09-21 23:12:54', '2021-09-21 21:51:19', NULL);
INSERT INTO `admin_import` VALUES (20, 'Q6bkRj2RmVZq', 'admin-main', '1.0', 'admin/console/tag', 'GET,POST', 0, 1, 1, 1, '2021-12-18 23:16:26', '2021-09-21 22:24:26', NULL);

-- ----------------------------
-- Table structure for admin_menu
-- ----------------------------
DROP TABLE IF EXISTS `admin_menu`;
CREATE TABLE `admin_menu`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `permit_key` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `module` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `parent_id` int(11) NOT NULL DEFAULT 1 COMMENT '上级权限ID 1：表示不属于任何一级',
  `label` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '显示名称 或接口备注',
  `icon` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `hide_in_menu` tinyint(1) NOT NULL DEFAULT 1 COMMENT '菜单是否显示 1：显示 ，2：不显示',
  `is_home_page` tinyint(4) NOT NULL DEFAULT 2 COMMENT '是否首页',
  `manage_import_permit` tinyint(3) UNSIGNED NOT NULL DEFAULT 1 COMMENT '是否有管理接口的权限 1:管理 2：不管理',
  `domain` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `url_path` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'URL路径后缀',
  `other_value` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '其他参数',
  `sort_value` int(11) NOT NULL DEFAULT 0 COMMENT '排序值 值越大越靠前',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `permit_key`(`permit_key`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 89 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '菜单表数据' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_menu
-- ----------------------------
INSERT INTO `admin_menu` VALUES (1, 'backend', 'platform', -1, '汽车', 'md-car', 0, 2, 0, '', '', '', 1000, '2021-05-21 00:50:18', '2021-02-28 21:18:21', NULL);
INSERT INTO `admin_menu` VALUES (2, 'style_admin_model', 'backend', 9, '车型车系', 'md-car', 2, 2, 2, '', '', '', 1, '2021-06-27 20:02:26', '2020-09-27 23:28:57', NULL);
INSERT INTO `admin_menu` VALUES (4, 'travel_cricle', 'backend', 1, '游记圈子', 'ios-book', 2, 2, 2, '', '', '', 70000000, '2021-04-16 23:41:01', '2020-12-09 23:27:10', NULL);
INSERT INTO `admin_menu` VALUES (9, 'car_data', 'backend', 1, '汽车数据', 'md-car', 2, 2, 2, '', '', '', 80000000, '2021-09-21 19:42:47', '2020-09-27 23:28:57', NULL);
INSERT INTO `admin_menu` VALUES (10, 'style_manager', 'backend', 2, '车系管理', 'md-train', 2, 2, 1, '', '', '', 0, '2021-03-10 00:42:01', '2021-03-06 20:19:20', NULL);
INSERT INTO `admin_menu` VALUES (11, 'model_manager', 'backend', 2, '车型管理', 'ios-car', 2, 2, 1, '', '', '', 0, '2021-03-10 00:42:41', '2021-03-06 20:22:00', NULL);
INSERT INTO `admin_menu` VALUES (12, 'system_manager', 'system', 49, '系统管理', 'md-settings', 2, 2, 2, '', '', '', 0, '2021-06-02 22:27:51', '2021-03-07 19:29:46', NULL);
INSERT INTO `admin_menu` VALUES (13, 'permit_index', 'system', 12, '权限', 'ios-keypad', 2, 2, 2, '', '', '', 0, '2021-04-17 00:00:42', '2021-03-07 19:32:04', NULL);
INSERT INTO `admin_menu` VALUES (14, 'user_list', 'system', 58, '系统用户', 'md-contact', 2, 2, 1, '', '', '', 0, '2021-04-17 00:01:02', '2021-03-07 19:32:35', NULL);
INSERT INTO `admin_menu` VALUES (15, 'console_group_list', 'system', 13, '管理员组', 'md-bulb', 2, 2, 1, '', '', '', 100, '2021-12-06 22:51:47', '2021-03-07 22:37:54', NULL);
INSERT INTO `admin_menu` VALUES (16, 'admin_group_permit', 'system', 13, '权限管理', 'ios-cash-outline', 2, 2, 1, '', '', '', 80, '2021-12-06 22:51:56', '2021-03-07 22:38:36', NULL);
INSERT INTO `admin_menu` VALUES (17, 'admin_group_apps', 'system', 52, '服务管理', 'ios-briefcase', 2, 2, 1, '', '', '', 0, '2021-05-15 21:59:21', '2021-03-07 22:47:12', NULL);
INSERT INTO `admin_menu` VALUES (18, 'system_index', 'system', 52, '系统设置', 'md-settings', 2, 2, 1, '', '', '', 0, '2021-05-15 21:59:02', '2021-03-07 23:35:41', NULL);
INSERT INTO `admin_menu` VALUES (19, 'car_media', 'backend', 9, '车型素材', '', 2, 2, 1, '', '', '', 0, '2021-03-08 00:25:24', '2021-03-08 00:25:24', NULL);
INSERT INTO `admin_menu` VALUES (28, 'home_index', 'backend', 1, '首页', 'md-home', 1, 1, 1, '', '', '', 90000001, '2021-05-12 22:21:05', '2021-04-15 23:21:25', NULL);
INSERT INTO `admin_menu` VALUES (29, 'Ky6wWjPM83Rm', 'backend', 1, '广告系统', 'md-attach', 2, 2, 2, '', '', '', 60000000, '2021-04-17 00:00:01', '2021-04-15 23:47:25', NULL);
INSERT INTO `admin_menu` VALUES (30, 'RL7bWm5WGV24', 'backend', 1, '用户管理', 'ios-man', 2, 2, 2, '', '', '', 50000000, '2021-04-17 00:00:04', '2021-04-15 23:49:36', NULL);
INSERT INTO `admin_menu` VALUES (31, 'site_recommend', 'backend', 29, '首页推荐', 'logo-angular', 2, 2, 1, '', '', '', 0, '2021-04-16 00:12:26', '2021-04-16 00:08:39', NULL);
INSERT INTO `admin_menu` VALUES (32, 'advertise_list', 'backend', 29, '广告管理', 'md-boat', 2, 2, 1, '', '', '', 0, '2021-04-16 00:12:45', '2021-04-16 00:09:17', NULL);
INSERT INTO `admin_menu` VALUES (33, 'user_index', 'backend', 30, '用户管理', 'md-contact', 2, 2, 1, '', '', '', 0, '2021-04-16 00:17:51', '2021-04-16 00:17:51', NULL);
INSERT INTO `admin_menu` VALUES (34, 'user_permit', 'backend', 30, '用户权限', 'md-bulb', 2, 2, 1, '', '', '', 0, '2021-04-16 00:18:44', '2021-04-16 00:18:44', NULL);
INSERT INTO `admin_menu` VALUES (35, 'post_list', 'backend', 4, '游记文章', 'ios-book', 2, 2, 1, '', '', '', 0, '2021-04-16 00:19:58', '2021-04-16 00:19:58', NULL);
INSERT INTO `admin_menu` VALUES (36, 'post_trash', 'backend', 35, '回收站', 'md-trash', 1, 2, 1, '', '', '', 0, '2021-04-16 00:53:57', '2021-04-16 00:20:30', NULL);
INSERT INTO `admin_menu` VALUES (37, 'cate_list', 'backend', 4, '分类列表', 'md-locate', 2, 2, 1, '', '', '', 0, '2021-04-16 00:21:22', '2021-04-16 00:21:22', NULL);
INSERT INTO `admin_menu` VALUES (38, 'data_tag_list', 'backend', 4, '标签列表', 'md-share', 2, 2, 1, '', '', '', 0, '2021-04-16 00:23:50', '2021-04-16 00:23:50', NULL);
INSERT INTO `admin_menu` VALUES (39, 'model_arg_edit', 'backend', 11, '车型参数编辑', 'ios-barcode', 1, 2, 1, '', '', '', 0, '2021-05-15 23:22:06', '2021-04-16 00:27:17', NULL);
INSERT INTO `admin_menu` VALUES (40, 'YmG4zXlzLp5l', 'backend', 1, '公共接口', '', 1, 2, 1, '', '', '', 90000000, '2021-04-16 23:19:55', '2021-04-16 23:19:49', NULL);
INSERT INTO `admin_menu` VALUES (46, 'user', 'platform', -1, '用户', 'md-contact', 0, 2, 0, '', '', '', 20000, '2021-05-21 00:50:30', '2021-05-15 20:50:10', NULL);
INSERT INTO `admin_menu` VALUES (47, 'boOQN3agzEKk', 'user', 46, '首页', 'md-home', 1, 1, 1, '', '', '', 90000001, '2021-05-15 20:50:10', '2021-05-15 20:50:10', NULL);
INSERT INTO `admin_menu` VALUES (48, 'B1epMQXBzrKl', 'user', 46, '公共接口', '', 1, 2, 1, '', '', '', 0, '2021-05-15 20:50:10', '2021-05-15 20:50:10', NULL);
INSERT INTO `admin_menu` VALUES (49, 'system', 'platform', -1, '系统', 'md-settings', 0, 2, 0, '', '', '', 30000, '2021-05-21 00:19:20', '2021-05-15 20:51:24', NULL);
INSERT INTO `admin_menu` VALUES (50, 'system_home_index', 'system', 49, '首页', 'md-home', 2, 1, 1, '', '', '', 90000000, '2021-05-15 20:51:24', '2021-05-15 20:51:24', NULL);
INSERT INTO `admin_menu` VALUES (51, '7lDBNwj5MrQ9', 'system', 49, '公共接口', '', 1, 2, 1, '', '', '', 0, '2021-05-15 20:51:24', '2021-05-15 20:51:24', NULL);
INSERT INTO `admin_menu` VALUES (52, 'xBZANrP7zpK8', 'system', 12, '系统控制', 'md-build', 2, 2, 2, '', '', '', 0, '2021-06-02 22:28:39', '2021-05-15 21:58:35', NULL);
INSERT INTO `admin_menu` VALUES (53, 'crontab', 'system', 12, '定时任务', 'md-hand', 2, 2, 1, '', '', '', 0, '2021-06-02 22:28:56', '2021-05-22 11:32:49', NULL);
INSERT INTO `admin_menu` VALUES (54, 'crontab_edit', 'system', 53, '定时任务编辑', '', 1, 2, 1, '', '', '', 0, '2021-05-22 12:04:02', '2021-05-22 11:36:13', NULL);
INSERT INTO `admin_menu` VALUES (55, 'group_permit_set', 'system', 15, '组权限设置', '', 1, 2, 1, '', '', '', 0, '2021-05-29 00:42:49', '2021-05-29 00:42:41', NULL);
INSERT INTO `admin_menu` VALUES (56, 'admin_app_edit', 'system', 17, '服务编辑', '', 1, 2, 1, '', '', '', 0, '2021-05-30 00:26:35', '2021-05-30 00:26:35', NULL);
INSERT INTO `admin_menu` VALUES (57, 'system_import', 'system', 13, '接口列表', 'md-book', 2, 2, 1, '', '', '', 0, '2021-06-03 00:01:43', '2021-06-03 00:00:54', NULL);
INSERT INTO `admin_menu` VALUES (58, 'rA29zJRvzq3R', 'system', 49, '用户管理', 'ios-contacts', 2, 2, 1, '', '', '', 70000000, '2021-10-24 23:33:24', '2021-10-24 23:32:13', NULL);
INSERT INTO `admin_menu` VALUES (59, 'user_tag', 'system', 58, '用户标签', 'md-apps', 2, 2, 1, '', '', '', 0, '2021-11-07 18:15:14', '2021-11-07 18:13:28', NULL);
INSERT INTO `admin_menu` VALUES (60, 'tag_config', 'system', 61, '标签配置', 'ios-bookmarks', 2, 2, 1, '', '', '', 0, '2021-11-07 18:31:32', '2021-11-07 18:29:43', NULL);
INSERT INTO `admin_menu` VALUES (61, 'tag_manager', 'system', 12, '标签管理', 'ios-bookmarks', 2, 2, 1, '', '', '', 1000, '2021-11-07 22:17:08', '2021-11-07 22:14:07', NULL);
INSERT INTO `admin_menu` VALUES (62, 'tag_group', 'system', 61, '标签组', 'md-bookmarks', 2, 2, 1, '', '', '', 0, '2021-11-07 22:16:26', '2021-11-07 22:15:53', NULL);
INSERT INTO `admin_menu` VALUES (63, 'admin_user', 'system', 13, '管理员', 'ios-contacts', 2, 2, 1, '', '', '', 90, '2021-12-06 22:53:04', '2021-12-06 22:50:24', NULL);
INSERT INTO `admin_menu` VALUES (64, 'papers', 'system', 61, '资质管理', 'md-paper', 2, 0, 1, '', '', '', 0, '2022-03-03 10:07:45', '2022-03-03 10:06:23', NULL);
INSERT INTO `admin_menu` VALUES (65, 'papers_list', 'system', 64, '资质列表', '', 2, 0, 1, '', '', '', 0, '2022-03-03 10:09:04', '2022-03-03 10:08:07', NULL);
INSERT INTO `admin_menu` VALUES (66, 'papers_relate', 'system', 64, '资质关系', '', 2, 0, 1, '', '', '', 0, '2022-03-03 10:09:14', '2022-03-03 10:08:34', NULL);
INSERT INTO `admin_menu` VALUES (67, 'papers_relate_config', 'system', 66, '资质关系编辑', '', 1, 2, 1, '', '', '', 0, '2022-03-06 17:36:50', '2022-03-06 17:36:50', NULL);
INSERT INTO `admin_menu` VALUES (68, 'static_data', '', 1, '运营统计', '', 1, 2, 1, '', '', '', 0, '2023-01-03 17:19:35', '2023-01-03 17:19:35', NULL);
INSERT INTO `admin_menu` VALUES (69, 'help_document', 'system', 58, '帮助文档', 'md-help', 2, 2, 1, '', '', '', 0, '2023-01-03 17:22:33', '2023-01-03 17:22:33', NULL);
INSERT INTO `admin_menu` VALUES (70, 'help_edit', 'system', 69, '帮助文档编辑', '', 1, 2, 1, '', '', '', 0, '2023-01-03 17:23:07', '2023-01-03 17:23:07', NULL);
INSERT INTO `admin_menu` VALUES (71, 'ad', 'system', 12, '广告', 'md-globe', 2, 0, 1, '', '', '', 1100, '2023-01-03 21:40:21', '2023-01-03 17:39:08', NULL);
INSERT INTO `admin_menu` VALUES (72, 'ad_scene', 'system', 71, '广告场景', 'md-ionitron', 2, 0, 1, '', '', '', 0, '2023-01-03 17:43:39', '2023-01-03 17:41:30', NULL);
INSERT INTO `admin_menu` VALUES (73, 'ad_list', 'system', 71, '广告数据', 'md-microphone', 2, 0, 1, '', '', '', 0, '2023-01-03 17:43:50', '2023-01-03 17:42:44', NULL);
INSERT INTO `admin_menu` VALUES (74, 'ad_scene_edit', 'system', 72, '广告场景编辑', '', 1, 2, 1, '', '', '', 0, '2023-01-03 17:52:32', '2023-01-03 17:52:32', NULL);
INSERT INTO `admin_menu` VALUES (75, 'ad_data_edit', 'system', 73, '广告数据编辑', '', 1, 2, 1, '', '', '', 0, '2023-01-03 21:31:00', '2023-01-03 21:31:00', NULL);
INSERT INTO `admin_menu` VALUES (76, 'mall_manager', 'system', 49, '商品管理', 'ios-contacts', 2, 0, 1, '', '', '', 90000000, '2023-01-03 21:50:55', '2023-01-03 21:45:04', NULL);
INSERT INTO `admin_menu` VALUES (78, 'spu_list', 'system', 76, 'SPU管理', '', 2, 0, 1, '', '', '', 10000, '2023-01-03 21:47:51', '2023-01-03 21:46:08', NULL);
INSERT INTO `admin_menu` VALUES (79, 'freight', 'system', 76, '物流管理', '', 2, 2, 1, '', '', '', 0, '2023-01-03 21:46:35', '2023-01-03 21:46:35', NULL);
INSERT INTO `admin_menu` VALUES (80, 'gift_list', 'system', 76, '赠品管理', '', 2, 0, 1, '', '', '', 8000, '2023-01-03 21:48:14', '2023-01-03 21:47:03', NULL);
INSERT INTO `admin_menu` VALUES (81, 'sku_list', 'system', 76, 'SKU管理', '', 2, 0, 1, '', '', '', 9000, '2023-01-03 21:48:02', '2023-01-03 21:47:38', NULL);
INSERT INTO `admin_menu` VALUES (82, 'activity', 'system', 83, '平台活动', '', 2, 0, 1, '', '', '', 0, '2023-01-03 21:50:37', '2023-01-03 21:49:04', NULL);
INSERT INTO `admin_menu` VALUES (83, 'marketing', 'system', 49, '营销管理', 'ios-contacts', 2, 0, 1, '', '', '', 80000000, '2023-01-03 21:51:02', '2023-01-03 21:50:00', NULL);
INSERT INTO `admin_menu` VALUES (84, 'coupon_list', 'system', 83, '优惠券列表', '', 2, 2, 1, '', '', '', 0, '2023-01-03 21:51:25', '2023-01-03 21:51:25', NULL);
INSERT INTO `admin_menu` VALUES (85, 'order_manager', 'system', 49, '订单管理', 'ios-contacts', 2, 0, 1, '', '', '', 85000000, '2023-01-03 21:52:35', '2023-01-03 21:52:22', NULL);
INSERT INTO `admin_menu` VALUES (86, 'shop_order', 'system', 85, '店铺管理', '', 2, 2, 1, '', '', '', 0, '2023-01-03 21:52:57', '2023-01-03 21:52:57', NULL);
INSERT INTO `admin_menu` VALUES (87, 'user_order', 'system', 85, '用户订单', '', 2, 2, 1, '', '', '', 0, '2023-01-03 21:53:10', '2023-01-03 21:53:10', NULL);
INSERT INTO `admin_menu` VALUES (88, 'order_static', 'system', 83, '订单统计', '', 2, 2, 1, '', '', '', 0, '2023-01-03 21:54:22', '2023-01-03 21:54:22', NULL);

-- ----------------------------
-- Table structure for admin_menu_import
-- ----------------------------
DROP TABLE IF EXISTS `admin_menu_import`;
CREATE TABLE `admin_menu_import`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `menu_id` int(11) NOT NULL DEFAULT 0,
  `menu_module` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '菜单所属模块',
  `import_id` int(11) NOT NULL DEFAULT 0,
  `default_open` tinyint(1) NOT NULL DEFAULT 2 COMMENT '默认开启菜单就开启本接口权限 1-开启 不为1-不开启（ 默认值2不开启）',
  `import_app_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '接口应用路径',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `menu_id`(`menu_id`, `import_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 189 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '菜单关联接口' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_menu_import
-- ----------------------------
INSERT INTO `admin_menu_import` VALUES (1, 50, 'system', 1, 2, 'admin-main', '2021-06-06 00:49:13', '2021-06-27 10:51:05', NULL);
INSERT INTO `admin_menu_import` VALUES (2, 50, 'system', 3, 2, 'admin-main', '2021-06-22 00:06:03', '2021-06-27 10:51:05', NULL);
INSERT INTO `admin_menu_import` VALUES (8, 50, 'system', 4, 2, 'admin-main', '2021-06-22 00:12:51', '2021-06-27 10:51:05', NULL);
INSERT INTO `admin_menu_import` VALUES (10, 50, 'system', 5, 2, 'admin-main', '2021-06-22 00:12:54', '2021-06-27 10:51:05', NULL);
INSERT INTO `admin_menu_import` VALUES (14, 50, 'system', 11, 2, 'admin-main', '2021-06-22 00:13:01', '2021-06-27 10:51:16', NULL);
INSERT INTO `admin_menu_import` VALUES (15, 50, 'system', 10, 2, 'admin-main', '2021-06-22 00:13:01', '2021-06-27 10:51:16', NULL);
INSERT INTO `admin_menu_import` VALUES (16, 50, 'system', 9, 2, 'admin-main', '2021-06-22 00:13:03', '2021-06-27 10:51:16', NULL);
INSERT INTO `admin_menu_import` VALUES (17, 50, 'system', 8, 2, 'admin-main', '2021-06-22 00:13:03', '2021-06-27 10:51:16', NULL);
INSERT INTO `admin_menu_import` VALUES (18, 50, 'system', 7, 2, 'admin-main', '2021-06-22 00:13:05', '2021-06-27 10:51:16', NULL);
INSERT INTO `admin_menu_import` VALUES (34, 50, 'system', 6, 2, 'admin-main', '2021-06-26 22:34:31', '2021-06-27 10:51:05', NULL);
INSERT INTO `admin_menu_import` VALUES (62, 50, 'system', 12, 2, 'admin-main', '2021-06-27 10:51:20', '2021-06-27 10:51:20', NULL);
INSERT INTO `admin_menu_import` VALUES (63, 50, 'system', 13, 2, 'admin-main', '2021-06-27 10:51:20', '2021-06-27 10:51:20', NULL);
INSERT INTO `admin_menu_import` VALUES (64, 14, 'system', 1, 2, 'admin-main', '2021-06-27 20:21:59', '2021-06-27 20:22:05', '2021-06-27 20:22:06');
INSERT INTO `admin_menu_import` VALUES (65, 14, 'system', 3, 2, 'admin-main', '2021-06-27 20:22:00', '2021-06-27 20:22:06', '2021-06-27 20:22:07');
INSERT INTO `admin_menu_import` VALUES (66, 14, 'system', 4, 2, 'admin-main', '2021-06-27 20:22:00', '2021-06-27 20:22:08', '2021-06-27 20:22:08');
INSERT INTO `admin_menu_import` VALUES (67, 14, 'system', 5, 2, 'admin-main', '2021-06-27 20:22:01', '2021-06-27 20:22:08', '2021-06-27 20:22:09');
INSERT INTO `admin_menu_import` VALUES (72, 11, 'backend', 8, 2, 'admin-main', '2021-10-31 02:11:34', '2021-10-31 02:12:09', '2021-10-31 02:12:10');
INSERT INTO `admin_menu_import` VALUES (80, 28, 'backend', 1, 2, 'admin-main', '2021-11-20 22:44:56', '2021-12-13 00:02:56', NULL);
INSERT INTO `admin_menu_import` VALUES (81, 58, 'system', 1, 2, 'admin-main', '2021-11-20 23:15:31', '2021-11-20 23:15:31', NULL);
INSERT INTO `admin_menu_import` VALUES (82, 58, 'system', 3, 2, 'admin-main', '2021-11-20 23:15:33', '2021-11-20 23:15:33', NULL);
INSERT INTO `admin_menu_import` VALUES (83, 58, 'system', 5, 2, 'admin-main', '2021-11-21 21:09:32', '2021-11-21 21:09:32', NULL);
INSERT INTO `admin_menu_import` VALUES (84, 58, 'system', 6, 2, 'admin-main', '2021-11-21 21:09:33', '2021-11-21 21:09:33', NULL);
INSERT INTO `admin_menu_import` VALUES (85, 58, 'system', 7, 1, 'admin-main', '2021-11-21 21:09:34', '2021-11-22 11:58:48', '2021-11-22 11:58:48');
INSERT INTO `admin_menu_import` VALUES (86, 14, 'system', 6, 2, 'admin-main', '2021-11-21 21:09:37', '2021-11-21 21:09:37', NULL);
INSERT INTO `admin_menu_import` VALUES (87, 14, 'system', 7, 2, 'admin-main', '2021-11-21 21:09:38', '2021-11-21 21:09:38', NULL);
INSERT INTO `admin_menu_import` VALUES (98, 10, 'backend', 1, 2, 'admin-main', '2021-12-02 23:20:40', '2021-12-17 22:48:46', NULL);
INSERT INTO `admin_menu_import` VALUES (101, 10, 'backend', 8, 1, 'admin-main', '2021-12-02 23:21:02', '2021-12-02 23:21:02', NULL);
INSERT INTO `admin_menu_import` VALUES (125, 28, 'backend', 3, 2, 'admin-main', '2021-12-13 00:01:00', '2021-12-13 00:02:57', NULL);
INSERT INTO `admin_menu_import` VALUES (128, 28, 'backend', 5, 2, 'admin-main', '2021-12-13 00:01:03', '2021-12-13 00:02:58', NULL);
INSERT INTO `admin_menu_import` VALUES (130, 28, 'backend', 6, 2, 'admin-main', '2021-12-13 00:01:05', '2021-12-13 00:02:59', NULL);
INSERT INTO `admin_menu_import` VALUES (132, 28, 'backend', 7, 2, 'admin-main', '2021-12-13 00:01:07', '2021-12-13 00:03:00', NULL);
INSERT INTO `admin_menu_import` VALUES (188, 40, 'backend', 1, 1, 'admin-main', '2021-12-18 23:34:42', '2021-12-18 23:34:42', NULL);

-- ----------------------------
-- Table structure for admin_user
-- ----------------------------
DROP TABLE IF EXISTS `admin_user`;
CREATE TABLE `admin_user`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_hid` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '用户ID',
  `real_name` varchar(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '姓名',
  `mobile` varchar(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '手机号',
  `email` varbinary(100) NOT NULL DEFAULT '' COMMENT '邮箱',
  `flag_admin` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否超级管理员',
  `can_not_use` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否可用0-可用 1-不可用',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '表示未删除',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `user_hid`(`user_hid`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '后台用户表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_user
-- ----------------------------
INSERT INTO `admin_user` VALUES (1, '1', '长江', '15108352617', '', 1, 0, '2022-01-02 20:10:54', '2022-01-02 20:10:54', NULL);
INSERT INTO `admin_user` VALUES (2, '21', '', '13500000020', '', 0, 0, '2022-01-11 23:19:54', '2022-01-11 23:21:46', NULL);

-- ----------------------------
-- Table structure for admin_user_group
-- ----------------------------
DROP TABLE IF EXISTS `admin_user_group`;
CREATE TABLE `admin_user_group`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `group_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '组ID',
  `user_hid` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户ID',
  `is_super_admin` tinyint(3) UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否是超级管理员 0-否,1-是',
  `is_admin_group` tinyint(3) UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否是后台管理员组 0-否,1-是',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_gid_hid`(`group_id`, `user_hid`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户组用户列表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_user_group
-- ----------------------------
INSERT INTO `admin_user_group` VALUES (1, 1, '1', 1, 1, '2022-01-02 20:11:30', '2022-01-02 20:11:30', NULL);
INSERT INTO `admin_user_group` VALUES (2, 6, '21', 0, 0, '2022-01-11 23:19:54', '2022-01-11 23:19:54', NULL);

-- ----------------------------
-- Table structure for admin_user_group_import
-- ----------------------------
DROP TABLE IF EXISTS `admin_user_group_import`;
CREATE TABLE `admin_user_group_import`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `group_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '管理员组ID',
  `app_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '应用名',
  `module` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '接口所属菜单所属系统',
  `menu_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '接口所属菜单界面',
  `import_permit_key` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '接口唯一Key',
  `default_open` tinyint(4) NOT NULL DEFAULT 0 COMMENT '0-默认不随菜单开启 1-随菜单开启权限',
  `import_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '接口ID',
  `menu_permit_key` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '菜单唯一Key',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_uid`(`group_id`, `menu_id`, `import_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 92 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户组所具备的权限' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_user_group_import
-- ----------------------------
INSERT INTO `admin_user_group_import` VALUES (79, 1, 'admin-main', 'system', 50, '87LVpwopYvEy', 0, 13, 'pv51WeL9N4oD', '2022-01-17 23:18:54', '2022-01-17 23:18:54', NULL);
INSERT INTO `admin_user_group_import` VALUES (90, 1, 'admin-main', 'backend', 40, 'M0Xqxvnx1J9B', 0, 1, 'YmG4zXlzLp5l', '2022-01-17 23:19:19', '2022-01-17 23:19:19', NULL);

-- ----------------------------
-- Table structure for admin_user_group_menu
-- ----------------------------
DROP TABLE IF EXISTS `admin_user_group_menu`;
CREATE TABLE `admin_user_group_menu`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `group_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '管理员组ID',
  `module` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '菜单所属系统',
  `menu_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '菜单ID',
  `menu_permit_key` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '菜单唯一Key',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_uid`(`group_id`, `module`, `menu_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 48 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户组所具备的权限' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_user_group_menu
-- ----------------------------
INSERT INTO `admin_user_group_menu` VALUES (39, 1, 'system', 50, 'pv51WeL9N4oD', '2022-01-17 23:18:54', '2022-01-17 23:18:54', NULL);
INSERT INTO `admin_user_group_menu` VALUES (40, 1, 'system', 49, 'system', '2022-01-17 23:18:54', '2022-01-17 23:18:54', NULL);
INSERT INTO `admin_user_group_menu` VALUES (41, 1, 'system', 51, '7lDBNwj5MrQ9', '2022-01-17 23:18:54', '2022-01-17 23:18:54', NULL);
INSERT INTO `admin_user_group_menu` VALUES (42, 1, 'user', 47, 'boOQN3agzEKk', '2022-01-17 23:19:17', '2022-01-17 23:19:17', NULL);
INSERT INTO `admin_user_group_menu` VALUES (43, 1, 'user', 46, 'user', '2022-01-17 23:19:17', '2022-01-17 23:19:17', NULL);
INSERT INTO `admin_user_group_menu` VALUES (44, 1, 'user', 48, 'B1epMQXBzrKl', '2022-01-17 23:19:17', '2022-01-17 23:19:17', NULL);
INSERT INTO `admin_user_group_menu` VALUES (45, 1, 'backend', 28, 'home_index', '2022-01-17 23:19:19', '2022-01-17 23:19:19', NULL);
INSERT INTO `admin_user_group_menu` VALUES (46, 1, 'backend', 1, 'backend', '2022-01-17 23:19:19', '2022-01-17 23:19:19', NULL);
INSERT INTO `admin_user_group_menu` VALUES (47, 1, 'backend', 40, 'YmG4zXlzLp5l', '2022-01-17 23:19:19', '2022-01-17 23:19:19', NULL);

-- ----------------------------
-- Table structure for recommend_ad_data
-- ----------------------------
DROP TABLE IF EXISTS `recommend_ad_data`;
CREATE TABLE `recommend_ad_data`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `scene_id` bigint(20) NULL DEFAULT 0,
  `data_type` int(6) NOT NULL DEFAULT 0,
  `data_id` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '数据ID',
  `status` tinyint(2) NOT NULL DEFAULT 1 COMMENT '状态',
  `weight` bigint(20) NULL DEFAULT 0 COMMENT '排序权重',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_pk`(`scene_id`, `data_type`, `data_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '广告数据表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of recommend_ad_data
-- ----------------------------

SET FOREIGN_KEY_CHECKS = 1;
