
# Sales Analytics Backend
A Go-based backend system for performing sales data analytics using PostgreSQL and CSV import.
Provides APIs for analyzing revenue, customer segmentation, profit margins, and more.

## Prerequisites

- **Go:** v1.18+
- **PostgreSQL:** v12+
- **Git**

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/sales-analytics/sales-analytics.git
cd sales-analytics/SalesAnalytics
```

### 2. Configure the Database

Update `dbconfig.toml` with your PostgreSQL credentials:

```toml
[database]
PostgresServer = "localhost"
PostgresPort = 5432
PostgresUser = "root"
PostgresPassword = "root"
PostgresDatabase = "postgres"
PostgresDBType = "postgres"
```

### 3. Initialize & Run

```bash
go mod tidy
go run main.go
```

App will start on: `http://localhost:8080`


## CSV Upload

Place your sales CSV file in the `data/` directory. The system loads it automatically when the application starts or when a refresh is triggered.



## API Reference

| Endpoint                         | Method | Params                        |  Header       | Description                                                      
|----------------------------------|--------|-------------------------------|---------------|-------------------------------------------                       
| `/api/revenue/dateRange`         | GET    | `start`, `end` (YYYY-MM-DD)   |   None        | Total revenue for a date range                                  
| `/api/revenue/product`           | GET    | `start`, `end` (YYYY-MM-DD)   |   None        | Revenue by product                                                
| `/api/revenue/category`          | GET    | `start`, `end` (YYYY-MM-DD)   |   None        | Revenue by category                                              
| `/api/revenue/region`            | GET    | `start`, `end` (YYYY-MM-DD)   |   None        | Revenue by region                                                

| `/api/nProducts/overall`         | GET    | `start`, `end` (YYYY-MM-DD)   |  NVALUE = ?   | Top N products overall on quantity sold (Within a date range)     
| `/api/nProducts/category`        | GET    | `start`, `end` (YYYY-MM-DD)   |  NVALUE = ?   | Top N products category on quantity sold (Within a date range)    
| `/api/nProducts/region`          | GET    | `start`, `end` (YYYY-MM-DD)   |  NVALUE = ?   | Top N products region on quantity sold   (Within a date range)    

| `/api/customers/dateRange`       | GET    | `start`, `end` (YYYY-MM-DD)   |   None        | Total Number of Customers (Within a date range)                  
| `/api/customers/orders`          | GET    | `start`, `end` (YYYY-MM-DD)   |   None        | Total Number of Orders (Within a date range)                  
| `/api/customers/avgOrderValue`   | GET    | `start`, `end` (YYYY-MM-DD)   |   None        | Average Order Value (Within a date range)                        

| `/api/calculations/profitMargin` | GET    | `start`, `end` (YYYY-MM-DD)   |   None        | Profit Margin by Product (Within a date range)                  
| `/api/calculations/clv`          | GET    | `start`, `end` (YYYY-MM-DD)   |   None        | Customer Lifetime Value (Within a date range)                     
| `/api/calculations/customerSeg`  | GET    | `start`, `end` (YYYY-MM-DD)   |   None        | Customer Segmentation (Within a date range)                        

| `/api/refresh`                   | POST   | None                          |   None        | Reloads CSV data to DB                                          

## Data Refresh
- The refresh API reloads data from CSV to DB.
- Logs are stored in the `refresh_logs` table.

## Logging
Data refresh attempts and their status are logged in the `refresh_logs` table.

```sql
SELECT * FROM refresh_logs ORDER BY timestamp DESC;
```

## Database details in separate file
database_schema.sql

## Author
Maintained by [Kavin kishore.S]. 


## Sample API Requests 
1. API: Total revenue for a date range 
Method: GET
Route: /api/revenue/dateRange?start=YYYY-MM-DD&end=YYYY-MM-DD
Host: localhost:8080

sample request:
http://localhost:8080/api/revenue/dateRange?start=2024-01-03&end=2024-02-28

**Sample Response** (Success):
```json
{
    "status": "S",
    "msg": "Data fetched successfully!!",
    "total_revenue": 1478.97
}
```

**Sample Response** (error):
```json
{
    "status": "E",
    "msg": "Error while fetching"
}
```

2. API: Revenue by product 
Method: GET
Route: /api/revenue/product?start=YYYY-MM-DD&end=YYYY-MM-DD
Host: localhost:8080

sample request:
http://localhost:8080/api/revenue/product?start=2024-01-03&end=2024-02-28

**Sample Response** (Success):
```json
{
    "status": "S",
    "msg": "Data fetched successfully!!",
    "total_product_revenue": [
        {
            "product_name": "iPhone 15 Pro",
            "revenue": 1299
        },
        {
            "product_name": "Levi's 501 Jeans",
            "revenue": 179.97
        }
    ]
}
```

**Sample Response** (error):
```json
{
    "status": "E",
    "msg": "Error while fetching"
}
```


3. API: Revenue by category 
Method: GET
Route: /api/revenue/category?start=YYYY-MM-DD&end=YYYY-MM-DD
Host: localhost:8080

sample request:
http://localhost:8080/api/revenue/category?start=2024-01-03&end=2024-02-28

**Sample Response** (Success):
```json
{
    "status": "S",
    "msg": "Data fetched successfully!!",
    "total_category_revenue": [
        {
            "category": "Clothing",
            "revenue": 179.97
        },
        {
            "category": "Electronics",
            "revenue": 1299
        }
    ]
}
```

**Sample Response** (error):
```json
{
    "status": "E",
    "msg": "Error while fetching"
}
```

4. API: Revenue by region   
Method: GET
Route: /api/revenue/region?start=YYYY-MM-DD&end=YYYY-MM-DD
Host: localhost:8080

sample request:
http://localhost:8080/api/revenue/region?start=2024-01-03&end=2024-02-28

**Sample Response** (Success):
```json
{
    "status": "S",
    "msg": "Data fetched successfully!!",
    "total_region_revenue": [
        {
            "region": "Asia",
            "revenue": 179.97
        },
        {
            "region": "Europe",
            "revenue": 1299
        }
    ]
}
```

**Sample Response** (error):
```json
{
    "status": "E",
    "msg": "Error while fetching"
}
```


5. API: Top N products overall on quantity sold (Within a date range) 
Method: GET
Route: /api/nProducts/overall?start=YYYY-MM-DD&end=YYYY-MM-DD
Host: localhost:8080
Header: NVALUE = ?

sample request:
http://localhost:8080/api/nProducts/overall?start=2024-01-03&end=2024-02-28

**Sample Response** (Success):
```json
{
    "status": "S",
    "msg": "Data fetched successfully!!",
    "top_N_Products": [
        {
            "product_name": "Levi's 501 Jeans",
            "quantity": 3
        },
        {
            "product_name": "iPhone 15 Pro",
            "quantity": 1
        }
    ]
}
```

**Sample Response** (error):
```json
{
    "status": "E",
    "msg": "Error while fetching"
}
```


6. API: Top N products category on quantity sold (Within a date range)
Method: GET
Route: /api/nProducts/category?start=YYYY-MM-DD&end=YYYY-MM-DD
Host: localhost:8080
Header: NVALUE = ?

sample request:
http://localhost:8080/api/nProducts/category?start=2024-01-03&end=2024-02-28

**Sample Response** (Success):
```json
{
    "status": "S",
    "msg": "Data fetched successfully!!",
    "top_N_Products": [
        {
            "category": "Clothing",
            "quantity": 3
        },
        {
            "category": "Electronics",
            "quantity": 1
        }
    ]
}
```

**Sample Response** (error):
```json
{
    "status": "E",
    "msg": "Error while fetching"
}
```


7. API: Top N products region on quantity sold   (Within a date range) 
Method: GET
Route: /api/nProducts/region?start=YYYY-MM-DD&end=YYYY-MM-DD
Host: localhost:8080
Header: NVALUE = ?

sample request:
http://localhost:8080/api/nProducts/region?start=2024-01-03&end=2024-02-28

**Sample Response** (Success):
```json
{
    "status": "S",
    "msg": "Data fetched successfully!!",
    "top_N_Products": [
        {
            "region": "Asia",
            "quantity": 3
        },
        {
            "region": "Europe",
            "quantity": 1
        }
    ]
}
```

**Sample Response** (error):
```json
{
    "status": "E",
    "msg": "Error while fetching"
}
```

8. API: Total Number of Customers (Within a date range) 
Method: GET
Route: /api/customers/dateRange?start=YYYY-MM-DD&end=YYYY-MM-DD
Host: localhost:8080

sample request:
http://localhost:8080/api/customers/dateRange?start=2024-01-03&end=2024-02-28

**Sample Response** (Success):
```json
{
    "status": "S",
    "msg": "Data fetched successfully!!",
    "customer_analysis": {
        "total_customers": 2
    }
}

**Sample Response** (error):
```json
{
    "status": "E",
    "msg": "Error while fetching"
}
```


9. API: Total Number of Orders (Within a date range) 
Method: GET
Route: /api/customers/orders?start=YYYY-MM-DD&end=YYYY-MM-DD
Host: localhost:8080

sample request:
http://localhost:8080/api/customers/orders?start=2024-01-03&end=2024-02-28

**Sample Response** (Success):
```json
{
    "status": "S",
    "msg": "Data fetched successfully!!",
    "customer_analysis": {
        "total_orders": 2
    }
}

**Sample Response** (error):
```json
{
    "status": "E",
    "msg": "Error while fetching"
}
```


10. API: Average Order Value (Within a date range)   
Method: GET
Route: /api/customers/avgOrderValue?start=YYYY-MM-DD&end=YYYY-MM-DD
Host: localhost:8080

sample request:
http://localhost:8080/api/customers/avgOrderValue?start=2024-01-03&end=2024-02-28

**Sample Response** (Success):
```json
{
    "status": "S",
    "msg": "Data fetched successfully!!",
    "customer_analysis": {
        "avg_order_value": 749.59
    }
}

**Sample Response** (error):
```json
{
    "status": "E",
    "msg": "Error while fetching"
}
```


11. API: Profit Margin by Product (Within a date range)
Method: GET
Route: /api/calculations/profitMargin?start=YYYY-MM-DD&end=YYYY-MM-DD
Host: localhost:8080

sample request:
http://localhost:8080/api/calculations/profitMargin?start=2024-01-03&end=2024-02-28

**Sample Response** (Success):
```json
{
    "status": "S",
    "msg": "Data fetched successfully!!",
    "calculations": [
        {
            "name": "iPhone 15 Pro",
            "profitMargin": 98.85
        },
        {
            "name": "Levi's 501 Jeans",
            "profitMargin": 97.11
        }
    ]
}
```

**Sample Response** (error):
```json
{
    "status": "E",
    "msg": "Error while fetching"
}
```


12. API: Customer Lifetime Value (Within a date range)  
Method: GET
Route: /api/calculations/clv?start=YYYY-MM-DD&end=YYYY-MM-DD
Host: localhost:8080

sample request:
http://localhost:8080/api/calculations/clv?start=2024-01-03&end=2024-02-28

**Sample Response** (Success):
```json
{
    "status": "S",
    "msg": "Data fetched successfully!!",
    "calculations": [
        {
            "name": "Emily Davis",
            "customer_lifetime_value": 1299
        },
        {
            "name": "John Smith",
            "customer_lifetime_value": 179.97
        }
    ]
}
```

**Sample Response** (error):
```json
{
    "status": "E",
    "msg": "Error while fetching"
}
```

13. API:  Customer Segmentation (Within a date range)  
Method: GET
Route: /api/calculations/customerSeg?start=YYYY-MM-DD&end=YYYY-MM-DD
Host: localhost:8080

sample request:
http://localhost:8080/api/calculations/customerSeg?start=2024-01-03&end=2024-02-28

**Sample Response** (Success):
```json
{
    "status": "S",
    "msg": "Data fetched successfully!!",
    "calculations": [
        {
            "name": "Emily Davis",
            "total_orders": 1,
            "total_spent": 1299,
            "avg_order_value": 1299,
            "last_order_date": "2024-01-03T00:00:00Z"
        },
        {
            "name": "John Smith",
            "total_orders": 1,
            "total_spent": 179.97,
            "avg_order_value": 179.97,
            "last_order_date": "2024-02-28T00:00:00Z"
        }
    ]
}
```

**Sample Response** (error):
```json
{
    "status": "E",
    "msg": "Error while fetching"
}
```

14. API:  Reloads and refreshes CSV data to DB      
Method: POST
Route: /api/refresh
Host: localhost:8080

sample request:
http://localhost:8080/api/refresh

**Sample Response** (Success):
```json
{
    "status": "S",
    "msg": "File loaded successfully!!"
}
```

**Sample Response** (error):
```json
{
    "status": "E",
    "msg": "Error while fetching"
}
```