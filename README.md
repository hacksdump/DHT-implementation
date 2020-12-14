# DHT Implementation

A Simple demonstration to understand the algorithms and operation of Distributed Hash Tables.

## Build instructions

```shell
go build DHT-implementation
./DHT-implementation
```

Run as many nodes as you want. The port number will be self-assigned. Just keep `8000` free from other applications.

## API routes

1. "/"  
    a. Use `GET` method to retrieve values. Pass this as JSON body.
   ```json
   {
       "key": "someQueryKey"
   }
   ```
   b. Use `PUT` method to retrieve values. Pass this as JSON body.
   ```json
   {
       "key":   "someKey",
       "value": "someValue"
   }
   ```
   c. Use `DELETE` method to delete data. Pass this as JSON body.   
   ```json
   {
       "key": "someKeyToDelete"
   }
   ```

2. "/data"  
    `GET` method to retrieve all data. No body required.
   
3. "/node"  
    `GET` method to get all node metadata. No body required.  
   Example: 
   ```
   GET localhost:8002/node
   ```
   
Rest of the routes are internal to the function of the application and are being used for inter-node communication.