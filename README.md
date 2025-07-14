# JFRI

```
       ___         ___           ___                 
      /\  \       /\  \         /\  \          ___   
      \:\  \     /::\  \       /::\  \        /\  \  
  ___ /::\__\   /:/\:\  \     /:/\:\  \       \:\  \ 
 /\  /:/\/__/  /::\~\:\  \   /::\~\:\  \      /::\__\
 \:\/:/  /    /:/\:\ \:\__\ /:/\:\ \:\__\  __/:/\/__/
  \::/  /     \/__\:\ \/__/ \/_|::\/:/  / /\/:/  /   
   \/__/           \:\__\      |:|::/  /  \::/__/    
                    \/__/      |:|\/__/    \:\__\    
                               |:|  |       \/__/    
                                \|__|              

```

A simple Go program that runs scripts scattered around the file system.
Scripts can be configured in ~/.config/jfri/jfri.conf allowing the user to select one. JFRI ensures it's executable, and then runs it.

## Installation
1. Clone this repository:
   ```bash
   git clone https://github.com/PalmieriClaudio/JFRI.git
   cd jfri
   ```
2. (Optional) Build the program:
   ```bash
   go build -o jfri main.go
   ```
3. (Optional) Move the binary to a location in your `$PATH`:
   ```bash
   mv jfri /usr/local/bin/
   ```
If you prefer not to build the program, the linux executable comes already built in the repository

## Usage
   (Optional) Create the jfri.conf file in the ~/.config/jfri/ directory, add entries in the format:
   ```
   run /path/to/script1.sh
   name name of script 2
   run script2 
   ```
   If the config file is not created prior to starting JFRI, it will be created at launch, and the only available option will be to configure it.
   By default, selecting the configuration option will try and open the configuration file with nvim, if nvim is not available it will default to nano.

   The name line will define the name for the run command in the next line.
   The name tag is not necessary, though it's helpful for commands, since those will not have a sensible name by default.
   If the run line contains a script in the format "run /path/to/file/script.sh" the name will be the script name by default (script in this example).
   The name of scripts can still be overwritten with a "name" tag

   Scripts that are not executable will be made executable automatically before being run.
   Only inline commands are supported currently. This is because the script was made for my convenience and I didn't have a need yet for multi-line scripts.
   If you have a need to run multi-line scripts, this is something that could wasily be setup in the source, otherwhise a .sh containining the multi-line command should serve the purpose of running the command.

Note, this program is by no means complex, I built it for personal use, independent of similar tools possibly being available, and while it's for my convenience, I thought others could find it useful. 

Note #2. The name of the program is pronounced like Jeffrey.

:q!
