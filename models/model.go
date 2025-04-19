package models

type RESPONSE struct {
	Status           string          `json:"status"`
	Msg              string          `json:"msg"`
	TotRevenue       *float64        `json:"total_revenue,omitempty"`
	TotProdRevenue   *[]REVPRODUCTS  `json:"total_product_revenue,omitempty"`
	TotCategRevenue  *[]REVCATEGORY  `json:"total_category_revenue,omitempty"`
	TotRegionRevenue *[]REVREGION    `json:"total_region_revenue,omitempty"`
	TopNProducts     *[]TOPNPRODUCTS `json:"top_N_Products,omitempty"`
	Customers        *CUSTOMERS      `json:"customer_analysis,omitempty"`
	Calculations     *[]CALCUALTIONS `json:"calculations,omitempty"`
}

type REVPRODUCTS struct {
	Name    string  `json:"product_name"`
	Revenue float64 `json:"revenue"`
}

type REVCATEGORY struct {
	Category string  `json:"category"`
	Revenue  float64 `json:"revenue"`
}

type REVREGION struct {
	Region  string  `json:"region"`
	Revenue float64 `json:"revenue"`
}

type TOPNPRODUCTS struct {
	Name     *string `json:"product_name,omitempty"`
	Category *string `json:"category,omitempty"`
	Region   *string `json:"region,omitempty"`
	Qty      int     `json:"quantity"`
}

type CUSTOMERS struct {
	TotCustomers  *int     `json:"total_customers,omitempty"`
	TotOrders     *int     `json:"total_orders,omitempty"`
	AvgOrderValue *float64 `json:"avg_order_value,omitempty"`
}

type CALCUALTIONS struct {
	Name         *string  `json:"name,omitempty"`
	ProfitMargin *float64 `json:"profitMargin,omitempty"`
	CLV          *float64 `json:"customer_lifetime_value,omitempty"`
	// CustomerSeg  *float64 `json:"avg_order_value,omitempty"`
	Total_orders    *int     `json:"total_orders,omitempty"`
	Total_spent     *float64 `json:"total_spent,omitempty"`
	Avg_order_value *float64 `json:"avg_order_value,omitempty"`
	Last_order_date *string  `json:"last_order_date,omitempty"`
}
