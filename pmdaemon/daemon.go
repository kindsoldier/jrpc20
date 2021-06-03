/*
 * Copyright 2019 Oleg Borodin  <borodin@unix7.org>
 */

package pmdaemon

import (
    "errors"
    "fmt"
    "io/fs"
    "io"
    "log"
    "os"
    "os/signal"
    "os/user"
    "path/filepath"
    "strconv"
    "syscall"
    "time"

)

type Daemon struct {
    username        string
    logFilename     string
    pidFilename     string
    debug           bool
    foreground      bool
}

const (
    piddirMode  fs.FileMode = 0750
    pidfileMode fs.FileMode = 0640

    logdirMode  fs.FileMode = 0750
    logfileMode fs.FileMode = 0640
)

func NewDaemon(logFilename, pidFilename string, debug, foreground bool) *Daemon {
    curUser, _ := user.Current()
    username := curUser.Username
    return &Daemon{
        username:           username,
        logFilename:        logFilename,     
        pidFilename:        pidFilename,
        debug:              debug,
        foreground:         foreground,
    }
}

func (this *Daemon) Daemonize() error {
    var err error

    if !this.foreground{
        forkProcess()
    }

    err = saveProcessID(this.pidFilename, piddirMode, pidfileMode)
    if err != nil {
        return errors.New(fmt.Sprintf("unable save process id: %s\n", err))
    }

    user, err := user.Lookup(this.username)
    if err != nil {
        return errors.New(fmt.Sprintf("user lookup error: %s\n", err))
    }
    uid, err := strconv.Atoi(user.Uid)

    /* Change effective user ID */
    //if uid != 0 {
    //    err = syscall.Setuid(uid)
    //    if err != nil {
    //        return errors.New(fmt.Sprintf("set process user id error: %s\n", err))
    //    }
    //    if syscall.Getuid() != uid {
    //        return errors.New(fmt.Sprintf("set process user id error: %s\n", err))
    //    }
    //}

    _, err = redirectLog(uid, this.logFilename, logdirMode, logfileMode, false)
    if err != nil {
        return errors.New(fmt.Sprintf("unable redirect log to message file: %s\n", err))
    }

    if !this.foreground {
        _, err := redirectIO();
        if err != nil {
            return errors.New(fmt.Sprintf("unable redirect stdio: %s\n", err))
        }
    }
    return nil
}

func SetSignalHandler() {
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGSTOP, syscall.SIGTERM, syscall.SIGQUIT)

    go func() {
        for {
            log.Printf("signal handler start")
            sig := <- sigs
            log.Printf("receive signal %s", sig.String())

            switch sig {
                case syscall.SIGINT, syscall.SIGTERM, syscall.SIGSTOP:
                    log.Printf("exit process by signal %s", sig.String())
                    time.Sleep(time.Millisecond * 100)
                    os.Exit(0)

                case syscall.SIGHUP:
                    log.Printf("restart program")
                    forkProcess()
            }
        }
    }()
}

func saveProcessID(filename string, piddirMode fs.FileMode, pidfileMode fs.FileMode) error {
    var err error

    err = os.MkdirAll(filepath.Dir(filename), piddirMode)
    if err != nil {
        return errors.New(fmt.Sprintf("unable create rundir: %s\n", err))
    }
    pid := os.Getpid()
    file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, pidfileMode)
    if err != nil {
         return err
    }
    defer file.Close()
    _, err = file.WriteString(strconv.Itoa(pid))
    if err != nil {
        return err
    }
    file.Sync()
    return nil
}

func redirectLog(uid int, filename string, logdirMode fs.FileMode, logfileMode fs.FileMode, shortFile bool) (*os.File, error) {
    var err error
    
    err = os.MkdirAll(filepath.Dir(filename), logdirMode)
    if err != nil {
        return nil, err
    }
    err = os.Chown(filepath.Dir(filename), uid, os.Getgid())
    if err != nil {
        return nil, err
    }

    file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, logfileMode)
    if err != nil {
        return nil, err
    }
    writer := io.MultiWriter(os.Stdout, file)
    //writer := io.Writer(file)
    if shortFile {
        log.SetFlags(log.LstdFlags | log.Lshortfile)
    } else {
        log.SetFlags(log.LstdFlags)
    }
    log.SetOutput(writer)
    return file, nil
}

func redirectIO() (*os.File, error) {
    file, err := os.OpenFile("/dev/null", os.O_RDWR, 0)
    if err != nil {
        return nil, err
    }
    syscall.Dup2(int(file.Fd()), int(os.Stdin.Fd()))
    syscall.Dup2(int(file.Fd()), int(os.Stdout.Fd()))
    syscall.Dup2(int(file.Fd()), int(os.Stderr.Fd()))
    return file, nil
}

func forkProcess() error {
    if _, exists := os.LookupEnv("GOGOFORK"); !exists {
        os.Setenv("GOGOFORK", "yes")

        cwd, err := os.Getwd()
        if err != nil {
            return err
        }

        procAttr := syscall.ProcAttr{}
        procAttr.Files = []uintptr{ uintptr(syscall.Stdin), uintptr(syscall.Stdout), uintptr(syscall.Stderr) }
        procAttr.Env = os.Environ()
        procAttr.Dir = cwd
        syscall.ForkExec(os.Args[0], os.Args, &procAttr)
        os.Exit(0)
    }
    _, err := syscall.Setsid()
    if err != nil {
        return err
    }
    os.Chdir("/")
    return nil
}


