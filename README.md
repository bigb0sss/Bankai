<p align="center">
    <br>
        <img src=img/bankai.jpg >
    <br>
</p>

# Bankai
Another Go shellcode loader designed to work with Cobalt Strike raw payload. I created this project to mainly educate myself learning Go and directly executing shellcode into the target Windows system using various techniques. 

<b>Encryption</b> - I implemented a simple payload encryption process (IV --> AES --> XOR --> Base64) that I learned while studying [SLAE32](https://bigb0ss.medium.com/expdev-custom-go-crypter-fb8f9bac0fe8). This is mainly for protecting Cobalt Strike payload when it's moved over to the target host. The final payload will include a decrypt function within.

<b>Templates</b> - Templates are the skeleton scripts to generate a final payload per each technique. 

## Installation
```
git clone https://github.com/bigb0sss/bankai.git
GO111MODULE=off go build bankai.go
```

## Usage & Example
Generate a Cobalt Strike payload:

<p align="center">
    <br>
        <img src=img/cobalt.png>
    <br>
</p>

```
./bankai -h                       

     _                 _         _ 
    | |               | |       (_)
    | |__   __ _ _ __ | | ____ _ _ 
    | '_ \ / _' | '_ \| |/ / _' | |
    | |_) | (_| | | | |   < (_| | |
    |_.__/ \__,_|_| |_|_|\_\__,_|_|
                        [bigb0ss]

    [INFO] Another Go Shellcode Loader

 
    Required:
    -i            Binary File (e.g., beacon.bin)
    -o            Payload Output (e.g, payload.exe)
    -t            Payload Template (e.g., win32_VirtualProtect.tmpl)
    -a            Arch (32|64)
    
    Optional:
    -h            Print this help menu
    -p            PID

    Templates:                                     Last update: 06/02/21
    +-----------------------------------------------+------------------+
    | Techniques                                    | Bypass Defender  |
    +-----------------------------------------------+------------------+
    | win32_VirtualProtect.tmpl                     |        No        |
    +-----------------------------------------------+------------------+
    | win64_CreateFiber.tmpl                        |        No        |
    +-----------------------------------------------+------------------+
    | win64_CreateRemoteThreadNative.tmpl           |        Yes       | 
    +-----------------------------------------------+------------------+
    | win64_CreateThread.tmpl                       |        No        | 
    +-----------------------------------------------+------------------+
    | win64_EtwpCreateEtwThread.tmpl                |        No        | 
    +-----------------------------------------------+------------------+
    | win64_Syscall.tmpl                            |        No        | 
    +-----------------------------------------------+------------------+

    Example:

    ./bankai -i beacon.bin -o payload.exe -t win64_CreateThread.tmpl -a 64

    [INFO] Key: SymE9GQBtyHL4IAq5Pm6r3b8I7PJB9l0
    [INFO] AES encrpyting the payload...
    [INFO] Arch: x64 (64-bit)
    [INFO] Template: win64_CreateThread.tmpl
    [INFO] InputFile: beacon.bin
    [INFO] OutputFile: payload.exe
```

## Credits / Acknowledgments / References
All of the work is inspired and done by the following researchers/projects:
* [go-shellcode](https://github.com/brimstone/go-shellcode) by brimstone
* [go-shellcode](https://github.com/Ne0nd0g/go-shellcode) by Ne0nd0g
* [GoPurple](https://github.com/sh4hin/GoPurple) by sh4hin
* Go Template - https://dev.to/kirklewis/go-text-template-processing-181d

## Todo
* Add more shellcode injection technique templates
* Add [AlternativeShellcodeExec](https://github.com/S4R1N/AlternativeShellcodeExec) techniques that Ali and Alfaro found





