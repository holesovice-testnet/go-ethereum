## Install Holešovice as a Devnet

This is a simple explainer for how to launch an instance of Holešovice, as a private devnet, on a computer running Linux.

It has been tested and confirmed as working on Ubuntu on `amd64` (Lenovo Thinkpad T490s), and Ubuntu on `arm64` (Raspberry Pi 4B), but will likely also work on many other environments.

### Install golang

Install golang by following [these golang installation instructions](https://go.dev/doc/install).

### Install pre-requisites

```
sudo apt install make git
```

### Install Holešovice Project

```
mkdir ~/holesovice-testnet
git clone https://github.com/holesovice-testnet/go-ethereum.git
cd go-ethereum
make geth
nano genesis.json
```

(optional) edit the `genesis.json` file, to add allocations of ETH to specific addresses.

See the appendices below, for more details on how to edit `genesis.json`.

Then save and close and continue.

### Initialise Local Devnet

```
./build/bin/geth --verbosity 7 init genesis.json --datadir=/home/ubuntu/.verkle
```

### Run Local Devnet

```
./build/bin/geth --mine --miner.etherbase=0xBfd826F430C41Db28e158642d8f3Db5E52CDd8a5 --miner.threads=1 --datadir=/home/ubuntu/.verkle --http --http.addr 0.0.0.0 --http.api=net,eth --http.corsdomain="*" --http.vhosts="*"
```

Note: you probably want to replace the `--miner.etherbase` address to one that you control.

:boom::boom::boom: Congratulations, you are now running a single-node _Execution Layer_ "network".

Things should look something like this now:

![image](https://user-images.githubusercontent.com/2212651/174289245-de828989-b5b6-46e9-a5a3-e065d230133c.png)

### Connecting to the Devnet

The node presents a standard RPC endpoint, which you can (for example), connect MetaMask to with the parameters below:

If running the node on _the same computer_ as your MetaMask:

```
Network Name: Holešovice
RPC URL: http://127.0.0.1:8545
Chain ID: 17000
```

Like so:

![image](https://user-images.githubusercontent.com/2212651/174290109-7fc52dc7-21f5-4e26-ba60-04adb7be1740.png)

If running the node on _a different computer_ than your MetaMask:

```
Network Name: Holešovice
RPC URL: http://{ip-address-of-node}:8545
Chain ID: 17000
```

Like so:

![image](https://user-images.githubusercontent.com/2212651/174290258-c57790fb-165e-433e-9ae6-67e5eac65233.png)

## Appendices

### Editing `genesis.json`

In order to launch with a non-zero amount of ETH associated with a given address, you can edit the following section in `genesis.json`:

```
  "coinbase": "0x0000000000000000000000000000000000000000",
  "alloc": {
    "0x0000000000000000000000000000000000000000": {
      "balance": "0x40000000000000000000"
    }
  },
  "number": "0x0",
```

An example of what to change it to, using a sample address, is found below:
```
  "coinbase": "0x0000000000000000000000000000000000000000",
  "alloc": {
    "0x0000000000000000000000000000000000000000": {
      "balance": "0x40000000000000000000"
    },
    "0xBfd826F430C41Db28e158642d8f3Db5E52CDd8a5": {
      "balance": "0x40000000000000000000"
    }
  },
  "number": "0x0",
```
To update yourself, replace the `0xBfd826F430C41Db28e158642d8f3Db5E52CDd8a5` with your own address.
