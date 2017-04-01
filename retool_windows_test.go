// +build windows

package main

// Go builds files on windows with an '.exe' suffix, so we need a few pieces of
// special logic to make sure things work there.
const osBinSuffix = ".exe"
