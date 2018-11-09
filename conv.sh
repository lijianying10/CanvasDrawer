#!/bin/bash
export MAGICK_HOME="$1/ImageMagick-7.0.8"
export DYLD_LIBRARY_PATH="$1/ImageMagick-7.0.8/lib"
$1/ImageMagick-7.0.8/bin/convert $2 $3
