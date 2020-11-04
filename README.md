# Some Coding Challenge

## Geolocation Service

### Overview
You're provided with a CSV file (`data_dump.csv`) that contains raw geolocation data; the goal is to develop a service that imports such data and expose it via an API.

```
ip_address,country_code,country,city,latitude,longitude,mystery_value
200.106.141.15,SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346
160.103.7.140,CZ,Nicaragua,New Neva,-68.31023296602508,-37.62435199624531,7301823115
70.95.73.73,TL,Saudi Arabia,Gradymouth,-49.16675918861615,-86.05920084416894,2559997162
,PY,Falkland Islands (Malvinas),,75.41685191518815,-144.6943217219469,0
125.159.20.54,LI,Guyana,Port Karson,-78.2274228596799,-163.26218895343357,1337885276
```

### Requirements
1. Develop a library with two main features:
    * a service that parses the CSV file containing the raw data and persists it in a database;
    * an interface to provide access to the geolocation data (model layer);
1. Develop a REST API that uses the aforementioned library to expose the geolocation data

In doing so:
* define a data format suitable for the data contained in the CSV file;
* sanitize the entries (the file comes from an unreliable source; this means that the entries can be duplicated, may miss some value, the value can not be in the correct format or completely bogus);
* at the end of the import process, return some statistics about the time elapsed, as well as the number of entries accepted/discarded;
* the library should be configurable by an external configuration (particularly with regards to the DB configuration);
* the API layer should implement a single endpoint that, given an IP address, returns information about the IP address' location (i.e. country, city);
* the endpoint should be developed according to the HTTP/1.1 standard;

### Expected outcome and shipping:
* a library that packages the import service and the interface for accessing the geolocation data;
* the REST API application (that uses the aforementioned library) should be Dockerised and the Dockerfile should be included in the solution;
* deploy the project on a cloud platform of your choice (e.g. AWS, Heroku, etc):
    * run a container for the API layer;
    * run any other container that you think necessary;
    * have a database prepared with the already imported data

### Notes
* the file's contents are fake, you don't have to worry about data correctness
* in production the import service would run as part of a scheduled/cron job, but we don't want that part implemented as part of this exercise
* for local/development run a DB container can be included
* you can structure the repository as you see it fit


## Solution

* It took me 3 evenings; around 10h total.
* CSV parser is single threaded and not optimized; Waiting for DB response stalls the pipeline. I've shared some ideas in comments (split files on new lines). I guess this is important for you, I should have implemented it to actually be more performant... at this stage, I think the fastest way to handle that is working on partitioned files, correcting/dropping row and saving correct CSV file in a way it is directly importable by postgres.
* Having both parser and REST service in the same application is unusual; here, the logic switches to the CSV parsing depending on the command line argument ("parse").
* Geolocation coordinates -> all real world details are ignored, blindly implemented those as two float64s.
* mystery_value exact datatype is unknown; selected 64-bit integer.
* Lookup endpoint looks like this: https://fh.rozestwinski.com/api/v1/geolocate/142.33.158.193
* Records are being updated if having the same IP address, in order of CSV rows.
* Project layout is unusual; GOPATH is managed by makefile, so dependencies are manually vendored. There was a flux in golang vendoring, nowadays it should be rather done using go modules vendoring, which is reasonably left-pad-proof.
* I felt the need for 'unit testing' only the upsert postgres capabilities to ensure it's correct behavior.
* On actual validation scenario: due to upserts, one must take rows from the back of the CSV file and check them with the REST endpoint.
* I am using on-host postgresql, as I think it is the right decision for the data layer. I've been using postgres db in docker for testing; and was thinking to provide some docker-compose template.
* Application can be configured by environment variables.
* Imported data could be repaired more, e.g. by country lookup when having only country code. Cities may be compared to list of known cities. IP format could be verified. Coordinates range bounded. Unicode can be problematic for some data.
* Please note that the hosted server is on a smallest Digital Ocean VM as possible (to save energy and limit our carbon footprint ;)) - please let me know if more than functional testing is going to be performed there, I'll assign more resources to the VM in that case.
* Additional endpoints and some extra middleware functionality is provided for the HTTP server, as I like this design and decided to share it.
* CSV filename is hardcoded. Directory as environment variable makes more sense. But then I'd need to mark file as already processed, etc..
* Bad decision: started with REST, implemented CSV parsing last - parsing was slow and it stalled my development pipeline as the time was running out.

Parsing result:
Duration: 3h20m49.092272335s
InsertedOrUpdated: 983021
Failed: 16980

avg performance[rows/s]:  82.99260996058257
select count(*) from geoentry; -> 916324



