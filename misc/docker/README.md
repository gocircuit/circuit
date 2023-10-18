# IN PROGRESS

## To preparing a Docker image with a circuit inside

    sudo docker build .

## Running containers

**Note**: Use the `-l` option so that `bash` is in "login shell mode". This will make bash execute `/etc/profile`, where environment variables are exported for convenience.

	sudo docker run -i -t circuit-img /bin/bash -l

If you don't add an entry command to `docker run` the container will start the circuit automatically. (you won't be able to interact with the container through a `tty`) By default, the script for automatic start of the circuit is using the address `228.8.8.8:8788` for UDP multicast channel. If you want to change this address, you can mount a file containing the chosen address to `/go/util/addr`, like this:

    echo "229.9.9.9:9899" > ~/addr
    sudo docker run -v ~/addr:/go/util/addr:ro circuit-img

## Example usage

Start a circuit host inside one Docker image:

    sudo docker run circuit-img &

Start second container, but this time enter `bash` instead:

    sudo docker run -i -t circuit-img /bin/bash -l

From inside the second container:

    start-circuit.sh &
    circuit ls -discover $ADDRESS /

## Utility
From inside the container you can use the following commands:

`start-circuit.sh`: will start the circuit (on the default address if `/go/util/addr` doesn't exist or on the address inside it otherwise).

`$ADDRESS`: This variable can be used to conveniently get the __default__ address.
