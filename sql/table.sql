drop table if exists asset;
create table asset (
    id bigint(20) unsigned not null auto_increment, # ID
    imgsrc varchar(512) not null,         			# img path 
    videosrc varchar(1024) not null,         		# video path 
    primary key(id),
	unique(imgsrc)
) engine=InnoDB default charset=utf8;

alter table asset add count int(10) unsigned default 0 not null;
alter table asset add rand_count int(64) unsigned default 0 not null;

create table ip_addr (
    id bigint(20) unsigned not null auto_increment, # ID
    ip_addr varchar(512) not null,         			# ip addr 
    time int(10) not null,         					# 时间 
    primary key(id)
) engine=InnoDB default charset=utf8;


