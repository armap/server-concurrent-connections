### **FEEDER SERVER**

- A feeder-server ("Application") that opens a socket and restricts input to at most 5 concurrent clients (providers). 
- Clients will connect to the Application and write any string (product sku) of 9 characters in a specific format, and then close the connection. 
- The Application must write a de-duplicated list of these numbers to a log file in no particular order.

**Folder Structure**
- cmd: The cmd directory is a Go convention for when an application produce multiple binaries. It contains a subdirectory for each main package and its binary. In our case, one for the feeder-server and another for the client-test.
- server: Contains the server implementation and its test. 

**How to Build and Run Feeder-Server and Client-Test:**
1. Open a Terminal on directory feeder-server/cmd/feederserver
2. Execute:  go build
3. Execute:  ./feederserver
4. Open another Terminal on directory feeder-server/cmd/clienttest
5. Execute:  go build
6. Execute:  ./clienttest
7. Press INTRO from Client-Test to terminate the Server

**How to test:**
1. Open a Terminal on directory feeder-server/
2. Execute:  go test ./...



