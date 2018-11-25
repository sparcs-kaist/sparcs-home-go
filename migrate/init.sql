CREATE SCHEMA sparcs_home;

CREATE TABLE IF NOT EXISTS Album
(
    id int PRIMARY KEY AUTO_INCREMENT,
    year int NOT NULL,
    title varchar(255) NOT NULL,
    date date
);

CREATE TABLE IF NOT EXISTS Photo
(
    id int PRIMARY KEY AUTO_INCREMENT,
    album_id int NOT NULL,
    path varchar(255) NOT NULL,
    CONSTRAINT Photo_Album_id_fk FOREIGN KEY (album_id) REFERENCES Album (id) ON DELETE CASCADE
);

CREATE TABLE `Seminar` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL,
  `speaker` varchar(255) NOT NULL,
  `date` date DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8
  COLLATE utf8_unicode_ci;

CREATE TABLE `SeminarResource` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `seminar_id` int(11) NOT NULL,
  `path` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `SeminarResource___fk_seminar` (`seminar_id`),
  CONSTRAINT `SeminarResource___fk_seminar` FOREIGN KEY (`seminar_id`) REFERENCES `Seminar` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8
  COLLATE utf8_unicode_ci;

-- INSERT INTO `sparcs_home`.`Album` (`year`, `title`, `date`) VALUES (2018, 'Test', '2018-11-17');
-- INSERT INTO `sparcs_home`.`Album` (`year`, `title`, `date`) VALUES (2017, 'Past', '2017-11-08');
-- INSERT INTO `sparcs_home`.`Photo` (`album_id`, `path`) VALUES (1, '/Users/Youngkyu/go/src/github.com/sparcs-home-go/static/test1.jpg');
-- INSERT INTO `sparcs_home`.`Photo` (`album_id`, `path`) VALUES (2, '/Users/Youngkyu/go/src/github.com/sparcs-home-go/static/test2.jpg');