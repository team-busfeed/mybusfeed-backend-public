package main

import (
	"encoding/json"
	"fmt"
	"math"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	
	"go.mongodb.org/mongo-driver/bson"
	
	"context"
	"log"
)

type DatamallResponse struct {
	OdataMetadata string `json:"odata.metadata"`
	BusStopCode   string `json:"BusStopCode"`
	Services      []struct {
		ServiceNo string `json:"ServiceNo"`
		Operator  string `json:"Operator"`
		NextBus   struct {
			OriginCode       string    `json:"OriginCode"`
			DestinationCode  string    `json:"DestinationCode"`
			EstimatedArrival time.Time `json:"EstimatedArrival"`
			Latitude         string    `json:"Latitude"`
			Longitude        string    `json:"Longitude"`
			VisitNumber      string    `json:"VisitNumber"`
			Load             string    `json:"Load"`
			Feature          string    `json:"Feature"`
			Type             string    `json:"Type"`
		} `json:"NextBus"`
		NextBus2 struct {
			OriginCode       string    `json:"OriginCode"`
			DestinationCode  string    `json:"DestinationCode"`
			EstimatedArrival time.Time `json:"EstimatedArrival"`
			Latitude         string    `json:"Latitude"`
			Longitude        string    `json:"Longitude"`
			VisitNumber      string    `json:"VisitNumber"`
			Load             string    `json:"Load"`
			Feature          string    `json:"Feature"`
			Type             string    `json:"Type"`
		} `json:"NextBus2"`
		NextBus3 struct {
			OriginCode       string    `json:"OriginCode"`
			DestinationCode  string    `json:"DestinationCode"`
			EstimatedArrival time.Time `json:"EstimatedArrival"`
			Latitude         string    `json:"Latitude"`
			Longitude        string    `json:"Longitude"`
			VisitNumber      string    `json:"VisitNumber"`
			Load             string    `json:"Load"`
			Feature          string    `json:"Feature"`
			Type             string    `json:"Type"`
		} `json:"NextBus3"`
	} `json:"Services"`
}

type DatamallResponseBus struct {
	OdataMetadata string `json:"odata.metadata"`
	BusStopCode   string `json:"BusStopCode"`
	Services      []struct {
		ServiceNo string `json:"ServiceNo"`
		Operator  string `json:"Operator"`
		NextBus   struct {
			OriginCode       string    `json:"OriginCode"`
			DestinationCode  string    `json:"DestinationCode"`
			EstimatedArrival time.Time `json:"EstimatedArrival"`
			Latitude         string    `json:"Latitude"`
			Longitude        string    `json:"Longitude"`
			VisitNumber      string    `json:"VisitNumber"`
			Load             string    `json:"Load"`
			Feature          string    `json:"Feature"`
			Type             string    `json:"Type"`
		} `json:"NextBus"`
		NextBus2 struct {
			OriginCode       string    `json:"OriginCode"`
			DestinationCode  string    `json:"DestinationCode"`
			EstimatedArrival time.Time `json:"EstimatedArrival"`
			Latitude         string    `json:"Latitude"`
			Longitude        string    `json:"Longitude"`
			VisitNumber      string    `json:"VisitNumber"`
			Load             string    `json:"Load"`
			Feature          string    `json:"Feature"`
			Type             string    `json:"Type"`
		} `json:"NextBus2"`
		NextBus3 struct {
			OriginCode       string    `json:"OriginCode"`
			DestinationCode  string    `json:"DestinationCode"`
			EstimatedArrival time.Time `json:"EstimatedArrival"`
			Latitude         string    `json:"Latitude"`
			Longitude        string    `json:"Longitude"`
			VisitNumber      string    `json:"VisitNumber"`
			Load             string    `json:"Load"`
			Feature          string    `json:"Feature"`
			Type             string    `json:"Type"`
		} `json:"NextBus3"`
	} `json:"Services"`
}

type BusStopResponse struct {
	OdataMetadata string `json:"odata.metadata"`
	Value         []struct {
		BusstopNo   string  `json:"BusStopCode"`
		RoadName    string  `json:"RoadName"`
		Description string  `json:"Description"`
		Latitude    float64 `json:"Latitude"`
		Longitude   float64 `json:"Longitude"`
	} `json:"value"`
}

type BusStop struct {
	Value []struct {
		BusstopNo   string  `json:"BusStopCode"`
		RoadName    string  `json:"RoadName"`
		Description string  `json:"Description"`
		Latitude    float64 `json:"Latitude"`
		Longitude   float64 `json:"Longitude"`
	} `json:"value"`
}

type DatamallBusesResponse struct {
	Services []struct {
		ServiceNo string `json:"serviceNo"`
	} `json:"services"`
}

type BusesResponse struct {
	Services []struct {
		ServiceNo string `json:"service_no"`
	} `json:"services"`
}

func healthcheck(c *fiber.Ctx) {
	c.Set("Content-type", "application/json; charset=utf-8")
	fmt.Println("Successful heartbeat")
	mapD := map[string]string{"type": "success", "status": "healthy"}
	json, _ := json.Marshal(mapD)
	c.Status(200).Send(json)
}

func getBusesRoute(c *fiber.Ctx) {
	buses := getBusesAPI(c.Params("busStopNo"))

	if len(buses) == 0 {
		c.Status(400).JSON(&fiber.Map{
			"message": "Invalid request",
		})
	} else {
		c.Status(200).JSON(&fiber.Map{
			"services": buses,
		})
	}

}

func getBusesAPI(busStopNo string) (buses []string) {
	client := &http.Client{}

	// Parameters
	parameters := url.Values{}
	parameters.Add("BusStopCode", busStopNo)

	req, err := http.NewRequest("GET", "http://datamall2.mytransport.sg/ltaodataservice/BusArrivalv2?"+parameters.Encode(), nil)
	if err != nil {
		fmt.Print(err.Error())
	}

	// Headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("AccountKey", "6FebnHNxTv+PGH/NDKkf/Q==")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Print(err.Error())
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	var result DatamallBusesResponse

	json.Unmarshal(data, &result)

	for _, service := range BusesResponse(result).Services {
		buses = append(buses, service.ServiceNo)
	}

	return buses
}

func getListOfBusstopNo(c *fiber.Ctx) {
	c.Set("Content-type", "application/json; charset=utf-8")

	const dbName = "BusFeedHardware"
	const collectionName = "BusStopLocation"

	collection, err := getMongoDbCollection(dbName, collectionName)
	if err != nil {
		c.Status(500).Send(err)
		return
	}


	latitude := c.Params("latitude")
	longitude := c.Params("longitude")

	// var location_info []location_info

	// client := &http.Client{}
	var filter bson.M = bson.M{}

	var BusStops []bson.M

	cur, err := collection.Find(context.Background(), filter)

	if err != nil {
		c.Status(500).Send(err)
		return
	}

	defer cur.Close(context.Background())

	cur.All(context.Background(), &BusStops)
	
	if len(BusStops) == 0 {
		outputJSON := map[string]string{"Message": "Fail to retrieve bus stop information"}
		json, _ := json.Marshal(outputJSON)
		c.Status(500).Send(json)
	} else {
		var busInfo string
		var myarray = []map[string]string{}

		for _, data := range BusStops {


			lat := data["latitude"]
			latstr := fmt.Sprint(lat)

			long := data["longitude"]
			longstr := fmt.Sprint(long)

			description := data["description"]

			// Convert the latitude & longitude to a whole number to manipulate the first 3dp
			// DB Location Data
			newArrayLat := strings.Split(latstr, ".")
			dplat := newArrayLat[0] + newArrayLat[1][:3]
			newArrayLong := strings.Split(longstr, ".")
			dplong := newArrayLong[0] + newArrayLong[1][:3]

			// Input Data and modify the input value to get the upper and lower range
			searchArrayLat := strings.Split(latitude, ".")
			currentLat := searchArrayLat[1][:3]
			latData, _ := strconv.Atoi(currentLat)
			upperRangeLat := latData + 3
			lowerRangeLat := latData - 3
			searchdplatupper := searchArrayLat[0] + fmt.Sprint(upperRangeLat)
			searchdplatlower := searchArrayLat[0] + fmt.Sprint(lowerRangeLat)

			searchArraylong := strings.Split(longitude, ".")
			currentLong := searchArraylong[1][:3]
			longData, _ := strconv.Atoi(currentLong)
			upperRangeLong := longData + 3
			lowerRangeLong := longData - 3

			searchdplongupper := searchArraylong[0] + fmt.Sprint(upperRangeLong)
			searchdplonglower := searchArraylong[0] + fmt.Sprint(lowerRangeLong)

			if dplat >= searchdplatlower && dplat <= searchdplatupper && dplong >= searchdplonglower && dplong <= searchdplongupper {

				// Calculate the euclidean distance
				searchLat, _ := strconv.ParseFloat(dplat, 64)
				searchLong, _ := strconv.ParseFloat(dplong, 64)
				actualLat, _ := strconv.ParseFloat(searchArrayLat[0]+currentLat, 64)
				actualLong, _ := strconv.ParseFloat(searchArraylong[0]+currentLong, 64)

				difference := int(math.Abs(math.Sqrt(math.Pow((searchLat-actualLat), 2) + math.Pow((searchLong-actualLong), 2))))
				differenceNew := fmt.Sprint(difference)
				result := data["busstop_no"]
				busStopNum := fmt.Sprint(result)
				desc := fmt.Sprint(description)
				desc = strings.ReplaceAll(desc, "%20", " ")

				temp_map := map[string]string{"busstop_number": busStopNum, "busstop_name": desc, "busstop_lat": fmt.Sprint(lat), "busstop_lng": fmt.Sprint(long), "compare": differenceNew}
				myarray = append(myarray, temp_map)

				busInfo += busStopNum + ","
			}
		}

		// Sorting function to sort the distance based on asc order
		sort.Slice(myarray, func(i, j int) bool {
			compare1, _ := strconv.Atoi(myarray[i]["compare"])
			compare2, _ := strconv.Atoi(myarray[j]["compare"])
			return compare1 < compare2
		})
		fmt.Println(myarray)

		if busInfo == "" {
			didNotFindResponse := map[string]string{"type": "success", "status": "not_found", "message": "No nearby bus stops."}
			responseJSON, _ := json.Marshal(didNotFindResponse)
			c.Status(200).Send(responseJSON)
		} else {
			json, _ := json.Marshal(myarray)
			c.Status(200).Send(json)
		}
	}
}

// Get bus stop information based on busstop description or busstop number
func getBusStopInformation(c *fiber.Ctx) {
	bustopNum := c.Params("busstopnum")

	const dbName = "BusFeedHardware"
	const collectionName = "BusStopLocation"

	collection, err := getMongoDbCollection(dbName, collectionName)
	if err != nil {
		c.Status(500).Send(err)
		return
	}

	var filter bson.M = bson.M{}

	var BusStops []bson.M

	cur, err := collection.Find(context.Background(), filter)

	if err != nil {
		c.Status(500).Send(err)
		return
	}

	defer cur.Close(context.Background())

	cur.All(context.Background(), &BusStops)



	if len(BusStops) == 0 {
		outputJSON := map[string]string{"type": "success", "status": "not_found", "message": "Incorrect busstop number or description"}
		json, _ := json.Marshal(outputJSON)
		c.Status(404).Send(json)
	} else {
		var myarray []interface{}
		for _, data := range BusStops {
			description := data["description"]
			descriptionString := fmt.Sprint(description)
			busNumber := data["busstop_no"]
			busNumberString := fmt.Sprint(busNumber)
			latitude := data["latitude"]
			longitude := data["longitude"]
			if strings.Contains(busNumberString, bustopNum) {
				outputJSON := map[string]string{"busstop_number": busNumberString, "busstop_name": descriptionString, "busstop_lat": fmt.Sprint(latitude), "busstop_lng": fmt.Sprint(longitude)}
				myarray = append(myarray, outputJSON)
			}
		}

		if len(myarray) == 0 {
			for _, data := range BusStops {
				description := data["description"]
				descString := fmt.Sprint(description)
				busNumber := data["busstop_no"]
				busNumString := fmt.Sprint(busNumber)
				latitude := data["latitude"]
				longitude := data["longitude"]
				bustopNum = strings.Replace(bustopNum, "%20", " ", -1)
				descString = strings.Replace(descString, "%20", " ", -1)
				if strings.Contains(strings.ToLower(descString), strings.ToLower(bustopNum)) {
					outputJSON := map[string]string{"busstop_number": busNumString, "busstop_name": descString, "busstop_lat": fmt.Sprint(latitude), "busstop_lng": fmt.Sprint(longitude)}
					myarray = append(myarray, outputJSON)
				}
			}
		}
		if len(myarray) > 0 {
			if len(myarray) > 10 {
				json, _ := json.Marshal(myarray[:10])
				c.Status(200).Send(json)
			} else {
				json, _ := json.Marshal(myarray)
				c.Status(200).Send(json)
			}

		} else {
			outputJSON := map[string]string{"type": "success", "status": "not_found", "message": "Incorrect busstop number or description"}
			json, _ := json.Marshal(outputJSON)
			c.Status(404).Send(json)
		}
	}
}

// based on bus stop number return all the bus stop information with the lat and long
func returnBusStopInformation(c *fiber.Ctx) {
	c.Set("Content-type", "application/json; charset=utf-8")

	var BusStop map[string]interface{}

	if err := c.BodyParser(&BusStop); err != nil {
		c.Status(503).Send(err)
	} else {
		var allBusStop []interface{}
		for _, t := range BusStop {
			// Convert the output to json string an manipulate the data
			s := fmt.Sprint(t)
			data := s[1 : len(s)-1]
			if data == "" {
				didNotFindResponse := map[string]string{"type": "success", "status": "not_found", "message": "Incorrect busstop number or description"}
				json, _ := json.Marshal(didNotFindResponse)
				c.Status(500).Send(json)
			} else {
				list_data := strings.Split(data, " ")
				// var location_info []location_info
				const dbName = "BusFeedHardware"
				const collectionName = "BusStopLocation"
			
				collection, err := getMongoDbCollection(dbName, collectionName)
				if err != nil {
					c.Status(500).Send(err)
					return
				}
			
				var filter bson.M = bson.M{}
			
				var BusStops []bson.M
			
				cur, err := collection.Find(context.Background(), filter)
			
				if err != nil {
					c.Status(500).Send(err)
					return
				}
			
				defer cur.Close(context.Background())
			
				cur.All(context.Background(), &BusStops)
				for _, busData := range list_data {
					fmt.Println(busData)
					for _, busstop := range BusStops {
						descString := fmt.Sprint(busstop["description"])
						descString = strings.Replace(descString, "%20", " ", -1)
						outputJSON := map[string]string{"busstop_number": fmt.Sprint(busstop["busstop_no"]), "busstop_name": descString, "busstop_lat": fmt.Sprint(busstop["latitude"]), "busstop_lng": fmt.Sprint(busstop["longitude"])}
						if (fmt.Sprint(busstop["busstop_no"]) == busData) {
							allBusStop = append(allBusStop, outputJSON)
						}
					}
				}
			
				if len(allBusStop) == 0 {
					didNotFindResponse := map[string]string{"type": "success", "status": "not_found", "message": "Incorrect busstop number or description"}
					json, _ := json.Marshal(didNotFindResponse)
					c.Status(500).Send(json)
				} else {
					json, _ := json.Marshal(allBusStop)
					c.Status(200).Send(json)
				}
			}
		}
	}
}




func setupRoutes(app *fiber.App) {
	app.Get("/demand/healthcheck", healthcheck)
	app.Get("/demand/bus-stop/:busStopNo", getBusesRoute)
	app.Get("/location/getListOfBusStop/:latitude-:longitude", getListOfBusstopNo)
	app.Get("/location/getBusStopInformation/:busstopnum", getBusStopInformation)
	app.Post("/location/returnBusStopInformation", returnBusStopInformation)
}

//GetMongoDbConnection get connection of mongodb
func GetMongoDbConnection() (*mongo.Client, error) {

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb+srv://BusFeedUser:SMUbusfeed!08@busfeed.jsui8.mongodb.net/test"))

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	return client, nil
}

func getMongoDbCollection(DbName string, CollectionName string) (*mongo.Collection, error) {
	client, err := GetMongoDbConnection()

	if err != nil {
		return nil, err
	}

	collection := client.Database(DbName).Collection(CollectionName)

	return collection, nil
}

func newApp() *fiber.App {
	app := fiber.New()
	app.Use(cors.New())
	setupRoutes(app)
	return app
}

func main() {
	port := os.Getenv("PORT")
	err := newApp().Listen(port)
	if err != nil {
		panic(err)
	}
}
