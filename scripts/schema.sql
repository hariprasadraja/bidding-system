-- bidding_app.`User` definition
CREATE TABLE `User` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `role` int DEFAULT NULL,
  `email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `password` varchar(200) NOT NULL,
  `date_created` datetime DEFAULT NULL,
  `date_modified` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


--  Drop table

--  DROP TABLE bidding_app.Auction;

-- bidding_app.Auction definition
CREATE TABLE `Auction` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `start_time` datetime NOT NULL,
  `end_time` datetime NOT NULL,
  `start_amount` float NOT NULL,
  `currency` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Create a First Admin Role User  // password: 1234
INSERT INTO bidding_app.`User` (name,`role`,email,password,date_created,date_modified) VALUES
('hari prasad',0,'hariprasadmails@gmail.com','A6xnQhbz4Vx2HuGl4lXwZ5U2I8iziLRFnhP5eNfIRvQ=','2020-09-13 15:29:42','2020-09-13 15:29:42')
;
-- bidding_app.Bid definition

CREATE TABLE `Bid` (
  `auction_id` bigint unsigned NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `amount` float NOT NULL,
  KEY `Bid_FK` (`auction_id`),
  KEY `Bid_FK_1` (`user_id`),
  CONSTRAINT `Bid_FK` FOREIGN KEY (`auction_id`) REFERENCES `Auction` (`id`),
  CONSTRAINT `Bid_FK_1` FOREIGN KEY (`user_id`) REFERENCES `User` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;