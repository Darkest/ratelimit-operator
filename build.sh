export IMAGE=tanyavoteva/ratelimit-operator:v0.0.1
echo $IMAGE

echo "operator-sdk build ..."
operator-sdk build $IMAGE

echo "docker push ..."
docker push $IMAGE