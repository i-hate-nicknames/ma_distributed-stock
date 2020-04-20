# Distributed stock
This is an implementation of homework assigned at Masters Academy golang course
Distributed stock consists of two main parts: stock service and warehouse services.
Warehouse services have each their own inventory and start independently, without knowing
anything about other services. Each warehouse sends invitations over UDP for the stock
service to find and contact them later.
Stock service receives invitations and keeps a catalog of warehouses.
Additionally, stock service runs a web service that allows orders for items to be placed,
and processes these orders by retrieving necessary items from warehouses.

# Building
Run `./scripts/build.sh` to build both warehouse and stock binaries

# Running
Run `./bin/stock` to run a stock and `./bin/warehouse <port>` to run a warehouse instance.

# Using
Stock web service operates on 8001. All the interaction should happen through the web service of
the stock.
So far three methods are implemented:

1. POST /submit with body json encoded object of items, e.g. {"items": [1, 2, 3]}
2. GET /order/:id - check status of an order
3. POST /cancel/:id - cancel an order

Additionally, scripts that use httpie can be used to make queries:
1. `./scripts/submit_order.sh [1,2]` to submit an order for items 1 and 2
2. `./scripts/get_order 1` to get order with id 1