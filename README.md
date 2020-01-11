# prologix

Go library to communicate with either a [Prologix][prologix-web] GPIB-ETHERNET
(GPIB-LAN) controller or a GPIB-USB (HPIB-USB) controller.

[![GoDoc][godoc badge]][godoc link]
[![Go Report Card][report badge]][report card]
[![License Badge][license badge]][LICENSE.txt]

## Overview

[Prologix][prologix-web] offers two GPIB controllers that enable a computer to
communicate over GPIB using either an Ethernet or USB interface. Both the
[GPIB-ETHERNET controller][gpib-ethernet] and the [GPIB-USB
controller][gpib-usb] can operate either in controller mode or device mode. In
controller mode, the GPIB-ETHERNET or GPIB-USB is the Controller-In-Charge
(CIC) on the GPIB bus. In device mode, the GPIB-ETHERNET or GPIB-USB acts a
GPIB device that is either a talker or a listener.

For more information, please see the User Manual and FAQ for either the
[GPIB-ETHERNET controller][gpib-ethernet] or the [GPIB-USB
controller][gpib-usb].

## Status

- GPIB-USB VCP: Partially implemented
- GPIB-USB Direct Driver: Not implemented
- GPIB-ETHERNET: Not implemented


## GPIB-USB

The GPIB-USB controller communicates with a computer either directly using the
D2XX driver or as a Virtual COM Port (VCP) using the FTDI FT245R driver.

### GPIB-USB VCP Driver Installation

The appropriate VCP driver for your operating system can be downloaded from the
[FTDI VCP Drivers webpage][ftdi-vcp-drivers]. Alternatively, on macOS you can use
[Homebrew][] to install the VCP driver as follows:

```bash
$ brew cask install ftdi-vcp-driver
```

### GPIB-USB D2XX Direct Driver Installation

The appropriate D2XX Direct Driver for your operating system can be downloaded
from the [FTDI D2XX Direct Drivers webpage][ftdi-d2xx-drivers]. Alternatively,
on macOS you can use [Homebrew][] to install the D2XX direct driver as follows:

```bash
$ brew install libftdi
```


## Contributing

To contribute, please fork the repository, create a feature branch, and then
submit a [pull request][].

## License

[prologix][prologix] is released under the MIT license. Please see the
[LICENSE.txt][] file for more information.


[ftdi-d2xx-drivers]: https://www.ftdichip.com/Drivers/D2XX.htm
[ftdi-vcp-drivers]: https://www.ftdichip.com/Drivers/VCP.htm
[godoc badge]: https://godoc.org/github.com/gotmc/prologix?status.svg
[godoc link]: https://godoc.org/github.com/gotmc/prologix
[gpib-ethernet]: http://prologix.biz/gpib-ethernet-controller.html
[gpib-usb]: http://prologix.biz/gpib-usb-controller.html
[homebrew]: https://brew.sh/
[LICENSE.txt]: https://github.com/gotmc/prologix/blob/master/LICENSE.txt
[license badge]: https://img.shields.io/badge/license-MIT-blue.svg
[prologix]: https://github.com/gotmc/prologix
[prologix-web]: http://prologix.biz/
[pull request]: https://help.github.com/articles/using-pull-requests
[report badge]: https://goreportcard.com/badge/github.com/gotmc/prologix
[report card]: https://goreportcard.com/report/github.com/gotmc/prologix
