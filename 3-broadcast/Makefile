MAELSTROM_DIR=~/Repo/github.com/jepsen-io/maelstrom
MODULE=3-broadcast

${MODULE}:
	go build -o bin/${MODULE} .

test: test1

test1:
	${MAELSTROM_DIR}/maelstrom test \
		-w broadcast \
		--bin ./bin/${MODULE} \
		--node-count 1 \
		--rate 10 \
		--time-limit 20

debug:
	${MAELSTROM_DIR}/maelstrom serve