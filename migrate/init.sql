CREATE SCHEMA sparcs_home;

CREATE TABLE IF NOT EXISTS Album
(
    id int PRIMARY KEY AUTO_INCREMENT,
    year int NOT NULL,
    title varchar(100) NOT NULL,
    date date
);

CREATE TABLE IF NOT EXISTS Photo
(
    id int PRIMARY KEY AUTO_INCREMENT,
    album_id int NOT NULL,
    path varchar(200) NOT NULL,
    CONSTRAINT Photo_Album_id_fk FOREIGN KEY (album_id) REFERENCES Album (id) ON DELETE CASCADE
);

-- INSERT INTO `sparcs_home`.`Album` (`year`, `title`, `date`) VALUES (2018, 'Test', '2018-11-17');
-- INSERT INTO `sparcs_home`.`Album` (`year`, `title`, `date`) VALUES (2017, 'Past', '2017-11-08');
-- INSERT INTO `sparcs_home`.`Photo` (`album_id`, `path`) VALUES (1, '/Users/Youngkyu/go/src/github.com/sparcs-home-go/static/test1.jpg');
-- INSERT INTO `sparcs_home`.`Photo` (`album_id`, `path`) VALUES (2, '/Users/Youngkyu/go/src/github.com/sparcs-home-go/static/test2.jpg');