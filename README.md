# PRODUCT CATALOGUE MANAGEMENT API #

- This project provides a backend API for managing products, suppliers, categories, stock movements, and users.<br>  
It includes two roles:
- system_admin  
- supplier_admin  (also known as users)

# SYSTEM ADMIN RESPONSIBILITIES #

- Create new suppliers and update existing supplier details.
- Get/Read all suppliers or Get suppliers by their ID.
- Delete any supplier by their ID from the database.
- Create or Update categories.
- Get/Read all categories or Get categories by their ID.
- Delete categories by their ID.
- Create supplier_admin user accounts from the supplier created. Can Create multiple supplier admins (users) for one supplier.
- Get/Read all users or Get users by their ID.
- Delete any supplier_admin user by their ID.
- Access the entire dashboard including total suppliers, total categories, total products, and low-stock alerts.
- View stock movements for all suppliers and all products.
- Cannot Create,Update or Delete products or cannot create stock movements because system_admin is not tied to any supplier.

# SUPPLIER ADMIN RESPONSIBILITIES #

- Get/Read suppliers or Get suppliers by their ID but these suppliers belong to the same company
- Create products for their supplier (supplier_id is taken automatically from token).
- Can update their own products but cannot change the product_supplier_id.
- Can delete only their own products.
- Can get/view only the products that belong to the supplier-admin (users) of their company.
- Can create stock movements (IN/OUT) for their products only.
- Cannot create stock movements for products belonging to other suppliers.
- Cannot reduce stock below the configured LOW_STOCK_ALERT value.
- Can view stock movement history for their supplier’s and suppliers belonging to same company products only.
- Cannot create,update or delete suppliers, categories, or users because these are system-level actions.
- Can view Dashboard

# IMPORTANT ROLE SEPARATION #

- System admin controls the entire platform (suppliers, categories, users, full dashboard).
- Supplier admin controls only their supplier’s or supplier admins belonging to same company's products or and stock (restricted by supplier_id in JWT token).
- This ensures data isolation and prevents suppliers from accessing or modifying each other’s data.

# DATABASE SETUP #

- CREATE DATABASE product_catalogue;
- You do not need to create any tables.
- All tables and sample data are automatically created by migration/init.sql when running the project.

# ENVIRONMENT VARIABLES (.env) #<br>

Create a ".env" file in the project root with the following:<br>

DSN="root:"your-password"@tcp(127.0.0.1:3306)/product_catalogue?parseTime=true"<br>
JWT_SECRET_KEY=<your-generated-secret">  //To generate secret key run go run tools/secret_key/main.go<br>
PORT=8080<br>
LOW_STOCK_ALERT=10<br>

# RUN THE PROJECT #

1. Run the server:
   go run main.go
2. API will run on:
   <http://localhost:8080>

# DEFAULT SYSTEM ADMIN LOGIN #

- A default system admin user is auto-inserted from init.sql:<br>

Email: <yuvrajbisht41@gmail.com>
Password: Yuvraj@2411

# SYSTEM_ADMIN API TESTING GUIDE #

## LOGIN ##

API -> POST to <http://localhost:8080/auth/login>  <br>
Login as system_admin (enter with email password given above)<br>
JSON REQUIRED -> (email, password)<br>

curl -X POST <http://localhost:8080/auth/login> \ <br>
  -H "Content-Type: application/json" \ <br>
  -d '{<br>
        "email": "your_email_here",<br>
        "password": "your_password_here"<br>
      }'<br>

## CREATE SUPPLIER ##

API -> POST to <http://localhost:8080/suppliers> <br>
Create Supplier <br>
JSON REQUIRED -> (name, contact_info, email, company) <br>
Headers -> Authorization: Bearer TOKEN <br>

curl -X POST <http://localhost:8080/suppliers> \ <br>
     -H "Authorization: Bearer TOKEN" \ <br>
     -H "Content-Type: application/json" \ <br>
     -d '{"name":"ABC Pvt Ltd","contact_info":"9876543210","email":"<abc@gmail.com>","company":"ABC"}'<br>

## GET SUPPLIERS ##

API -> GET <http://localhost:8080/suppliers>  <br>
View all suppliers <br>
Headers -> Authorization: Bearer TOKEN <br>

curl -X GET <http://localhost:8080/suppliers> \ <br>
     -H "Authorization: Bearer TOKEN" <br>

## GET SUPPLIER BY ID ##

API -> GET <http://localhost:8080/suppliers/:id>   <br>
View supplier by ID<br>
Headers -> Authorization: Bearer TOKEN <br>

- Add any int in place of :id<br>

curl -X GET <http://localhost:8080/suppliers/:id> \ <br>
     -H "Authorization: Bearer TOKEN" <br>

## UPDATE SUPPLIER ##

API -> PUT <http://localhost:8080/suppliers/:id>  <br>
Update supplier<br>
Headers -> Authorization: Bearer TOKEN <br>
JSON REQUIRED -> (name, contact_info, email, company)<br>

- can update any number of fields or all fields<br>
- Add any int in place of :id<br>

curl -X PUT <http://localhost:8080/suppliers/:id> \ <br>
     -H "Authorization: Bearer TOKEN" \ <br>
     -H "Content-Type: application/json" \ <br>
     -d '{"company":"Updated Pvt Ltd"}' <br>

## DELETE SUPPLIER ##

API -> DELETE <http://localhost:8080/suppliers/:id>  <br>
Delete supplier<br>
Headers -> Authorization: Bearer TOKEN <br>

- Add any int in place of :id<br>

curl -X DELETE <http://localhost:8080/suppliers/:id> \ <br>
     -H "Authorization: Bearer TOKEN" <br>

## CREATE SUPPLIER ADMIN (USER) ##

API -> POST <http://localhost:8080/users/supplier-admin>  <br>
Create supplier_admin user<br>
Headers -> Authorization: Bearer TOKEN <br>
JSON REQUIRED -> (name, email, password, supplier_id)<br>

curl -X POST <http://localhost:8080/users/supplier-admin> \ <br>
     -H "Authorization: Bearer TOKEN" \ <br>
     -H "Content-Type: application/json" \ <br>
     -d '{"name":"any name","email":"<abc@gmail.com>","password":"abc@1234","supplier_id":any integer}' <br>

## GET USERS ##

API -> GET <http://localhost:8080/users/supplier-admin>  <br>
List supplier admins<br>
Headers -> Authorization: Bearer TOKEN <br>

curl -X GET <http://localhost:8080/users/supplier-admin> \ <br>
     -H "Authorization: Bearer TOKEN" <br>

## GET USER BY ID ##

API -> GET <http://localhost:8080/users/supplier-admin/:id>  <br>
Get supplier_admin by ID<br>
Headers -> Authorization: Bearer TOKEN <br>

- Add any int in place of :id<br>

curl -X GET <http://localhost:8080/users/supplier-admin/:id> \ <br>
     -H "Authorization: Bearer TOKEN" <br>

## DELETE USER ##

API -> GET <http://localhost:8080/users/supplier-admin/:id>  <br>
Get supplier_admin by ID<br>
Headers -> Authorization: Bearer TOKEN <br>

- Add any int in place of :id <br>

curl -X GET <http://localhost:8080/users/supplier-admin/:id> \ <br>
     -H "Authorization: Bearer TOKEN" <br>

## CREATE CATEGORY ##

API -> POST <http://localhost:8080/categories>  <br>
Create category <br>
Headers -> Authorization: Bearer TOKEN <br>
JSON REQUIRED -> category_name <br>
OPTIONAL -> category_description<br>

curl -X POST <http://localhost:8080/categories> \ <br>
     -H "Authorization: Bearer TOKEN" \ <br>
     -H "Content-Type: application/json" \ <br>
     -d '{"category_name":"ABC", "category_description":"any string"}' //description should contain atleast one letter and it is not mandatory<br>

## GET CATEGORY ##

API -> GET <http://localhost:8080/categories>  <br>
List categories<br>
Headers -> Authorization: Bearer TOKEN <br>

curl -X GET <http://localhost:8080/categories> \ <br>
     -H "Authorization: Bearer TOKEN" <br>

## GET CTAEGORY BY ID ##

API -> GET <http://localhost:8080/categories/:id>  <br>
Get category by ID<br>
Headers -> Authorization: Bearer TOKEN <br>

- Add any int in place of :id<br>

curl -X GET <http://localhost:8080/categories/:id> \ <br>
     -H "Authorization: Bearer TOKEN" <br>

## UPDATE CATEGORY ##

API -> PUT <http://localhost:8080/categories/:id>  <br>
Update category<br>
Headers -> Authorization: Bearer TOKEN <br>
JSON REQUIRED -> category_name, category_description<br>

- can update any number of fields or all fields<br>
- Add any int in place of :id <br>

curl -X PUT <http://localhost:8080/categories/:id> \ <br>
     -H "Authorization: Bearer TOKEN" \ <br>
     -H "Content-Type: application/json" \ <br>
     -d '{"category_name":"abc" "category_description":"Updated"}' <br>

## DELETE CATEGORY ##

API -> DELETE <http://localhost:8080/categories/:id>  <br>
Delete category<br>
Headers -> Authorization: Bearer TOKEN <br>

- Add any int in place of :id<br>

curl -X DELETE <http://localhost:8080/categories/:id> \ <br>
     -H "Authorization: Bearer TOKEN" <br>

## GET PRODUCTS ##

API -> GET <http://localhost:8080/products> <br>
View all products<br>
Headers -> Authorization: Bearer TOKEN <br>

curl -X GET <http://localhost:8080/products> \ <br>
     -H "Authorization: Bearer TOKEN" <br>

## GET PRODUCTS BY ID ##

API -> GET <http://localhost:8080/products/:id>  <br>
View product by ID<br>
Headers -> Authorization: Bearer TOKEN <br>

- Add any int in place of :id <br>

curl -X GET <http://localhost:8080/products/:id> \ <br>
     -H "Authorization: Bearer TOKEN" <br>

## GET STOCK_MOVEMENTS ##

API -> GET <http://localhost:8080/stock_movements>  <br>
View all stock movements <br>
Headers -> Authorization: Bearer TOKEN <br>

curl -X GET <http://localhost:8080/stock_movements> \ <br>
     -H "Authorization: Bearer TOKEN" <br>

## GET STOCK_MOVEMENTS BY ID ##

API -> GET to <http://localhost:8080/stock_movements?product_id=:id> <br>
View stock movements for specific product <br>
Headers -> Authorization: Bearer TOKEN  <br>

curl -X GET "<http://localhost:8080/stock_movements?product_id=:id>" \ <br>
     -H "Authorization: Bearer TOKEN" <br>

## GET DASHBOARD ##

API -> GET <http://localhost:8080/dashboard>  <br>
Dashboard summary<br>
Headers -> Authorization: Bearer TOKEN <br>

curl -X GET <http://localhost:8080/dashboard> \ <br>
     -H "Authorization: Bearer TOKEN" <br>

# SUPPLIER ADMIN API REQUESTS #

(Requires login as supplier_admin (user))<br>

## LOGIN ##

API -> POST <http://localhost:8080/auth/login>  <br>
Login as supplier_admin<br>
JSON REQUIRED -> (email, password)<br>

curl -X POST <http://localhost:8080/auth/login> \ <br>
     -H "Content-Type: application/json" \ <br>
     -d '{"email":"<supplier@gmail.com>","password":"Supplier@123"}' <br>

## GET SUPPLIERS ##

API -> GET to <http://localhost:8080/suppliers>  <br>
View suppliers (supplier_admin can view suppliers of their company)<br>
Headers -> Authorization: Bearer TOKEN <br>

curl -X GET <http://localhost:8080/suppliers> \ <br>
     -H "Authorization: Bearer TOKEN" <br>

## GET SUPPLIER BY ID ##

API -> GET to <http://localhost:8080/suppliers/:id>  <br>
View a supplier by ID (only if it belongs to the supplier_admin's company)  <br>
Headers -> Authorization: Bearer TOKEN <br>

- Add any int in place of :id <br>

curl -X GET <http://localhost:8080/suppliers/:id> \ <br>
     -H "Authorization: Bearer TOKEN" <br>

## GET CATEGORY ##

API -> GET to <http://localhost:8080/categories>  <br>
List all categories  <br>
Headers -> Authorization: Bearer TOKEN <br>

curl -X GET <http://localhost:8080/categories> \ <br>
     -H "Authorization: Bearer TOKEN" <br>

## GET CATEGORY BY ID ##

API -> GET to <http://localhost:8080/categories/:id>  <br>
Get category details by ID<br>
Headers -> Authorization: Bearer TOKEN <br>

- Add any int in place of :id<br>

curl -X GET <http://localhost:8080/categories/:id> \ <br>
     -H "Authorization: Bearer TOKEN" <br>

## GET PRODUCTS ##

API -> GET to <http://localhost:8080/products>  <br>
List products visible to this user (supplier_admin sees products for their supplier belonging to same company)<br>
Headers -> Authorization: Bearer TOKEN <br>

curl -X GET <http://localhost:8080/products> \ <br>
     -H "Authorization: Bearer TOKEN" <br>

## GET PRODUCTS BY ID ##

API -> GET to <http://localhost:8080/products/:id>  <br>
Get product by ID (only if product belongs to supplier who is in same company)<br>
Headers -> Authorization: Bearer TOKEN <br>

- Add any int in place of :id <br>

curl -X GET <http://localhost:8080/products/:id> \ <br>
     -H "Authorization: Bearer TOKEN" <br>

## CREATE PRODUCT ##

API -> POST to <http://localhost:8080/products>  <br>
Create a new product for the logged-in supplier (supplier_id is taken from token)<br>
JSON REQUIRED -> (product_name,product_category_id, product_cost)<br>
JSON OPTIONAL -> (product_description, discount_type, discount_value)<br>  
Note -> If discount_type is provided, discount_value is required. If discount_type == "percent" then discount_value must be <= 100. Discount type can only be flat or percent.

curl -X POST <http://localhost:8080/products> \ <br>
     -H "Authorization: Bearer TOKEN" \ <br>
     -H "Content-Type: application/json" \ <br>
     -d '{ <br>
           "product_name":"Screen", <br>
           "product_description":"Big LED", <br>
           "product_cost":4999.00, <br>
           "product_category_id":2, <br>
           "discount_type":"percent", <br>
           "discount_value":10 <br>
         }' <br>

## UPDATE PRODUCT ##

API -> PUT to <http://localhost:8080/products/:id>  <br>
Update an existing product (supplier_admin can update only their products)<br>
JSON OPTIONAL -> (product_name, product_cost, product_category_id, discount_type, discount_value, product_description)<br>
Notes -> You can update any field. If discount_type is present you must include discount_value (<=100 for percent). discount_type allowed values: "percent" or "flat".<br>

- Add any int in place of :id<br>

curl -X POST <http://localhost:8080/products/:id> \ <br>
     -H "Authorization: Bearer TOKEN" \ <br>
     -H "Content-Type: application/json" \ <br>
     -d '{ <br>
           "product_name":"Earpods", <br>
           "product_cost":5499.00, <br>
           "discount_type":"flat", <br>
           "discount_value":300 <br>
         }' <br>

## GET STOCK_MOVEMENTS ##

API -> GET to <http://localhost:8080/stock_movements> <br>
View stock movements (supplier_admin sees only movements for the products created by supplier admin of same company)<br>
Headers -> Authorization: Bearer TOKEN  <br>

curl -X GET "<http://localhost:8080/stock_movements>" \ <br>
     -H "Authorization: Bearer TOKEN" <br>

## GET STOCK_MOVEMENTS BY ID ##

API -> GET to <http://localhost:8080/stock_movements?product_id=:id> <br>
View stock movements for specific product (supplier_admin sees only movements for product created by supplier admin of same company)<br>
Headers -> Authorization: Bearer TOKEN  <br>

curl -X GET "<http://localhost:8080/stock_movements?product_id=:id>" \ <br>
     -H "Authorization: Bearer TOKEN" <br>

## CREATE STOCK_MOVEMENTS ##

API -> POST to <http://localhost:8080/stock_movements> <br>
Create stock movement for a product owned by the supplier admin<br>
JSON REQUIRED -> (product_id, quantity, movement_type)<br>  
JSON OPTIONAL -> (reason)<br>  
Notes → movement_type must be "IN" or "OUT". Quantity must be > 0. OUT will be rejected if it would drive stock below LOW_STOCK_ALERT (10 for this project).<br>

curl -X POST <http://localhost:8080/stock_movements> \ <br>
     -H "Authorization: Bearer TOKEN" \ <br>
     -H "Content-Type: application/json" \ <br>
     -d '{ <br>
           "product_id":any integer, <br>
           "quantity": any integer, <br>
           "movement_type":"IN", <br>
           "reason":"Restock batch abc" <br>
         }'<br>

## GET DASHBOARD ##

API -> GET to <http://localhost:8080/dashboard>  <br>
View supplier-specific dashboard (totals and low-stock for the supplier logged in and other supplier admins of same company )<br>
Headers -> Authorization: Bearer TOKEN <br>

curl -X GET <http://localhost:8080/dashboard> \ <br>
     -H "Authorization: Bearer TOKEN" <br>

## INFO ##

1. Supplier_admin must not include product_supplier_id in create/update bodies — the server uses the supplier_id from the JWT<br>
