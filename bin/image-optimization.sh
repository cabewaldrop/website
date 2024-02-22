#!/bin/bash

# Based on Chris Titus' image optimization script https://christitus.com/script-for-optimizing-images/

# Dependencies
# - img-optimize - https://virtubox.github.io/img-optimize/
# - imagemagick
# - jpegoptim
# - optipng

FOLDER="/Users/cabewaldrop/input-images"

# max width
WIDTH=800

# max height
HEIGHT=600

#resize png or jpg to either height or width, keeps proportions using imagemagick
find ${FOLDER} -iname '*.jpg' -o -iname '*.png' -exec convert \{} -verbose -resize $WIDTHx$HEIGHT\> /Users/cabewaldrop/\{} \;
# img-optimize --std --path ${FOLDER}
