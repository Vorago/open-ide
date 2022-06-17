# About

This project solves a basic problem of navigating in a bunch of projects and opening them in IDE of your preference, which happens IntelliJ IDEA in my case.
I use [i3wm](https://i3wm.org/), but same functionality can be adapted for other desktops as well.

# Demo

![demo](demonstration.gif)

# Requirements

* rofi
  * for project incremental searching
* i3
  * for focusing already opened project
* xdotool
  * for searching through opened projects

# Usage

```sh
go build
./open-ide --depth 3 --codeDir /home/user/code/ --ideCommand /opt/idea/idea
```
