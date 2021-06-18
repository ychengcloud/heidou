# 数据库

准备数据库

## ERD

数据实体关系模型定义

![erd](/diagrams/out/ecommerce/ecommerce.png)

## 数据字典

设计应用的数据字典 `xiaokeai.sql` 如下：

```sql
-- Create a database
CREATE DATABASE IF NOT EXISTS `xiaokeai` DEFAULT CHARACTER SET = `utf8mb4`;

USE `xiaokeai`;
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` char(36) NOT NULL,
  `name` varchar(32) NOT NULL DEFAULT '' COMMENT '用户名',
  `description` varchar(255) NOT NULL DEFAULT '' COMMENT '描述',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT = 'User';

DROP TABLE IF EXISTS `group`;
CREATE TABLE `group` (
  `id` char(36) NOT NULL,
  `name` varchar(32) NOT NULL DEFAULT '' COMMENT '用户组名',
  `description` varchar(255) NOT NULL DEFAULT '' COMMENT '描述',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT = 'group';

DROP TABLE IF EXISTS `user_group`;
CREATE TABLE `user_group` (
  `id` char(36) NOT NULL,
  `user_id` char(36) NOT NULL DEFAULT '' COMMENT '用户ID',
  `group_id` char(36) NOT NULL DEFAULT '' COMMENT '用户组ID',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_group_id` (`group_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT = 'user_group';

DROP TABLE IF EXISTS `product`;
CREATE TABLE `product` (
  `id` char(36) NOT NULL,
  `name` varchar(32) NOT NULL DEFAULT '' COMMENT '商品名',
  `description` varchar(255) NOT NULL DEFAULT '' COMMENT '描述',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT = 'product';

DROP TABLE IF EXISTS `order`;
CREATE TABLE `order` (
  `id` char(36) NOT NULL,
  `user_id` char(36) NOT NULL DEFAULT '' COMMENT '用户ID',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT = 'order';

DROP TABLE IF EXISTS `order_item`;
CREATE TABLE `order_item` (
  `id` char(36) NOT NULL,
  `product_id` char(36) NOT NULL DEFAULT '' COMMENT '商品ID',
  `order_id` char(36) NOT NULL DEFAULT '' COMMENT '订单ID',
  PRIMARY KEY (`id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_order_id` (`order_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT = 'order_item';

```

## 建表

```bash
mysql -h <host> -u <user> -p < xiaokeai.sql
```