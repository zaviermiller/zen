# Zen
Give yourself peace of mind.

![Usage example](/public/zen_example.gif)

## What is Zen?
Zen is a tool I built to help give me peace of mind when it comes to my CS labs. This semester, for my first lab at least, I was given an example executable of the working lab. After finishing the lab, I figured I could use these provided examples to create a test program that compares my lab binary to the correct one. If I know my program gives a pretty similar (if not the same) output, I can rest easy knowing I'll get an A on my lab :]

### Linux Installation (UTK Lab Machines)
The installation process is pretty easy. First, download the binary from the [releases page](https://github.com/zaviermiller/zen/releases), and put it into a new folder in your home directory called `bin` (`mkdir ~/bin`). Next, run `vim ~/.zshrc` to open your zsh config and add the following line:
`alias zen=$HOME/bin/zen`

And that's it! Now just use `zen` in the terminal.

### How it works
Basically, after installing just run `zen [correct lab binary path] [your lab binary path] [options that need to be passed to each program]` (IN THAT ORDER) and then go through the lab program like you would if you were using it. Then, when you're finished, exit the program and watch the differences roll in.