# selfupdate golang binary

Demo application for self updating golang binary

Instructions:
 * Use package.go Helper method to generate a .json file that will have meta data of executable binary
 * Use update.go Helper method to actually update the download binary with current executable

Steps:
 * once you have generated package.json upload it to target server, e.g. https://X.X/somepath/package.json
 * run a goroutine which will poll this address ever X minutes, e.g. 15
 * at every poll, unmarshal json and compare json version with current executable version
 * if its different(or greater), download the binary from https://X.X/somepath/myapplication to local machine
 * run UpdatePackage in update.go with downloaded binary path and current applicaiton path as arguments
