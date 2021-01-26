## Zen
Give yourself peace of mind.

![Usage example](/public/zen.gif)

### What is Zen?
Zen is a tool I built to help give me peace of mind when it comes to my CS labs. This semester, for my first lab at least, I was given an example executable of the working lab. After finishing the lab, I figured I could use these provided examples to create a test program that compares my lab binary to the correct one. If I know my program gives a pretty similar (if not the same) output, I can rest easy knowing I'll get an A on my lab :]

## Installation
The installation process is pretty easy. First, either download and compile the source, or if you are on Linux, download the binary from the dist folder.
Once you have the binary, put it somewhere like your home directory. Finally, add an `alias` to your `.zshrc/.bashrc` with the following line (or something similar): `alias zen=$HOME/zen`

### How it works
Basically, after installing just run `zen [correct lab binary path] [your lab binary path] [options that need to be passed to each program]` (IN THAT ORDER) and then go through the lab program like you would if you were using it. Then, when you're finished, exit the program and watch the differences roll in.