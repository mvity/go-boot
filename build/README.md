# `/build`

打包和持续集成所需的文件。

* build/ci：存放持续集成的配置和脚本，如果持续集成平台(例如 Travis CI)对配置文件有路径要求，则可将其 link 到指定位置。
* build/package：存放 AMI、Docker、系统包（deb、rpm、pkg）的配置和脚本等。