# URL Shortener service task

## Task requirements

* A RESTful web service that allows users to POST a URL and receive a short URL in return, or GET a short URL and be redirected to the original long URL.
* A web app for users to access the service's features.

## Task details

* the service should be easy to deploy in the cloud (AWS, Heroku, or similar)
* use a language of your choice to build the service; feel free to use additional libraries as needed
* use a database of your choice to store the service's state
* take as much time as you need to complete the assignment, but somewherebetween 3 and 6 hours of time sounds about right
* if you have any questions about the assignment, please contact me.

## Some assumptions

* URL will be stored forever (no expiration of any short url)
* Each short URL will be based on `<BASE_URL>/<URL_KEY>`
  * `base_url` will be the hostname of the service
  * `url_key` will be alphanumeric string of not more than 7 characters

## Some design options

1. We will use BASE62 encoding to generate the `URL_KEY`, so this will allow of total 7^62 short URLs stored in the system ( BASE62 encoding is using a-z, A-Z, 0-9)
2. We will use one table for storing generated short URLs
   1. Each short URL will have key as primary column
   2. an owner id for tracking who can update/delete the record
   3. original URL
   4. timestamp of creating the record
3. TODO KEY generation
4. We will use one table for username/password authentication, each authenticated user will be able to create new short URLs, and delete/update short URL created by the same user
5. For data store, we can multiple of options

## Out of scope, but also good to consider

* We can assume that for each created short URL there will be multiple ( more ten 100) request to read this short URL, so a caching system could be a huge benefit for the performance of the system.
* History of updated/created/deleted records.
* Metrics/tracing data for monitoring and troubleshooting the system. [Opentelemetry](https://opentelemetry.io/) could be used for both.
