package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"sort"
	"time"

	AthenaCall "my-lambda-app/athena"
	dynamo "my-lambda-app/dynamo"
	"my-lambda-app/env"
	helperFunctions "my-lambda-app/helperFunctions"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Portfolio struct {
	Symbol string
	Amount int
}

func HandleBackTest() events.APIGatewayProxyResponse {
	// Loop through back data Data
	// Define back test Data Loop
	client := dynamo.GetClient()
	// date := "2017-02-14"
	// exampleArray := []string{"2017-05-15", "2017-08-14", "2017-11-14", "2018-02-14"}
	// counter := 0
	// Stopping 02-14 Not all data is up to 5-13
	// TODO: Dynamicaly Set the array

	// predictionArray0 := []string{"2012-02-14", "2012-05-14", "2012-08-14", "2012-11-14", "2013-02-14", "2013-05-14", "2013-08-14", "2013-11-14",
	// 	"2014-02-14", "2014-05-14", "2014-08-14", "2014-11-14", "2015-02-13", "2015-05-14", "2015-08-14", "2015-11-13",
	// 	"2016-02-12", "2016-05-13", "2016-08-15", "2016-11-14", "2017-02-14", "2017-05-15", "2017-08-14", "2017-11-14",
	// 	"2018-02-14", "2018-05-14", "2018-08-14", "2018-11-14", "2019-02-14", "2019-05-14", "2019-08-14", "2019-11-14",
	// 	"2020-02-14", "2020-05-14", "2020-08-14", "2020-11-13", "2021-02-12", "2021-05-14", "2021-08-13", "2021-11-15",
	// 	"2022-02-14", "2022-05-13", "2022-08-15", "2022-11-14", "2023-02-14", "2023-05-12", "2023-08-14", "2023-11-14",
	// 	"2024-02-14", "2024-05-13", "2024-07-23", "2024-10-21", "2025-01-19"}

	// Prediction for 2016-2019 (jan1 -jan1 )
	// predictionArray1 := []string{
	// 	"2016-02-12", "2016-05-13", "2016-08-15", "2016-11-14", "2017-02-14", "2017-05-15", "2017-08-14", "2017-11-14",
	// 	"2018-02-14", "2018-05-14", "2018-08-14", "2018-11-14", "2019-02-14", "2019-05-14", "2019-08-14", "2019-11-14",
	// 	"2020-02-14"}

	// Prediction for 2021-2024 (jan1 -jan1 )
	// predictionArray2 := []string{"2021-02-12", "2021-05-14", "2021-08-13", "2021-11-15",
	// 	"2022-02-14", "2022-05-13", "2022-08-15", "2022-11-14", "2023-02-14", "2023-05-12", "2023-08-14", "2023-11-14",
	// 	"2024-02-14", "2024-05-13", "2024-07-23", "2024-10-21", "2025-01-19"}

	predictionArray := []string{
		"2022-02-14", "2022-05-13", "2022-08-15", "2022-11-14", "2023-02-14", "2023-05-12", "2023-08-14", "2023-11-14",
		"2024-02-14", "2024-05-13", "2024-07-23", "2024-10-21", "2025-01-19"}

	// Loop over data ( build first with just one data set)
	// call the SQL athena function to get actual values
	// TODO: Loop every 3 intergers start Date End Date - no longer Hard Coded
	// TODO: Need to gather
	// TODO compare dates and if future date is greater then current date, move date itterator
	predictionIterator := 0
	predictionString := predictionArray[predictionIterator]

	dateArray := helperFunctions.ReverseSlice(env.DateArray)
	currentPredictionFlag := false

	// Convert current date TIme string into a Date time
	predictionTime, err := time.Parse("2006-01-02", predictionString)
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < len(dateArray); i += 5 {

		dateStringDaily := dateArray[i]

		// Convert current date TIme string into a Date time
		dateTimeDaily, err := time.Parse("2006-01-02", dateStringDaily)
		if err != nil {
			fmt.Println(err)
		}

		// Compare times to adjust
		if dateTimeDaily.Unix() >= predictionTime.Unix() {
			predictionIterator += 1
			predictionString = predictionArray[predictionIterator]

			predictionTime, err = time.Parse("2006-01-02", predictionString)
			if err != nil {
				fmt.Println(err)
			}

		}
		currentStockPrices := AthenaCall.SQL_date_Price(dateArray[i])

		// Call Get Current Portfolio distribution
		currentPortfolio := dynamo.GetPortfolio("historicalTestPortfolio", client)
		totalCash := 0.0

		// Is there a currentPortfolio
		if len(*currentPortfolio) == 0 {
			// Get current cash by accessing cashTable
			currentCash := dynamo.GetCashTotal("historicalCashBackTest", client)
			totalCash = currentCash.Amount
		} else {
			currentCash := dynamo.GetCashTotal("historicalCashBackTest", client)
			totalCash = currentCash.Amount + CalculateCashTotal(currentPortfolio, currentStockPrices)
		}

		// Get Predictions  in the range of Dynamo table (Eliminate all stocks that are below .5 bias)
		predictions := TradingWeight(*currentStockPrices, predictionArray, predictionIterator, client, context.TODO())

		// Call UpdatePortfolio - creates a new portfolio
		updatedPortfolio, cashUpdate := UpdatePortfolio(totalCash, &predictions, currentPredictionFlag, client)

		// rectify_portfolio - buy and sell data into new table - save logs in portfolio transaction
		buys, sells := RectifyPortfolio(dateArray[i], currentPortfolio, &updatedPortfolio)

		// save new Portfolio distribution in dynamo table
		stringPortfolio := ConvertPortToString(updatedPortfolio)
		err = dynamo.PutCurrentPortfolio("historicalTestPortfolio", stringPortfolio, client)
		if err != nil {
			log.Fatalf("Failed to PutCurrentPortfolio: %v", err)
		}

		// Update Transaction Trade Logs
		// Update Cash Table
		dynamo.PutCurrentCash("historicalCashBackTest", cashUpdate, client)
		if err != nil {
			log.Fatalf("BatchWriteTransactionHistory Sells: %v", err)
		}

		// Update Sells
		if len(*sells) != 0 {
			dynamo.BatchWriteTransactionHistory(*sells, client)
			if err != nil {
				log.Fatalf("BatchWriteTransactionHistory Sells: %v", err)
			}
		}

		// Update Buys
		dynamo.BatchWriteTransactionHistory(*buys, client)
		if err != nil {
			log.Fatalf("BatchWriteTransactionHistory Sells: %v", err)
		}

	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       jsonMessage("BackTest process started"),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}
}

func CalculateCashTotal(currentPortfolio *[]dynamo.PortfolioDistribution, currentStockPrices *[]map[string]string) float64 {
	totalCash := 0.0
	portfolioMap := make(map[string]float64)
	for i := 0; i < len(*currentPortfolio); i++ {
		portfolioMap[(*currentPortfolio)[i].Symbol] = (*currentPortfolio)[i].Amount
	}

	for _, record := range *currentStockPrices {
		amount, ok := portfolioMap[record["symbol"]]
		if !ok {
			continue
		}

		fmt.Println("In portfolio:", record["symbol"])
		totalCash = totalCash + (amount * helperFunctions.FloatConverter(record["close"]))
	}

	return totalCash

}

func TradingWeight(currentStocks []map[string]string, dates []string, counter int, client *dynamodb.Client, ctx context.Context) map[string]dynamo.PredictionData {
	predictedData := make(map[string]dynamo.PredictionData)
	for i := 0; i < 4; i++ {
		predictions, err := dynamo.ScanPredictions(dates[counter+i], client, ctx)
		if err != nil {
			fmt.Println("Query API call failed:", err)
			return predictedData
		}
		for _, prediction := range *predictions {
			// Median Distance
			if prediction.Bias < 0.45 {
				continue
			}
			// // Mean Distance
			// if prediction.MeanDistance < 0 {
			// 	continue
			// }
			// // Model Is not accurate enough
			// if prediction.DefaultError > 30.0 {
			// 	continue
			// }

			data, ok := predictedData[prediction.Symbol]
			if !ok {
				priceModifier := .1 * prediction.Price
				data = dynamo.PredictionData{
					Symbol:         prediction.Symbol,
					PredictedPrice: priceModifier,
					Bias:           prediction.Bias,
					MedianDistance: prediction.MedianDistance,
				}
			}
			priceModifier := .1 * prediction.Price
			// 1 year out (best Prediction)
			if i == 3 {
				priceModifier = .7 * prediction.Price
			}
			data.PredictedPrice += priceModifier
			predictedData[prediction.Symbol] = data
		}
	}

	for symbol, data := range predictedData {
		predictionBias := data.Bias
		// if predictionBias < .5 {
		// 	predictionBias += .5
		// }
		// Date modifier medianDistance
		data.PredictedPrice += (predictionBias * 2 * data.MedianDistance)
		predictedData[symbol] = data

	}

	for _, stock := range currentStocks {
		// Find Current Stock in map
		data, ok := predictedData[stock["symbol"]]
		if !ok {
			continue
		}
		originalPrice := helperFunctions.FloatConverter(stock["close"])
		data.OriginalPrice = originalPrice
		fmt.Printf("Retrieved stock: %v\n", stock)
		MovingAverage200Days := helperFunctions.FloatConverter(stock["movingaverage200days"])
		MovingAverage50Days := helperFunctions.FloatConverter(stock["movingaverage50days"])
		// Golden Cross boost
		if MovingAverage200Days > MovingAverage50Days {
			data.PredictedPrice = data.PredictedPrice * .85
		} else {
			data.PredictedPrice = data.PredictedPrice * 1.15
		}

		data.GrowthRate = data.PredictedPrice / originalPrice
		predictedData[stock["symbol"]] = data

	}
	return predictedData
}

func UpdatePortfolio(totalCash float64, predictions *map[string]dynamo.PredictionData, currentPredictionFlag bool, client *dynamodb.Client) ([]Portfolio, float64) {
	// 16 + 14 + 14 + 12 + 10 + 10 + 8 + 6 + 6 + 4
	weights := []float64{.16, .14, .14, .12, .10, .10, .08, .06, .06, .04}
	var newPortfolio []Portfolio
	cashBalance := totalCash
	allPredictions := make([]dynamo.PredictionData, 0, len(*predictions))

	// Extracting all PredictionData from the map
	for _, pd := range *predictions {
		allPredictions = append(allPredictions, pd)
	}

	// Sorting the slice based on the GrowthRate field
	sort.Slice(allPredictions, func(i, j int) bool {
		// sort descending
		return allPredictions[i].GrowthRate > allPredictions[j].GrowthRate
	})

	fmt.Println("hit Update Portfolio")
	if currentPredictionFlag {
		fmt.Println("hit Update CurrentPredictionFlag")
		err := dynamo.BatchWritePredictionWeights(allPredictions, client)
		if err != nil {
			log.Fatalf("BatchWritePredictionWeights Sells: %v", err)
		}
	}

	// Take top 10 share distribution
	portfolioMaxBalance := min(10, (len(*predictions)))
	for i := 0; i < portfolioMaxBalance; i++ {
		shares := (totalCash * weights[i]) / allPredictions[i].OriginalPrice
		floorNumber := math.Floor(shares) // Take the floor of the number
		shareNumber := int(floorNumber)   // Convert the floor number to int

		// Update Cash
		cashLoss := floorNumber * allPredictions[i].OriginalPrice
		cashBalance += (-cashLoss)

		// Make sure Cash never goes Negative when buying stock
		if cashBalance < 0 {
			cashBalance += cashLoss
			fmt.Println("shares")
			fmt.Println("shares")
			return newPortfolio, cashBalance
		}

		temp := Portfolio{
			Symbol: allPredictions[i].Symbol,
			Amount: shareNumber,
		}

		newPortfolio = append(newPortfolio, temp)
	}
	return newPortfolio, cashBalance

}

func RectifyPortfolio(date string, oldPortfolio *[]dynamo.PortfolioDistribution, currentPortfolio *[]Portfolio) (*[]dynamo.Orders, *[]dynamo.Orders) {
	// Create New map
	tradeMap := make(map[string]int)
	var buys []dynamo.Orders
	var sells []dynamo.Orders

	// Loop Over Current Portfolio - add positive weights to Map
	for _, new := range *currentPortfolio {
		tradeMap[new.Symbol] = new.Amount
	}
	// Loop over Old Portfolio and subtract from Map
	for _, old := range *oldPortfolio {
		tradeMap[old.Symbol] += int(-old.Amount)
	}

	// Loop over Map
	// append to Sell Orders if Negative and Buy Orders if Positive
	for key, value := range tradeMap {
		temp := dynamo.Orders{
			Symbol:    key,
			Date:      date,
			OrderType: "Buy",
			Amount:    value,
		}
		if value > 0 {
			buys = append(buys, temp)
		} else {
			temp.OrderType = "Sell"
			sells = append(sells, temp)
		}
	}

	return &buys, &sells

}

func ConvertPortToString(portfolio []Portfolio) []string {
	var stringPortfolio []string
	for _, item := range portfolio {
		temp := fmt.Sprintf("%s:%d", item.Symbol, item.Amount)

		stringPortfolio = append(stringPortfolio, temp)
	}
	return stringPortfolio
}
