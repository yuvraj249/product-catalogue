create database if not exists product_catalogue;
use product_catalogue;

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


-- example input

-- insert ignore into suppliers (name, contact_info, email, company)
-- values
-- ('Messi', 'messi@suppliers.com', 'messi@suppliers.com', 'messi Corp'),
-- ('Ronaldo', 'beta@suppliers.com', 'ronaldo@suppliers.com', 'ronaldo Group'),
-- ('Neymar', 'gamma@suppliers.com', 'neymar@suppliers.com', 'neymar Industries');

-- insert ignore into categories (category_name, category_description)
-- values
-- ('Electronics', 'Devices, gadgets, and electrical equipment'),
-- ('Clothing', 'Men and Women fashion wear'),
-- ('Home Appliances', 'Kitchen and home utility products'),
-- ('Books', 'Educational and fictional books');

-- insert ignore into users (name, email, password_hash, role, supplier_id)
-- values
-- ('Messi Admin', 'messi@suppliers.com', '$2a$10$v0UwrHMxpNFgkzG7pZ.7/uUU814u7udGO8jNgMCbzH2NsSvdX0P5y', 'supplier_admin', 1),
-- ('Ronaldo Admin', 'ronaldo@suppliers.com', '$2a$10$v0UwrHMxpNFgkzG7pZ.7/uUU814u7udGO8jNgMCbzH2NsSvdX0P5y', 'supplier_admin', 2),
-- ('Neymar Admin', 'neymar@suppliers.com', '$2a$10$v0UwrHMxpNFgkzG7pZ.7/uUU814u7udGO8jNgMCbzH2NsSvdX0P5y', 'supplier_admin', 3);


-- insert ignore into products (product_name, product_description, product_cost, product_category_id, product_supplier_id, discount_type, discount_value)
-- values
-- ('Wireless Mouse', 'Bluetooth-enabled ergonomic mouse', 499.00, 1, 1, 'percent', 10),
-- ('Laptop Backpack', 'Waterproof backpack with multiple compartments', 1599.00, 2, 2, 'flat', 100),
-- ('Air Conditioner', '1.5 Ton Split AC, Energy efficient', 35999.00, 3, 3, 'percent', 5),
-- ('Smartphone', '5G-enabled Android smartphone', 18999.00, 1, 1, 'flat', 500),
-- ('Fiction Novel', 'Bestselling mystery novel', 499.00, 4, 2, 'percent', 15);

-- insert ignore into stock_movements (product_id, quantity, movement_type, reason, performed_by)
-- values
-- (1, 50, 'IN', 'Initial stock', 2),
-- (2, 30, 'IN', 'Initial stock', 3),
-- (3, 20, 'IN', 'Initial stock', 4),
-- (1, 5, 'OUT', 'Sold to customer', 2),
-- (4, 40, 'IN', 'New shipment arrived', 2),
-- (5, 10, 'OUT', 'Damaged returns', 3);





