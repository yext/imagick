#!/bin/sh
#
# This script can be used on either Linux or MacOS to build ImageMagick and
# dependent static libraries.
#
# Notes on diagnosing problems:
#
# - If ImageMagick's ./configure shows something unexpectedly disabled, e.g. tiff:
#   (a) it can't find the delegate library, e.g. libtiff.a
#   (b) a dependency of that library can not be found, e.g. tiff depends on lzma
#   => Look in ./config.log to see exactly what the error was.
#
# - If you change an option to ./configure, but it doesn't seem to be taking
#   effect: "make clean"
#
# After running to completion, you can find all headers & libraries in the work
# directory.
#
# Delegates: mpeg jng jpeg lcms png ps tiff webp zlib

set -e

IMAGEMAGICK_URL=https://www.yext-static.com/cms/imagemagick/ImageMagick-6.8.8-7.tar.xz

OS=$(uname -s)
function _tar_xz() {
    if [ "$OS" == "Darwin" ]; then
        tar zxf "$@"
    else
        tar Jxf "$@"
    fi
}

# Create a temp directory to work in
IMBUILD=`echo ~/imagemagick-build`
mkdir $IMBUILD
echo "Work directory: $IMBUILD"

# Download ImageMagick
cd $IMBUILD
curl -O $IMAGEMAGICK_URL
_tar_xz *.xz
rm *.xz
cd ImageMagick*

##############################
# Delegates
##############################

export CPPFLAGS=-I$IMBUILD/include
export LDFLAGS=-L$IMBUILD/lib

echo CONFIGURE: zlib
curl -LO http://www.imagemagick.org/download/delegates/zlib-1.2.11.tar.xz
_tar_xz zlib*.xz && rm zlib*.xz && mv zlib* zlib && cd zlib
./configure --prefix=$IMBUILD --static
make install
cd ..

echo CONFIGURE: png
curl -LO http://www.imagemagick.org/download/delegates/libpng-1.6.31.tar.xz
_tar_xz libpng*.xz && rm libpng*.xz && mv libpng* png && cd png
./configure --prefix=$IMBUILD --disable-shared --disable-dependency-tracking
make install
cd ..

echo CONFIGURE: webp
curl -LO http://www.imagemagick.org/download/delegates/libwebp-0.6.0.tar.gz
tar zxf libwebp*.gz && rm libwebp*.gz && mv libwebp* webp && cd webp
./configure --prefix=$IMBUILD --disable-shared --disable-dependency-tracking
make install
cd ..

echo CONFIGURE: jpeg
curl -LO http://www.imagemagick.org/download/delegates/jpegsrc.v9b.tar.gz
tar zxf jpeg*.gz && rm jpeg*.gz && mv jpeg* jpeg && cd jpeg
./configure --prefix=$IMBUILD --disable-shared --disable-dependency-tracking
make install
cd ..

echo CONFIGURE: tiff
curl -LO http://www.imagemagick.org/download/delegates/tiff-4.0.8.tar.gz
tar zxf tiff*.gz && rm tiff*.gz && mv tiff* tiff && cd tiff
./configure --prefix=$IMBUILD --disable-shared --disable-dependency-tracking --disable-lzma
make install
cd ..

echo CONFIGURE: lcms2
curl -LO http://www.imagemagick.org/download/delegates/lcms2-2.8.tar.gz
tar zxf lcms2*.gz && rm lcms2*.gz && mv lcms* lcms && cd lcms
./configure --prefix=$IMBUILD --disable-shared --disable-dependency-tracking
make install
cd ..

##############################
# Build ImageMagick
##############################

./configure                           \
    --prefix $IMBUILD                 \
    --enable-static                   \
    --disable-shared                  \
    --disable-installed               \
    --disable-dependency-tracking     \
    --enable-delegate-build           \
    --without-frozenpaths             \
    --disable-docs                    \
    --enable-hdri=no                  \
    --without-modules                 \
    --disable-openmp                  \
    --with-threads                    \
    --without-magick-plus-plus        \
                                      \
    --with-bzlib=no                   \
    --with-autotrace=no               \
    --with-djvu=no                    \
    --with-dps=no                     \
    --with-fftw=no                    \
    --with-fpx=no                     \
    --with-fontconfig=no              \
    --with-freetype=no                \
    --with-gslib=no                   \
    --with-gvc=no                     \
    --with-jbig=no                    \
    --with-jpeg=yes                   \
    --with-lcms=yes                   \
    --with-lqr=no                     \
    --with-lzma=no                    \
    --with-openexr=no                 \
    --with-openjp2=no                 \
    --with-pango=no                   \
    --with-perl=no                    \
    --with-png=yes                    \
    --with-rsvg=no                    \
    --with-tiff=yes                   \
    --with-webp=yes                   \
    --with-wmf=no                     \
    --with-x=no                       \
    --with-xml=no                     \
    --with-zlib=no

make install
