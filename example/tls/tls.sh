#!/bin/bash

set -e

# Generate server certificate
cfssl selfsign 127.0.0.1 csr.json | cfssljson -bare server
