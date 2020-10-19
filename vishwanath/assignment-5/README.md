# Assignment-5


## Qemu-Img utility

This CLI utility is mainly used for building, managing, conversion and modification of images in an offline way. All the formats supported by QEMU can be handled by this utility. However, we should not use this utility while the image is being used by any process or VM as it can corrupt the image. This is sort of like the ".iso" image in a compressed format. It acts like a disk image utility for various file systems.

This is the syntax: `qemu-img command --flag`

These are the supported image formats by qemu-img:
* raw
* qcow2
* bochs
* cloop
* cow
* dmg
* nbd
* parallels
* qcow
* vdi
* vmdk
* vpc
* vvfat

Commands for qemu-img: (tried these too)

1. `qemu-img create [size] -f <fmt> <filename>` : This creates qemu images with the given name, size and format provided. We also have the ability to add further options to this command using `-o` flag.

2. `qemu-img amend -p [options] <filename>` : This changes the image format of the given filename with respect to the options provided.

3. `qemu-img bench -c [count] -s [buffer-size] <filename>` : This is used for benchmark of the image created with sequential I/O reads and writes as mentioned with `-w` flag. The count tells how many I/O operations must be performed with the size referred in buffer-size.

4. `qemu-img check -r <filename>` : This performs a consistency check on the given image as disk image. It outputs in the format QFMT which is usually a JSON format. The The JSON output is an object of QAPI type ImageCheck. However, when `-r` is specfied, qemu-img tries to repair any inconsistencies found during the check. If everything is fine, it returns with 0 meaning that the image is consistent and can be used (error-free). 2 tells that image is corrupted.

5. `qemu-img compare <file1> <file2>` : It checks if the given two images have the same content or not. Check can be performed with different formats or settings.

6. `qemu-img convert -l [snapshot-parameter] -o [output-format] <filename> <output-filename>` : This is used to convert a given disk image or a snapshot to a disk image while changing its format. However, it can be used for compression too by using `-c` flag. Only the formats `qcow` and `qcow2` support compression. The compression is read-only. It means that if a compressed sector is rewritten, then it is rewritten as uncompressed data. This conversion can also be useful in a use case where the image needs to be of smaller size than original (`qcow format`)

7. `qemu-img dd -f [fmt] -o [output-fmt] if=INPUT of=OUTPUT` : This copies the input file to output file while converting it from previous format to the given format.

8. `qemu-img info <filename>` : This gives entire information about the image including its size and the displayed size. We can get a JSON output using `--output json` flag.

### Useful Flags:
* -D is for diff
* -O is output format
* -B is backing path
* -f is for First image format
* -F is for Second image format

### Reason for using qemu-img

Qemu-img is directly compatible to most used image formats and performs operation on them effectively using same commands. It has the option of conversion, compression, management, testing, creation, updation, resizing and operations considering backing chain of images and snapshots. The package installation is easy. It ccan also be used for formatting guest images, additional storage devices and network storage.


## Qcow2 image format

`qcow` is a file format for disk image files used by QEMU, a hosted VM monitor. It stands for "QEMU Copy On Write" and uses a disk storage optimization strategy that delays allocation of storage until it is actually needed. Files in qcow format can contain a variety of disk images which are generally associated with specific guest OSs. Three versions of the format exist: qcow, qcow2 and qcow3 which use the .qcow, .qcow2 and .qcow3 file extensions.

It works like git commits.

One of the main characteristics of qcow disk images is that files with this format can grow as data is added. This allows for smaller file sizes than raw disk images, which allocate the whole image space to a file, even if parts of it are empty. This is particularly useful for file systems that do not support sparse files, such as FAT32.

`qemu-img` command allows to inspect, check, create, convert, resize and take snapshot of qcow images.

The qcow format also allows storing changes made to a read-only base image on a separate qcow file by using copy on write. This new qcow file contains the path to the base image to be able to refer back to it when required. When a particular piece of data has to be read from this new image, the content is retrieved from it if it is new and was stored there. If it is not, the data is fetched from the base image.

`qcow2` is an updated version of the q`cow` format and it supports AES encryption. The difference from the original version is that qcow2 supports multiple snapshots using a newer, more flexible model for storing them. A qcow2 image file is organized in units of constant size, which are called (host) clusters. A cluster is the unit in which all allocations are done, both for actual guest data and for image metadata.  All numbers in qcow2 are stored in Big Endian byte order.

If the image has a backing file then the backing file name should be stored in the remaining space between the end of the header extension area and the end of the first cluster.

There are certain options which can be used by `qemu-img` that are suitable while performing operations on `qcow2` images:

* `compat` : Determines the qcow2 version to use.
* `backing_file` : File name of a base image.
* `backing_fmt` : Image format of the base image.
* `encryption` : If this option is set to on, the image is encrypted with 128-bit AES-CBC.
* `cluster_size` : Changes the qcow2 cluster size (must be between 512 and 2M). 
* `lazy_refcounts` : If this option is set to on, reference count updates are postponed with the goal of avoiding metadata I/O and improving performance.

## Difference between Virtual Size and Disk Size

The Virtual Size and Disk Size refer to the static and dynamic size of the image (or container too). However, they both show the used size only.

**The actual difference is that when the image is kept as it is and will be loaded while starting a VM/container from it, that original size will be its Disk Size. It is a read-only layer which will be shared to every instance that will be run from it. Now when the image is loaded and an instance starts from it, a writable-layer is written on top of it as a representation and basic dependencies for the instance. This also holds any changes required while running it. This newly added size which has writable-layer is known as Virtual Size.**

## Restoring data from qcow2 images

1. Converting qcow2 image to raw image: `qemu-img convert -f qcow2 -O raw srcQcow2Path destPath`
2. Reading the raw file: `LESSOPEN= less file.raw`

## GuestFS for creating images

1. From Filesystem: `Disk_create()` will create a blank disk image. `Copy_in()` from `guestfs.go` will be responsible to copy local files or directories into an image.
2. From Block Device: `qemu-img convert -O %s %s -B %s %s", format, srcPath, prevQcow2Path, intermediateQcow2Path`.

## Backing Chain & Diff in Image

Modern disk image formats allow users to create an overlay on top of an existing image which will be the target of the new guest writes. This allows us to do snapshots of the disk state of a VM efficiently. This is also known as Backing Chain.

To get the entire backing chain in of one image: `qemu-img info --backing-chain $FILENAME`

Creating image with backing chain: `"qemu-img convert -O %s %s -B %s %s", format, srcPath, prevQcow2Path, intermediateQcow2Path`

Explicitly putting a chain to an image: `sudo qemu-img rebase -f qcow2 -b $NEW_BACKING_FILE $QCOW2_FILE_TO_CHANGE`

Creating a `Diff` image from 2 sources: `"qemu-img convert -f qcow2 -O qcow2 -D  %s -F qcow2 %s %s", prevPvQcow2Path, intermediateQcow2Path, finalQcow2Path` 


## Used Links:

* https://docs.fedoraproject.org/en-US/Fedora/18/html/Virtualization_Administration_Guide/sect-Virtualization-Tips_and_tricks-Using_qemu_img.html
* https://www.qemu.org/docs/master/interop/qemu-img.html
* https://en.wikipedia.org/wiki/Qcow
* https://git.qemu.org/?p=qemu.git;a=blob;f=docs/interop/qcow2.txt
* https://stackoverflow.com/questions/37966973/what-is-the-difference-between-the-size-and-the-virtual-size-of-the-docker-image
* https://unix.stackexchange.com/questions/208105/open-a-raw-file-as-text-in-less
* https://blog.programster.org/qemu-img-cheatsheet