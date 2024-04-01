# Go concurrency

## Original assignment

1. Cars arrive at the gas station and wait in the queue for free station
1. Total number of cars and their arrival time is configurable
1. There are 4 types of stations: gas, diesel, LPG, electric
1. Count of stations and their serve time is configurable as interval (e.g. 2–5s) and can be different for each type
1. Each station can serve only one car at a time, serving time is chosen randomly from station's interval
1. After the car is served, it goes to the cash register
1. Count of cash registers and their handle time is configurable
1. After the car is handled (random time from register handle time range) by the cash register, it leaves the station
1. Program collects statistics about the time spent in the queue, time spent at the station and time spent at the cash register for every car
1. Program prints the aggregate statistics at the end of the simulation

## ⚠ Implementation details

* All time values in `config.yaml` are in milliseconds by default; this could be changed in `random.go`, if really necessary
* Before car enters the _stand queue_ for its specific fuel type, it has to go through the _shared queue_
  * Shared queue is meant as a road between e.g. the highway and the main gas station area
  * If shared queue is full, spawning of new cars is paused
* Every stand/register/(shared) queue has its maximum length
  * If queue for the next part of the process is full, car waits and blocks the previous part
  * e.g.: Car _c_ finished refueling at stand _x_ and wants to go to the cash register _y_. _y_'s queue is full though, so _c_ waits until there is room for it. While waiting it still blocks _x_. Car that is behind _c_ in _x_'s queue cannot start refueling until _c_ leaves to the _y_'s queue.
* Queue lengths in `config.yaml` are actually Go channel capacities, meaning [one more car can fit into it](https://stackoverflow.com/a/25539742)
* Car always chooses stand/register with shortest queue (where there are currently least other cars)
* Times in `output.yaml` count time spent in queue, doing some action, but also waiting for queue (and blocking)
  * e.g. `register:handle_time` = time spent after refueling waiting __for__ _y_'s queue + time spent waiting __in__ _y_'s queue + time spent paying

## Example `config.yaml`

```yaml
cars:
  count: 1000
  arrival_time_min: 1
  arrival_time_max: 2
  shared_queue_length_max: 5
stations:
  gas:
    count: 4
    serve_time_min: 2
    serve_time_max: 5
    queue_length_max: 3
  diesel:
    count: 2
    serve_time_min: 3
    serve_time_max: 6
    queue_length_max: 3
  lpg:
    count: 1
    serve_time_min: 4
    serve_time_max: 7
    queue_length_max: 3
  electric:
    count: 1
    serve_time_min: 5
    serve_time_max: 10
    queue_length_max: 3
registers:
  count: 2
  handle_time_min: 1
  handle_time_max: 3
  queue_length_max: 1
```

### Example `output.yaml`

```yaml
sharedqueue:
    total_cars: 1000
    total_time: 105.6µs
    avg_time: 105ns
    max_time: 63.9µs
stations:
    gas:
        total_cars: 273
        total_time: 4.2492077s
        avg_time: 15.564863ms
        max_time: 16.1471ms
    diesel:
        total_cars: 253
        total_time: 3.927832s
        avg_time: 15.525027ms
        max_time: 16.1233ms
    lpg:
        total_cars: 234
        total_time: 3.6457938s
        avg_time: 15.580315ms
        max_time: 16.1354ms
    electric:
        total_cars: 240
        total_time: 3.7343759s
        avg_time: 15.559899ms
        max_time: 16.1336ms
registers:
    total_cars: 1000
    total_time: 15.5572461s
    avg_time: 15.557246ms
    max_time: 16.1471ms
```
