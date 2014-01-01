# anolis/irc
---
anolis/irc is a library used to connecting to and interacting with the IRC protocol

### TODO
- better documentation
- logging
- better error handling
- rate limit connection (to stop excess flood errors)
- split long lines (rfc2812 states max length: 512 (510+\r\n))
- more commands
  - action
  - mode
  - whois
  - whowas
- mIRC color support

### reference:
* [RFC-2812](https://tools.ietf.org/html/rfc2812) IRCv2