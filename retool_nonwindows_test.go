// +build !windows

package main

// Go builds files on windows with an '.exe' suffix. Everywhere else, there's no
// suffix.
const osBinSuffix = ""
