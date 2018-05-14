/*
SQLyog v10.2 
MySQL - 5.5.5-10.2.7-MariaDB-10.2.7+maria~jessie-log : Database - 17wan
*********************************************************************
*/

/*!40101 SET NAMES utf8 */;

/*!40101 SET SQL_MODE=''*/;

/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
CREATE DATABASE /*!32312 IF NOT EXISTS*/`17wan` /*!40100 DEFAULT CHARACTER SET utf8mb4 */;

USE `17wan`;

/*Data for the table `go_config_android_trade_item` */

insert  into `go_config_android_trade_item`(`item_id`,`describe`,`money`,`diamond`) values ('gold_100000','100000钻石',100000,100000),('gold_12800','12800钻石',12800,12800),('gold_3000','3000钻石',3000,3000),('gold_300000','300000钻石',300000,300000),('gold_32800','32800钻石',32800,32800),('gold_500000','500000钻石',500000,500000),('gold_600','600钻石',100,600),('gold_64800','64800钻石',64800,64800);

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
