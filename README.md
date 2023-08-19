# dollar-exchange-rate

### What
Simple Go server-client project that get the dollar exchange rate in BRL and save it in a database.

### Server
Exposes the endpoint:
```
curl http://localhost:8080/cotacao
```

Response: 
```
{
    "exchangeRate": 4.9695
}
```

Run it:
```
go run cmd/server/main.go
```

It also saves the exchange rate in a [sqlite](https://www.sqlite.org/index.html) database (dolar_exchange_rate.db).

To query it:
```
sqlite3 dolar_exchange_rate.db

.tables

select * from dolar_exchange_rate;
```

### Client
It makes a request to the server above and saves the response in a text file (cotacao.txt).

Run it:
```
go run cmd/client/main.go
```
