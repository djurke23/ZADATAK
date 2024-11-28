# Aplikacija za Upravljanje Korisnicima


## Frontend:

 - URL: http://localhost:4200/
 - Folder: user-management-app
 - Pokretanje frontend-a u terminalu:
#### cd user-management-app  
#### ng serve  

## Backend:

- URL: http://localhost:8080/
- Folder: rest-app
- Pokretanje backend-a u terminalu:
#### cd rest-app  
#### go run main.go  


### ***Napomena: Ako se pojave problemi sa SSL-om ukucati u terminalu: ***

#### $env:DB_USER = "postgres"  
#### $env:DB_PASSWORD = "admin"  
#### $env:DB_NAME = "user_management"

## Baza Podataka:
Baza pod nazivom user_management kreirana je u pgAdmin4.

## Funkcionalnosti

#### 1. Prijava (Login):
Autentifikacija korisnika i generisanje JWT tokena.
#### 2. Upravljanje korisnicima:

a) Prikaz korisnika: GET /users/

b)Dodavanje korisnika: POST /users/

c)Ažuriranje korisnika: PUT /users/{id}

d) Brisanje korisnika: DELETE /users/{id}


## API

### LOGIN

$headers = @{ 
    "Content-Type" = "application/json" 
} 

$body = @{ 
    nickname = "janed" 
    password = "password123" 
} | ConvertTo-Json -Depth 10 

$response = Invoke-WebRequest -Uri http://localhost:8080/login -Method POST -Headers $headers -Body $body

$jwtToken = $response.Content | ConvertFrom-Json | Select-Object -ExpandProperty token

$authHeaders = @{ "Authorization" = "Bearer $jwtToken" } 


###  PRIKAZ 

$response = Invoke-WebRequest -Uri http://localhost:8080/users/ -Method GET -Headers $authHeaders 

$response.Content  

### DODAVANJE 

$body = @{ 

    firstName = "" 
    lastName = "" 
    nickname = "" 
    password = "" 
} | ConvertTo-Json -Depth 10 

$response = Invoke-WebRequest -Uri http://localhost:8080/users/ -Method POST -Headers $authHeaders -Body $body  


### AZURIRANJE

$body = @{ 

    nickname = "" 
    password = "" 
} | ConvertTo-Json -Depth 10 

$response = Invoke-WebRequest -Uri http://localhost:8080/users/4 -Method PUT -Headers $authHeaders -Body $body  



### BRISANJE 

$response = Invoke-WebRequest -Uri http://localhost:8080/users/4 -Method DELETE -Headers $authHeaders  


 

