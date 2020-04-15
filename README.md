# cbgo

cbgo implements Go bindings for [CoreBluetooth](https://developer.apple.com/documentation/corebluetooth?language=objc) central functionality.

## Documentation

For documentation, see the [CoreBluetooth docs](https://developer.apple.com/documentation/corebluetooth?language=objc).

Examples are in the `examples` directory.

## Scope

cbgo implements all central functionality that is supported in macOS 10.13.

## Naming

Function and type names in cbgo are intended to match the corresponding CoreBluetooth functionality as closely as possible.  There are a few (consistent) deviations:

* All cbgo identifiers start with a capital letter to make them public.
* Named arguments in CoreBluetooth functions are eliminated.

## Issues

cbgo makes no attempt to release CoreBluetooth objects allocated in the objective C code.  I don't anticipate this causing any issues during typical usage.  This could become noticeable if your process has a very long lifetime (months) or it interacts with hundreds of thousands of peripherals.
