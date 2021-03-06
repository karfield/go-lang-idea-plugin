# Google Go language plugin for IntelliJ Idea

Google go language plugin is an attempt to build an outstanding IDE for
[Google Go language](http://golang.org) using Intellij IDEA.

[![Build Status](https://travis-ci.org/go-lang-plugin-org/go-lang-idea-plugin.png?branch=master)](https://travis-ci.org/go-lang-plugin-org/go-lang-idea-plugin)
[![Coverity Scan Build Status](https://scan.coverity.com/projects/3280/badge.svg)](https://scan.coverity.com/projects/3280)

[![Gitter](https://badges.gitter.im/Join Chat.svg)](https://gitter.im/go-lang-plugin-org/go-lang-idea-plugin?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

## TOC

+ [Current WIP](#current-wip)
+ [Links](#links)
+ [What it does](#what-it-does)
+ [How to use it](#how-to-use-it)
+ [Tutorials](#tutorials)
+ [Bugs](#bugs)
+ [EAP downloads](#eap-download)
+ [Authors](#authors)
+ [Contributing](#contributing)
+ [Thank you](#thank-you)

## Current WIP

We are working towards 0.9.16 release (https://github.com/go-lang-plugin-org/go-lang-idea-plugin/milestones/0.9.16) and
the main things we want to fix are:

+ [ ] fix default imports of packages where the last part of the import path differs from the package name
+ [ ] add a supported alpha/beta channel so we don't have to push people to dropbox
+ [ ] improve the documentation of what the plugin can/can't do
    + [ ] improve the usage docs for what it can do
+ [ ] various type inference fixes

## Links

+ Project pages on github:
    + <http://github.com/go-lang-plugin-org/go-lang-idea-plugin>
    + <http://wiki.github.com/go-lang-plugin-org/go-lang-idea-plugin/> (wiki)
+ Google Go language page:
    + <http://golang.org/>
+ Intellij IDEA homepage:
    + <http://www.jetbrains.com/idea/>
    + <http://plugins.intellij.net/plugin/5047> (plugin page at the intellij plugin repository).

## What it does

* Basic language parsing and highlighting
* Code folding
* Brace matching
* Comment/Uncomment (Single/Multiple line) support
* Go SDK (work with the latest release and on windows)
* File type icon
* Go application file and library generation.
* Auto completion of sdk package names and/or local application packages.
* Compilation of the go applications (supported semantics are similar to those of gobuild)
* Go To definition (for types) works across files and Go SDK
* Code formatting - experimental (disabled)
* Type name completion
* ColorsAndSettings page with a new color scheme
* Automatically add new line at end of file
* Force UTF-8 encoding for go files
* Go module type
* Go SDK indexing mode

## How to use it

* Download and install Intellij IDEA (Ultimate or Community edition) or any other IntelliJ IDE
that is build on the IDEA 133.326+ platform
* Open the Plugins installation page: File -> Settings -> Plugins -> Browse Repositories...
* Search for "golang"
* Right click on the proper plugin and install
* Download [latest release of the Google Go language](http://golang.org/doc/install.html).
* Build and install it.
* Open IDEA and create an empty Go project.
* Go to File -> Project Structure and select SDKs entry in the left column of the new window
* Add a new Google Go SDK by clicking the plus sign an choosing the appropriate SDK type. (See MacOS note below)
* After the SDK is defined go to the Modules entry and add a new google facet to your default module.
Select the proper sdk for the module.
* If you have only one directory in the ``` GOPATH ``` and you are creating a project inside that path
when you are working with packages that are part of the project you must still specify the whole import
path for them, not the relative one. Example:
    - ``` GOPATH ``` is: ``` /home/florin/go ```
    - the correct way to setup a project called ``` demogo ``` is: ``` /home/florin/go/src/github.com/dlsniper/demogo/ ```
    - new package is: ``` /home/florin/go/src/github.com/dlsniper/demogo/newpack ```
    - the correct import statement is: ``` github.com/dlsniper/demogo/newpack ``` not ``` newpack ```


Now you are ready to play with golang.

_NOTE:_ In MacOS you may need to press <kbd>&#8984;</kbd>+<kbd>&#8679;</kbd>+<kbd>G</kbd> in Finder dialog and enter `/usr/local/go` if Go in installed 
using package installer)

## Tutorials

* [Philip Andrew's Intellij Go Plugin tutorial](http://webapp.org.ua/dev/intellij-idea-and-go-plugin/)

## Bugs

If you found a bug, please report it at the Google Go plugin project's tracker
on GitHub: <http://github.com/go-lang-plugin-org/go-lang-idea-plugin/issues>.

When reporting a bug, please include the following:
- IDEA version
- OS version
- JDK version

Please ensure that the bug is not reported already.

## EAP download

There will always be at least one EAP downloadable on [@dlsniper's dropbox](https://www.dropbox.com/sh/kzcmavr2cmqqdqw/j8wjp8SdNH) built from master.
The folder may contain other builds as well such as issue prefixed builds that attempt to fix a specific issue and have not been merged to master yet.
They are meant to be a convenient way to quickly test upcoming changes without the hasle of building the plugin but it also means that you should take them with a grain of salt, they are after all, beta builds.
Also, after you've done testing, please remember to check the official plugin repository in order to ensure you always have the latest stable build or that there's no new release as automatic updates won't work with EAP releases of the plugin.

## Contributing

If you want to contribute to this effort you should read the following page:
[how to contribute](https://github.com/go-lang-plugin-org/go-lang-idea-plugin/blob/master/contributing.md).

## Authors

+ [Mihai Claudiu Toader (mtoader@gmail.com)](http://redeul.ro)
    + Language parsing and PSI building
    + Sdk indexing
    + Reference Resolving & Completions

+ [leojay (@leojay)](https://github.com/leojay)
    + folding improvement
    + refactorings
        + introduce variable
    + inspections
        + unused imports

+ [Alexandre Normand (@alexandre-normand)](https://github.com/alexandre-normand)
    + gomake integration

+ [Jhonny (@khronnuz)](https://github.com/khronnuz)
    +  gofmt integration
    +  simple Structure View for Go files.

+ [Florin Patan (@dlsniper)](https://github.com/dlsniper)
    +  Cardea support
    +  Non-IntelliJ IDEs plugin support
    +  GDB debug support integration

+ [Chris Spencer (@spencercw)](https://bitbucket.org/spencercw/ideagdb/overview)
    + IDEA GDB support

## Thank you

We would like to thank [YourKit](http://www.yourkit.com) for supporting the development
of this plugin by providing us licenses to their full-featured [Yourkit Java Profiler](http://www.yourkit.com/java/profiler/index.jsp).
YourKit, LLC is the creator of innovative and intelligent tools for profiling
Java and .NET applications. Take a look at YourKit's leading software products:

- [YourKit Java Profiler](http://www.yourkit.com/java/profiler/index.jsp)
- [YourKit .NET Profiler](http://www.yourkit.com/.net/profiler/index.jsp)
