MAELSTROM_DIR=~/Repo/github.com/jepsen-io/maelstrom
MODULE=1-echo

${MODULE}:
	go build -o bin/${MODULE} .

test:
	${MAELSTROM_DIR}/maelstrom test -w echo --bin ./bin/${MODULE} --node-count=1 --time-limit=10

debug:
	${MAELSTROM_DIR}/maelstrom serve