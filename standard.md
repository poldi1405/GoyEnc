# Standard as Specified by yenc.org/yenc-draft.1.3.txt

## Encoding

encoded := (input + 42) % 256

### Critical characters

ASCII 0x3D (`=`) indicates a critical character

such Characters are to be encoded using:

critical := (encoded + 64) % 256

Critical Characters (1.3)
	- 0x0   NULL
	- 0xA   LF
	- 0xD   CR
	- 0x3D  =

additionally for 1.2:
	- 0x9   TAB

> Careful writers of encoders will encode TAB (09h) SPACES (20h) 
> if they would appear in the first or last column of a line.
>
> Implementors who write directly to a TCP stream will care about 
> the doubling of dots in the first column - or also encode a DOT 
> in the first column.

if the escaping `=` happens to be the last character on the line the
then the escaped character (NULL for example) is to be written on the
same line leading to a total line-lenth of n+1

### Headers and Trailers

```
=ybegin line=128 size=123456 name=mybinary.dat
…
=yend size=123456
```

> (1.2) If [one of] the parameters "line=" "size=" "name=" [is] not present then
> the =ybegin might be part of a text-message with a discussion about
> yEnc. In such cases the decoder should assume that there is no binary. 

Must contain:
	- typical line length (might be +1 if critical character is present) (header)
	- size of the unencoded file
	- name of the fragmented file (header)

Can contain:
	- crc32 Checksum (trailer)

Notes:
	- the filename must be the last item on the header line
	- leading and trailing spaces will be cut 
	- filenames may contain non-US-ASCII-characters, control characters, and characters
   not supported by the current platform
	- the filename may be up to 256 characters long
	- the sizes of the header, trailer, and decoded binary have to be verified. if any one
   of them differs the fragment needs to be considered corrupt and a warning must be
   issued, the binary must be discarded (discarding is to be handled by clients using this
   library)
   
### Multipart Encoded Binaries

#### Headers and Trailers

```
=ybegin part=1 total=5 line=128 size=500000 name=mybinary.dat
=ypart begin=1 end=100000
…
=yend size=100000 part=1 pcrc32=abcdef12 

=ybegin part=5 line=128 size=500000 name=mybinary.dat
=ypart begin=400001 end=500000 
…
=yend size=100000 part=10 pcrc32=12a45c78 crc32=abcdef12 
```
