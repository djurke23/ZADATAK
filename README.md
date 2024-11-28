# Aplikacija za Upravljanje Korisnicima


frontend:
http://localhost:4200/

folder: user-management-app

pokretanje u terminalu -> 
cd rest-app
go run main.go

****  u slucaju da nece zbog SSL, onda kucati:  ****
$env:DB_USER = "postgres"  
$env:DB_PASSWORD = "admin"  
$env:DB_NAME = "user_management" 

backend:
http://localhost:8080/

folder:  rest-app

pokretanje u terminalu ->
cd user-management-app
ng serve


Tabela uradjena u pgadmin4, u bazi pod nazivom user_management.




FUNKCIONALNSOTI: 

LOGIN 
Prikaz korisnika              GET /users/ 
Dodavanje korisnika    POST /USERS/ 
Azuriranje korisnika     PUT /users/{id} 
Brisanje korisnika         DELETE /users/ {id} 



LOGIN:

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


PRIKAZ:

$response = Invoke-WebRequest -Uri http://localhost:8080/users/ -Method GET -Headers $authHeaders 
$response.Content 

 
DODAVANJE:

$body = @{ 
    firstName = "" 
    lastName = "" 
    nickname = "" 
    password = "" 
} | ConvertTo-Json -Depth 10 
$response = Invoke-WebRequest -Uri http://localhost:8080/users/ -Method POST -Headers $authHeaders -Body $body 


AZURIRANJE:

$body = @{ 
    nickname = "" 
    password = "" 
} | ConvertTo-Json -Depth 10 

$response = Invoke-WebRequest -Uri http://localhost:8080/users/4 -Method PUT -Headers $authHeaders -Body $body 


BRISANJE:

$response = Invoke-WebRequest -Uri http://localhost:8080/users/4 -Method DELETE -Headers $authHeaders 

 

