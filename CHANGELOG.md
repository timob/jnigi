# Change Log

## v2.0.0-alpha.1
- Add VERSION file
- Update dependencies
- Include more tests of package functions
- Update tests to use github.com/stretchr/testify
- Export ToGoArray, ToJavaArray methods
- Fixed bug with env.GetUTF8String()
- Updated JNI version constants
- Change to use int32 internally to represent Java int
- Fix loading DLL from paths with unicode characters on Windows
- Fixed Memory leak
- Add support for using a different class loader for FindClass, useful on Android
- 2021-12-05 Version 2: New idiomatic API. Converter interfaces. Add docs.

## v1
- 2020-12-09 Add go.mod file, updated import path to tekao.net/jnigi.
- 2020-08-21 Add ExceptionHandler interface for handling Java exceptions. Add 3 general handlers DescribeExceptionHandler (default), ThrowableToStringExceptionHandler and ThrowableErrorExceptionHandler.
- 2020-08-11 Add DestroyJavaVM support, JNI_VERSION_1_8 const
- 2019-05-29 Better multiplatform support, dynamic loading of JVM library.
- 2016-08-01 Initial version.
