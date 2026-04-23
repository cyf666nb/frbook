-- 全市高校二手图书共享流转平台 数据库设计
-- Database: bookshare

CREATE DATABASE IF NOT EXISTS bookshare DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE bookshare;

-- ============================================
-- 1. 用户表
-- ============================================
CREATE TABLE `users` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    `account` VARCHAR(50) NOT NULL COMMENT '账号(手机号或自定义昵称)',
    `password` VARCHAR(255) NOT NULL COMMENT '密码(加密存储)',
    `birth_date` DATE NOT NULL COMMENT '出生日期',
    `nickname` VARCHAR(50) DEFAULT NULL COMMENT '昵称',
    `avatar` VARCHAR(255) DEFAULT NULL COMMENT '头像URL',
    `school` VARCHAR(100) DEFAULT NULL COMMENT '所在学校',
    `campus` VARCHAR(100) DEFAULT NULL COMMENT '校区',
    `grade` VARCHAR(20) DEFAULT NULL COMMENT '年级',
    `interest_tags` JSON DEFAULT NULL COMMENT '兴趣标签(JSON数组)',
    `credit_score` INT DEFAULT 100 COMMENT '信用评分(初始100分)',
    `status` TINYINT DEFAULT 1 COMMENT '状态: 0-禁用 1-正常',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_account` (`account`),
    KEY `idx_school` (`school`),
    KEY `idx_credit_score` (`credit_score`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- ============================================
-- 2. 图书表
-- ============================================
CREATE TABLE `books` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '图书ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '发布者用户ID',
    `isbn` VARCHAR(20) DEFAULT NULL COMMENT 'ISBN编号',
    `title` VARCHAR(255) NOT NULL COMMENT '书名',
    `author` VARCHAR(255) DEFAULT NULL COMMENT '作者',
    `publisher` VARCHAR(100) DEFAULT NULL COMMENT '出版社',
    `cover_image` VARCHAR(255) DEFAULT NULL COMMENT '封面图URL',
    `description` TEXT DEFAULT NULL COMMENT '图书描述',
    `category` VARCHAR(50) DEFAULT NULL COMMENT '分类',
    `mode` TINYINT NOT NULL COMMENT '流转模式: 1-租赁 2-出售 3-赠送',
    `status` TINYINT DEFAULT 1 COMMENT '状态: 0-下架 1-上架 2-已交易',
    
    -- 租赁相关字段
    `daily_rent` DECIMAL(10,2) DEFAULT NULL COMMENT '日租金',
    `weekly_rent` DECIMAL(10,2) DEFAULT NULL COMMENT '周租金',
    `deposit` DECIMAL(10,2) DEFAULT NULL COMMENT '押金金额',
    `min_rent_days` INT DEFAULT 1 COMMENT '最短租期(天)',
    
    -- 出售相关字段
    `sell_price` DECIMAL(10,2) DEFAULT NULL COMMENT '出售价格',
    
    -- 赠送模式价格锁定为0
    
    -- 通用字段
    `images` JSON NOT NULL COMMENT '实物图片(JSON数组)',
    `pickup_location` VARCHAR(255) DEFAULT NULL COMMENT '线下取书首选位置',
    `view_count` INT DEFAULT 0 COMMENT '浏览次数',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_isbn` (`isbn`),
    KEY `idx_mode` (`mode`),
    KEY `idx_status` (`status`),
    KEY `idx_created_at` (`created_at`),
    FULLTEXT KEY `ft_title_author` (`title`, `author`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='图书表';

-- ============================================
-- 3. 租赁订单表
-- ============================================
CREATE TABLE `rental_orders` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '订单ID',
    `order_no` VARCHAR(32) NOT NULL COMMENT '订单编号',
    `book_id` BIGINT UNSIGNED NOT NULL COMMENT '图书ID',
    `book_title` VARCHAR(255) NOT NULL COMMENT '图书名称(冗余)',
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT '出租人ID',
    `renter_id` BIGINT UNSIGNED NOT NULL COMMENT '承租人ID',
    
    -- 租赁信息
    `daily_rent` DECIMAL(10,2) NOT NULL COMMENT '日租金',
    `deposit` DECIMAL(10,2) NOT NULL COMMENT '押金',
    `rent_days` INT NOT NULL COMMENT '租赁天数',
    `total_rent` DECIMAL(10,2) NOT NULL COMMENT '总租金',
    `total_amount` DECIMAL(10,2) NOT NULL COMMENT '总金额(押金+首期租金)',
    
    -- 时间相关
    `start_time` DATETIME DEFAULT NULL COMMENT '租赁开始时间(出租人确认取书后)',
    `end_time` DATETIME DEFAULT NULL COMMENT '租赁结束时间',
    `actual_return_time` DATETIME DEFAULT NULL COMMENT '实际归还时间',
    
    -- 状态: 0-待支付 1-已支付待确认取书 2-租赁中 3-待验收 4-已完成 5-已取消 6-逾期中
    `status` TINYINT DEFAULT 0 COMMENT '订单状态',
    
    -- 金额结算
    `overdue_fee` DECIMAL(10,2) DEFAULT 0.00 COMMENT '逾期费用',
    `refund_amount` DECIMAL(10,2) DEFAULT NULL COMMENT '退还金额',
    `settled_amount` DECIMAL(10,2) DEFAULT NULL COMMENT '结算给出租人的金额',
    
    -- 评价
    `owner_rating` TINYINT DEFAULT NULL COMMENT '出租人评分(1-5)',
    `renter_rating` TINYINT DEFAULT NULL COMMENT '承租人评分(1-5)',
    `owner_comment` VARCHAR(500) DEFAULT NULL COMMENT '出租人评价',
    `renter_comment` VARCHAR(500) DEFAULT NULL COMMENT '承租人评价',
    
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_order_no` (`order_no`),
    KEY `idx_book_id` (`book_id`),
    KEY `idx_owner_id` (`owner_id`),
    KEY `idx_renter_id` (`renter_id`),
    KEY `idx_status` (`status`),
    KEY `idx_start_time` (`start_time`),
    KEY `idx_end_time` (`end_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='租赁订单表';

-- ============================================
-- 4. 出售订单表
-- ============================================
CREATE TABLE `sell_orders` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '订单ID',
    `order_no` VARCHAR(32) NOT NULL COMMENT '订单编号',
    `book_id` BIGINT UNSIGNED NOT NULL COMMENT '图书ID',
    `book_title` VARCHAR(255) NOT NULL COMMENT '图书名称(冗余)',
    `seller_id` BIGINT UNSIGNED NOT NULL COMMENT '卖家ID',
    `buyer_id` BIGINT UNSIGNED NOT NULL COMMENT '买家ID',
    
    -- 价格信息
    `price` DECIMAL(10,2) NOT NULL COMMENT '出售价格',
    `platform_fee` DECIMAL(10,2) DEFAULT 0.00 COMMENT '平台抽成',
    `settle_amount` DECIMAL(10,2) DEFAULT NULL COMMENT '结算给卖家的金额',
    
    -- 状态: 0-待支付 1-已支付待发货 2-已发货待收货 3-已完成 4-已取消
    `status` TINYINT DEFAULT 0 COMMENT '订单状态',
    
    -- 发货信息
    `delivery_type` TINYINT DEFAULT NULL COMMENT '交付方式: 1-快递 2-线下自提',
    `delivery_no` VARCHAR(50) DEFAULT NULL COMMENT '快递单号',
    `delivery_company` VARCHAR(50) DEFAULT NULL COMMENT '快递公司',
    `delivery_time` DATETIME DEFAULT NULL COMMENT '发货时间',
    `receive_time` DATETIME DEFAULT NULL COMMENT '收货时间',
    
    -- 自动结算时间
    `auto_settle_time` DATETIME DEFAULT NULL COMMENT '自动结算时间(发货后7天)',
    
    -- 评价
    `seller_rating` TINYINT DEFAULT NULL COMMENT '卖家评分(1-5)',
    `buyer_rating` TINYINT DEFAULT NULL COMMENT '买家评分(1-5)',
    `seller_comment` VARCHAR(500) DEFAULT NULL COMMENT '卖家评价',
    `buyer_comment` VARCHAR(500) DEFAULT NULL COMMENT '买家评价',
    
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_order_no` (`order_no`),
    KEY `idx_book_id` (`book_id`),
    KEY `idx_seller_id` (`seller_id`),
    KEY `idx_buyer_id` (`buyer_id`),
    KEY `idx_status` (`status`),
    KEY `idx_auto_settle_time` (`auto_settle_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='出售订单表';

-- ============================================
-- 5. 赠送记录表
-- ============================================
CREATE TABLE `gift_records` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '记录ID',
    `record_no` VARCHAR(32) NOT NULL COMMENT '记录编号',
    `book_id` BIGINT UNSIGNED NOT NULL COMMENT '图书ID',
    `book_title` VARCHAR(255) NOT NULL COMMENT '图书名称(冗余)',
    `giver_id` BIGINT UNSIGNED NOT NULL COMMENT '赠送人ID',
    `receiver_id` BIGINT UNSIGNED NOT NULL COMMENT '接收人ID',
    
    -- 状态: 0-待确认 1-已确认待交付 2-已完成 3-已取消
    `status` TINYINT DEFAULT 0 COMMENT '状态',
    
    -- 交付信息
    `delivery_type` TINYINT DEFAULT NULL COMMENT '交付方式: 1-快递 2-线下自提',
    `delivery_time` DATETIME DEFAULT NULL COMMENT '交付时间',
    
    -- 评价
    `giver_rating` TINYINT DEFAULT NULL COMMENT '赠送人评分(1-5)',
    `receiver_rating` TINYINT DEFAULT NULL COMMENT '接收人评分(1-5)',
    `giver_comment` VARCHAR(500) DEFAULT NULL COMMENT '赠送人评价',
    `receiver_comment` VARCHAR(500) DEFAULT NULL COMMENT '接收人评价',
    
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_record_no` (`record_no`),
    KEY `idx_book_id` (`book_id`),
    KEY `idx_giver_id` (`giver_id`),
    KEY `idx_receiver_id` (`receiver_id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='赠送记录表';

-- ============================================
-- 6. 交易流水表
-- ============================================
CREATE TABLE `transactions` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '流水ID',
    `transaction_no` VARCHAR(32) NOT NULL COMMENT '流水号',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `related_user_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '关联用户ID',
    
    -- 业务类型: 1-租赁押金 2-租赁租金 3-租赁押金退还 4-购买支付 5-购买结算 6-平台抽成
    `type` TINYINT NOT NULL COMMENT '业务类型',
    
    -- 关联订单
    `order_type` TINYINT NOT NULL COMMENT '订单类型: 1-租赁 2-出售',
    `order_id` BIGINT UNSIGNED NOT NULL COMMENT '订单ID',
    `order_no` VARCHAR(32) NOT NULL COMMENT '订单编号',
    
    -- 金额
    `amount` DECIMAL(10,2) NOT NULL COMMENT '金额(正数收入,负数支出)',
    `balance_before` DECIMAL(10,2) DEFAULT NULL COMMENT '变动前余额',
    `balance_after` DECIMAL(10,2) DEFAULT NULL COMMENT '变动后余额',
    
    -- 状态: 0-待处理 1-成功 2-失败
    `status` TINYINT DEFAULT 1 COMMENT '状态',
    `remark` VARCHAR(255) DEFAULT NULL COMMENT '备注',
    
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_transaction_no` (`transaction_no`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_order_id` (`order_id`),
    KEY `idx_type` (`type`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='交易流水表';

-- ============================================
-- 7. 用户钱包表
-- ============================================
CREATE TABLE `wallets` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '钱包ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `balance` DECIMAL(10,2) DEFAULT 0.00 COMMENT '可用余额',
    `frozen_balance` DECIMAL(10,2) DEFAULT 0.00 COMMENT '冻结余额(押金/担保金额)',
    `total_income` DECIMAL(10,2) DEFAULT 0.00 COMMENT '累计收入',
    `total_expense` DECIMAL(10,2) DEFAULT 0.00 COMMENT '累计支出',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户钱包表';

-- ============================================
-- 8. 用户信用记录表
-- ============================================
CREATE TABLE `credit_records` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '记录ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `score_change` INT NOT NULL COMMENT '分数变化(正数加分,负数扣分)',
    `score_before` INT NOT NULL COMMENT '变动前分数',
    `score_after` INT NOT NULL COMMENT '变动后分数',
    `reason` VARCHAR(255) NOT NULL COMMENT '原因',
    `related_order_type` TINYINT DEFAULT NULL COMMENT '关联订单类型',
    `related_order_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '关联订单ID',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户信用记录表';

-- ============================================
-- 9. 图书收藏表
-- ============================================
CREATE TABLE `book_favorites` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '收藏ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `book_id` BIGINT UNSIGNED NOT NULL COMMENT '图书ID',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_book` (`user_id`, `book_id`),
    KEY `idx_book_id` (`book_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='图书收藏表';

-- ============================================
-- 10. 验证码表
-- ============================================
CREATE TABLE `verification_codes` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `target` VARCHAR(50) NOT NULL COMMENT '目标(手机号/邮箱)',
    `code` VARCHAR(10) NOT NULL COMMENT '验证码',
    `type` TINYINT NOT NULL COMMENT '类型: 1-注册 2-登录 3-重置密码',
    `expires_at` DATETIME NOT NULL COMMENT '过期时间',
    `used` TINYINT DEFAULT 0 COMMENT '是否已使用: 0-否 1-是',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_target_type` (`target`, `type`),
    KEY `idx_expires_at` (`expires_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='验证码表';

-- ============================================
-- 初始化数据
-- ============================================
-- 插入测试用户(密码为 bcrypt 加密的 "123456")
INSERT INTO `users` (`account`, `password`, `birth_date`, `nickname`, `school`, `campus`) VALUES
('13800138000', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt6Z5EH', '2000-01-01', '测试用户1', '北京大学', '海淀校区'),
('testuser', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt6Z5EH', '1999-05-15', '测试用户2', '清华大学', '五道口校区');

-- 为测试用户创建钱包
INSERT INTO `wallets` (`user_id`, `balance`) VALUES (1, 1000.00), (2, 500.00);
