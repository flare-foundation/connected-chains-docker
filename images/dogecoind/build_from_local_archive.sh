docker build --build-arg SOURCE_ZIP=dogecoin-1.14.7-flare.zip --build-arg SOURCE_FOLDER=dogecoin-1.14.7-flare --progress=plain -f Dockerfile-wls -t flarefoundation/dogecoin:1.14.7-wls .