# psql-example

## Requirements
This project requires golang v1.17 and docker to run successfully.

## Running
1. To run this project simply run `make`.
    ```bash
    make 
    ```
    The makefile will create a postgres db in docker and run the app on the host machine directly.
    [`./services.sql`](./services.sql) contains the database schema and some test data.

## API

### GET /service?search=<>&sort=<>&page_size=<>&page_number=<>
This endpoint returns a list of services. The list can be filtered, paginated, and sorted.

#### Results
This api returns a JSON object of the following shape for `200`:
```json
{
    "services": [
        {
            "Name": "modern",
            "Description": "a fine description",
            "Versions": ["v1", "v2"],
        },
        {
            "Name": "legacy",
            "Description": "a description",
            "Versions": ["v100"],
        }
    ]
}
```

If this api has an error, the result looks like the following:
```json 
{
    "error": "something went wrong"
}
```

#### Params

##### search
This is the string to find within the Name or Description field of a Service.

##### sort
* "asc" to to return results sorted by name ascending.
* "desc" to to return results sorted by name descending.

##### page_size
This is the maximum number of results to return. A value of 0 represents all.

##### page_number
This is the page of results to return. This is ignored when `page_size` = 0

### GET /service/\<name\>
This endpoint returns a specific services based on the passed name.

#### Results
This api returns a JSON object of the following shape for `200`:
```json
{
    "service": {
        "Name": "modern",
        "Description": "a fine description",
        "Versions": ["v1", "v2"],
    }
}
```

If this api has an error, the result looks like the following:
```json 
{
    "error": "something went wrong"
}
```

## Design 

### Sorting & Pagination in application vs using database features.
Sorting & Pagination are implemented using database features so the application does not use as much memory.
The tradeoff here is the DB will be put under more load to perform these operations but since dbs typically have good caching properties, it should use less resources globally.

### Filtering in application vs in database.
Filtering is done in the application. 
This allows flexibility and clarity in what filters can be created. 
This also allows for filters that databases might not support.

### Model interface
An interface was created for operations regarding the data. 
This allows the backend to be abstracted. 
We can change from postgres to mongodb without having to rewrite routes. 

### Database
An external database is used to allow this application to scale horizontally.
A SQL database was chosen because the data can be relational (services relate to version).

## Testing
This application has unit and integration tests that can be run locally.
```bash
make test
```