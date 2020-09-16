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
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;



--  Drop table

--  DROP TABLE bidding_app.Auction;
CREATE TABLE `Auction` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `start_time` datetime NOT NULL,
  `end_time` datetime NOT NULL,
  `start_amount` bigint unsigned NOT NULL,
  `end_amount` double NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


-- Create a First Admin Role User  // password: 1234
INSERT INTO bidding_app.`User` (name,`role`,email,password,date_created,date_modified) VALUES
('hari prasad',0,'hariprasadmails@gmail.com','47DEQpj8HBSa+/TImW+5JCeuQeRkm5NMpJWZG3hSuFU=','2020-09-13 15:29:42','2020-09-13 15:29:42')
;

CREATE TABLE bidding_app.Bid (
	auction_id BIGINT UNSIGNED NOT NULL,
	user_id BIGINT UNSIGNED NOT NULL,
	amount FLOAT NOT NULL,
	CONSTRAINT Bid_FK FOREIGN KEY (auction_id) REFERENCES bidding_app.Auction(id),
	CONSTRAINT Bid_FK_1 FOREIGN KEY (user_id) REFERENCES bidding_app.`User`(id)
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8mb4
COLLATE=utf8mb4_0900_ai_ci;
