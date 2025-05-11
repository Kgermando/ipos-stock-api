package dashboard

import (
	"sort"
	"time"
)

// CalculateRevenueGrowth calculates the revenue growth percentage between two periods.
func CalculateRevenueGrowth(currentRevenue, previousRevenue float64) float64 {
	if previousRevenue == 0 {
		return 0
	}
	return ((currentRevenue - previousRevenue) / previousRevenue) * 100
}

// CalculateAverageDailySales calculates the average daily sales over a given period.
func CalculateAverageDailySales(totalRevenue float64, startDate, endDate time.Time) float64 {
	days := endDate.Sub(startDate).Hours() / 24
	if days == 0 {
		return 0
	}
	return totalRevenue / days
}

// CalculateSalesConversionRate calculates the sales conversion rate.
func CalculateSalesConversionRate(totalSales, totalLeads float64) float64 {
	if totalLeads == 0 {
		return 0
	}
	return (totalSales / totalLeads) * 100
}

// CalculateSalesByCategory calculates sales grouped by product category.
func CalculateSalesByCategory(salesData map[string]float64) map[string]float64 {
	// salesData is a map where the key is the category and the value is the sales amount.
	return salesData
}

// CalculateSalesByRegion calculates sales grouped by region.
func CalculateSalesByRegion(salesData map[string]float64) map[string]float64 {
	// salesData is a map where the key is the region and the value is the sales amount.
	return salesData
}

// GetTopProductsByRevenue returns the top N products by revenue.
func GetTopProductsByRevenue(productSales map[string]float64, topN int) []string {
	type product struct {
		Name  string
		Sales float64
	}
	var products []product
	for name, sales := range productSales {
		products = append(products, product{Name: name, Sales: sales})
	}

	// Sort products by sales in descending order
	sort.Slice(products, func(i, j int) bool {
		return products[i].Sales > products[j].Sales
	})

	// Extract the top N product names
	var topProducts []string
	for i := 0; i < topN && i < len(products); i++ {
		topProducts = append(topProducts, products[i].Name)
	}
	return topProducts
}

// CalculateContributionPercentage calculates the percentage contribution of each product to total sales.
func CalculateContributionPercentage(productSales map[string]float64, totalRevenue float64) map[string]float64 {
	contribution := make(map[string]float64)
	if totalRevenue == 0 {
		return contribution
	}
	for product, sales := range productSales {
		contribution[product] = (sales / totalRevenue) * 100
	}
	return contribution
}

// CalculateCumulativeRevenue calculates the cumulative revenue over a period.
func CalculateCumulativeRevenue(dailySales []float64) []float64 {
	cumulative := make([]float64, len(dailySales))
	var total float64
	for i, sales := range dailySales {
		total += sales
		cumulative[i] = total
	}
	return cumulative
}

// CalculateHourlySales calculates the sales curve over 24 hours for a given day.
func CalculateHourlySales(hourlySalesData map[int]float64) []float64 {
	// hourlySalesData is a map where the key is the hour (0-23) and the value is the sales amount.
	hourlySales := make([]float64, 24)
	for hour, sales := range hourlySalesData {
		if hour >= 0 && hour < 24 {
			hourlySales[hour] = sales
		}
	}
	return hourlySales
}

// CalculateTotalSales calculates the total sales from a list of sales amounts.
func CalculateTotalSales(sales []float64) float64 {
	var total float64
	for _, sale := range sales {
		total += sale
	}
	return total
}

// CalculateTotalTransactions calculates the total number of transactions.
func CalculateTotalTransactions(transactions []float64) int {
	return len(transactions)
}

// CalculateAverageSalesPerTransaction calculates the average sales per transaction.
func CalculateAverageSalesPerTransaction(totalRevenue float64, totalTransactions int) float64 {
	if totalTransactions == 0 {
		return 0
	}
	return totalRevenue / float64(totalTransactions)
}

// CalculateSuccessfulTransactionRate calculates the percentage of successful transactions.
func CalculateSuccessfulTransactionRate(successfulTransactions, totalTransactions int) float64 {
	if totalTransactions == 0 {
		return 0
	}
	return (float64(successfulTransactions) / float64(totalTransactions)) * 100
}
