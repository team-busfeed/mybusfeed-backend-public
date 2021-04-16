# BusFeed Backend Services

## Overview

This is the Busfeed react native mobile platform backend servies. It's a simple service that lies between the frontend mobile application and LTA datamall (Bus information data source). This backend service is code with Go and deployed on heroku (https://mybusfeed.herokuapp.com). 

## Dependencies

```
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
```

## Bus Inforamtion Restful APIs

### Healthcheck
This is to check the status of the endpoint to ensrue that it is up and running.

**GET /demand/healthcheck**

Sample Response:

```json
{
    "status": "healthy",
    "type": "success"
}
```

### List of Bus Services 
This is to check the list of bus services in a bus stop

**GET /demand/bus-stop/{busStopNo}**

Sample Request:

```
https://mybusfeed.herokuapp.com/demand/bus-stop/63031
```

Sample Response:

```json
{
    "services": [
        "80"
    ]
}
```

### Retrieve all Bus Stop information based on lat/long provided
This endpoint will return all the bus stop number around 200-300m radius (Configurable) from the current lat/long provided. The bus stop will be sorted based on ascending order based on distance from the current point.

**GET location/getListOfBusStopNo/{latitude}-{longitude}**

Sample Request:

```
https://mybusfeed.herokuapp.com/location/getListOfBusStopNo/1.2345312-3.4654623
```

Sample Response:

```json
[
    {
        "busstop_lat": "1.29684825487648",
        "busstop_lng": "103.852535916540067",
        "busstop_name": "Hotel Grand Pacific",
        "busstop_number": "1012",
        "compare": 1
    },
    {
        "busstop_lat": "1.29770970610083",
        "busstop_lng": "103.8532247463225",
        "busstop_name": "St. Joseph's Ch",
        "busstop_number": "1013",
        "compare": 2.5
    },
    {
        "busstop_lat": "1.29698951191332",
        "busstop_lng": "103.85302201172507",
        "busstop_name": "Bras Basah Cplx",
        "busstop_number": "1019",
        "compare": 5.5
    },
    {
        "busstop_lat": "1.2966729849642",
        "busstop_lng": "103.85441422464267",
        "busstop_name": "Opp Natl Lib",
        "busstop_number": "1029",
        "compare": 5.6
    },
]
```

### Retrieve Bus Stop Information

This endpoint will retrieve all the bus stop information where the input matches the bus stop number or description. - Works like tha search function

**GET location/getBusStopInformation/{busstopnum}**

Sample Request:

```
https://mybusfeed.herokuapp.com/location/getBusStopInformation/4121
```

Sample Response:

```json
[
    {
        "busstop_lat": "1.29615544092335",
        "busstop_lng": "103.84965029974611",
        "busstop_name": "SMU",
        "busstop_number": "4121"
    },
    {
        "busstop_lat": "1.26591350024271",
        "busstop_lng": "103.81925660621779",
        "busstop_name": "Aft HarbourFront Stn",
        "busstop_number": "14121"
    },
    {
        "busstop_lat": "1.37769846227758",
        "busstop_lng": "103.75093552549583",
        "busstop_name": "Blk 113",
        "busstop_number": "44121"
    },
]
```

### Retrieve bus stop information based on bus stop number input

This endpoint will retrieve all the bus stop information (Lat/Long) based on the list of bus stop number.

**POST /location/returnBusStopInformation**

Sample Request:

```
https://mybusfeed.herokuapp.com/location/returnBusStopInformation
```

```json
{
    "Busstops": [4121, 62121]
}
```

Sample Response:

```json
[
    {
        "busstop_lat": "1.29615544092335",
        "busstop_lng": "103.84965029974611",
        "busstop_name": "SMU",
        "busstop_number": "4121"
    },
    {
        "busstop_lat": "1.34682778451192",
        "busstop_lng": "103.87156809881778",
        "busstop_name": "Aft Wolskel Rd",
        "busstop_number": "62121"
    }
]
```

## Deployment

To run the backend services locally. 

**run go main.go**

To deploy the backend services contact TeamBusFeed - hello@mybusfeed.com


