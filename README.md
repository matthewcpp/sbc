# sbc
Reimplementation of N64 SGI Audio Tools [Sequence Bank Compiler](https://ultra64.ca/files/documentation/online-manuals/functions_reference_manual_2.0i/tools/sbc.html)

### Usage

```shell
sbc [-o output_file] seqfile0 seqfile1 seqfile2 
```

### Creating a Sequence Bank
1. Convert your sequences to MIDI type 0 using [midicvt](https://github.com/lambertjamesd/midicvt) by lambertjamesd.
1. Create a sequence bank using this tool.
1. Set the sequence Bank in your game:  
With nusys: [nuAuSeqPlayerSeqSet](https://ultra64.ca/files/documentation/online-manuals/man-v5-1/nusystem/nu_f/audio_sgi/nuAuSeqPlayerSeqSet.htm)  
With al lib: [alSeqFileNew](https://ultra64.ca/files/documentation/online-manuals/functions_reference_manual_2.0i/al/alSeqFileNew.html)