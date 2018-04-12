# logrus-formatter
A custom text formater for logrus with debug info(filename, function name and line no.)
## Enable Debug Info
### Code
```
var log = logrus.New()
log.Out = os.Stdout
log.Formatter = &Formatter{Debug: true}
log.Info("This is a info")
log.Error("This is a error")
log.Warnf("This is a warnf")
log.Warningf("This is a warningf")
```
### Output
```
debug="[formatter_test.go][logrus-formatter.TestFormatter][16]" time="2018-04-12T11:58:37+08:00" level=info msg="This is a info"
debug="[formatter_test.go][logrus-formatter.TestFormatter][17]" time="2018-04-12T11:58:37+08:00" level=error msg="This is a error"
debug="[formatter_test.go][logrus-formatter.TestFormatter][18]" time="2018-04-12T11:58:37+08:00" level=warning msg="This is awarnf"
debug="[formatter_test.go][logrus-formatter.TestFormatter][19]" time="2018-04-12T11:58:37+08:00" level=warning msg="This is awarningf"
```
