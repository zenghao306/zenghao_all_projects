/*
Navicat MySQL Data Transfer

Source Server         : 120.76.96.73 game
Source Server Version : 50505
Source Host           : 120.76.96.73:3306
Source Database       : games_temp

Target Server Type    : MYSQL
Target Server Version : 50505
File Encoding         : 65001

Date: 2017-06-07 14:18:40
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for `php_cash_moon_switch`
-- ----------------------------
DROP TABLE IF EXISTS `php_cash_moon_switch`;
CREATE TABLE `php_cash_moon_switch` (
  `id` int(10) NOT NULL,
  `channel_id` varchar(125) NOT NULL COMMENT '渠道id',
  `status` int(1) unsigned zerofill NOT NULL DEFAULT '0' COMMENT '状态，0关闭(不显示)，1开启（显示）',
  `the_time` int(10) NOT NULL COMMENT '新增/更新时间，时间戳',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of php_cash_moon_switch
-- ----------------------------
