CREATE TABLE `post` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `version` int(10) unsigned NOT NULL DEFAULT 1,
  `parent` int(10) unsigned NOT NULL,
  `title` varchar(255) NOT NULL DEFAULT '',
  `body` longtext NOT NULL DEFAULT '',
  `published` datetime,
  `created` timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4;