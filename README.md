# avitojob-trainee-task
My solution for the Avito trainee task https://github.com/avito-tech/job-backend-trainee-assignment/

# Quickstart 

With `docker-compose`
```shell
$ git clone https://github.com/tkrouty/avitojob-trainee-task.git
$ cd avitojob-trainee-task
$ docker-compose up -d
```

# Usage

Using HTTPie (https://httpie.org/) to send requests

Request: Edit balance: 

```shell
$ http POST localhost:8000/edit_balance/1 sum:=100
```
Response 
```
{
    "message": "transaction completed"
}

```

Request: Show balance

```shell
$ http GET localhost:8000/show_balance/1
```

Response 
```
{
    "balance": 100,
    "currency": "RUB",
    "user_id": "1"
}
 ```
Request: Show balance in custom currency
```shell
$ http GET localhost:8000/show_balance/1 currency==USD
```
Response 
```
{
    "balance": 2.54206276,
    "currency": "USD",
    "user_id": "1"
}
 ```

Request: Transfer money from user 1 to user 2
```shell
$ http POST localhost:8000/transfer source_id=1, target_id=2, sum:=20
```

Response 
```
{
    "message": "transaction is successful"
}
```

Request: Show transaction history for user 1
```shell
$ http GET localhost:8000/show_history/1
```

Response 
```
{
    "transaction_history": [
        {
            "source_id": "",
            "sum": 100,
            "target_id": "1",
            "transaction_id": 4,
            "transaction_time": "2020-09-28T12:19:18Z"
        },
        {
            "source_id": "1",
            "sum": 20,
            "target_id": "2",
            "transaction_id": 6,
            "transaction_time": "2020-09-28T12:19:59Z"
        }
    ],
    "user_id": "1"
}
```
Transaction history is sorted by sum and by date.

# Running tests

```shell
$ docker exec -it gin_app ./run_tests.sh
```
