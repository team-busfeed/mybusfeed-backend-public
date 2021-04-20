package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
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

	latitude := c.Params("latitude")
	longitude := c.Params("longitude")

	var allBusStop [][]string
	// var location_info []location_info

	client := &http.Client{}

	// Parameters
	for i := 0; i < 5000+1; i += 500 {
		req, err := http.NewRequest("GET", "http://datamall2.mytransport.sg/ltaodataservice/BusStops?$skip="+fmt.Sprint(i), nil)
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

		var result BusStopResponse

		json.Unmarshal(data, &result)

		for _, busstop := range result.Value {
			var temp []string
			temp = append(temp, busstop.BusstopNo, busstop.Description, fmt.Sprint(busstop.Latitude), fmt.Sprint(busstop.Longitude))
			allBusStop = append(allBusStop, temp)
		}

	}

	if len(allBusStop) == 0 {
		outputJSON := map[string]string{"Message": "Fail to retrieve bus stop information"}
		json, _ := json.Marshal(outputJSON)
		c.Status(500).Send(json)
	} else {
		var busInfo string
		var myarray = []map[string]string{}

		for _, data := range allBusStop {

			lat := data[2]
			latstr := fmt.Sprint(lat)

			long := data[3]
			longstr := fmt.Sprint(long)

			description := data[1]

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
				result := data[0]
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

	// Load dataset into a list
	var allBusStop [][]string

	client := &http.Client{}

	// Parameters
	for i := 0; i < 5000+1; i += 500 {
		req, err := http.NewRequest("GET", "http://datamall2.mytransport.sg/ltaodataservice/BusStops?$skip="+fmt.Sprint(i), nil)
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

		var result BusStopResponse

		json.Unmarshal(data, &result)

		for _, busstop := range result.Value {
			var temp []string
			temp = append(temp, busstop.BusstopNo, busstop.Description, fmt.Sprint(busstop.Latitude), fmt.Sprint(busstop.Longitude))
			allBusStop = append(allBusStop, temp)
		}
	}

	if len(allBusStop) == 0 {
		outputJSON := map[string]string{"type": "success", "status": "not_found", "message": "Incorrect busstop number or description"}
		json, _ := json.Marshal(outputJSON)
		c.Status(404).Send(json)
	} else {
		var myarray []interface{}
		for _, data := range allBusStop {
			description := data[1]
			busNumber := data[0]
			latitude := data[2]
			longitude := data[3]
			if strings.Contains(busNumber, bustopNum) {
				outputJSON := map[string]string{"busstop_number": busNumber, "busstop_name": description, "busstop_lat": fmt.Sprint(latitude), "busstop_lng": fmt.Sprint(longitude)}
				myarray = append(myarray, outputJSON)
			}
		}

		if len(myarray) == 0 {
			for _, data := range allBusStop {
				description := data[1]
				busNumber := data[0]
				latitude := data[2]
				longitude := data[3]
				bustopNum = strings.Replace(bustopNum, "%20", " ", -1)
				if strings.Contains(strings.ToLower(description), strings.ToLower(bustopNum)) {
					print("HERE HER HERE")
					outputJSON := map[string]string{"busstop_number": busNumber, "busstop_name": description, "busstop_lat": fmt.Sprint(latitude), "busstop_lng": fmt.Sprint(longitude)}
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
				var allBusStop []interface{}
				// var location_info []location_info

				client := &http.Client{}

				// Parameters
				for i := 0; i < 5000+1; i += 500 {
					req, err := http.NewRequest("GET", "http://datamall2.mytransport.sg/ltaodataservice/BusStops?$skip="+fmt.Sprint(i), nil)
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

					var result BusStopResponse

					json.Unmarshal(data, &result)

					for _, busstop := range result.Value {
						outputJSON := map[string]string{"busstop_number": busstop.BusstopNo, "busstop_name": busstop.Description, "busstop_lat": fmt.Sprint(busstop.Latitude), "busstop_lng": fmt.Sprint(busstop.Longitude)}
						for _, busData := range list_data {
							if len(busData) == 4 {
								busData = "0" + busData
							}
							if busstop.BusstopNo == busData {
								allBusStop = append(allBusStop, outputJSON)
							}
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
