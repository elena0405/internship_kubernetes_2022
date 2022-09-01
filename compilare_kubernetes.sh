#!/bin/bash

systemctl stop kubelet
./build/run.sh make DBG=1 kubelet
sudo cp _output/dockerized/bin/linux/amd64/kubelet /usr/bin/
systemctl restart kubelet
systemctl status kubelet

