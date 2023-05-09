# OpenMSR

OpenMSR is a cross-platform GUI application for the MSR605 and MSR605X.
The pkg directory also has two Go libraries,
[libmsr](https://pkg.go.dev/github.com/egginabucket/openmsr/pkg/libmsr)
and [libtracks](https://pkg.go.dev/github.com/egginabucket/openmsr/pkg/libtracks)
which are great if you want to develop your own applications / commands. 

## Features

- Read, write (WIP), and erase on individual tracks
- BPI, BPC, and parity selection for each track
- IEC 7813 and AAMVA parsing with Luhn validation
- Hi/lo coercitivity selection

## Installation

Clone and cd into the repository:

`git clone https://github.com/egginabucket/openmsr.git && cd openmsr`

Build for your OS (requires [Go](https://go.dev/doc/install))

`make build`

An executable should appear in a `bin` folder.

To run without root privileges, Linux users need to add a udev rule:

`# cp 50-msr605x.rules /etc/udev/rules.d/`
`# udevadm control --reload-rules && udevadm trigger`

To use the MSR605, make sure your user has access to the serial ports
(`dialout` group for Debian-based, `uucp` for Arch).

## Limitations

- Writing currently has some issues. I'll fix it soon!
- The GUI looks different depending on your system theme; see [andlabs/ui](https://github.com/andlabs/ui)
- Refreshing the list of devices just adds onto it.
- Saving / opening files currently doesn't work, as I want to make it compatible with Deftun's MSRX software.
