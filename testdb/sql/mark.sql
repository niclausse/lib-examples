CREATE TABLE `tblCoursePacks` (
    `id` bigint(11) unsigned NOT NULL AUTO_INCREMENT,
    `deleted` int(10) NOT NULL DEFAULT '0' COMMENT '删除状态 0.未删除 1.已删除',
    `grade_id` int(10) NOT NULL COMMENT '年级',
    `subject_id` int(10) NOT NULL COMMENT '科目',
    `course_count` int(10) NOT NULL DEFAULT '0' COMMENT '课程数量',
    `use_year` int(10) NOT NULL DEFAULT '0',
    `create_time` int(11) NOT NULL COMMENT '创建时间',
    `update_time` int(11) NOT NULL COMMENT '更新时间',
    `name` varchar(100) NOT NULL DEFAULT '' COMMENT '课程包名称',
    PRIMARY KEY (`id`),
    KEY `idx_grade_subject` (`grade_id`,`subject_id`),
    KEY `idx_create_time` (`create_time`),
    KEY `idx_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=1214 DEFAULT CHARSET=utf8mb4 COMMENT='课程包列表';

CREATE TABLE `tblCoursePackRelation` (
    `id` bigint(11) unsigned NOT NULL AUTO_INCREMENT,
    `deleted` int(10) NOT NULL DEFAULT '0' COMMENT '删除状态 0.未删除 1.已删除',
    `course_pack_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '课程包的id',
    `course_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '课程的id',
    `create_time` int(11) NOT NULL COMMENT '创建时间',
    `update_time` int(11) NOT NULL COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_course` (`course_id`),
    KEY `idx_pack_course` (`course_pack_id`,`course_id`,`deleted`)
) ENGINE=InnoDB AUTO_INCREMENT=1399 DEFAULT CHARSET=utf8mb4 COMMENT='课程包列表';

CREATE TABLE `tblCourse` (
    `id` bigint(11) unsigned NOT NULL AUTO_INCREMENT,
    `deleted` int(10) NOT NULL DEFAULT '0' COMMENT '删除状态 0.未删除 1.已删除',
    `book_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '书的id',
    `tiku_book_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '书id-题库--已废弃',
    `grade_id` int(10) NOT NULL COMMENT '年级',
    `subject_id` int(10) NOT NULL COMMENT '科目',
    `textbook_version` int(10) NOT NULL DEFAULT '0' COMMENT '教材版本',
    `name` varchar(100) NOT NULL DEFAULT '' COMMENT '课程名称',
    `cover_img` varchar(500) NOT NULL DEFAULT '' COMMENT '封面图片',
    `ti_edit_status` int(11) NOT NULL DEFAULT '0' COMMENT '题目补充状态 0未补充 1已补充',
    `type` int(10) unsigned NOT NULL DEFAULT '1' COMMENT '1.练习册 21.周测',
    `source` int(10) unsigned NOT NULL DEFAULT '10' COMMENT '10-mis后台创建 21-试题篮自动创建',
    `create_time` int(11) NOT NULL COMMENT '创建时间',
    `update_time` int(11) NOT NULL COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_create_time` (`create_time`),
    KEY `idx_grade_subject` (`grade_id`,`subject_id`),
    KEY `idx_book_id` (`book_id`),
    KEY `idx_source_type` (`source`,`type`)
) ENGINE=InnoDB AUTO_INCREMENT=67239 DEFAULT CHARSET=utf8mb4 COMMENT='课程信息';