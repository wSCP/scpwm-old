# keter

## Name

keter - A simple X hotkey daemon


## Synopsis

keter [OPTIONS] [EXTRA_CONFIG...]


## Options

--help

Print the synopsis to standard output and exit.

--version

Print the version information to standard output and exit. (Version information
is derived from data passed to the "ldflags" flag with [go build](https://golang.org/pkg/go/build/))

--timeout

Timeout in seconds for the recording of chord chains. Default is 2 seconds.

--config

Read the keter main configuration from the given file

--verbose

Verbose logging of any messages, default is false.


## Behavior

keter is a daemon that listens to keyboard events and execute commands.

It reads its configuration file from $XDG_CONFIG_HOME/keter/keterrc by default, or from the given file if the --config option is used.

If keter receives a SIGUSR1 (resp. SIGUSR2) signal, it will reload its configuration file (resp. toggle the grabbing state of all its bindings).

The commands are executed via [exec.Cmd](https://golang.org/pkg/os/exec/) (you can use environment variables).


## Configuration

Each line of the configuration file is interpreted as so:

    If it is empty or starts with #, it is ignored.

    If it starts with a space, it is read as a command.

    Otherwise, it is read as a hotkey.

General syntax:

```
HOTKEY
    COMMAND

HOTKEY      := CHORD_1 ; CHORD_2 ; … ; CHORD_n
CHORD_i     := [MODIFIERS_i +] [@|;|~]KEYSYM_i
MODIFIERS_i := MODIFIER_i1 + MODIFIER_i2 + … + MODIFIER_ik
```

The valid modifier names are: super, hyper, meta, alt, control, ctrl, shift, mode_switch, lock, mod1, mod2, mod3, mod4, mod5 and any.

The keysym names are given by the output of **xev**.

Hotkeys and commands can be spread across multiple lines by ending each partial line with a backslash character.

When multiple chords are separated by semicolons, the hotkey is a chord chain: the command will only be executed after receiving each chord of the chain in consecutive order.

<!--
The colon character can be used instead of the semicolon to indicate that the chord chain shall not be aborted when the chain tail is reached.

If a command starts with a semicolon, it will be executed synchronously, otherwise asynchronously.

The Escape key can be used to abort a chord chain.
-->

If @ is added at the beginning of the keysym, the command will be run on key release events, otherwise on key press events.

<!--
If ~ is added at the beginning of the keysym, the captured event will be replayed for the other clients.
-->

Mouse hotkeys can be defined by using one of the following special keysym names: button1, button2, button3, …, button24.

The hotkey and the command may contain sequences of the form {STRING_1,…,STRING_N}.

In addition, the sequences can contain ranges of the form A-Z where A and Z are alphanumeric characters.

<!--
The underscore character represents an empty sequence element.
-->


###EXAMPLE BINDINGS

```
XF86Audio{Prev,Next}
    mpc -q {prev,next}

@XF86LaunchA
    scrot -s -e 'image_viewer $f'

super + shift + equal
    sxiv -rt "$HOME/image"

XF86LaunchB
    xdotool selectwindow | xsel -bi

super + {h,j,k,l}
    scpc window -f {left,down,up,right}

super + alt + {0-9}
    mpc -q seek {0-9}0%

super + {alt,ctrl,alt + ctrl} + XF86Eject
    sudo systemctl {suspend,reboot,poweroff}

super + button{1-3}
    scpc pointer -g {move,resize_side,resize_corner}

super + @button{1-3}
    scpc pointer -u

super + o ; {e,w,m}
    {gvim,firefox,thunderbird}

super + m ; {h,j,k,l}
    xdo move {-x -5,-y +5,-y -5,-x +5}

super + alt + control + {h,j,k,l} ; {0-9}
    scpc window -e {left,down,up,right} 0.{0-9}

super + alt + p
    scpc config focus_follows_pointer {true,false}
```

<!--
super + {_,shift + }{h,j,k,l}
    scpc window {-f,-s} {left,down,up,right}

{_,shift + ,super + }XF86MonBrightness{Down,Up}
    bright {-1,-10,min,+1,+10,max}

~button1
    scpc pointer -g focus
-->
