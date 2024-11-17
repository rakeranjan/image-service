# backend-service

The service is using aws stack , for testing/development we are using localstack

Command to start testing
`docker-compose up`

# user-service
on localhost:8002
command to create a user

```curl --location 'localhost:8002/v1/user' \
--header 'Authorization;' \
--header 'User-Agent: Apidog/1.0.0 (https://apidog.com)' \
--header 'Content-Type: application/json' \
--data-raw '{
    "firstName": "a",
    "lastName": "b",
    "userName": "cat.ran",
    "password": "cat.ran@",
    "phoneNumber": "1234567890"
}'
```

# image-service

on localhost:8003

POST ismage
```
curl --location 'localhost:8003/v1/image' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzE0MzE0MDMsImZpcnN0TmFtZSI6ImEiLCJpZCI6ImI3MTg3NjVlLWI5NjEtNDYxOS05NWJmLTg5MDEwMjE2NzdjZSIsImxhc3ROYW1lIjoiYiIsInBob25lTnVtYmVyIjoiMTIzNDU2Nzg5MCIsInVzZXJOYW1lIjoiY2F0LnJhbiJ9.3d2RxaoexMZzdtnEZ_gdBd6IHDLloLmjgmvPsqQaGuA' \
--header 'User-Agent: Apidog/1.0.0 (https://apidog.com)' \
--form 'file=@"/Users/rakeshranjan/Documents/BetterWorld.jpg"'
```

Once the image is uploaded to processing bucket , it pushes message to pricessing queue

GetImage by imageId
```
curl --location 'localhost:8003/v1/image/acf11d9e-cab1-4f45-afc9-0f4cad165985' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzE0MzE0MDMsImZpcnN0TmFtZSI6ImEiLCJpZCI6ImI3MTg3NjVlLWI5NjEtNDYxOS05NWJmLTg5MDEwMjE2NzdjZSIsImxhc3ROYW1lIjoiYiIsInBob25lTnVtYmVyIjoiMTIzNDU2Nzg5MCIsInVzZXJOYW1lIjoiY2F0LnJhbiJ9.3d2RxaoexMZzdtnEZ_gdBd6IHDLloLmjgmvPsqQaGuA'
```

Get All images of a user
```
curl --location 'localhost:8003/v1/images' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzE0MzE0MDMsImZpcnN0TmFtZSI6ImEiLCJpZCI6ImI3MTg3NjVlLWI5NjEtNDYxOS05NWJmLTg5MDEwMjE2NzdjZSIsImxhc3ROYW1lIjoiYiIsInBob25lTnVtYmVyIjoiMTIzNDU2Nzg5MCIsInVzZXJOYW1lIjoiY2F0LnJhbiJ9.3d2RxaoexMZzdtnEZ_gdBd6IHDLloLmjgmvPsqQaGuA'
```

Update image 
```
curl --location --request PUT 'localhost:8001/v1/image/e2f37355-646a-4cec-8f81-c0c06d099ac3' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzE4OTk3NTgsImZpcnN0TmFtZSI6ImEiLCJpZCI6IjkyOTEyZTM4LWRhZTgtNDk3MC04MDQ3LTI3YjJkMTViMjY4YiIsImxhc3ROYW1lIjoiYiIsInBob25lTnVtYmVyIjoiMTIzNDU2Nzg5MCIsInVzZXJOYW1lIjoiY2F0LnJhbiJ9.WjorCqDJxX4xUfPyiJnLgpSALr0RSc6svNaQTYzuvts' \
--form 'file=@"/Users/rakeshranjan/Documents/aidash.jpg"'
```

DeleteImage
```
curl --location --request DELETE 'localhost:8001/v1/image/e2f37355-646a-4cec-8f81-c0c06d099ac3' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzE4OTk3NTgsImZpcnN0TmFtZSI6ImEiLCJpZCI6IjkyOTEyZTM4LWRhZTgtNDk3MC04MDQ3LTI3YjJkMTViMjY4YiIsImxhc3ROYW1lIjoiYiIsInBob25lTnVtYmVyIjoiMTIzNDU2Nzg5MCIsInVzZXJOYW1lIjoiY2F0LnJhbiJ9.WjorCqDJxX4xUfPyiJnLgpSALr0RSc6svNaQTYzuvts'
```



# image-analysis
It reads the message from pricessing queue, 
Step 1: Process the image
Step 2: Persist the analysis data 
Step 3: Move the image from processing bucket to procressed bucket
Step 4: Push message to processed queue

# notification-service
fetched the message from processed queue, and send notification to user that their image as been processed.