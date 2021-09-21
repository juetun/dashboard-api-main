/*
 Navicat Premium Data Transfer

 Source Server         : 本机器
 Source Server Type    : MySQL
 Source Server Version : 80023
 Source Host           : localhost:3306
 Source Schema         : juetun_admin

 Target Server Type    : MySQL
 Target Server Version : 80023
 File Encoding         : 65001

 Date: 21/09/2021 10:19:57
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for admin_app
-- ----------------------------
DROP TABLE IF EXISTS `admin_app`;
CREATE TABLE `admin_app` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `unique_key` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '系统唯一KEY',
  `port` int unsigned NOT NULL DEFAULT '80' COMMENT '端口号',
  `hosts` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'host配置json字符串',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '系统名称',
  `desc` varchar(300) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '系统描述',
  `is_stop` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '是否停用 1:使用 2:停用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of admin_app
-- ----------------------------
BEGIN;
INSERT INTO `admin_app` VALUES (1, 'api-user', 8093, '{\"dev\":\"dev-api.juetun.com\",\"pre\":\"pre-api.juetun.com\",\"release\":\"api.juetun.com\",\"test\":\"test-api.juetun.com\"}', '用户', '用户用户', 1, '2021-05-23 18:53:52', '2021-06-01 22:20:52', NULL);
INSERT INTO `admin_app` VALUES (2, 'api-upload', 8090, '{\"dev\":\"dev-api.juetun.com\",\"pre\":\"pre-api.juetun.com\",\"release\":\"api.juetun.com\",\"test\":\"test-api.juetun.com\"}', '上传', '', 2, '2021-05-23 22:07:10', '2021-06-01 22:20:42', NULL);
INSERT INTO `admin_app` VALUES (3, 'admin-car', 8091, '{\"dev\":\"dev-api.juetun.com\",\"pre\":\"pre-api.juetun.com\",\"release\":\"api.juetun.com\",\"test\":\"test-api.juetun.com\"}', '汽车', '', 1, '2021-05-27 00:20:27', '2021-05-28 23:45:32', NULL);
INSERT INTO `admin_app` VALUES (4, 'admin-scheduler', 8193, '{\"dev\":\"dev-api.juetun.com\",\"pre\":\"pre-api.juetun.com\",\"release\":\"api.juetun.com\",\"test\":\"test-api.juetun.com\"}', '定时任务', '定时任务调度器', 1, '2021-05-27 00:27:20', '2021-05-28 23:45:38', NULL);
INSERT INTO `admin_app` VALUES (5, 'web', 8092, '{\"dev\":\"dev-api.juetun.com\",\"pre\":\"pre-api.juetun.com\",\"release\":\"api.juetun.com\",\"test\":\"test-api.juetun.com\"}', '网站界面', '', 1, '2021-05-27 00:27:58', '2021-05-28 23:45:19', NULL);
INSERT INTO `admin_app` VALUES (6, 'admin-main', 8089, '{\"dev\":\"dev-api.juetun.com\",\"pre\":\"pre-api.juetun.com\",\"release\":\"api.juetun.com\",\"test\":\"test-api.juetun.com\"}', '后台管理', '后台管理代码', 1, '2021-05-27 23:59:05', '2021-05-28 23:45:23', NULL);
COMMIT;

-- ----------------------------
-- Table structure for admin_group
-- ----------------------------
DROP TABLE IF EXISTS `admin_group`;
CREATE TABLE `admin_group` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '组名',
  `parent_id` int NOT NULL DEFAULT '0' COMMENT '上级用户组',
  `group_code` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '机构号',
  `is_super_admin` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否是超级管理员 0:否    1：是',
  `is_admin_group` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否为管理员组',
  `last_child_group_code` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of admin_group
-- ----------------------------
BEGIN;
INSERT INTO `admin_group` VALUES (1, '系统基础组', 0, 'A2', 1, 1, 'A2ED2', '2021-05-15 00:18:49', '2021-05-15 00:38:23', NULL);
INSERT INTO `admin_group` VALUES (2, '管理员', 1, 'A2EDD', 0, 0, '', '2021-06-01 22:54:31', '2021-06-01 22:54:31', NULL);
INSERT INTO `admin_group` VALUES (3, '管理组2', 1, 'A2ED1', 0, 0, '', '2021-06-01 23:05:16', '2021-09-05 21:10:43', NULL);
COMMIT;

-- ----------------------------
-- Table structure for admin_import
-- ----------------------------
DROP TABLE IF EXISTS `admin_import`;
CREATE TABLE `admin_import` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `permit_key` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `app_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '应用英文名',
  `app_version` varchar(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '1.0' COMMENT '引用版本号',
  `url_path` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'URL路径后缀',
  `request_method` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '请求方法 (GET)',
  `sort_value` int NOT NULL DEFAULT '0' COMMENT '排序值 值越大越靠前',
  `default_open` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '1:开启权限 2：不开启权限',
  `need_login` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '是否需要登录 2:不需要;1:需要',
  `need_sign` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '是否需要签名验证 1:需要 ;2:不需要',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='接口表数据';

-- ----------------------------
-- Records of admin_import
-- ----------------------------
BEGIN;
INSERT INTO `admin_import` VALUES (1, 'M0Xqxvnx1J9B', 'admin-main', '1.0', 'admin/permit/menu', 'GET', 1, 2, 0, 0, '2021-06-09 23:15:52', '2021-06-01 22:45:05', NULL);
INSERT INTO `admin_import` VALUES (3, '7OWZR3qpL9a0', 'admin-main', '1.0', 'admin/export/list', 'GET', 2, 2, 1, 2, '2021-06-28 21:05:15', '2021-06-03 22:55:53', NULL);
INSERT INTO `admin_import` VALUES (4, '438VRNq5zlGg', 'admin-main', '1.0', 'admin/permit/menu1', 'GET', 2, 1, 1, 1, '2021-06-10 22:13:29', '2021-06-03 22:56:35', NULL);
INSERT INTO `admin_import` VALUES (5, 'Y3VdpOqR20Nn', 'admin-main', '1.0', 'admin/permit/admin_user', 'POST', 2, 1, 1, 1, '2021-06-15 23:47:10', '2021-06-03 22:58:27', NULL);
INSERT INTO `admin_import` VALUES (6, 'qyE757npKLwb', 'admin-main', '1.0', 'admin/permit/admin_user_group_release', 'POST', 2, 1, 1, 1, '2021-06-10 22:16:07', '2021-06-03 23:00:16', NULL);
INSERT INTO `admin_import` VALUES (7, 'dDqz5z4xjPGo', 'admin-main', '1.0', 'admin/permit/admin_user_group_add', 'POST', 0, 1, 1, 1, '2021-06-15 23:47:07', '2021-06-03 23:01:03', NULL);
INSERT INTO `admin_import` VALUES (8, 'XE1959m5jg2l', 'admin-main', '1.0', 'admin/permit/admin_user_delete', 'POST', 2, 1, 1, 1, '2021-06-15 23:47:05', '2021-06-03 23:01:35', NULL);
INSERT INTO `admin_import` VALUES (9, 'AYq45PZxj6aG', 'admin-main', '1.0', 'admin/permit/admin_group', 'GET', 7, 1, 1, 1, '2021-06-10 22:16:09', '2021-06-03 23:02:16', NULL);
INSERT INTO `admin_import` VALUES (10, 'BejyRD1xG0L3', 'admin-main', '1.0', 'admin/console/home', 'GET', 0, 1, 1, 1, '2021-06-09 23:51:40', '2021-06-09 23:51:26', NULL);
INSERT INTO `admin_import` VALUES (11, 'Pq3d5rV5k2Nj', 'admin-main', '1.0', 'admin/console/cate', 'GET', 0, 1, 1, 1, '2021-06-15 23:47:03', '2021-06-09 23:52:42', NULL);
INSERT INTO `admin_import` VALUES (12, 'K4Nd5mBR9E8n', 'admin-main', '1.0', 'admin/console/cate/edit/:id', 'GET', 0, 1, 1, 1, '2021-06-15 23:47:02', '2021-06-09 23:53:07', NULL);
INSERT INTO `admin_import` VALUES (13, '87LVpwopYvEy', 'admin-main', '1.0', 'admin/console/cate/:id', 'PUT', 0, 2, 1, 1, '2021-06-28 21:30:05', '2021-06-09 23:53:25', NULL);
INSERT INTO `admin_import` VALUES (14, '87LVpwopYvEZ', 'admin-main', '1.0', 'admin/permit/get_app_config', 'GET', 0, 1, 0, 1, '2021-09-05 11:52:34', '2021-09-05 11:52:34', NULL);
COMMIT;

-- ----------------------------
-- Table structure for admin_menu
-- ----------------------------
DROP TABLE IF EXISTS `admin_menu`;
CREATE TABLE `admin_menu` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `permit_key` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `module` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
  `parent_id` int NOT NULL DEFAULT '1' COMMENT '上级权限ID 1：表示不属于任何一级',
  `label` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '显示名称 或接口备注',
  `icon` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
  `hide_in_menu` tinyint(1) NOT NULL DEFAULT '1' COMMENT '菜单是否显示 1：显示 ，2：不显示',
  `manage_import_permit` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '是否有管理接口的权限 1:管理 2：不管理',
  `domain` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
  `url_path` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT 'URL路径后缀',
  `other_value` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '其他参数',
  `sort_value` int NOT NULL DEFAULT '0' COMMENT '排序值 值越大越靠前',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `permit_key` (`permit_key`)
) ENGINE=InnoDB AUTO_INCREMENT=58 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='菜单表数据';

-- ----------------------------
-- Records of admin_menu
-- ----------------------------
BEGIN;
INSERT INTO `admin_menu` VALUES (1, 'backend', 'platform', -1, '汽车', 'md-car', 0, 0, '', '', '', 0, '2021-05-21 00:50:18', '2021-02-28 21:18:21', NULL);
INSERT INTO `admin_menu` VALUES (2, 'style_admin_model', 'backend', 9, '车型车系', 'md-car', 2, 2, '', '', '', 1, '2021-06-27 20:02:26', '2020-09-27 23:28:57', NULL);
INSERT INTO `admin_menu` VALUES (4, 'travel_cricle', 'backend', 1, '游记圈子', 'ios-book', 2, 2, '', '', '', 70000000, '2021-04-16 23:41:01', '2020-12-09 23:27:10', NULL);
INSERT INTO `admin_menu` VALUES (9, 'car_data', 'backend', 1, '汽车数据111111', 'md-car', 2, 2, '', '', '', 80000000, '2021-09-05 21:48:38', '2020-09-27 23:28:57', NULL);
INSERT INTO `admin_menu` VALUES (10, 'style_manager', 'backend', 2, '车系管理', 'md-train', 2, 1, '', '', '', 0, '2021-03-10 00:42:01', '2021-03-06 20:19:20', NULL);
INSERT INTO `admin_menu` VALUES (11, 'model_manager', 'backend', 2, '车型管理', 'ios-car', 2, 1, '', '', '', 0, '2021-03-10 00:42:41', '2021-03-06 20:22:00', NULL);
INSERT INTO `admin_menu` VALUES (12, 'system_manager', 'system', 49, '系统管理', 'md-settings', 2, 2, '', '', '', 0, '2021-06-02 22:27:51', '2021-03-07 19:29:46', NULL);
INSERT INTO `admin_menu` VALUES (13, 'permit_index', 'system', 12, '权限', 'ios-keypad', 2, 2, '', '', '', 0, '2021-04-17 00:00:42', '2021-03-07 19:32:04', NULL);
INSERT INTO `admin_menu` VALUES (14, 'console_user_list', 'system', 13, '管理员', 'md-contact', 2, 1, '', '', '', 0, '2021-04-17 00:01:02', '2021-03-07 19:32:35', NULL);
INSERT INTO `admin_menu` VALUES (15, 'console_group_list', 'system', 13, '管理员组', 'md-bulb', 2, 1, '', '', '', 0, '2021-03-07 22:37:54', '2021-03-07 22:37:54', NULL);
INSERT INTO `admin_menu` VALUES (16, 'admin_group_permit', 'system', 13, '权限管理', 'ios-cash-outline', 2, 1, '', '', '', 0, '2021-03-07 22:38:36', '2021-03-07 22:38:36', NULL);
INSERT INTO `admin_menu` VALUES (17, 'admin_group_apps', 'system', 52, '服务管理', 'ios-briefcase', 2, 1, '', '', '', 0, '2021-05-15 21:59:21', '2021-03-07 22:47:12', NULL);
INSERT INTO `admin_menu` VALUES (18, 'system_index', 'system', 52, '系统设置', 'md-settings', 2, 1, '', '', '', 0, '2021-05-15 21:59:02', '2021-03-07 23:35:41', NULL);
INSERT INTO `admin_menu` VALUES (19, 'car_media', 'backend', 9, '车型素材', '', 2, 1, '', '', '', 0, '2021-03-08 00:25:24', '2021-03-08 00:25:24', NULL);
INSERT INTO `admin_menu` VALUES (28, 'home_index', 'backend', 1, '首页', 'md-home', 2, 1, '', '', '', 90000001, '2021-05-12 22:21:05', '2021-04-15 23:21:25', NULL);
INSERT INTO `admin_menu` VALUES (29, 'Ky6wWjPM83Rm', 'backend', 1, '广告系统', 'md-attach', 2, 2, '', '', '', 60000000, '2021-04-17 00:00:01', '2021-04-15 23:47:25', NULL);
INSERT INTO `admin_menu` VALUES (30, 'RL7bWm5WGV24', 'backend', 1, '用户管理', 'ios-man', 2, 2, '', '', '', 50000000, '2021-04-17 00:00:04', '2021-04-15 23:49:36', NULL);
INSERT INTO `admin_menu` VALUES (31, 'site_recommend', 'backend', 29, '首页推荐', 'logo-angular', 2, 1, '', '', '', 0, '2021-04-16 00:12:26', '2021-04-16 00:08:39', NULL);
INSERT INTO `admin_menu` VALUES (32, 'advertise_list', 'backend', 29, '广告管理', 'md-boat', 2, 1, '', '', '', 0, '2021-04-16 00:12:45', '2021-04-16 00:09:17', NULL);
INSERT INTO `admin_menu` VALUES (33, 'user_index', 'backend', 30, '用户管理', 'md-contact', 2, 1, '', '', '', 0, '2021-04-16 00:17:51', '2021-04-16 00:17:51', NULL);
INSERT INTO `admin_menu` VALUES (34, 'user_permit', 'backend', 30, '用户权限', 'md-bulb', 2, 1, '', '', '', 0, '2021-04-16 00:18:44', '2021-04-16 00:18:44', NULL);
INSERT INTO `admin_menu` VALUES (35, 'post_list', 'backend', 4, '游记文章', 'ios-book', 2, 1, '', '', '', 0, '2021-04-16 00:19:58', '2021-04-16 00:19:58', NULL);
INSERT INTO `admin_menu` VALUES (36, 'post_trash', 'backend', 35, '回收站', 'md-trash', 1, 1, '', '', '', 0, '2021-04-16 00:53:57', '2021-04-16 00:20:30', NULL);
INSERT INTO `admin_menu` VALUES (37, 'cate_list', 'backend', 4, '分类列表', 'md-locate', 2, 1, '', '', '', 0, '2021-04-16 00:21:22', '2021-04-16 00:21:22', NULL);
INSERT INTO `admin_menu` VALUES (38, 'tag_list', 'backend', 4, '标签列表', 'md-share', 2, 1, '', '', '', 0, '2021-04-16 00:23:50', '2021-04-16 00:23:50', NULL);
INSERT INTO `admin_menu` VALUES (39, 'model_arg_edit', 'backend', 11, '车型参数编辑', 'ios-barcode', 1, 1, '', '', '', 0, '2021-05-15 23:22:06', '2021-04-16 00:27:17', NULL);
INSERT INTO `admin_menu` VALUES (40, 'YmG4zXlzLp5l', 'backend', 1, '公共接口', '', 1, 1, '', '', '', 90000000, '2021-04-16 23:19:55', '2021-04-16 23:19:49', NULL);
INSERT INTO `admin_menu` VALUES (46, 'user', 'platform', -1, '用户', 'md-contact', 0, 0, '', '', '', 0, '2021-05-21 00:50:30', '2021-05-15 20:50:10', NULL);
INSERT INTO `admin_menu` VALUES (47, 'boOQN3agzEKk', 'user', 46, '首页', 'md-home', 0, 1, '', '', '', 90000001, '2021-05-15 20:50:10', '2021-05-15 20:50:10', NULL);
INSERT INTO `admin_menu` VALUES (48, 'B1epMQXBzrKl', 'user', 46, '公共接口', '', 1, 1, '', '', '', 0, '2021-05-15 20:50:10', '2021-05-15 20:50:10', NULL);
INSERT INTO `admin_menu` VALUES (49, 'system', 'platform', -1, '系统', 'md-settings', 0, 0, '', '', '', 0, '2021-05-21 00:19:20', '2021-05-15 20:51:24', NULL);
INSERT INTO `admin_menu` VALUES (50, 'pv51WeL9N4oD', 'system', 49, '首页', 'md-home', 0, 1, '', '', '', 90000000, '2021-05-15 20:51:24', '2021-05-15 20:51:24', NULL);
INSERT INTO `admin_menu` VALUES (51, '7lDBNwj5MrQ9', 'system', 49, '公共接口', '', 1, 1, '', '', '', 0, '2021-05-15 20:51:24', '2021-05-15 20:51:24', NULL);
INSERT INTO `admin_menu` VALUES (52, 'xBZANrP7zpK8', 'system', 12, '系统控制', 'md-build', 2, 2, '', '', '', 0, '2021-06-02 22:28:39', '2021-05-15 21:58:35', NULL);
INSERT INTO `admin_menu` VALUES (53, 'crontab', 'system', 12, '定时任务', 'md-hand', 2, 1, '', '', '', 0, '2021-06-02 22:28:56', '2021-05-22 11:32:49', NULL);
INSERT INTO `admin_menu` VALUES (54, 'crontab_edit', 'system', 53, '定时任务编辑', '', 1, 1, '', '', '', 0, '2021-05-22 12:04:02', '2021-05-22 11:36:13', NULL);
INSERT INTO `admin_menu` VALUES (55, 'group_permit_set', 'system', 15, '组权限设置', '', 1, 1, '', '', '', 0, '2021-05-29 00:42:49', '2021-05-29 00:42:41', NULL);
INSERT INTO `admin_menu` VALUES (56, 'admin_app_edit', 'system', 17, '服务编辑', '', 1, 1, '', '', '', 0, '2021-05-30 00:26:35', '2021-05-30 00:26:35', NULL);
INSERT INTO `admin_menu` VALUES (57, 'system_import', 'system', 13, '接口列表', 'md-book', 2, 1, '', '', '', 0, '2021-06-03 00:01:43', '2021-06-03 00:00:54', NULL);
COMMIT;

-- ----------------------------
-- Table structure for admin_menu_import
-- ----------------------------
DROP TABLE IF EXISTS `admin_menu_import`;
CREATE TABLE `admin_menu_import` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `menu_id` int NOT NULL DEFAULT '0',
  `menu_module` varchar(50) NOT NULL DEFAULT '' COMMENT '菜单所属模块',
  `import_id` int NOT NULL DEFAULT '0',
  `import_app_name` varchar(50) NOT NULL DEFAULT '' COMMENT '接口应用路径',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `menu_id` (`menu_id`,`import_id`)
) ENGINE=InnoDB AUTO_INCREMENT=72 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='菜单关联接口';

-- ----------------------------
-- Records of admin_menu_import
-- ----------------------------
BEGIN;
INSERT INTO `admin_menu_import` VALUES (1, 50, 'system', 1, 'admin-main', '2021-06-06 00:49:13', '2021-06-27 10:51:05', NULL);
INSERT INTO `admin_menu_import` VALUES (2, 50, 'system', 3, 'admin-main', '2021-06-22 00:06:03', '2021-06-27 10:51:05', NULL);
INSERT INTO `admin_menu_import` VALUES (8, 50, 'system', 4, 'admin-main', '2021-06-22 00:12:51', '2021-06-27 10:51:05', NULL);
INSERT INTO `admin_menu_import` VALUES (10, 50, 'system', 5, 'admin-main', '2021-06-22 00:12:54', '2021-06-27 10:51:05', NULL);
INSERT INTO `admin_menu_import` VALUES (14, 50, 'system', 11, 'admin-main', '2021-06-22 00:13:01', '2021-06-27 10:51:16', NULL);
INSERT INTO `admin_menu_import` VALUES (15, 50, 'system', 10, 'admin-main', '2021-06-22 00:13:01', '2021-06-27 10:51:16', NULL);
INSERT INTO `admin_menu_import` VALUES (16, 50, 'system', 9, 'admin-main', '2021-06-22 00:13:03', '2021-06-27 10:51:16', NULL);
INSERT INTO `admin_menu_import` VALUES (17, 50, 'system', 8, 'admin-main', '2021-06-22 00:13:03', '2021-06-27 10:51:16', NULL);
INSERT INTO `admin_menu_import` VALUES (18, 50, 'system', 7, 'admin-main', '2021-06-22 00:13:05', '2021-06-27 10:51:16', NULL);
INSERT INTO `admin_menu_import` VALUES (34, 50, 'system', 6, 'admin-main', '2021-06-26 22:34:31', '2021-06-27 10:51:05', NULL);
INSERT INTO `admin_menu_import` VALUES (62, 50, 'system', 12, 'admin-main', '2021-06-27 10:51:20', '2021-06-27 10:51:20', NULL);
INSERT INTO `admin_menu_import` VALUES (63, 50, 'system', 13, 'admin-main', '2021-06-27 10:51:20', '2021-06-27 10:51:20', NULL);
INSERT INTO `admin_menu_import` VALUES (64, 14, 'system', 1, 'admin-main', '2021-06-27 20:21:59', '2021-06-27 20:22:05', '2021-06-27 20:22:06');
INSERT INTO `admin_menu_import` VALUES (65, 14, 'system', 3, 'admin-main', '2021-06-27 20:22:00', '2021-06-27 20:22:06', '2021-06-27 20:22:07');
INSERT INTO `admin_menu_import` VALUES (66, 14, 'system', 4, 'admin-main', '2021-06-27 20:22:00', '2021-06-27 20:22:08', '2021-06-27 20:22:08');
INSERT INTO `admin_menu_import` VALUES (67, 14, 'system', 5, 'admin-main', '2021-06-27 20:22:01', '2021-06-27 20:22:08', '2021-06-27 20:22:09');
COMMIT;

-- ----------------------------
-- Table structure for admin_user
-- ----------------------------
DROP TABLE IF EXISTS `admin_user`;
CREATE TABLE `admin_user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_hid` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '用户ID',
  `real_name` varchar(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '姓名',
  `mobile` varchar(15) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '手机号',
  `flag_admin` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否超级管理员',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '表示未删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_hid` (`user_hid`)
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='后台用户表';

-- ----------------------------
-- Records of admin_user
-- ----------------------------
BEGIN;
INSERT INTO `admin_user` VALUES (1, '5', '长江', '', 0, '2020-09-20 17:30:30', '2020-09-20 17:30:30', NULL);
INSERT INTO `admin_user` VALUES (2, '7', 'asdfasd', '', 0, '2020-10-28 23:57:06', '2020-10-29 23:19:01', NULL);
INSERT INTO `admin_user` VALUES (3, '1', '测试', '15108352617', 0, '2020-10-28 23:57:39', '2020-10-29 00:01:16', '2021-03-10 00:09:33');
INSERT INTO `admin_user` VALUES (6, '8', '昵称', '', 0, '2020-10-29 00:01:40', '2020-10-29 00:01:40', NULL);
INSERT INTO `admin_user` VALUES (7, '9', '昵称', '', 0, '2020-10-29 00:01:47', '2020-10-29 00:01:47', '2020-10-29 22:38:25');
INSERT INTO `admin_user` VALUES (8, '11', '2', '', 0, '2020-10-29 00:04:58', '2020-10-29 00:04:58', '2020-11-05 00:08:42');
INSERT INTO `admin_user` VALUES (9, '12', '22321', '', 0, '2020-10-29 00:05:00', '2020-10-29 00:05:00', NULL);
INSERT INTO `admin_user` VALUES (10, '13', '12', '', 0, '2020-10-29 00:05:05', '2020-10-29 00:05:05', '2021-05-21 23:42:09');
INSERT INTO `admin_user` VALUES (11, '14', '123', '', 0, '2020-10-29 00:05:08', '2020-10-29 00:05:08', '2021-06-03 23:01:18');
INSERT INTO `admin_user` VALUES (12, '15', '1236', '', 0, '2020-10-29 00:05:13', '2020-10-29 00:05:13', '2021-06-01 22:58:30');
INSERT INTO `admin_user` VALUES (13, '16', '1238', '', 0, '2020-10-29 00:05:19', '2020-10-29 00:05:19', '2021-05-15 18:18:10');
INSERT INTO `admin_user` VALUES (14, '17', '345', '', 0, '2020-10-29 00:05:30', '2020-10-29 00:05:30', '2021-04-22 23:50:52');
INSERT INTO `admin_user` VALUES (15, '18', '123', '', 0, '2020-10-29 00:05:34', '2020-10-29 00:05:34', '2021-04-22 23:50:44');
INSERT INTO `admin_user` VALUES (16, '19', '666', '', 0, '2020-10-29 00:05:39', '2020-10-29 00:05:39', '2021-04-22 23:50:42');
COMMIT;

-- ----------------------------
-- Table structure for admin_user_group
-- ----------------------------
DROP TABLE IF EXISTS `admin_user_group`;
CREATE TABLE `admin_user_group` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `group_id` int NOT NULL COMMENT '用户组ID',
  `user_hid` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '用户HID',
  `deleted_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `group_id` (`group_id`,`user_hid`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户组';

-- ----------------------------
-- Records of admin_user_group
-- ----------------------------
BEGIN;
INSERT INTO `admin_user_group` VALUES (2, 2, '7', '2021-06-02 23:41:29', '2021-06-02 23:34:31', '2021-06-02 21:35:22');
INSERT INTO `admin_user_group` VALUES (3, 1, '5', '2021-06-13 21:59:44', '2021-06-13 21:59:40', '2021-06-02 23:33:11');
INSERT INTO `admin_user_group` VALUES (15, 3, '8', NULL, '2021-09-05 21:53:27', '2021-09-05 21:53:27');
INSERT INTO `admin_user_group` VALUES (16, 2, '12', NULL, '2021-09-05 21:53:42', '2021-09-05 21:53:42');
COMMIT;

-- ----------------------------
-- Table structure for admin_user_group_permit
-- ----------------------------
DROP TABLE IF EXISTS `admin_user_group_permit`;
CREATE TABLE `admin_user_group_permit` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `group_id` int NOT NULL DEFAULT '0' COMMENT '用户组ID',
  `path_type` enum('page','api') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'page' COMMENT '当前连接的属性 page:页面,api:接口',
  `module` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `menu_id` int NOT NULL DEFAULT '0' COMMENT '权限ID（admin_menu表主键）',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `group_id` (`group_id`,`path_type`,`menu_id`)
) ENGINE=InnoDB AUTO_INCREMENT=165 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户组所具有的权限';

-- ----------------------------
-- Records of admin_user_group_permit
-- ----------------------------
BEGIN;
INSERT INTO `admin_user_group_permit` VALUES (1, 2, 'page', 'backend', 28, '2021-06-01 22:54:50', '2021-06-01 22:54:50', NULL);
INSERT INTO `admin_user_group_permit` VALUES (2, 2, 'page', 'backend', 40, '2021-06-01 22:54:50', '2021-06-01 22:54:50', NULL);
INSERT INTO `admin_user_group_permit` VALUES (3, 3, 'api', 'system', 1, '2021-06-27 10:59:42', '2021-06-27 10:59:45', '2021-06-27 20:10:59');
INSERT INTO `admin_user_group_permit` VALUES (5, 3, 'api', 'system', 3, '2021-06-27 10:59:45', '2021-06-27 10:59:45', '2021-06-27 20:10:59');
INSERT INTO `admin_user_group_permit` VALUES (6, 3, 'page', 'user', 47, '2021-06-27 18:16:55', '2021-06-27 19:59:19', '2021-06-27 20:36:48');
INSERT INTO `admin_user_group_permit` VALUES (7, 3, 'page', 'user', 48, '2021-06-27 18:16:55', '2021-06-27 19:59:19', '2021-06-27 20:36:48');
INSERT INTO `admin_user_group_permit` VALUES (8, 3, 'api', 'system', 4, '2021-06-27 18:17:00', '2021-06-28 23:25:34', NULL);
INSERT INTO `admin_user_group_permit` VALUES (9, 3, 'api', 'system', 5, '2021-06-27 18:17:00', '2021-06-28 23:26:13', NULL);
INSERT INTO `admin_user_group_permit` VALUES (10, 3, 'api', 'system', 11, '2021-06-27 18:17:00', '2021-06-27 19:59:22', NULL);
INSERT INTO `admin_user_group_permit` VALUES (11, 3, 'api', 'system', 10, '2021-06-27 18:17:00', '2021-06-27 19:59:22', NULL);
INSERT INTO `admin_user_group_permit` VALUES (12, 3, 'api', 'system', 9, '2021-06-27 18:17:00', '2021-06-27 19:59:22', '2021-06-27 20:11:10');
INSERT INTO `admin_user_group_permit` VALUES (13, 3, 'api', '', 8, '2021-06-27 18:17:00', '2021-06-29 00:27:17', NULL);
INSERT INTO `admin_user_group_permit` VALUES (14, 3, 'api', '', 7, '2021-06-27 18:17:00', '2021-06-29 00:27:12', NULL);
INSERT INTO `admin_user_group_permit` VALUES (15, 3, 'api', '', 6, '2021-06-27 18:17:00', '2021-06-29 00:27:11', NULL);
INSERT INTO `admin_user_group_permit` VALUES (16, 3, 'api', '', 12, '2021-06-27 18:17:00', '2021-06-29 00:27:05', NULL);
INSERT INTO `admin_user_group_permit` VALUES (17, 3, 'api', 'system', 13, '2021-06-27 18:17:00', '2021-06-27 19:59:22', '2021-06-27 20:11:17');
INSERT INTO `admin_user_group_permit` VALUES (18, 3, 'page', 'system', 50, '2021-06-27 18:17:00', '2021-06-27 19:59:22', NULL);
INSERT INTO `admin_user_group_permit` VALUES (19, 3, 'page', 'system', 51, '2021-06-27 18:17:00', '2021-06-27 19:59:22', '2021-06-27 20:12:23');
INSERT INTO `admin_user_group_permit` VALUES (20, 3, 'page', 'backend', 28, '2021-06-27 18:17:18', '2021-06-27 20:14:45', NULL);
INSERT INTO `admin_user_group_permit` VALUES (21, 3, 'page', 'backend', 40, '2021-06-27 18:17:18', '2021-06-27 20:16:19', NULL);
INSERT INTO `admin_user_group_permit` VALUES (22, 3, 'page', 'system', 12, '2021-06-27 18:26:38', '2021-06-27 18:26:38', '2021-06-27 20:12:23');
INSERT INTO `admin_user_group_permit` VALUES (23, 3, 'page', 'system', 13, '2021-06-27 18:26:38', '2021-06-27 18:26:38', '2021-06-27 20:12:23');
INSERT INTO `admin_user_group_permit` VALUES (24, 3, 'page', 'system', 14, '2021-06-27 18:26:38', '2021-06-27 18:26:38', NULL);
INSERT INTO `admin_user_group_permit` VALUES (25, 3, 'page', 'system', 15, '2021-06-27 18:26:38', '2021-06-27 18:26:38', '2021-06-27 20:12:23');
INSERT INTO `admin_user_group_permit` VALUES (26, 3, 'page', 'system', 55, '2021-06-27 18:26:38', '2021-06-27 18:26:38', '2021-06-27 20:12:23');
INSERT INTO `admin_user_group_permit` VALUES (27, 3, 'page', 'system', 16, '2021-06-27 18:26:38', '2021-06-27 18:26:38', NULL);
INSERT INTO `admin_user_group_permit` VALUES (28, 3, 'page', 'system', 57, '2021-06-27 18:26:38', '2021-06-27 18:26:38', NULL);
INSERT INTO `admin_user_group_permit` VALUES (29, 3, 'page', 'system', 52, '2021-06-27 18:26:38', '2021-06-27 18:26:38', NULL);
INSERT INTO `admin_user_group_permit` VALUES (30, 3, 'page', 'system', 17, '2021-06-27 18:26:38', '2021-06-27 18:26:38', NULL);
INSERT INTO `admin_user_group_permit` VALUES (31, 3, 'page', 'system', 56, '2021-06-27 18:26:38', '2021-06-27 18:26:38', NULL);
INSERT INTO `admin_user_group_permit` VALUES (32, 3, 'page', 'system', 18, '2021-06-27 18:26:38', '2021-06-27 18:26:38', NULL);
INSERT INTO `admin_user_group_permit` VALUES (33, 3, 'page', 'system', 53, '2021-06-27 18:26:38', '2021-06-27 18:26:38', NULL);
INSERT INTO `admin_user_group_permit` VALUES (34, 3, 'page', 'system', 54, '2021-06-27 18:26:38', '2021-06-27 18:26:38', NULL);
INSERT INTO `admin_user_group_permit` VALUES (35, 3, 'page', 'backend', 9, '2021-06-27 18:28:47', '2021-06-27 18:28:47', '2021-06-27 18:29:27');
INSERT INTO `admin_user_group_permit` VALUES (36, 3, 'page', 'backend', 2, '2021-06-27 18:28:47', '2021-06-27 20:14:38', '2021-06-27 20:14:40');
INSERT INTO `admin_user_group_permit` VALUES (37, 3, 'page', 'backend', 10, '2021-06-27 18:28:47', '2021-06-27 20:14:38', '2021-06-27 20:14:40');
INSERT INTO `admin_user_group_permit` VALUES (38, 3, 'page', 'backend', 11, '2021-06-27 18:28:47', '2021-06-27 20:16:15', '2021-06-27 20:16:20');
INSERT INTO `admin_user_group_permit` VALUES (39, 3, 'page', 'backend', 39, '2021-06-27 18:28:47', '2021-06-27 20:16:15', '2021-06-27 20:16:20');
INSERT INTO `admin_user_group_permit` VALUES (40, 3, 'page', 'backend', 19, '2021-06-27 18:28:47', '2021-06-27 18:28:47', '2021-06-27 18:29:27');
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
