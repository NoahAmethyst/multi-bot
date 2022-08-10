# generate contract go file
# should have contract abi file
abigen --abi="../resources/contract/pancakeswap.abi.json" --type="Loot" --pkg=contract --out="../contract/Pancakeswap.go"
