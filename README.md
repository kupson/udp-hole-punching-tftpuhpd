# TFTP Udp Hole Punching Helper Server

See [this project](https://github.com/kupson/udp-hole-punching) for explanation.

## Protocol

This server listens on the standard TFTP 69/udp port and waits for GET requests with filenames matching the regex `v1_tftp_udp_\d+`.

Details:

- `v1` - protocol version
- `_tftp_udp_` - literal string
- `\d+` - local socket port number or `0`

Request and response lengths are approximately equal in bytes, limiting the server’s usefulness in UDP amplification attacks.

## Usage

It's recommended to run this daemon from systemd with limited priviliges using socket activation.

- install proper binary from `build/*/tftpudpd` into `/usr/local/bin/tftpudpd`
- copy `tftpudpd.service` and `tftpudpd.socket` from `examples/systemd/` directory into `/etc/systemd/system`
- enable & start both `tftpudpd.service` and `tftpudpd.socket` with `systemctl` command
