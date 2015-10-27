# EUCLID

Key Features

    Configured and controlled through messages.

    Multiple monitors support.

    EWMH support.

    Hybrid tiling.


#Synopsis

euclid [-h|-v|-c CONFIG_PATH]


#Description


euclid is a split container partioning window manager, the main tiling manager for [scpwm](https://github.com/thrisp/scpwm) 

It is controlled and configured via [scpc](https://github.com/thrisp/scpwm/tree/master/scpc).


#Options
<!---
-h

    Print the synopsis and exit.

-v 

    verbose logging of messages to stdout

-version

    Print the version and exit.

-c CONFIG_PATH

    Use the given configuration file.
-->

#Configuration


euclid has only two sources of information: the X events it receives and the messages it receives from scpc.

The default configuration file is $XDG_CONFIG_HOME/euclid/euclidrc.

Keyboard and pointer bindings are defined with [keter](https://github.com/thrisp/scpwm/tree/master/keter).

<!---Example configuration files can be found in the examples directory.-->


#Splitting Modes

TBD


#Containers

Each monitor contains at least one desktop.

Each desktop contains at most one tree.


#Selectors


Selectors are used to select a target window, desktop, or monitor. A selector can either describe the target relatively or name it globally.

Descriptive (relative) selectors consist of a primary selector and any number of non-conflicting modifiers as follows:

PRIMARY_SELECTOR[.MODIFIER]\*

For obvious reasons, neither desktop nor monitor names may be valid descriptive selectors

###Window

Select a window

WINDOW_SELECTOR := <window_id> | (DIR|CYCLE_DIR|biggest|last|focused|older|newer)[.floating|.tiled][.like|.unlike][.manual|.automatic][.urgent][.local][.unfocused]

####Window States

floating

    Is above any tiled window and can be moved/resized freely. Although it doesn’t occupy any tiling space, it is still part of the window tree.

pseudo_tiled

    Has an unrestricted size while being centered in its tiling space.

fullscreen

    Fills its monitor rectangle, is above all the other windows and has no borders.

locked

    Ignores the close message.

sticky

    Stays in the focused desktop of its monitor.

private

    Tries to keep the same tiling position/size.

###Desktop

Select a desktop

DESKTOP_SELECTOR := <desktop_name> | [MONITOR_SEL:]^<n> | (CYCLE_DIR|last|[MONITOR_SELECTOR:]focused|older|newer)[.occupied|.free][.urgent][.local]

###Monitor

Select a monitor

MONITOR_SELECTOR := <monitor_name> | ^<n> | (DIR|CYCLE_DIR|last|primary|focused|older|newer)[.occupied|.free]

#Commands

###Window

######General Syntax

window *[WINDOW_SELECTOR] OPTIONS*

###Desktop

######General Syntax

desktop *[DESKTOP_SELECTOR] OPTIONS*

###Monitor

######General Syntax

monitor *[MONITOR_SELECTOR] OPTIONS*

###Query

######General Syntax

query *OPTIONS*

not yet implemented

###Pointer

######General Syntax

not yet implemented

###Control

######General Syntax

not yet implemented

###Restore

######General Syntax

not yet implemented

###Rule

######General Syntax

rule *OPTIONS*

######Options

*-a, --add <class_name>|<instance_name>| [-o|--one-shot] [monitor=MONITOR_SELECTOR|desktop=DESKTOP_SELECTOR|window=WINDOW_SELECTOR] [(floating|fullscreen|pseudo_tiled|locked|sticky|private|center|follow|manage|focus|border)=(on|off)] [split_dir=DIR] [split_ratio=RATIO]*

    Create a new rule.

*-r, --remove ^<n>|head|tail|<class_name>|<instance_name>*

    Remove the given rules.

*-l, --list [<class_name>|<instance_name>]*

    List the rules.

###Config

######General Syntax

config *[-m MONITOR_SELECTOR |-d DESKTOP_SELECTOR |-w WINDOW_SELELECTOR] <key> [<value>]*

Get or set the value of <key>.


#Settings

tbd

<!--
#Environment Variables
SCPWM_SOCKET

    The path of the socket used for the communication between scpc and euclid. 

    If it isn’t defined, then the following path is used:

    /tmp/scpwm<host_name>_<display_number>_<screen_number>-socket.
-->
