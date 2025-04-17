#!/bin/bash
(cd ./inventoryapi/ && go run cmd/api/main.go) &
(cd ./ts-inventory-ui/ && npm run dev)
