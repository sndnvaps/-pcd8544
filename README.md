# rpi_cpu_infoscreen
A small program to drive the pcd8544 cpu info screen on a raspberry pi


## Introduction
This is a simple small program to drive a product called the RPi 1.6" CPUInfo Screen
This product basically contains a Nokia 5110 screen with a pcd8544 controller.


## The hardware

The pinout of this circuit is as follows

| Function |    PIN   | GPIO     |
|----------|----------|----------|
|CLK       |        11|        17|
|DIN       |        12|        18|
|DC        |        13|        27|
|CE/CS     |        15|        22|
|Vcc       |         1|          |
|Gnd       |         6|          |
|Rst       |        16|        23|
|BackLgt   |         7|         4|


## Prerequisites
Install go compiler
  because the program write with golang


git clone this code
```
git clone https://github.com/sndnvaps/pcd8544.git
```

## How to build

enter the directory where you installed this source
Enter the src directory and type
```
go build -o rpi_cpuinfo_screen
```
After a little while you should have a program called 

rpi_cpuinfo_screen

## How to use

Now you can enter
```
_rpi_cpuinfo_screen
```
on the command line

There is also a little script called update that fetches some more or less usefull stuff to put on your shiny new display.

Enjoy!