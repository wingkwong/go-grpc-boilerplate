CREATE TABLE `Foo` (
  `ID` bigint(20) NOT NULL AUTO_INCREMENT,
  `Title` varchar(200),
  `Desc` varchar(1024),
  `CreatedBy` varchar(1024),
  `UpdatedBy` varchar(1024),
  `CreatedAt` timestamp NOT NULL,
  `UpdatedAt` timestamp NOT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `ID_UNIQUE` (`ID`)
);