MAELSTROM_DIR=~/Repo/github.com/jepsen-io/maelstrom
OUT=./bin/maelstrom-node

${OUT}: *.go
	go build -o bin/maelstrom-node .

debug:
	${MAELSTROM_DIR}/maelstrom serve

# Tests are from guide https://fly.io/dist-sys

test-echo:
	${MAELSTROM_DIR}/maelstrom test -w echo --bin ${OUT} \
		--node-count 1 \
		--time-limit 10

test-unique-ids:
	${MAELSTROM_DIR}/maelstrom test -w unique-ids --bin ${OUT} \
		--node-count 3 \
		--rate 1000 \
		--time-limit 30 \
		--availability total \
		--nemesis partition

test-broadcast-a:
	${MAELSTROM_DIR}/maelstrom test -w broadcast --bin ${OUT} \
		--node-count 1 \
		--rate 10 \
		--time-limit 20

test-broadcast-b:
	${MAELSTROM_DIR}/maelstrom test -w broadcast --bin ${OUT} \
		--node-count 5 \
		--rate 10 \
		--time-limit 20

test-broadcast-c:
	${MAELSTROM_DIR}/maelstrom test -w broadcast --bin ${OUT} \
		--node-count 5 \
		--rate 10 \
		--time-limit 20 \
		--nemesis partition
