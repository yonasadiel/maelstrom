MAELSTROM_DIR=~/Repo/github.com/jepsen-io/maelstrom
OUT=./bin/maelstrom-node

${OUT}: *.go
	go build -o bin/maelstrom-node .

debug:
	${MAELSTROM_DIR}/maelstrom serve

# Tests are from guide https://fly.io/dist-sys

test-echo:
	MWORKLOAD=echo \
		${MAELSTROM_DIR}/maelstrom test -w echo --bin ${OUT} \
		--node-count 1 \
		--time-limit 10

test-unique-ids:
	MWORKLOAD=unique-ids \
		${MAELSTROM_DIR}/maelstrom test -w unique-ids --bin ${OUT} \
		--node-count 3 \
		--rate 1000 \
		--time-limit 30 \
		--availability total \
		--nemesis partition

test-broadcast-a:
	MWORKLOAD=broadcast \
		${MAELSTROM_DIR}/maelstrom test -w broadcast --bin ${OUT} \
		--node-count 1 \
		--rate 10 \
		--time-limit 20

test-broadcast-b:
	MWORKLOAD=broadcast \
		${MAELSTROM_DIR}/maelstrom test -w broadcast --bin ${OUT} \
		--node-count 5 \
		--rate 10 \
		--time-limit 20

test-broadcast-c:
	MWORKLOAD=broadcast \
		${MAELSTROM_DIR}/maelstrom test -w broadcast --bin ${OUT} \
		--node-count 5 \
		--rate 10 \
		--time-limit 20 \
		--nemesis partition

test-broadcast-d:
	MWORKLOAD=broadcast \
		${MAELSTROM_DIR}/maelstrom test -w broadcast --bin ${OUT} \
		--node-count 25 \
		--rate 100 \
		--time-limit 20 \
		--latency 100

test-grow-only-counter:
	MWORKLOAD=grow-only-counter \
		${MAELSTROM_DIR}/maelstrom test -w g-counter --bin ${OUT} \
		--node-count 3 \
		--rate 100 \
		--time-limit 20 \
		--nemesis partition

test-kafka-a:
	MWORKLOAD=kafka \
		${MAELSTROM_DIR}/maelstrom test -w kafka --bin ${OUT} \
		--node-count 1 \
		--rate 1000 \
		--time-limit 20 \
		--concurrency 2n

test-kafka-b:
	MWORKLOAD=kafka \
		${MAELSTROM_DIR}/maelstrom test -w kafka --bin ${OUT} \
		--node-count 1 \
		--rate 1000 \
		--time-limit 20 \
		--concurrency 2n
