**PRODUCT CATALOGUE MANAGEMENT API**

This project provides a backend API for managing products, suppliers, categories, stock movements, and users.  
It includes two roles:
- system_admin  
- supplier_admin  ( also known as users)

**SYSTEM ADMIN RESPONSIBILITIES**
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

**SUPPLIER ADMIN RESPONSIBILITIES**
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

**IMPORTANT ROLE SEPARATION**
- System admin controls the entire platform (suppliers, categories, users, full dashboard).
- Supplier admin controls only their supplier’s or supplier admins belonging to same company's products or and stock (restricted by supplier_id in JWT token).
- This ensures data isolation and prevents suppliers from accessing or modifying each other’s data.

**DATABASE SETUP**
CREATE DATABASE product_catalogue;

You do not need to create any tables.  
All tables and sample data are automatically created by migration/init.sql when running the project.


**ENVIRONMENT VARIABLES (.env)**
Create a ".env" file in the project root with the following:

DSN="root:<your password>@tcp(127.0.0.1:3306)/product_catalogue?parseTime=true"
JWT_SECRET_KEY=<your-generated-secret">  //To generate secret key run go run tools/secret_key/main.go
PORT=8080
LOW_STOCK_ALERT=10


**RUN THE PROJECT**
1. Run the server:
   go run main.go
2. API will run on:
   http://localhost:8080


**DEFAULT SYSTEM ADMIN LOGIN**
A default system admin user is auto-inserted from init.sql:

Email: yuvrajbisht41@gmail.com
Password: Yuvraj@2411


**SYSTEM_ADMIN API TESTING GUIDE**

**LOGIN** 
API -> POST to http://localhost:8080/auth/login  
Login as system_admin (enter with email password given above)
JSON REQUIRED -> (email, password)

curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
        "email": "your_email_here",
        "password": "your_password_here"
      }'


**CREATE SUPPLIER**
API -> POST to http://localhost:8080/suppliers
Create Supplier
JSON REQUIRED -> (name, contact_info, email, company)
Headers -> Authorization: Bearer <TOKEN>

curl -X POST http://localhost:8080/suppliers \
     -H "Authorization: Bearer <TOKEN>" \
     -H "Content-Type: application/json" \
     -d '{"name":"ABC Pvt Ltd","contact_info":"9876543210","email":"abc@gmail.com","company":"ABC"}'


**GET SUPPLIERS**
API -> GET http://localhost:8080/suppliers  
View all suppliers
Headers -> Authorization: Bearer <TOKEN>

curl -X GET http://localhost:8080/suppliers \
     -H "Authorization: Bearer <TOKEN>"


**GET SUPPLIER BY ID**
API -> GET http://localhost:8080/suppliers/:id  
View supplier by ID
Headers -> Authorization: Bearer <TOKEN>
- Add any int in place of :id

curl -X GET http://localhost:8080/suppliers/:id \
     -H "Authorization: Bearer <TOKEN>"


**UPDATE SUPPLIER**
API -> PUT http://localhost:8080/suppliers/:id  
Update supplier  
Headers -> Authorization: Bearer <TOKEN>
JSON REQUIRED -> (name, contact_info, email, company)
- can update any number of fields or all fields
- Add any int in place of :id

curl -X PUT http://localhost:8080/suppliers/:id \
     -H "Authorization: Bearer <TOKEN>" \
     -H "Content-Type: application/json" \
     -d '{"company":"Updated Pvt Ltd"}'


**DELETE SUPPLIER**
API -> DELETE http://localhost:8080/suppliers/:id  
Delete supplier
Headers -> Authorization: Bearer <TOKEN>
- Add any int in place of :id

curl -X DELETE http://localhost:8080/suppliers/:id \
     -H "Authorization: Bearer <TOKEN>"


**CREATE SUPPLIER ADMIN (USER)**
API -> POST http://localhost:8080/users/supplier-admin  
Create supplier_admin user  
Headers -> Authorization: Bearer <TOKEN>
JSON REQUIRED -> (name, email, password, supplier_id)

curl -X POST http://localhost:8080/users/supplier-admin \
     -H "Authorization: Bearer <TOKEN>" \
     -H "Content-Type: application/json" \
     -d '{"name":"any name","email":"abc@gmail.com","password":"abc@1234","supplier_id":any integer}'


**GET USERS**
API -> GET http://localhost:8080/users/supplier-admin  
List supplier admins
Headers -> Authorization: Bearer <TOKEN>

curl -X GET http://localhost:8080/users/supplier-admin \
     -H "Authorization: Bearer <TOKEN>"


**GET USER BY ID**
API -> GET http://localhost:8080/users/supplier-admin/:id  
Get supplier_admin by ID
Headers -> Authorization: Bearer <TOKEN>
- Add any int in place of :id

curl -X GET http://localhost:8080/users/supplier-admin/:id \
     -H "Authorization: Bearer <TOKEN>"


**DELETE USER**
API -> GET http://localhost:8080/users/supplier-admin/:id  
Get supplier_admin by ID
Headers -> Authorization: Bearer <TOKEN>
- Add any int in place of :id

curl -X GET http://localhost:8080/users/supplier-admin/:id \
     -H "Authorization: Bearer <TOKEN>"


**CREATE CATEGORY**
API -> POST http://localhost:8080/categories  
Create category  
Headers -> Authorization: Bearer <TOKEN>
JSON REQUIRED -> category_name  
OPTIONAL -> category_description

curl -X POST http://localhost:8080/categories \
     -H "Authorization: Bearer <TOKEN>" \
     -H "Content-Type: application/json" \
     -d '{"category_name":"ABC", "category_description":"any string"}' //description should contain atleast one letter and it is not mandatory


**GET CATEGORY**
API -> GET http://localhost:8080/categories  
List categories
Headers -> Authorization: Bearer <TOKEN>

curl -X GET http://localhost:8080/categories \
     -H "Authorization: Bearer <TOKEN>"


**GET CTAEGORY BY ID**
API -> GET http://localhost:8080/categories/:id  
Get category by ID
Headers -> Authorization: Bearer <TOKEN>
- Add any int in place of :id

curl -X GET http://localhost:8080/categories/:id \
     -H "Authorization: Bearer <TOKEN>"


**UPDATE CATEGORY**
API -> PUT http://localhost:8080/categories/:id  
Update category 
Headers -> Authorization: Bearer <TOKEN> 
JSON REQUIRED -> category_name, category_description
- can update any number of fields or all fields
- Add any int in place of :id

curl -X PUT http://localhost:8080/categories/:id \
     -H "Authorization: Bearer <TOKEN>" \
     -H "Content-Type: application/json" \
     -d '{"category_name":"abc" "category_description":"Updated"}' 


**DELETE CATEGORY**
API -> DELETE http://localhost:8080/categories/:id  
Delete category
Headers -> Authorization: Bearer <TOKEN>
- Add any int in place of :id

curl -X DELETE http://localhost:8080/categories/:id \
     -H "Authorization: Bearer <TOKEN>"


**GET PRODUCTS**
API -> GET http://localhost:8080/products  
View all products
Headers -> Authorization: Bearer <TOKEN>

curl -X GET http://localhost:8080/products \
     -H "Authorization: Bearer <TOKEN>"


**GET PRODUCTS BY ID**
API -> GET http://localhost:8080/products/:id  
View product by ID
Headers -> Authorization: Bearer <TOKEN>
- Add any int in place of :id

curl -X GET http://localhost:8080/products/:id \
     -H "Authorization: Bearer <TOKEN>"


**GET STOCK_MOVEMENTS**
API -> GET http://localhost:8080/stock_movements  
View all stock movements
Headers -> Authorization: Bearer <TOKEN>

curl -X GET http://localhost:8080/stock_movements \
     -H "Authorization: Bearer <TOKEN>"


**GET DASHBOARD**
API -> GET http://localhost:8080/dashboard  
Dashboard summary
Headers -> Authorization: Bearer <TOKEN>

curl -X GET http://localhost:8080/dashboard \
     -H "Authorization: Bearer <TOKEN>"




**SUPPLIER ADMIN API REQUESTS**  
(Requires login as supplier_admin (user))

**LOGIN**
API -> POST http://localhost:8080/auth/login  
Login as supplier_admin  
JSON REQUIRED -> (email, password)

curl -X POST http://localhost:8080/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email":"supplier@gmail.com","password":"Supplier@123"}'


**GET SUPPLIERS**
API -> GET to http://localhost:8080/suppliers  
View suppliers (supplier_admin can view suppliers of their company)  
Headers -> Authorization: Bearer <TOKEN>

curl -X GET http://localhost:8080/suppliers \
     -H "Authorization: Bearer <TOKEN>"


**GET SUPPLIER BY ID**
API -> GET to http://localhost:8080/suppliers/:id  
View a supplier by ID (only if it belongs to the supplier_admin's company)  
Headers -> Authorization: Bearer <TOKEN>
- Add any int in place of :id

curl -X GET http://localhost:8080/suppliers/:id \
     -H "Authorization: Bearer <TOKEN>"


**GET CATEGORY**
API -> GET to http://localhost:8080/categories  
List all categories  
Headers -> Authorization: Bearer <TOKEN>

curl -X GET http://localhost:8080/categories \
     -H "Authorization: Bearer <TOKEN>"


**GET CATEGORY BY ID**
API -> GET to http://localhost:8080/categories/:id  
Get category details by ID  
Headers -> Authorization: Bearer <TOKEN>
- Add any int in place of :id

curl -X GET http://localhost:8080/categories/:id \
     -H "Authorization: Bearer <TOKEN>"


**GET PRODUCTS**
API -> GET to http://localhost:8080/products  
List products visible to this user (supplier_admin sees products for their supplier belonging to same company)  
Headers -> Authorization: Bearer <TOKEN>

curl -X GET http://localhost:8080/products \
     -H "Authorization: Bearer <TOKEN>"


**GET PRODUCTS BY ID**
API -> GET to http://localhost:8080/products/:id  
Get product by ID (only if product belongs to supplier who is in same company)  
Headers -> Authorization: Bearer <TOKEN>
- Add any int in place of :id

curl -X GET http://localhost:8080/products/:id \
     -H "Authorization: Bearer <TOKEN>"


**CREATE PRODUCT**
API -> POST to http://localhost:8080/products  
Create a new product for the logged-in supplier (supplier_id is taken from token)  
JSON REQUIRED -> (product_name,product_category_id, product_cost)  
JSON OPTIONAL -> (product_description, discount_type, discount_value)  
Note -> If discount_type is provided, discount_value is required. If discount_type == "percent" then discount_value must be <= 100. Discount type can only be flat or percent.


curl -X POST http://localhost:8080/products \
     -H "Authorization: Bearer <TOKEN>" \
     -H "Content-Type: application/json" \
     -d '{
           "product_name":"Screen",
           "product_description":"Big LED",
           "product_cost":4999.00,
           "product_category_id":2,
           "discount_type":"percent",
           "discount_value":10
         }'


**UPDATE PRODUCT**
API -> PUT to http://localhost:8080/products/:id  
Update an existing product (supplier_admin can update only their products)  
JSON OPTIONAL -> (product_name, product_cost, product_category_id, discount_type, discount_value, product_description)  
Notes -> You can update any field. If discount_type is present you must include discount_value (<=100 for percent). discount_type allowed values: "percent" or "flat".
- Add any int in place of :id

curl -X POST http://localhost:8080/products/:id \
     -H "Authorization: Bearer <TOKEN>" \
     -H "Content-Type: application/json" \
     -d '{
           "product_name":"Earpods",
           "product_cost":5499.00,
           "discount_type":"flat",
           "discount_value":300
         }'


**GET STOCK_MOVEMENTS**
API -> GET to http://localhost:8080/stock_movements  
View stock movements (supplier_admin sees only movements for the products created by supplier admin of same company)  
Headers -> Authorization: Bearer <TOKEN>  

curl -X GET "http://localhost:8080/stock_movements" \
     -H "Authorization: Bearer <TOKEN>"



**CREATE STOCK_MOVEMENTS**
API -> POST to http://localhost:8080/stock_movements  
Create stock movement for a product owned by the supplier admin
JSON REQUIRED -> (product_id, quantity, movement_type)  
JSON OPTIONAL -> (reason)  
Notes → movement_type must be "IN" or "OUT". Quantity must be > 0. OUT will be rejected if it would drive stock below LOW_STOCK_ALERT (10 for this project).

curl -X POST http://localhost:8080/stock_movements \
     -H "Authorization: Bearer <TOKEN>" \
     -H "Content-Type: application/json" \
     -d '{
           "product_id":any integer,
           "quantity": any integer,
           "movement_type":"IN",
           "reason":"Restock batch abc"
         }'


**GET DASHBOARD**
API -> GET to http://localhost:8080/dashboard  
View supplier-specific dashboard (totals and low-stock for the supplier logged in and other supplier admins of same company )  
Headers -> Authorization: Bearer <TOKEN>

curl -X GET http://localhost:8080/dashboard \
     -H "Authorization: Bearer <TOKEN>"



**INFO**
1. Replace <TOKEN> with the JWT from /auth/login.
2. Supplier_admin must not include product_supplier_id in create/update bodies — the server uses the supplier_id from the JWT



















