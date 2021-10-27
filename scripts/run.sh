: ${IMAGE_NAME:=asssaf/automationhat:latest}
docker run --rm -it --privileged --device /dev/gpiomem "$IMAGE_NAME" $*
