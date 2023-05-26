MAELSTROM_DIR=~/Repo/github.com/jepsen-io/maelstrom
MODULE=2-unique-id

${MODULE}:
	go build -o bin/${MODULE} .

test:
	${MAELSTROM_DIR}/maelstrom test \
		-w unique-ids \
		--bin ./bin/${MODULE} \
		--node-count=3 \
		--rate=1000 \
		--time-limit=30 \
		--availability total \
		--nemesis partition

debug:
	${MAELSTROM_DIR}/maelstrom serve