package main


import (
    "os"
    "fmt"
    "runtime"
    "syscall"
)


func Daemon(pfile, lfile string, uid, gid, nochdir, noclose int) int {
    if syscall.Getppid() == 1 {
        return 0
    }

    ret, mret, errno := syscall.RawSyscall(syscall.SYS_FORK, 0, 0, 0)
    if errno != 0 {
        return -1
    }

    if mret < 0 {
        os.Exit(-1)
    }

    if runtime.GOOS == "darwin" && mret == 1 {
        ret = 0
    }

    if ret > 0 {
        os.Exit(0)
    }

    _ = syscall.Umask(0)
    s_ret, err := syscall.Setsid()
    if err != nil {
        panic(err)
    }

    if s_ret < 0 {
        return -1
    }

    err = WriteFile(pfile, fmt.Sprintf("%d\n", os.Getpid()))
    if err != nil {
        panic(err)
    }

    err = syscall.Setgid(gid)
    if err != nil {
        panic(err)
    }

    err = syscall.Setuid(uid)
    if err != nil {
        panic(err)
    }

    if nochdir == 0 {
        os.Chdir("/")
    }

    if noclose == 0 {
        fp, err := os.OpenFile(lfile, os.O_RDWR | os.O_APPEND, 0)
        if err != nil {
            panic(err)
        }

        fd := fp.Fd()
        syscall.Dup2(int(fd), int(os.Stdin.Fd()))
        syscall.Dup2(int(fd), int(os.Stdout.Fd()))
        syscall.Dup2(int(fd), int(os.Stderr.Fd()))
    }

    return 0
}
