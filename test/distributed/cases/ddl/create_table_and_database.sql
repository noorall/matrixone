drop database if exists test01;
create database test01;
drop database if exists test02;
create database test02;
drop database if exists test03;
create database test03;
drop database if exists test04;
create database test04;
drop database if exists test05;
create database test05;
show databases;

use test01;
drop table if exists test01;
create table test01(col1 int, col2 char);
drop table if exists test02;
create table test02(col1 int, col2 char);
drop table if exists test03;
create table test03(col1 int, col2 char);
drop table if exists test04;
create table test04(col1 int, col2 char);
drop table if exists test05;
create table test05(col1 int, col2 char);
drop table if exists test06;
create table test06(col1 int, col2 char);
show tables;

drop table test01;
drop table test02;
drop table test03;
drop table test04;
drop table test05;
drop table test06;
drop database test01;
drop database test02;
drop database test03;
drop database test04;
drop database test05;