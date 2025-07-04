#!/bin/bash
(cd ./inventoryapi/ && go run main.go) &
(cd ./ts-inventory-ui/ && npm run dev)
