
create table if not exists suppliers(
    supplier_id int auto_increment primary key,
    name varchar(225) not null,
    contact_info varchar(225) default null,
    email varchar(225) unique not null,
    company varchar(225) not null
);

create table if not exists categories(
   category_id int auto_increment primary key,
   category_name varchar(225) not null,
   category_description text
   
);

create table if not exists users(
    user_id int auto_increment primary key,
    name varchar(50) not null,
    email varchar(100) unique not null,
    password_hash varchar(255) not null,
    role enum('system_admin', 'supplier_admin') not null default 'supplier_admin',
    supplier_id int default null,
    constraint fk_supplier_id foreign key (supplier_id)references suppliers(supplier_id) on delete set null on update cascade
);

create table if not exists products(
    product_id int auto_increment primary key,
    product_name varchar(225) not null,
    product_description text,
    product_cost decimal(12,2) default 0,
    product_category_id int default null,
    product_supplier_id int default null,
    discount_type enum('percent', 'flat') default null,
    discount_value decimal(12,2) default 0,
    constraint fk_product_cat foreign key(product_category_id) references categories(category_id) on delete set null on update cascade,
    constraint fk_product_sup foreign key(product_supplier_id) references suppliers(supplier_id) on delete set null on update cascade
);

create table if not exists stock_movements(
    stock_id int auto_increment primary key,
    product_id int not null,
    quantity int not null default 1,
    movement_type enum('IN', 'OUT') not null,
    reason varchar(225),
    performed_by int not null,
    foreign key(product_id) references products(product_id) on delete cascade ,
    foreign key(performed_by) references users(user_id) on delete restrict  

); 

insert ignore into users(name,email,password_hash,role) values('Owner', 'yuvrajbisht41@gmail.com', '$2a$10$CCw/Xx/.lW1BCcc0MYIH5.xh2QJq7pBqMrWeE.WPxgRI4F8Af12s2', 'system_admin');


