package bundle

import "errors"

var ErrNotASharedObject = errors.New("this file is not a .so plugin")

var ErrNoBundleSymbol = errors.New("shared object does not contain a bundle")

var ErrPluginDoesntContainBundle = errors.New("bundle object couldn't be parsed as scry bundle")
