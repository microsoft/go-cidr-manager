# IPv4 CIDR Manager
The package `IPv4CIDR` contains utilities to perform the following operations:

1. Parse a string representing the CIDR block
    - Take a single IP address as input
    - Take a CIDR block in a standard notation where the `IP` part of the `IP/CIDR` range is the first IP address in the CIDR block
    - Take a non-standard CIDR block and enable a `standardize` flag to convert it to the standard notation
2. Split the CIDR block into two halves
3. Get the following information from the CIDR block
    - Convert to string
    - Get the IP part of the block representation
    - Get the CIDR mask part of the block representation
    - Get the nth IP address in range
    - Get the netmask
    - Get the size of the CIDR block

## To Use
Import the package into your code using:

    import "github.com/microsoft/go-cidr-manager/ipv4cidr"
