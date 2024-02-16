# Terminal Timer

A simple terminal-based timer application designed for productivity, written in Go. 



<p align="center">
  <img src="example.gif" alt="Terminal Timer Usage">
</p>



## Getting Started

To get a local copy up and running, you have two options: installing the precompiled binary directly or building the application from source. Follow these simple steps based on your preference.

### Prerequisites

- Go 1.15+ (if building from source)
- Git

### Option 1: Installing the Precompiled Binary

For users who prefer to install the precompiled binary directly:

1. Download the latest release from the [Releases page](https://github.com/cameroncuttingedge/terminal_timer/releases).
2. Extract the binary from the downloaded archive.
3. Move the binary to a directory in your `PATH` to make it executable from anywhere. This will depend on your OS

   For example, on Unix-like systems, (May need sudo):

   ```sh
   mv terminal-timer /usr/local/bin/
   ```
4. Verify the installation by running:

    ```sh
    terminal-timer -c    
    ```


### Option 2: Installing from Source Using `go install`

For users who prefer to install the application directly from the source using Go:

1. Ensure you have Go installed on your system. You can check by running `go version` in your terminal.

    ```sh
    go version
    ```

2. Use `go install` to download and install the binary directly from the repository:

    ```sh
    go install github.com/cameroncuttingedge/terminal-timer@main
    ```

    This command will install the latest version of `terminal-timer` to your `$GOPATH/bin` directory or `$GOBIN` if set.

3. Verify the installation by running:

    ```sh
    terminal-timer -c    
    ```

    If `$GOPATH/bin` or `$GOBIN` is not in your PATH, you can add it with the following command:

    ```sh
    echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bash_profile
    source ~/.bash_profile
    ```

    Replace `~/.bash_profile` with `~/.bashrc` or `~/.zshrc` depending on your shell and system configuration.

4. Now, you should be able to run `terminal-timer` from anywhere in your terminal.



### Option 3: Building from Source

For users who prefer to build the application from source:

1. Ensure you have Go installed on your system. You can check by running go version in your terminal.


   ```sh
   go version
    ```

2. Clone the repository:

    ```sh
    git clone https://github.com/yourusername/terminal-timer.git
    cd terminal-timer
    ```

3. Build the application:

   For example, on Unix-like systems, you might do:

   ```sh
   go build -o terminal-timer
    ```
4. Verify the installation by running:

    ```sh
    terminal-timer -c    
    ```

5. Optionally, move the timer executable to a directory in your PATH to run it from anywhere:
   
   ```sh
   mv terminal-timer /usr/local/bin/
    ```

6. Verify the installation by running:

    ```sh
    terminal-timer -c    
    ```


### Core Functionality

- **Timer (`-t`)**: Set the duration of your timer using the format `hh:mm`. For example, `-t 0:30` sets a 30-minute timer.
- **Alarm (`-a`)**: Specify an alarm time in a 24-hour format `hh:mm`. For instance, `-a 13:45` sets the alarm to trigger at 1:45 PM.
- **Reminder Message (`-r`)**: Customize the reminder message displayed when the timer ends. Default message is "Time is Up!".

### Customization Options

- **Font (`-f`)**: Set a new font for the timer display. Use `-f FontName` to change the font.
- **Preview Font (`-pf`)**: Preview how a font looks with `-pf FontName`.
- **List Valid Fonts (`-lf`)**: List all valid fonts available for the timer display.

### Sound Options

- **Preview Sound (`-ps`)**: Preview a specific sound with `-ps SoundName`.
- **List Valid Sounds (`-ls`)**: Display a list of all default sounds included in the application.
- **Set Sound (`-s`)**: Choose a new sound for timer notifications with `-s SoundName`.

### Logging and Configuration

- **Enable Logging (`-log`)**: Enable logging to a file for debugging or record-keeping purposes.
- **Show Current Config (`-c`)**: Display the current configuration for sound and font settings.

### Example Usage

To start a 25-minute timer with a custom reminder message and sound, you might use the following command:

```sh
terminal-timer -t 0:25 -r "Break Time!" 
```