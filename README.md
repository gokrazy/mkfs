# gokrazy mkfs

This program is intended to be run on gokrazy only, where it will create an ext4
file system on the perm partition and then reboot your system. If `/perm` is
already mounted, the program will exit without changing anything.

The gokrazy mkfs program includes a [frozen
copy](https://github.com/gokrazy/freeze) of the `mke2fs` program from the
`e2fsprogs` package from Debian.

## Usage

You can either add this program to your gokrazy instance:

```
gok add github.com/gokrazy/mkfs
```

…or, if you want to run it only once without otherwise including it in your
installation, you can use `gok run`:

```
git clone https://github.com/gokrazy/mkfs
cd mkfs
gok -i bakery run
```


